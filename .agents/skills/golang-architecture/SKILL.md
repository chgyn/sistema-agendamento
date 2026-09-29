---
name: golang-architecture
description: >-
  Diretrizes de arquitetura Clean/Hexagonal, boas práticas Go 1.24+, Gin Gonic,
  GORM, Asynq, concorrência, criptografia AES-256-GCM e Vue 3 para o projeto
  sistema-agendamento.
---

# 📐 Skill: Go (Golang) Clean Architecture & Best Practices Standard

Você é um especialista em engenharia de software com foco no ecossistema Go (1.24+), Clean Architecture, Gin Gonic, GORM, Asynq (Redis), PostgreSQL e Vue 3. Siga estritamente estas diretrizes ao planejar, revisar ou implementar código neste projeto.

---

## 1. Princípios de Go Moderno (Go 1.24+)

- **Tratamento Explícito de Erros:**
  - NUNCA ignore erros com `_`. Sempre trate ou propague: `if err != nil { return fmt.Errorf("contexto da falha: %w", err) }`.
  - Utilize `%w` com `fmt.Errorf` para permitir inspeção com `errors.Is()` e `errors.As()`.
- **Tipagem Forte e Struct Tags:**
  - Todas as estruturas devem ter tags explicitadas para GORM (`gorm:"..."`), JSON (`json:"..."`) e validação Gin (`binding:"..."`).
  - Use `json:"-"` para campos confidenciais (senhas, cookies não criptografados, tokens brutos).
- **Concorrência e Ciclo de Vida Seguros:**
  - Sempre propague `ctx context.Context` em operações I/O (repositórios, chamadas HTTP, workers).
  - Controle timeouts com `context.WithTimeout` para evitar bloqueios indefinidos em requisições externas.
  - Sincronize acessos concorrentes a recursos compartilhados com `sync.Mutex` ou `sync.RWMutex`.
  - Use `sync.WaitGroup` para coordenar o término de goroutines.

---

## 2. Clean Architecture & Estrutura de Pacotes

```text
backend/
├── cmd/
│   ├── api/main.go               # Entrypoint HTTP REST (Gin) + Cron Scheduler
│   └── worker/main.go            # Entrypoint do Worker Assíncrono (Asynq)
├── internal/
│   ├── domain/                   # Modelos puros e Interfaces (Ports)
│   │   ├── models.go             # Structs GORM, entidades e tags
│   │   └── ports.go              # DTOs e Interfaces de Repositórios e Provedores
│   ├── storage/repositories/     # Implementações GORM das interfaces de repositório
│   ├── api/                      # Handlers HTTP Gin e Roteador
│   ├── auth/                     # JWT, Middlewares e Hash de senhas (Bcrypt)
│   ├── crypto/                   # Criptografia AES-256-GCM para cookies e credenciais
│   ├── jobs/                     # Tasks Asynq, Handlers do Worker e Agendador Cron
│   └── integrations/             # Clientes externos (WUZAPI, Telegram, Gemini, Scraper)
```

### 2.1 Handlers Enxutos (Thin Handlers no Gin)
- Handlers HTTP têm a responsabilidade estrita de:
  1. Fazer o binding e validação do payload: `if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, ...); return }`
  2. Obter o usuário autenticado do contexto Gin: `userID, _ := c.Get("user_id")`
  3. Invocar a camada de repositório, serviço ou despachar a task para o Asynq.
  4. Retornar a resposta JSON padronizada (`c.JSON(http.StatusOK, gin.H{...})`).

### 2.2 Repositórios GORM & Prevenção de N+1
- Métodos de repositório devem sempre receber `ctx context.Context` e aplicar `.WithContext(ctx)`.
- **Prevenção de N+1:** Utilize `.Preload("NomeDoRelacionamento")` ao buscar entidades com associações (ex: `Product.AffiliateLink`, `ScrapingSource.Executions`).
- **Transações Atômicas:** Em operações que envolvem múltiplas mutações, use transações GORM:
  ```go
  tx := r.db.WithContext(ctx).Begin()
  defer func() {
      if r := recover(); r != nil {
          tx.Rollback()
      }
  }()
  if err := tx.Create(&entity).Error; err != nil {
      tx.Rollback()
      return err
  }
  return tx.Commit().Error
  ```

---

## 3. Tarefas Assíncronas (Asynq + Redis)

- **Desacoplamento:** Qualquer operação de longa duração (scraping de páginas, geração de copys com IA, disparos de mensagens, encurtamento de links) DEVE ser executada assincronamente via Asynq.
- **Tipagem de Tarefas:** Defina nomes constantes em `internal/jobs/tasks.go`:
  ```go
  const (
      TypeScrapeSource    = "scraping:run"
      TypeGenerateAICopy  = "ai:generate"
      TypeSendDispatch    = "dispatch:send"
  )
  ```
- **Payloads Limpos:** Envie apenas identificadores ou estruturas de dados estritamente necessárias no payload JSON da tarefa.

---

## 4. Frontend Vue 3 + Tailwind CSS v4 + Pinia

- **Composition API:** Utilize `<script setup>` em todos os componentes e telas (`frontend/src/pages/`).
- **Comunicação API:** Utilize o cliente Axios centralizado (`frontend/src/services/api.js`), que gerencia a injeção do cabeçalho `Authorization: Bearer <token>` e trata respostas `401 Unauthorized` com redirecionamento para o login.
- **Tailwind CSS v4:** Aplique classes utilitárias modernas, tema escuro consistente (dark slate / zinc), efeitos de glassmorphism (`backdrop-blur-md bg-slate-900/80 border border-slate-800`), e micro-animações em botões e cards.

---

## 5. 🚨 REGRA CRÍTICA DE ISOLAMENTO DE BANCO DE DADOS EM TESTES

- **PROIBIÇÃO ABSOLUTA**: **NUNCA** execute testes unitários ou comandos automatizados que alterem, limpem (`TRUNCATE`) ou executem drops de tabela no banco de dados principal de desenvolvimento (`DB_NAME=gestao_afiliados`).
- **RISCO GRAVE**: Resetar a base principal apaga dados de desenvolvimento, fontes de scraping cadastradas e credenciais de integração do usuário.
- **EXECUÇÃO CORRETA DE TESTES**:
  - Em testes unitários, utilize **Mocks das interfaces** declaradas em `internal/domain/ports.go` (ex: `type MockUserRepository struct {}`).
  - Para testes de integração de banco de dados, utilize **EXCLUSIVAMENTE** uma base de testes dedicada (`DB_NAME=gestao_afiliados_test`) ou banco SQLite em memória.
  - **Comando Permitido**: `go test -v ./internal/...`
  - **Exemplo ESTRITAMENTE PROIBIDO**: Rodar scripts com `db.Migrator().DropTable(...)` apontando para a base `gestao_afiliados`.
