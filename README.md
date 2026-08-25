# Sistema Multi-Tenant de Agendamento (Barbearias & Salões de Beleza)

Plataforma completa de agendamento de atendimentos desenvolvida com **Go (Gin Gonic + GORM + PostgreSQL + Redis + Asynq)** no backend e **Vue 3 (Vite + Tailwind CSS v4 + Pinia + Lucide Icons)** no frontend, seguindo os princípios de **Clean Architecture**, **Isolamento Rigoroso Multi-Tenant**, **Planos Comerciais** e **Integração Recorrente com Asaas (v3)**.

---

## 🚀 Novidades & Módulos Implementados

1. **💳 Planos de Assinatura & Monetização SaaS:**
   - CRUD completo de planos comerciais restrito ao **Administrador Geral** (`ADMIN_GLOBAL`).
   - Suporte a múltiplas periodicidades (`MONTHLY`, `QUARTERLY`, `SEMIANNUALLY`, `YEARLY`), limites de profissionais/serviços e recursos configuráveis.
2. **🔄 Integração com Gateway Asaas (API v3):**
   - Criação automática de clientes e assinaturas recorrentes durante o cadastro do estabelecimento.
   - Processamento resiliente de Webhooks com filas **Asynq + Redis** (`critical`) para confirmação de pagamentos, geração de faturas e bloqueios por inadimplência.
3. **🛡️ Gatekeeper de Acesso por Status de Assinatura:**
   - Middleware `RequireActiveSubscription` protegendo rotas operacionais do tenant (`/appointments`, `/services`, `/professionals`, `/customers`).
4. **📊 Visualização & Gestão da Situação da Assinatura dos Estabelecimentos:**
   - Tabela de estabelecimentos com identificadores de plano, status financeiro (Ativa, Pendente, Vencida, Cancelada) e filtros compostos.
   - Modal de detalhamento com resumo do plano, identificadores Asaas (`AsaasSubscriptionID`, `AsaasCustomerID`), histórico de faturas e ajuste manual de status.
5. **🔒 Prevenção Concorrente de Dupla Reserva:**
   - Transações atômicas com `SELECT ... FOR UPDATE` no PostgreSQL e validação matemática de sobreposição de intervalos temporais.

---

## 🏗️ Estrutura do Projeto

```text
sistema-agendamento/
├── backend/
│   ├── cmd/
│   │   ├── api/                  # Inicializa a API REST Gin
│   │   └── worker/               # Inicializa o Worker Asynq (Redis)
│   ├── internal/
│   │   ├── config/               # Variáveis de ambiente (.env)
│   │   ├── domain/               # Modelos, DTOs e Portas de Interfaces
│   │   ├── handler/              # Handlers HTTP Gin (Públicos, Admin e Webhooks)
│   │   ├── integrations/asaas/   # Cliente HTTP nativo da API v3 do Asaas
│   │   ├── middleware/           # JWT Auth, Scoper de Tenant e Gatekeeper de Assinatura
│   │   ├── queue/                # Handlers e Produtores de tarefas do Asynq
│   │   ├── repository/postgres/  # Repositórios GORM com Preload e Locks
│   │   └── service/              # Regras de Negócio, Disponibilidade e Assinaturas
│   ├── pkg/                      # Database, Hash bcrypt, JWT e Respostas JSON
│   ├── .env.example
│   └── Dockerfile.backend
├── frontend/
│   ├── src/
│   │   ├── components/layout/    # Layout administrativo com Sidebar dinâmica
│   │   ├── router/               # Rotas Vue Router com guards de perfil
│   │   ├── services/             # Cliente Axios com interceptors
│   │   ├── stores/               # Stores Pinia (Auth, Agendamento)
│   │   └── views/
│   │       ├── admin/            # Dashboard, Agenda, Planos, Assinaturas, Estabelecimentos
│   │       ├── auth/             # Login e Cadastro com Seleção de Planos
│   │       └── public/           # Wizard de Agendamento Público (/agendamento/:slug)
│   └── Dockerfile.frontend
├── docs/
│   ├── MODULO_PLANOS_E_ASSINATURAS_ASAAS.md
│   └── ARQUITETURA_E_API.md
├── docker-compose.yml
└── README.md
```

---

## ⚡ Como Executar a Aplicação

### Opção 1: Docker Compose (Recomendado)

Suba toda a infraestrutura com um único comando:

```bash
docker compose up --build
```

Após a inicialização:
- **Frontend SPA:** `http://localhost:5173`
- **Painel Administrativo:** `http://localhost:5173/login`
- **Página Pública Barbearia Dom Navalha:** `http://localhost:5173/agendamento/dom-navalha`
- **API REST Backend:** `http://localhost:8080/api/v1/health`

---

### Opção 2: Execução Local (Standalone)

#### 1. Backend API & Worker
```bash
cd backend
cp .env.example .env
go run cmd/api/main.go
```
Em outro terminal:
```bash
cd backend
go run cmd/worker/main.go
```

#### 2. Frontend
```bash
cd frontend
npm install
npm run dev
```

---

## 🔑 Credenciais de Demonstração (Seeding Automático)

### 👑 1. Administrador Geral da Plataforma (Admin Global)
- **E-mail:** `admin@plataforma.com`
- **Senha:** `admin123`
- **Acesso:** Gestão global de estabelecimentos, planos de assinatura, dashboards agregados e controle de faturas.

### 💈 2. Estabelecimento Demonstração 1 (Barbearia Dom Navalha)
- **Slug Público:** `/agendamento/dom-navalha`
- **E-mail do Administrador:** `admin@domnavalha.com`
- **Senha:** `admin123`
- **Plano Vinculado:** Plano Profissional (Assinatura Ativa)

### 💇 3. Estabelecimento Demonstração 2 (Salão Bella Vista)
- **Slug Público:** `/agendamento/bella-vista`
- **E-mail do Administrador:** `admin@bellavista.com`
- **Senha:** `admin123`

---

## 🧪 Execução de Testes Automatizados

```bash
# Testes do Backend (Unitários, Concorrência, Multi-Tenant e Assinaturas)
cd backend
go test -v ./...

# Validação de Tipos e Build do Frontend
cd frontend
npm run build
```

---

## 📚 Documentação Técnica Detalhada

- [Módulo de Planos & Assinaturas Asaas](file:///home/charles/projetos/sistema-agendamento/docs/MODULO_PLANOS_E_ASSINATURAS_ASAAS.md)
- [Arquitetura, Concorrência & Catálogo de APIs](file:///home/charles/projetos/sistema-agendamento/docs/ARQUITETURA_E_API.md)
