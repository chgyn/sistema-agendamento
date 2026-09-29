# Instalação em VPS Linux

Este guia descreve como instalar o Sistema de Agendamento em uma **VPS Linux** com Docker, Caddy, domínio próprio e HTTPS automático.

Fluxo esperado:

**VPS Linux → `install.sh` → Docker/Compose → Aplicação → Caddy → Domínio + HTTPS**

Para operação diária (logs, atualização, backup e troubleshooting), use [DEPLOY.md](./DEPLOY.md).

---

## 1. Requisitos mínimos da VPS

| Recurso | Mínimo | Recomendado |
|---|---|---|
| vCPU | 2 | 2–4 |
| RAM | 2 GB | 4 GB |
| Disco | 20 GB SSD | 40 GB SSD |
| Rede | IPv4 público | IPv4 + IPv6 |
| Portas | 22, 80 e 443 liberadas no firewall do provedor | Idem |

O build das imagens (Go + Node) consome memória. Com 2 GB o processo pode ser lento ou falhar; prefira 4 GB.

---

## 2. Sistemas operacionais suportados

Instalação automática do Docker (via `install.sh`):

- Ubuntu 22.04 LTS
- Ubuntu 24.04 LTS
- Debian 12

Outras distros Linux (AlmaLinux, Rocky, etc.) são suportadas **se o Docker Engine e o plugin Compose já estiverem instalados**. Nesse caso:

```bash
sudo ./install.sh --skip-docker-install --domain app.seudominio.com.br --email admin@seudominio.com.br
```

Arquiteturas: `amd64` (`x86_64`) e `arm64` (`aarch64`).

---

## 3. Como obter / clonar o projeto

Na VPS:

```bash
sudo apt-get update && sudo apt-get install -y git
git clone <URL_DO_REPOSITORIO> sistema-agendamento
cd sistema-agendamento
```

O script **deve** ser executado na raiz do repositório (onde está o `docker-compose.yml`).

---

## 4. Como executar o script de instalação

```bash
sudo chmod +x install.sh scripts/*.sh
sudo ./install.sh
```

O instalador pergunta o **domínio** e o **e-mail** do Let's Encrypt, gera senhas, instala Docker se necessário, sobe os containers e valida o healthcheck da API.

### Modo não interativo

```bash
sudo ./install.sh \
  --domain app.seudominio.com.br \
  --email admin@seudominio.com.br \
  --non-interactive
```

### HTTP apenas (IP, sem certificado)

Útil para testes quando ainda não há DNS:

```bash
sudo ./install.sh --http-only --non-interactive
```

### Firewall UFW (opcional)

Por padrão o script **não** altera o firewall (evita cortar o SSH). Para liberar 22/80/443 no UFW:

```bash
sudo ./install.sh --domain app.seudominio.com.br --email admin@seudominio.com.br --configure-firewall
```

Mesmo assim, libere as mesmas portas no **painel do provedor** (Security Group / Firewall).

O script é **idempotente**: executá-lo de novo não apaga volumes nem sobrescreve segredos já gravados no `.env`.

---

## 5. Informações que precisam ser configuradas

| Item | Obrigatório | Como informar |
|---|---|---|
| Domínio | Sim (produção) | `--domain` ou prompt do instalador |
| E-mail ACME | Sim (HTTPS) | `--email` |
| Chaves Asaas | Recomendado | `.env` → `ASAAS_API_KEY`, `ASAAS_WEBHOOK_SECRET` |
| DNS A/AAAA | Sim (HTTPS) | Painel do registrador de domínio |

Gerados automaticamente pelo instalador (se ainda forem `CHANGE_ME`):

- `POSTGRES_PASSWORD`
- `REDIS_PASSWORD`
- `JWT_SECRET`
- `ENCRYPTION_KEY` (64 hex)
- `WUZAPI_ADMIN_TOKEN`

---

## 6. Variáveis de ambiente

O arquivo canônico de produção é **`.env` na raiz** (nunca versionado). O modelo está em `.env.example`.

```bash
# já criado pelo install.sh; edição posterior:
sudo nano .env
sudo docker compose up -d
```

Variáveis principais:

```text
DOMAIN=app.seudominio.com.br
ACME_EMAIL=admin@seudominio.com.br
APP_PUBLIC_URL=https://app.seudominio.com.br
CORS_ORIGIN=https://app.seudominio.com.br

POSTGRES_PASSWORD=...
REDIS_PASSWORD=...
JWT_SECRET=...
ENCRYPTION_KEY=...          # 64 caracteres hexadecimais
WUZAPI_ADMIN_TOKEN=...

ASAAS_API_KEY=...
ASAAS_BASE_URL=https://api.asaas.com/api/v3   # produção
ASAAS_WEBHOOK_SECRET=...
```

Chaves Asaas que começam com `$` devem ser escritas com `$$` no `.env` para o Compose não interpolar, por exemplo:

```text
ASAAS_API_KEY=$$aact_prod_sua_chave
```

Depois de alterar o `.env`:

```bash
sudo docker compose up -d
```

O arquivo `backend/.env` é apenas para execução **local sem Docker**. Em VPS use o `.env` da raiz.

---

## 7. Como configurar o domínio no DNS

1. Obtenha o IPv4 da VPS (`curl -4 ifconfig.me` ou o painel do provedor).
2. No DNS do domínio, crie um registro **A**:
   - Nome/host: `app` (ou `@` se for o domínio raiz)
   - Valor: IP da VPS
   - TTL: 300 (enquanto testa)
3. Opcional IPv6: registro **AAAA**.
4. Aguarde a propagação (`dig +short app.seudominio.com.br`).
5. Aponte o domínio **antes** ou logo após o `install.sh`. O Caddy só emite o certificado quando o hostname resolve para a VPS e as portas 80/443 estão abertas.

Não use Cloudflare com proxy laranja (proxied) até validar o certificado; comece em **DNS only** (nuvem cinza).

---

## 8. Como o Caddy funciona neste projeto

O serviço `caddy` (perfil Compose `prod`) é o **único** ponto de entrada público:

- Portas **80** e **443** (e 443/UDP para HTTP/3)
- HTTPS automático via Let's Encrypt
- Renovação automática (volume `caddy_data`)
- Redirecionamento HTTP → HTTPS
- Headers: HSTS, `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`

Roteamento (`Caddyfile`):

- `https://$DOMAIN/api/*` → container `backend:8080`
- `https://$DOMAIN/` → container `frontend:80` (SPA Vue)

PostgreSQL, Redis, WUZAPI e a API **não** são publicados na internet.

Para trocar o domínio depois:

```bash
sudo nano .env   # DOMAIN, ACME_EMAIL, APP_PUBLIC_URL, CORS_ORIGIN
sudo docker compose up -d caddy backend
```

Atualize o DNS para o mesmo IP. O Caddy solicitará um certificado novo.

Modo HTTP (`Caddyfile.http`, `HTTP_ONLY=true`): escuta apenas `:80`, sem Let's Encrypt.

---

## 9. Verificar o status dos containers

```bash
cd /caminho/sistema-agendamento
sudo docker compose ps
```

Esperado: `postgres`, `redis`, `backend`, `worker`, `wuzapi`, `frontend` e `caddy` em execução; `backend` e `frontend` **healthy**.

---

## 10. Visualizar logs

```bash
sudo docker compose logs -f --tail=100
sudo docker compose logs -f backend worker
sudo docker compose logs -f caddy
```

---

## 11. Reiniciar os serviços

```bash
sudo docker compose restart
sudo docker compose restart backend worker
```

---

## 12. Atualizar a aplicação

Ver o procedimento completo em [DEPLOY.md](./DEPLOY.md#12-como-atualizar-a-aplicação).

Resumo:

```bash
cd /caminho/sistema-agendamento
sudo ./scripts/backup.sh
git pull
sudo docker compose up -d --build
```

---

## 13. Rollback

Ver [DEPLOY.md](./DEPLOY.md#13-rollback). Em resumo: volte o código (`git checkout <tag>`) e, se necessário, restaure o dump com `./scripts/restore.sh`.

---

## 14. Backup e recuperação

```bash
sudo ./scripts/backup.sh
sudo ./scripts/restore.sh backups/<DATA>/postgres.sql.gz
```

Detalhes em [DEPLOY.md](./DEPLOY.md#14-backup-e-recuperação).

---

## 15. Troubleshooting

Problemas frequentes (DNS, certificado, memória no build, `.env`) estão em [DEPLOY.md](./DEPLOY.md#15-troubleshooting).

---

## Desenvolvimento local (sem Caddy)

```bash
cp .env.example .env
# substitua todos os CHANGE_ME por senhas locais
# defina COMPOSE_PROFILES=  (vazio) para não subir o Caddy
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

Acessos locais:

- Frontend: http://localhost:5173
- API: http://localhost:8080/api/v1/health
- Postgres host: `localhost:5433`
- WUZAPI: http://localhost:8081

---

## Credenciais iniciais (seed)

Troque imediatamente após o primeiro login.

| Perfil | E-mail | Senha |
|---|---|---|
| Admin global | `admin@plataforma.com` | `admin123` |
| Demo Dom Navalha | `admin@domnavalha.com` | `admin123` |
| Demo Bella Vista | `admin@bellavista.com` | `admin123` |
