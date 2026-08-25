# 🏛️ Arquitetura, Concorrência e Catálogo de APIs da Plataforma

Documentação técnica unificada sobre a arquitetura Clean Architecture, modelo de isolamento multi-tenant, motor de concorrência com locks pessimistas, mensageria Asynq e catálogo completo de APIs.

---

## 1. 🏗️ Arquitetura do Backend (Go 1.24+)

O projeto adota os princípios de **Clean / Hexagonal Architecture**:

```text
[ HTTP Controllers / Handlers (Gin) ]
                 │
                 ▼
[ Service Layer (Business Logic & Availability Engine) ]
                 │
                 ▼
[ Domain Layer (Models, Ports & Pure Interfaces) ]
                 ▲
                 │
[ Repositories (PostgreSQL / GORM) ] & [ Integrations (Asaas, Asynq/Redis) ]
```

- **Domain (`internal/domain`):** Entidades puras e interfaces de contratos (`ports.go`). Não possui dependências externas.
- **Service (`internal/service`):** Regras de negócio, validações de disponibilidade, cálculo de horários livres, geração de agendamentos e orquestração de assinaturas.
- **Repository (`internal/repository/postgres`):** Implementação concreta de acesso ao banco com GORM, transações e `SELECT ... FOR UPDATE`.
- **Integrations (`internal/integrations/asaas`):** Clientes de APIs externas (Gateway Asaas).
- **Handlers (`internal/handler`):** Camada de apresentação HTTP com binding de DTOs e padronização JSON via `pkg/response`.
- **Queue (`internal/queue`):** Produção e consumo de tarefas assíncronas com Redis e Asynq.

---

## 2. 🛡️ Isolamento Multi-Tenant & Hierarquia de Acesso

O sistema suporta três níveis hierárquicos de acesso:
1. **`ADMIN_GLOBAL`:** Administrador da plataforma SaaS. Possui visão global de todos os estabelecimentos, métricas agregadas, gestão de planos comerciais e controle de faturamento.
2. **`ADMIN_TENANT`:** Administrador do estabelecimento. Gerencia profissionais, serviços, bloqueios de agenda, clientes e visualiza sua própria fatura/assinatura.
3. **`OPERATOR`:** Profissional ou atendente. Possui acesso estrito à agenda de atendimentos do seu estabelecimento.

### Scoper de Isolamento no Banco de Dados
Todas as consultas operacionais no PostgreSQL são filtradas obrigatoriamente pelo `tenant_id` injetado via JWT:
```go
db.Where("tenant_id = ?", tenantID)
```

---

## 3. ⚡ Motor de Concorrência & Prevenção de Dupla Reserva

Para garantir que dois agendamentos não sejam confirmados no mesmo horário e profissional simultaneamente:

1. **Abertura de Transação Atômica:** `tx := r.db.Begin()`
2. **Lock Pessimista de Linha:**
   ```sql
   SELECT * FROM professionals WHERE id = ? FOR UPDATE;
   ```
3. **Validação de Intervalos Sobrepostos:**
   ```sql
   SELECT COUNT(*) FROM appointments 
   WHERE professional_id = ? 
     AND status IN ('PENDING', 'CONFIRMED')
     AND start_at < ? AND end_at > ?;
   ```
4. **Decisão Atômica:** Se a contagem for maior que 0, a transação sofre `Rollback` imediato e retorna `ErrSlotAlreadyBooked` (HTTP 409 Conflict). Caso contrário, grava o agendamento e faz `Commit`.

---

## 4. 🔄 Tarefas Assíncronas (Asynq + Redis)

| Tipo de Tarefa | Fila | Payload | Comportamento |
|---|---|---|---|
| `task:appointment_created` | `critical` | `{ "appointment_id": "uuid" }` | Processa confirmações e notificações de agendamento |
| `task:appointment_reminder` | `default` | `{ "appointment_id": "uuid" }` | Dispara lembretes agendados com antecedência |
| `task:process_asaas_webhook` | `critical` | `{ "event": "PAYMENT_CONFIRMED", "payment": {...} }` | Atualiza assinaturas e faturas em segundo plano |

---

## 5. 🚀 Catálogo Completo de Endpoints REST

### 5.1 Público
- `GET /api/v1/health` -> Healthcheck do servidor.
- `GET /api/v1/public/plans` -> Catálogo de planos ativos.
- `GET /api/v1/public/:slug` -> Informações públicas do estabelecimento.
- `GET /api/v1/public/:slug/services` -> Serviços disponíveis.
- `GET /api/v1/public/:slug/services/:id/professionals` -> Profissionais do serviço.
- `GET /api/v1/public/:slug/availability` -> Horários livres em uma data.
- `POST /api/v1/public/:slug/appointments` -> Realizar agendamento público.
- `GET /api/v1/public/:slug/appointments/:id` -> Detalhes do agendamento.
- `POST /api/v1/public/:slug/appointments/:id/cancel` -> Cancelar agendamento.
- `POST /api/v1/webhooks/asaas` -> Webhook do gateway Asaas.

### 5.2 Autenticação
- `POST /api/v1/auth/register` -> Cadastro de novo estabelecimento com plano.
- `POST /api/v1/auth/login` -> Autenticação e emissão de JWT.
- `GET /api/v1/auth/me` -> Dados do usuário e estabelecimento logado.

### 5.3 Administração Global (`ADMIN_GLOBAL`)
- `GET /api/v1/admin/global-dashboard` -> Métricas globais da plataforma.
- `GET /api/v1/admin/tenants` -> Listagem de estabelecimentos com status de assinatura.
- `POST /api/v1/admin/tenants` -> Cadastro de estabelecimento.
- `GET /api/v1/admin/tenants/:id` -> Detalhes completos e histórico de faturas.
- `PUT /api/v1/admin/tenants/:id` -> Edição de dados do estabelecimento.
- `PATCH /api/v1/admin/tenants/:id/status` -> Ativar/inativar estabelecimento.
- `GET /api/v1/admin/plans` -> Listagem de planos comerciais.
- `POST /api/v1/admin/plans` -> Criação de plano.
- `PUT /api/v1/admin/plans/:id` -> Edição de plano.
- `PATCH /api/v1/admin/plans/:id/status` -> Alternar status de venda do plano.
- `DELETE /api/v1/admin/plans/:id` -> Exclusão protegida de plano.
- `PATCH /api/v1/admin/subscriptions/:id/status` -> Ajuste manual de status de assinatura.

### 5.4 Gestão do Estabelecimento (`ADMIN_TENANT` & `OPERATOR`)
- `GET /api/v1/admin/subscription` -> Dados da assinatura contratada e faturas.
- `GET /api/v1/admin/dashboard` -> KPIs do estabelecimento (receita, agendamentos, clientes).
- `GET /api/v1/admin/appointments` -> Listagem e agenda de atendimentos.
- `POST /api/v1/admin/appointments` -> Criar agendamento interno.
- `PATCH /api/v1/admin/appointments/:id/complete` -> Concluir atendimento.
- `PATCH /api/v1/admin/appointments/:id/cancel` -> Cancelar agendamento.
- `PATCH /api/v1/admin/appointments/:id/reschedule` -> Reagendar atendimento.
- `GET /api/v1/admin/services` -> Gestão de serviços.
- `GET /api/v1/admin/professionals` -> Gestão de profissionais e horários.
- `GET /api/v1/admin/customers` -> Histórico e CRM de clientes.
- `GET /api/v1/admin/settings` -> Configurações visuais e intervalos de atendimento.
