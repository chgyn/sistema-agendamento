#!/usr/bin/env bash
# Restaura um dump PostgreSQL gerado por scripts/backup.sh
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

FORCE=0
DUMP_PATH=""

get_env() {
  local key="$1"
  grep -E "^${key}=" .env | tail -n 1 | cut -d= -f2- || true
}

usage() {
  cat <<'EOF'
Uso: ./scripts/restore.sh <arquivo.sql.gz|arquivo.sql> [--force]

Restaura o banco da aplicação a partir de um dump. Os volumes do Redis,
WUZAPI e certificados do Caddy NÃO são revertidos por este script.

--force  não pede confirmação interativa
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --force)
      FORCE=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      DUMP_PATH="$1"
      shift
      ;;
  esac
done

if [[ -z "$DUMP_PATH" ]]; then
  usage
  exit 1
fi

if [[ ! -f "$DUMP_PATH" ]]; then
  echo "Arquivo não encontrado: $DUMP_PATH" >&2
  exit 1
fi

if [[ ! -f .env ]]; then
  echo "Arquivo .env não encontrado." >&2
  exit 1
fi

if docker compose version >/dev/null 2>&1; then
  DC="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  DC="docker-compose"
else
  echo "Docker Compose não encontrado." >&2
  exit 1
fi

POSTGRES_USER="$(get_env POSTGRES_USER)"
POSTGRES_PASSWORD="$(get_env POSTGRES_PASSWORD)"
POSTGRES_DB="$(get_env POSTGRES_DB)"

if [[ "$FORCE" -ne 1 ]]; then
  echo "ATENÇÃO: esta operação substitui os dados atuais do PostgreSQL."
  read -r -p "Digite 'restaurar' para continuar: " confirm
  if [[ "$confirm" != "restaurar" ]]; then
    echo "Cancelado."
    exit 1
  fi
fi

echo "Restaurando dump em ${POSTGRES_DB:-sistema_agendamento}..."

if [[ "$DUMP_PATH" == *.gz ]]; then
  gzip -dc "$DUMP_PATH" | $DC exec -T -e PGPASSWORD="${POSTGRES_PASSWORD:-}" postgres psql \
    -U "${POSTGRES_USER:-postgres}" \
    -d "${POSTGRES_DB:-sistema_agendamento}"
else
  $DC exec -T -e PGPASSWORD="${POSTGRES_PASSWORD:-}" postgres psql \
    -U "${POSTGRES_USER:-postgres}" \
    -d "${POSTGRES_DB:-sistema_agendamento}" < "$DUMP_PATH"
fi

echo "Restore concluído. Reiniciando API e worker..."
$DC restart backend worker
echo "Pronto."
