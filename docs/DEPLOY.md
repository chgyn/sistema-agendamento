# Deploy, operação e manutenção (VPS + Docker + Caddy)

Complemento de [INSTALL.md](./INSTALL.md). Use este documento no dia a dia da VPS.

Todos os comandos abaixo pressupõem o diretório raiz do projeto e, em geral, `sudo` quando o usuário não pertence ao grupo `docker`.

---

## 1. Arquitetura de produção

```text
Internet :80 / :443
        │
        ▼
     Caddy  (TLS Let's Encrypt, headers, gzip)
        │
        ├── /api/*  → backend:8080  (Gin, AutoMigrate + seed)
        └── /       → frontend:80   (Nginx SPA)
                        │
        postgres ◄──────┤
        redis    ◄──────┤ worker (Asynq)
        wuzapi   ◄──────┘ (rede interna somente)
```

Volumes persistentes:

| Volume | Conteúdo |
|---|---|
| `postgres_data` | Banco PostgreSQL |
| `redis_data` | Filas / AOF Redis |
| `wuzapi_data` | Sessões WhatsApp |
| `caddy_data` | Certificados ACME |
| `caddy_config` | Configuração automática do Caddy |

---

## 8. Caddy, domínio e HTTPS

O site address do Caddy é a variável `DOMAIN`. O e-mail ACME é `ACME_EMAIL`.

Arquivos:

- `Caddyfile` — HTTPS automático, redirect HTTP→HTTPS
- `Caddyfile.http` — somente porta 80 (`CADDYFILE_PATH=./Caddyfile.http`)

Trocar o domínio:

1. Atualize o DNS (registro A/AAAA).
2. Edite `.env`: `DOMAIN`, `ACME_EMAIL`, `APP_PUBLIC_URL`, `CORS_ORIGIN`.
3. `docker compose up -d caddy backend`

Webhook Asaas (cadastre no painel Asaas):

```text
https://SEU_DOMINIO/api/v1/webhooks/asaas
```

O webhook do WhatsApp (WUZAPI) usa a rede Docker (`http://backend:8080/api/v1/webhooks/whatsapp`) e **não** precisa ser público.

---

## 9. Status dos containers

```bash
docker compose ps
docker compose ps --format 'table {{.Name}}\t{{.Status}}\t{{.Ports}}'
```

Healthcheck da API (interno):

```bash
docker compose exec backend wget -qO- http://127.0.0.1:8080/api/v1/health
```

Pela internet (após DNS + TLS):

```bash
curl -fsS https://SEU_DOMINIO/api/v1/health
```

---

## 10. Logs

```bash
docker compose logs -f --tail=200
docker compose logs -f backend
docker compose logs --since 1h worker
docker compose logs -f caddy
docker compose logs --tail=200 wuzapi
```

---

## 11. Reinício

```bash
docker compose restart
docker compose restart backend worker frontend
docker compose up -d --force-recreate caddy
```

Parada completa (preserva volumes):

```bash
docker compose down
```

**Não** use `docker compose down -v` em produção: apaga banco, Redis, sessões WhatsApp e certificados.

---

## 12. Como atualizar a aplicação

1. Faça backup:

   ```bash
   ./scripts/backup.sh
   ```

2. Atualize o código:

   ```bash
   git fetch --all
   git pull
   ```

3. Recrie imagens e containers (volumes permanecem):

   ```bash
   docker compose up -d --build
   ```

4. Confirme:

   ```bash
   docker compose ps
   docker compose logs --tail=50 backend
   curl -fsS https://SEU_DOMINIO/api/v1/health
   ```

A API executa `AutoMigrate` na subida. Não há passo separado de migration.

Se o frontend não refletir mudança de `VITE_API_URL`, force rebuild sem cache:

```bash
docker compose build --no-cache frontend
docker compose up -d frontend
```

---

## 13. Rollback

Os volumes **não** são revertidos pelo `git checkout`. Escolha o cenário:

### Só código (bug na versão nova)

```bash
./scripts/backup.sh
git log --oneline -n 20
git checkout <commit-ou-tag-estavel>
docker compose up -d --build
```

### Código + dados (restore do banco)

1. Volte o código para a tag compatível com o dump.
2. Restaure:

   ```bash
   ./scripts/restore.sh backups/<DATA>/postgres.sql.gz
   ```

Sessões WUZAPI e certificados Caddy ficam nos volumes; em geral **não** precisam de rollback.

---

## 14. Backup e recuperação

### Backup

```bash
./scripts/backup.sh
```

Gera `backups/<timestamp>/`:

- `postgres.sql.gz` — dump SQL
- `env.copy` — cópia do `.env` (restrinja permissões; contém segredos)
- `compose-status.txt` — status dos containers no momento do backup

Copie o diretório para fora da VPS (S3, outro servidor, máquina local):

```bash
scp -r backups/<timestamp> usuario@backup-host:~/agendamento-backups/
```

Frequência sugerida: diária + imediata antes de cada atualização.

### Restore

```bash
./scripts/restore.sh backups/<timestamp>/postgres.sql.gz
```

Confirme digitando `restaurar`. Use `--force` apenas em automação.

O restore **não** reverte Redis, WUZAPI nem Caddy. Para desastre total: reinstale com `./install.sh` (preserva `.env` se existir), recrie a stack e restaure o dump.

---

## 15. Troubleshooting

### Certificado SSL não é emitido

- `dig +short SEU_DOMINIO` deve retornar o IP da VPS.
- Portas 80 e 443 abertas no provedor **e** no UFW.
- Nada além do Caddy escutando 80/443 (`ss -ltnp | grep -E ':80|:443'`).
- Logs: `docker compose logs caddy`.
- Cloudflare: modo DNS only.
- Let's Encrypt **não** emite certificado para IP puro — use um hostname.

### Site não abre / 502

```bash
docker compose ps
docker compose logs --tail=80 backend frontend caddy
docker compose exec backend wget -qO- http://127.0.0.1:8080/api/v1/health
```

502 costuma ser backend ainda inicializando (migrate/seed) ou unhealthy. Aguarde o healthcheck ou veja logs.

### Build falha por memória

Sintomas: `Killed`, `signal: killed`, `gocache`/`npm` abortando.

- Use VPS com 4 GB ou adicione swap (1–2 GB) temporariamente.
- Não compile outras coisas em paralelo.

### Containers “unhealthy”

```bash
docker compose ps
docker inspect agendamento_api --format '{{json .State.Health}}'
docker compose logs backend
```

Postgres recusa senha: o volume antigo foi criado com outra `POSTGRES_PASSWORD`. Não altere a senha no `.env` de uma instalação existente sem rotacionar no banco, ou suba um volume novo (isso **apaga** dados).

### Redis AUTH failed

`REDIS_PASSWORD` no `.env` deve ser o mesmo usado na criação do volume Redis. O instalador não rotaciona senha Redis automaticamente em reexecuções se o valor já não for `CHANGE_ME`.

### WhatsApp / WUZAPI

WUZAPI não tem porta pública. Pareamento é feito pelo painel da aplicação. Persistência: volume `wuzapi_data`. Se o QR nunca avança, `docker compose logs -f wuzapi backend worker`.

### Asaas não confirma pagamento

- URL do webhook: `https://SEU_DOMINIO/api/v1/webhooks/asaas`
- `ASAAS_WEBHOOK_SECRET` igual ao cadastrado no Asaas
- Chave com `$` inicial escrita como `$$` no `.env`
- Worker em execução: `docker compose logs -f worker`

### Reexecução do install.sh

Segura. Não recria segredos válidos nem apaga volumes. Útil após falha de rede no `docker pull`.

### Espaço em disco

```bash
docker system df
docker image prune
```

Não execute `docker volume prune` em produção sem conferir os nomes dos volumes.

---

## Checklist pós-install

1. `https://SEU_DOMINIO` abre o frontend.
2. Login com `admin@plataforma.com` e **troca da senha**.
3. `https://SEU_DOMINIO/api/v1/health` retorna `{"status":"ok",...}`.
4. DNS estável; cadeado TLS no navegador.
5. Webhook Asaas cadastrado (se for usar assinaturas).
6. Backup inicial copiado para fora da VPS.
