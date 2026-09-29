#!/usr/bin/env bash
# =============================================================================
# Sistema de Agendamento — instalador idempotente para VPS Linux + Docker
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

DOMAIN_ARG=""
EMAIL_ARG=""
NON_INTERACTIVE=0
CONFIGURE_FIREWALL=0
HTTP_ONLY=0
SKIP_DOCKER_INSTALL=0

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log()  { printf "%b\n" "${BLUE}[INFO]${NC} $*"; }
ok()   { printf "%b\n" "${GREEN}[OK]${NC} $*"; }
warn() { printf "%b\n" "${YELLOW}[AVISO]${NC} $*"; }
err()  { printf "%b\n" "${RED}[ERRO]${NC} $*" >&2; }

usage() {
  cat <<'EOF'
Uso: sudo ./install.sh [opções]

Opções:
  --domain DOMINIO          Domínio público (ex: app.seudominio.com.br)
  --email EMAIL             E-mail para Let's Encrypt (ACME)
  --http-only               Sobe o Caddy apenas em HTTP (:80), sem certificado
  --non-interactive         Não pergunta; exige --domain (ou --http-only)
  --configure-firewall      Libera 22/80/443 no UFW (se o ufw estiver instalado)
  --skip-docker-install     Não instala Docker; falha se não existir
  -h, --help                Exibe esta ajuda

O script é idempotente: pode ser executado novamente sem apagar volumes
nem sobrescrever segredos já definidos no arquivo .env.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --domain)
      DOMAIN_ARG="${2:-}"
      shift 2
      ;;
    --email)
      EMAIL_ARG="${2:-}"
      shift 2
      ;;
    --http-only)
      HTTP_ONLY=1
      shift
      ;;
    --non-interactive)
      NON_INTERACTIVE=1
      shift
      ;;
    --configure-firewall)
      CONFIGURE_FIREWALL=1
      shift
      ;;
    --skip-docker-install)
      SKIP_DOCKER_INSTALL=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      err "Opção desconhecida: $1"
      usage
      exit 1
      ;;
  esac
done

need_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    err "Execute com sudo: sudo ./install.sh"
    exit 1
  fi
}

have_cmd() {
  command -v "$1" >/dev/null 2>&1
}

upsert_env() {
  local key="$1"
  local value="$2"
  local file="${3:-.env}"
  local tmp
  tmp="$(mktemp)"
  if [[ ! -f "$file" ]]; then
    printf '%s=%s\n' "$key" "$value" > "$file"
    rm -f "$tmp"
    return
  fi
  awk -v k="$key" -v v="$value" '
    BEGIN { found=0 }
    index($0, k "=") == 1 { print k "=" v; found=1; next }
    { print }
    END { if (!found) print k "=" v }
  ' "$file" > "$tmp"
  cat "$tmp" > "$file"
  rm -f "$tmp"
}

get_env() {
  local key="$1"
  local file="${2:-.env}"
  [[ -f "$file" ]] || return 0
  grep -E "^${key}=" "$file" | tail -n 1 | cut -d= -f2- || true
}

is_placeholder() {
  local value="${1:-}"
  [[ -z "$value" || "$value" == "CHANGE_ME" || "$value" == "app.seudominio.com.br" || "$value" == "admin@seudominio.com.br" || "$value" == "https://app.seudominio.com.br" ]]
}

rand_hex() {
  local bytes="${1:-32}"
  openssl rand -hex "$bytes"
}

ensure_secret() {
  local key="$1"
  local bytes="${2:-32}"
  local current
  current="$(get_env "$key")"
  if is_placeholder "$current"; then
    upsert_env "$key" "$(rand_hex "$bytes")"
    ok "Segredo gerado: $key"
  else
    log "Mantendo $key já definido no .env"
  fi
}

detect_os() {
  if [[ ! -f /etc/os-release ]]; then
    err "Não foi possível detectar o sistema operacional (/etc/os-release)."
    exit 1
  fi
  # shellcheck disable=SC1091
  . /etc/os-release
  OS_ID="${ID:-unknown}"
  OS_VERSION="${VERSION_ID:-unknown}"
  log "Sistema detectado: ${PRETTY_NAME:-$OS_ID $OS_VERSION}"
}

check_architecture() {
  local arch
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64|aarch64|arm64)
      ok "Arquitetura suportada: $arch"
      ;;
    *)
      warn "Arquitetura $arch não foi testada. O Docker pode falhar."
      ;;
  esac
}

check_resources() {
  local mem_kb disk_kb
  mem_kb="$(awk '/MemTotal/ {print $2}' /proc/meminfo)"
  disk_kb="$(df -Pk "$SCRIPT_DIR" | awk 'NR==2 {print $4}')"

  if [[ "${mem_kb:-0}" -lt 1500000 ]]; then
    err "Memória insuficiente (${mem_kb} kB). Mínimo recomendado: 2 GB."
    exit 1
  fi
  if [[ "${mem_kb}" -lt 3500000 ]]; then
    warn "Menos de 4 GB de RAM: o build das imagens pode ser lento ou falhar."
  else
    ok "Memória: $((mem_kb / 1024)) MB"
  fi

  if [[ "${disk_kb:-0}" -lt 8000000 ]]; then
    err "Espaço em disco insuficiente. Livre pelo menos 10 GB no diretório do projeto."
    exit 1
  fi
  ok "Espaço livre: $((disk_kb / 1024)) MB"
}

port_listening() {
  local port="$1"
  if have_cmd ss; then
    ss -ltn 2>/dev/null | awk '{print $4}' | grep -Eq "[:.]${port}$"
    return $?
  fi
  if have_cmd lsof; then
    lsof -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1
    return $?
  fi
  return 1
}

check_ports() {
  local port
  for port in 80 443; do
    if port_listening "$port"; then
      if docker ps --format '{{.Names}} {{.Ports}}' 2>/dev/null | grep -q "agendamento_caddy"; then
        log "Porta $port já está em uso pelo Caddy deste projeto (reexecução)."
      else
        warn "Porta $port já está em uso. O Caddy precisa dela para HTTP/HTTPS."
      fi
    fi
  done
}

install_prereqs() {
  if have_cmd apt-get; then
    export DEBIAN_FRONTEND=noninteractive
    apt-get update -y
    apt-get install -y ca-certificates curl gnupg openssl git
  elif have_cmd dnf; then
    dnf install -y ca-certificates curl openssl git
  elif have_cmd yum; then
    yum install -y ca-certificates curl openssl git
  else
    warn "Gerenciador de pacotes não reconhecido. Certifique-se de ter curl, openssl e git."
  fi
}

compose_cmd() {
  if docker compose version >/dev/null 2>&1; then
    echo "docker compose"
  elif have_cmd docker-compose; then
    echo "docker-compose"
  else
    echo ""
  fi
}

install_docker() {
  if have_cmd docker && [[ -n "$(compose_cmd)" ]]; then
    ok "Docker e Compose já estão instalados."
    docker --version
    $(compose_cmd) version
    return
  fi

  if [[ "$SKIP_DOCKER_INSTALL" -eq 1 ]]; then
    err "Docker/Compose não encontrados e --skip-docker-install foi informado."
    exit 1
  fi

  case "$OS_ID" in
    ubuntu|debian)
      log "Instalando Docker Engine + Compose (script oficial)..."
      curl -fsSL https://get.docker.com | sh
      ;;
    *)
      if [[ -n "$(compose_cmd)" ]]; then
        return
      fi
      err "Instalação automática do Docker é suportada em Ubuntu/Debian."
      err "Instale o Docker manualmente e execute o script novamente."
      exit 1
      ;;
  esac

  systemctl enable --now docker >/dev/null 2>&1 || true

  if [[ -n "${SUDO_USER:-}" && "${SUDO_USER}" != "root" ]]; then
    usermod -aG docker "$SUDO_USER" || true
    log "Usuário ${SUDO_USER} adicionado ao grupo docker (novo login pode ser necessário)."
  fi

  if ! have_cmd docker || [[ -z "$(compose_cmd)" ]]; then
    err "Falha ao instalar Docker Compose. Verifique o log acima."
    exit 1
  fi
  ok "Docker instalado."
}

prompt_config() {
  local current_domain current_email
  current_domain="$(get_env DOMAIN)"
  current_email="$(get_env ACME_EMAIL)"

  if [[ -n "$DOMAIN_ARG" ]]; then
    DOMAIN_VALUE="$DOMAIN_ARG"
  elif ! is_placeholder "$current_domain"; then
    DOMAIN_VALUE="$current_domain"
  elif [[ "$HTTP_ONLY" -eq 1 ]]; then
    DOMAIN_VALUE="${current_domain:-}"
  elif [[ "$NON_INTERACTIVE" -eq 1 ]]; then
    err "Modo --non-interactive exige --domain ou --http-only."
    exit 1
  else
    read -r -p "Domínio público (ex: app.seudominio.com.br). Vazio para HTTP-only: " DOMAIN_VALUE
  fi

  if [[ -z "$DOMAIN_VALUE" ]]; then
    HTTP_ONLY=1
    DOMAIN_VALUE="localhost"
  fi

  if [[ -n "$EMAIL_ARG" ]]; then
    EMAIL_VALUE="$EMAIL_ARG"
  elif ! is_placeholder "$current_email"; then
    EMAIL_VALUE="$current_email"
  elif [[ "$HTTP_ONLY" -eq 1 ]]; then
    EMAIL_VALUE="admin@localhost"
  elif [[ "$NON_INTERACTIVE" -eq 1 ]]; then
    EMAIL_VALUE="admin@${DOMAIN_VALUE}"
  else
    read -r -p "E-mail para Let's Encrypt [${current_email:-admin@$DOMAIN_VALUE}]: " EMAIL_VALUE
    EMAIL_VALUE="${EMAIL_VALUE:-${current_email:-admin@$DOMAIN_VALUE}}"
  fi
}

prepare_env() {
  mkdir -p backups

  if [[ ! -f .env ]]; then
    cp .env.example .env
    ok "Arquivo .env criado a partir de .env.example"
  else
    log "Arquivo .env existente reutilizado (segredos preservados)."
  fi

  prompt_config

  ensure_secret POSTGRES_PASSWORD 24
  ensure_secret REDIS_PASSWORD 24
  ensure_secret JWT_SECRET 32
  ensure_secret WUZAPI_ADMIN_TOKEN 24

  local enc
  enc="$(get_env ENCRYPTION_KEY)"
  if is_placeholder "$enc" || [[ "${#enc}" -ne 64 ]]; then
    upsert_env ENCRYPTION_KEY "$(rand_hex 32)"
    ok "Segredo gerado: ENCRYPTION_KEY"
  else
    log "Mantendo ENCRYPTION_KEY já definido no .env"
  fi

  local pg_pass pg_user pg_db
  pg_pass="$(get_env POSTGRES_PASSWORD)"
  pg_user="$(get_env POSTGRES_USER)"
  pg_db="$(get_env POSTGRES_DB)"
  upsert_env DB_PASSWORD "$pg_pass"
  upsert_env POSTGRES_USER "${pg_user:-postgres}"
  upsert_env POSTGRES_DB "${pg_db:-sistema_agendamento}"
  upsert_env DB_USER "${pg_user:-postgres}"
  upsert_env DB_NAME "${pg_db:-sistema_agendamento}"
  upsert_env COMPOSE_PROFILES "prod"
  upsert_env APP_ENV "production"
  upsert_env TZ "America/Sao_Paulo"
  upsert_env WUZAPI_BASE_URL "http://wuzapi:8080"
  upsert_env REDIS_HOST "redis"
  upsert_env DB_HOST "postgres"
  upsert_env VITE_API_URL "/api/v1"
  upsert_env ACME_EMAIL "$EMAIL_VALUE"
  upsert_env DOMAIN "$DOMAIN_VALUE"

  if [[ "$HTTP_ONLY" -eq 1 ]]; then
    upsert_env HTTP_ONLY "true"
    upsert_env CADDYFILE_PATH "./Caddyfile.http"
    local public_ip=""
    public_ip="$(curl -4 -fsS --max-time 5 https://ifconfig.me 2>/dev/null || true)"
    if [[ -n "$public_ip" ]]; then
      upsert_env APP_PUBLIC_URL "http://${public_ip}"
    else
      upsert_env APP_PUBLIC_URL "http://${DOMAIN_VALUE}"
    fi
    upsert_env CORS_ORIGIN "*"
    warn "Modo HTTP-only: certificado SSL não será emitido."
  else
    upsert_env HTTP_ONLY "false"
    upsert_env CADDYFILE_PATH "./Caddyfile"
    upsert_env APP_PUBLIC_URL "https://${DOMAIN_VALUE}"
    upsert_env CORS_ORIGIN "https://${DOMAIN_VALUE}"
  fi
}

configure_firewall() {
  if [[ "$CONFIGURE_FIREWALL" -ne 1 ]]; then
    log "Firewall não alterado. Libere 22/tcp, 80/tcp e 443/tcp no provedor/UFW."
    return
  fi
  if ! have_cmd ufw; then
    warn "UFW não encontrado; pulando --configure-firewall."
    return
  fi
  ufw allow OpenSSH >/dev/null 2>&1 || ufw allow 22/tcp >/dev/null 2>&1 || true
  ufw allow 80/tcp >/dev/null 2>&1 || true
  ufw allow 443/tcp >/dev/null 2>&1 || true
  ufw --force enable >/dev/null 2>&1 || true
  ok "UFW: 22, 80 e 443 liberados."
}

start_stack() {
  local dc
  dc="$(compose_cmd)"
  log "Construindo e iniciando containers..."
  COMPOSE_PROFILES=prod $dc up -d --build --remove-orphans
}

wait_healthy() {
  local dc i
  dc="$(compose_cmd)"
  log "Aguardando containers saudáveis (até 3 minutos)..."
  for i in $(seq 1 36); do
    if $dc ps --format json >/dev/null 2>&1; then
      local unhealthy
      unhealthy="$($dc ps --format '{{.Name}} {{.Status}}' | grep -E 'unhealthy|Exit' || true)"
      if $dc ps --format '{{.Status}}' | grep -q 'healthy'; then
        if [[ -z "$unhealthy" ]] && $dc exec -T backend wget -qO- http://127.0.0.1:8080/api/v1/health >/dev/null 2>&1; then
          ok "API respondeu em /api/v1/health"
          return 0
        fi
      fi
    fi
    sleep 5
  done
  warn "Timeout aguardando healthcheck. Verifique: $(compose_cmd) ps && $(compose_cmd) logs --tail=80"
  return 1
}

print_summary() {
  local dc public_url
  dc="$(compose_cmd)"
  public_url="$(get_env APP_PUBLIC_URL)"

  echo
  echo "=============================================================="
  echo " Instalação do Sistema de Agendamento"
  echo "=============================================================="
  $dc ps
  echo
  echo " URL pública:          ${public_url}"
  echo " Domínio configurado:  $(get_env DOMAIN)"
  echo " Webhook Asaas:        ${public_url}/api/v1/webhooks/asaas"
  echo
  echo " Login inicial (troque a senha imediatamente):"
  echo "   E-mail: admin@plataforma.com"
  echo "   Senha:  admin123"
  echo
  echo " Comandos úteis:"
  echo "   $(compose_cmd) ps"
  echo "   $(compose_cmd) logs -f --tail=100"
  echo "   $(compose_cmd) restart"
  echo "   ./scripts/backup.sh"
  echo
  if [[ "$HTTP_ONLY" -ne 1 ]]; then
    echo " DNS: o registro A (e AAAA, se houver IPv6) de $(get_env DOMAIN)"
    echo " deve apontar para o IP desta VPS. O Caddy emite o certificado"
    echo " automaticamente quando o domínio resolver corretamente."
  fi
  echo "=============================================================="
}

main() {
  need_root

  if [[ ! -f docker-compose.yml ]]; then
    err "Execute este script na raiz do repositório (docker-compose.yml não encontrado)."
    exit 1
  fi

  detect_os
  check_architecture
  check_resources
  install_prereqs
  check_ports
  install_docker
  prepare_env
  configure_firewall
  start_stack
  wait_healthy || true
  print_summary
}

main "$@"
