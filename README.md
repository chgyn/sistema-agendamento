# Sistema Multi-Tenant de Agendamento (Barbearias & Salões de Beleza)

Plataforma completa de agendamento de atendimentos desenvolvida com **Go (Gin Gonic + GORM + PostgreSQL + Redis + Asynq)** no backend e **Vue 3 (Vite + Tailwind CSS + Pinia + Lucide Icons)** no frontend, seguindo os princípios de **Clean Architecture** e isolamento rigoroso **Multi-Tenant**.

---

## 🚀 Tecnologias Utilizadas

### Backend (Go)
- **Linguagem:** Go 1.23+
- **API REST:** Gin Gonic + CORS + Recovery & Logger Middlewares
- **Autenticação & Segurança:** JWT (JSON Web Tokens) com claims multi-tenant (`tenant_id`, `user_id`, `role`) e senhas criptografadas com `bcrypt`
- **Banco de Dados:** PostgreSQL 16 com ORM GORM, Auto-Migrations e Seeding automático de demonstração
- **Filas & Tarefas Assíncronas:** Redis 7 + Asynq (Workers dedicados para notificações, confirmações e lembretes)
- **Prevenção de Dupla Reserva:** Transações atômicas com row-level lock (`SELECT ... FOR UPDATE`) e validação de sobreposição temporal de intervalos no banco de dados

### Frontend (Vue.js)
- **Framework:** Vue 3 (Composition API com `<script setup lang="ts">`)
- **Build Tool:** Vite 5 + TypeScript
- **Estilização:** Tailwind CSS (Dark Mode nativo, Glassmorphism e gradientes modernos)
- **Gerenciamento de Estado:** Pinia
- **Roteamento:** Vue Router 4 com Navigation Guards de autenticação
- **Ícones & UI:** Lucide Icons (`lucide-vue-next`), `date-fns` (Localização pt-BR)

---

## 🏗️ Estrutura do Monorepo

```text
sistema-agendamento/
├── backend/
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go           # Inicializa a API REST Gin
│   │   └── worker/
│   │       └── main.go           # Inicializa os Workers Asynq (Redis)
│   ├── internal/
│   │   ├── config/               # Leitura de variáveis de ambiente (.env)
│   │   ├── domain/               # Modelos de domínio, entidades e erros
│   │   ├── handler/              # Handlers HTTP Gin (Públicos e Administrativos)
│   │   ├── middleware/           # JWT Auth, Scoper de Tenant e CORS
│   │   ├── queue/                # Handlers e Produtores de tarefas do Asynq
│   │   ├── repository/postgres/  # Acesso a dados GORM com isolamento de tenant e locks
│   │   └── service/              # Regras de negócio, Motor de Disponibilidade e Agendamento
│   ├── pkg/
│   │   ├── database/             # Conexão DB, Migrações e Seeding de Demonstração
│   │   ├── hash/                 # Criptografia de senhas bcrypt
│   │   ├── jwt/                  # Gerenciamento de tokens JWT
│   │   └── response/             # Formato padronizado de respostas JSON
│   ├── go.mod, go.sum
│   └── Dockerfile.backend
├── frontend/
│   ├── src/
│   │   ├── assets/               # Estilos Tailwind e Glassmorphism
│   │   ├── components/layout/    # Layout administrativo com Sidebar e Navegação
│   │   ├── router/               # Rotas com proteção de autenticação
│   │   ├── services/             # Cliente HTTP Axios com interceptor JWT
│   │   ├── stores/               # Stores Pinia (Auth e Wizard de Agendamento)
│   │   └── views/
│   │       ├── admin/            # Dashboard, Agenda/Calendário, Serviços, Profissionais, Clientes, Configurações
│   │       ├── auth/             # Login e Cadastro de Novo Estabelecimento (Onboarding)
│   │       ├── public/           # Wizard de Agendamento Público (/agendamento/:slug)
│   │       └── HomeView.vue      # Página Inicial com links e demonstrações de 1 clique
│   ├── package.json, vite.config.ts, tailwind.config.js
│   ├── nginx.conf
│   └── Dockerfile.frontend
├── docker-compose.yml             # Orquestrador local (Postgres 16, Redis 7, API, Worker, Front)
└── README.md
```

---

## ⚡ Como Executar a Aplicação

### Opção 1: Via Docker Compose (Recomendado)

Suba toda a infraestrutura (PostgreSQL, Redis, API Go, Worker Go e Frontend Nginx) com um único comando:

```bash
docker compose up --build
```

Após a inicialização:
- **Página Inicial & Portal de Demos:** `http://localhost:5173` ou `http://localhost`
- **Painel Administrativo:** `http://localhost:5173/login`
- **Página Pública Barbearia Dom Navalha:** `http://localhost:5173/agendamento/dom-navalha`
- **Página Pública Salão Bella Vista:** `http://localhost:5173/agendamento/bella-vista`
- **API REST Backend:** `http://localhost:8080/api/v1/health`

---

### Opção 2: Execução Local para Desenvolvimento (Standalone)

#### 1. Iniciar o Backend
```bash
cd backend
go run cmd/api/main.go
```
*(O backend possui fallback automático e seed inteligente de dados).*

Para rodar o Worker de tarefas em segundo plano (em outro terminal):
```bash
cd backend
go run cmd/worker/main.go
```

#### 2. Iniciar o Frontend
```bash
cd frontend
npm install
npm run dev
```

Acesse no navegador: `http://localhost:5173`

---

## 🔑 Credenciais e Dados de Demonstração (Seed)

O banco de dados é populado automaticamente na primeira execução com 2 estabelecimentos completos:

### 💈 1. Barbearia Dom Navalha
- **Slug Público:** `/agendamento/dom-navalha`
- **E-mail do Administrador:** `admin@domnavalha.com`
- **Senha:** `admin123`
- **Profissionais cadastrados:** Carlos Navalha (Master), Lucas Tesoura (Especialista), Diego Visagista
- **Serviços:** Corte Degradê, Barba Terapia com Toalha Quente, Combo VIP, Acabamento, Pigmentação de Barba

### 💇 2. Salão Bella Vista Studio
- **Slug Público:** `/agendamento/bella-vista`
- **E-mail do Administrador:** `admin@bellavista.com`
- **Senha:** `admin123`

---

## 🛡️ Prevenção Concorrente de Dupla Reserva

O sistema foi rigorosamente desenhado para impedir que dois clientes consigam reservar o mesmo profissional e horário simultaneamente:
1. Ao receber a solicitação de agendamento, abre-se uma **transação atômica** no banco de dados.
2. A linha do profissional é travada com `SELECT ... FOR UPDATE`.
3. É realizada a verificação de sobreposição:
   $$\text{existing.StartAt} < \text{new.EndAt} \quad \text{AND} \quad \text{existing.EndAt} > \text{new.StartAt}$$
4. Caso haja colisão com outro agendamento `CONFIRMED` ou `PENDING`, a transação sofre Rollback e a API responde com status `409 Conflict` (`ErrSlotAlreadyBooked`).
5. Caso esteja livre, o agendamento é persistido, as métricas do cliente são atualizadas e uma tarefa assíncrona é disparada para o **Asynq + Redis**.

### Executar Testes Automatizados de Concorrência:
```bash
cd backend
go test -v ./internal/service/...
```
*O teste `TestDoubleBookingConcurrencyValidation` dispara 20 goroutines simultâneas disputando o mesmo horário no mesmo milissegundo e valida que exatamente 1 agendamento é aprovado e os 19 restantes são bloqueados com segurança.*
