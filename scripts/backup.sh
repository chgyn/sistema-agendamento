#!/usr/bin/env bash
# Backup lógico do PostgreSQL e metadados de volumes Docker.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

get_env() {
  local key="$1"
  grep -E "^${key}=" .env | tail -n 1 | cut -d= -f2- || true
}

if [[ ! -f .env ]]; then
  echo "Arquivo .env não encontrado na raiz do projeto." >&2
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
DOMAIN="$(get_env DOMAIN)"
APP_PUBLIC_URL="$(get_env APP_PUBLIC_URL)"

STAMP="$(date +%Y%m%d-%H%M%S)"
BACKUP_DIR="${ROOT_DIR}/backups/${STAMP}"
mkdir -p "$BACKUP_DIR"

echo "Gerando dump PostgreSQL..."
$DC exec -T -e PGPASSWORD="${POSTGRES_PASSWORD:-}" postgres pg_dump \
  -U "${POSTGRES_USER:-postgres}" \
  -d "${POSTGRES_DB:-sistema_agendamento}" \
  --no-owner --no-acl \
  | gzip > "${BACKUP_DIR}/postgres.sql.gz"

echo "Copiando .env (guarde este diretório com acesso restrito)..."
cp .env "${BACKUP_DIR}/env.copy"
chmod 600 "${BACKUP_DIR}/env.copy"

{
  echo "timestamp=${STAMP}"
  echo "domain=${DOMAIN:-}"
  echo "app_public_url=${APP_PUBLIC_URL:-}"
  $DC ps
} > "${BACKUP_DIR}/compose-status.txt"

echo "Backup concluído em ${BACKUP_DIR}"
echo "Arquivos:"
ls -lh "$BACKUP_DIR"
