# 📅 Módulo de Integração Google Calendar (OAuth 2.0 Multi-Tenant BYOK)

Este documento descreve a arquitetura, modelo de dados, camada de segurança e criptografia, fluxo assíncrono via Asynq, catálogo de endpoints REST e interface administrativa da integração opcional e multi-tenant com o **Google Calendar (Google Agenda)** utilizando a filosofia **BYOK (*Bring Your Own Key*)**.

---

## 1. 📌 Visão Geral do Módulo

O módulo permite que cada estabelecimento comercial (`Tenant`) conecte sua própria agenda Google Calendar para sincronizar agendamentos de forma unidirecional (Agende Fácil → Google Calendar) e assíncrona via fila Asynq.

### Principais Pilares Arquiteturais:
1. **Filosofia BYOK (*Bring Your Own Key*):** A plataforma não utiliza uma conta centralizada ou Service Account compartilhada. Cada estabelecimento configura seu próprio `Client ID` e `Client Secret` obtidos no seu Google Cloud Console, garantindo isolamento total de cotas de API, limites de requisições e soberania de dados.
2. **Criptografia em Repouso (AES-256-GCM):** O `Client Secret`, `Access Token` e `Refresh Token` são gravados no banco de dados PostgreSQL cifrados com chave simétrica de 32 bytes (`ENCRYPTION_KEY`), vetor de inicialização (nonce) de 12 bytes gerado via `crypto/rand`, e marcados com `json:"-"` nos modelos para evitar vazamento em payloads de API ou logs.
3. **Segurança OAuth 2.0 com Anti-Tamper & Anti-Replay:**
   - O parâmetro `state` é assinado criptograficamente via **HMAC-SHA256** utilizando `JWT_SECRET`.
   - Um nonce aleatório é armazenado de forma atômica no Redis com TTL de 10 minutos e consumido exclusivamente via `GetDel` (uso único), inviabilizando ataques CSRF e replay.
   - Operação **fail-closed**: caso o Redis ou `OAuthStateStore` esteja inacessível, o fluxo OAuth é bloqueado preventivamente.
4. **Centralização e Unificação de Agendamentos:**
   - Todas as criações, reagendamentos e cancelamentos originados pela Web Pública (`/agendamento/:slug`), Painel Administrativo (`/admin/appointments`) e Atendimento IA via WhatsApp (`AIChatService`) convergem no `AppointmentService`.
   - O `AppointmentService` verifica preventivamente se a integração do tenant está ativa e conectada antes de enfileirar tarefas assíncronas no Asynq, garantindo higiene de fila.
5. **Reconciliação Idempotente pelo Estado Real no Banco:**
   - A tarefa assíncrona (`task:sync_google_calendar`) avalia o estado atual do agendamento persistido no PostgreSQL. Se `apt.Status == StatusCancelled`, efetua `DELETE` idempotente no Google Calendar (tratando `404 Not Found` e `410 Gone` como sucesso).
   - Para agendamentos confirmados ou pendentes, realiza `Upsert` idempotente via `UpdateEvent` ou `CreateEvent` utilizando `sendUpdates=none` para não disparar notificações abusivas.
   - Caso a autorização seja revogada ou expire (`invalid_grant`), o worker marca a conexão como `EXPIRED`, aborta retentativas desnecessárias com `asynq.SkipRetry` e expõe o erro no painel para re-autenticação.
6. **Interface Vue 3 (`GoogleCalendarView.vue`):**
   - Rota `/admin/google-calendar` com painel completo: inserção de credenciais BYOK, cópia facilitada da Redirect URI autorizada, geração de link de consentimento Google, seleção de calendário gravável (com fallback automático), status de conexão e painel de histórico/retry manual de falhas.

---

## 2. 🗄️ Modelo de Dados (GORM / PostgreSQL)

```mermaid
erDiagram
    TENANT ||--o| TENANT_GOOGLE_CALENDAR_CONFIG : "possui"
    TENANT ||--o{ APPOINTMENT_CALENDAR_SYNC : "possui"
    APPOINTMENT ||--o| APPOINTMENT_CALENDAR_SYNC : "vincula"

    TENANT_GOOGLE_CALENDAR_CONFIG {
        uuid id PK
        uuid tenant_id FK "unique"
        boolean is_enabled
        string status "DISCONNECTED | CONNECTED | EXPIRED | ERROR"
        string google_client_id
        text google_client_secret_encrypted "AES-256-GCM / json:-"
        string google_account_email
        text access_token_encrypted "AES-256-GCM / json:-"
        text refresh_token_encrypted "AES-256-GCM / json:-"
        datetime token_expires_at
        string calendar_id "default: primary"
        string calendar_summary
        datetime last_sync_at
        text last_error
        datetime created_at
        datetime updated_at
    }

    APPOINTMENT_CALENDAR_SYNC {
        uuid id PK
        uuid tenant_id FK "index / composite unique"
        uuid appointment_id FK "index / composite unique"
        string provider "GOOGLE_CALENDAR / composite unique"
        string external_calendar_id
        string external_event_id "index"
        string sync_status "PENDING | SYNCED | FAILED | DELETED"
        string last_action "CREATE | UPDATE | DELETE"
        int attempt_count
        text last_error
        datetime last_synced_at
        datetime created_at
        datetime updated_at
    }
```

---

## 3. 🔐 Arquitetura de Segurança, OAuth 2.0 & Criptografia

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Admin do Tenant
    participant Web as Painel Vue 3 (/admin/google-calendar)
    participant API as Backend (Gin / GoogleCalendarService)
    participant Redis as Redis (OAuthStateStore)
    participant Google as Google OAuth 2.0

    Admin->>Web: Salva Client ID e Client Secret (BYOK)
    Web->>API: PUT /api/v1/admin/integrations/google-calendar/settings
    API-->>Web: 200 OK (Credenciais cifradas com AES-256-GCM)

    Admin->>Web: Clica em "Conectar Google Calendar"
    Web->>API: GET /api/v1/admin/integrations/google-calendar/auth-url
    Note over API,Redis: Gera payload {tenant_id, user_id, nonce, exp}<br/>Assina HMAC-SHA256 e grava nonce no Redis (TTL 10m)
    API->>Redis: SET gcal:oauth:state:<nonce> <payload> EX 600
    API-->>Web: { "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?..." }
    Web->>Google: Redireciona usuário para consentimento Google

    Admin->>Google: Autoriza permissão de calendário
    Google->>API: GET /api/v1/integrations/google-calendar/callback?code=...&state=...
    Note over API,Redis: Valida assinatura HMAC-SHA256<br/>Executa GetDel atômico no Redis (uso único / anti-replay)
    API->>Redis: GETDEL gcal:oauth:state:<nonce>
    API->>Google: Troca code por tokens (access_token, refresh_token)
    Google-->>API: Tokens OAuth 2.0
    Note over API: Criptografa tokens com AES-256-GCM<br/>Consulta e-mail e lista calendários graváveis
    API-->>Web: Redireciona para /admin/google-calendar?status=success
```

---

## 4. ⚡ Fluxo de Sincronização Assíncrona & Reconciliação (Asynq)

```mermaid
flowchart TD
    A[Agendamento: Web Pública / Admin / WhatsApp IA] --> B[AppointmentService]
    B --> C[(PostgreSQL: Appointments)]
    B --> D{Integração Google ativa & conectada?}
    D -- Não --> E[Fim do Fluxo - Sem enqueue]
    D -- Sim --> F[Enfileira task:sync_google_calendar na fila default]
    F --> G[Worker Asynq]
    G --> H[GoogleCalendarService.SyncAppointment]
    H --> I[(Consulta Estado Atual no PostgreSQL)]
    I --> J{Status do Agendamento?}
    J -- CANCELLED --> K[DELETE Event no Google Calendar]
    K --> L[Trata 404/410 como sucesso idempotente]
    L --> M[(Atualiza AppointmentCalendarSync: DELETED)]
    J -- PENDING ou CONFIRMED --> N[Upsert Event no Google Calendar com sendUpdates=none]
    N --> O[(Atualiza AppointmentCalendarSync: SYNCED)]
    H -- Erro 401/invalid_grant --> P[Marca Config como EXPIRED e dispara asynq.SkipRetry]
```

---

## 5. 🔄 Tarefas Assíncronas (Asynq + Redis)

| Tipo de Tarefa | Fila | Payload | Comportamento |
|---|---|---|---|
| `task:sync_google_calendar` | `default` | `{"tenant_id": "uuid", "appointment_id": "uuid", "action": "CREATE\|UPDATE\|DELETE"}` | Sincroniza criação, reagendamento e cancelamento de eventos no Google Calendar do tenant de forma idempotente e não-bloqueante |

---

## 6. 🚀 Catálogo de Endpoints REST do Módulo

### 6.1 Callback Público
- `GET /api/v1/integrations/google-calendar/callback`
  - **Descrição:** Recebe redirecionamento da Google após o consentimento do usuário.
  - **Parâmetros Query:** `code`, `state` (e `error` em caso de recusa).
  - **Comportamento:** Valida a assinatura HMAC-SHA256, executa `GetDel` no Redis, troca o `code` pelos tokens OAuth, cifra com AES-256-GCM, obtém o e-mail da conta Google conectada e redireciona o navegador para `${FRONTEND_URL}/admin/google-calendar?status=success`.

### 6.2 Gestão da Integração (`ADMIN_TENANT` & `ADMIN_GLOBAL`)
*Protegidos por autenticação JWT e Gatekeeper de assinatura (`RequireActiveSubscription`).*

- `GET /api/v1/admin/integrations/google-calendar/status`
  - Retorna o status da conexão (`DISCONNECTED`, `CONNECTED`, `EXPIRED`, `ERROR`), `google_client_id`, flags booleanas (`has_client_secret`, `has_credentials`, `is_enabled`), e-mail da conta conectada, ID e nome do calendário ativo, Redirect URI e contador de falhas pendentes.
- `GET /api/v1/admin/integrations/google-calendar/auth-url`
  - Gera a URL de consentimento OAuth 2.0 personalizada com o `state` anti-CSRF e anti-replay.
- `GET /api/v1/admin/integrations/google-calendar/calendars`
  - Lista os calendários graváveis (`owner` ou `writer`) da conta Google conectada para seleção pelo tenant.
- `PUT /api/v1/admin/integrations/google-calendar/settings`
  - Permite configurar/atualizar `google_client_id`, `google_client_secret`, `calendar_id`, `calendar_summary` e habilitar/desabilitar a sincronização (`is_enabled`).
- `POST /api/v1/admin/integrations/google-calendar/disconnect`
  - Desconecta a conta Google (revoga o token na API do Google e limpa tokens/e-mail no banco), preservando as credenciais BYOK para facilitar reconexão futura.
- `GET /api/v1/admin/integrations/google-calendar/failures`
  - Retorna o histórico de agendamentos que falharam na sincronização com detalhes de erro e número de tentativas.
- `POST /api/v1/admin/integrations/google-calendar/syncs/:appointment_id/retry`
  - Reenfileira manualmente a sincronização de um agendamento específico no Asynq.

---

## 7. 💻 Interface Administrativa no Frontend (Vue 3)

A interface em `frontend/src/views/admin/GoogleCalendarView.vue` (rota `/admin/google-calendar`) foi desenhada com Tailwind CSS e Lucide Icons, oferecendo:
1. **Painel de Credenciais BYOK:** Campos para inserção de `Client ID` e `Client Secret`, botão para visualização de senha e botão para copiar com um clique a `Redirect URI` autorizada para inclusão no Google Cloud Console.
2. **Controle de Sessão & Conexão:** Botão de login OAuth ("Conectar com Google"), indicador visual de status (Conectado, Desconectado, Expirado ou Erro) e opção de desconectar.
3. **Seleção de Calendário:** Dropdown que carrega via API os calendários disponíveis na conta do usuário, com opção de fallback e preservação do calendário primário (`primary`).
4. **Histórico de Falhas & Ação de Retry:** Tabela dinâmica listando falhas recentes de sincronização com mensagem detalhada de erro, data/hora da tentativa e botão para retry instantâneo.

---

## 8. ⚙️ Variáveis de Ambiente e Configurações

| Variável | Padrão | Descrição |
|---|---|---|
| `APP_PUBLIC_URL` | `http://localhost:8080` | URL pública base da API backend para montagem de callbacks |
| `FRONTEND_URL` | `http://localhost:5173` | URL pública do frontend Vue 3 para onde o callback redireciona |
| `GOOGLE_OAUTH_REDIRECT_URI` | `${APP_PUBLIC_URL}/api/v1/integrations/google-calendar/callback` | Redirect URI padrão registrada na Google Cloud Console |
| `ENCRYPTION_KEY` | *(32 bytes hex)* | Chave mestra AES-256-GCM para cifrar secrets e tokens em repouso |
| `JWT_SECRET` | *(string secreta)* | Chave utilizada para assinar o parâmetro `state` via HMAC-SHA256 |
