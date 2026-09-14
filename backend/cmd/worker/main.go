package main

import (
	"log"

	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/integrations/ai"
	"github.com/sistema-agendamento/backend/internal/integrations/wuzapi"
	"github.com/sistema-agendamento/backend/internal/queue"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/database"
)

func main() {
	log.Println("👷 Inicializando Sistema de Agendamento — Worker Daemon (Asynq + Redis)...")

	cfg := config.Load()

	if err := runWorker(cfg); err != nil {
		log.Fatalf("❌ Erro fatal no Worker Asynq: %v", err)
	}
}

func runWorker(cfg *config.Config) error {
	// 1. Conecta ao banco de dados
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Worker falhou ao conectar ao banco de dados: %v", err)
	}

	repo := postgres.NewRepository(db)
	queueClient := queue.NewQueueClient(cfg)
	defer queueClient.Close()

	// 2. Inicializa integrações e serviços necessários para o atendimento da IA
	wuzapiClient := wuzapi.NewClient(cfg.WUZAPIBaseURL, cfg.WUZAPIAdminToken)
	dispatcher := service.NewWhatsAppMessageDispatcher(wuzapiClient)

	geminiProvider := ai.NewGeminiProvider()
	openaiProvider := ai.NewOpenAIProvider()

	availService := service.NewAvailabilityService(repo)
	aptService := service.NewAppointmentService(repo, queueClient)

	aiChatService := service.NewAIChatService(
		repo,
		availService,
		aptService,
		geminiProvider,
		openaiProvider,
		dispatcher,
		cfg.EncryptionKey,
	)

	// 3. Inicia servidor consumidor de filas com o processador de WhatsApp injetado
	_, err = queue.StartWorkerServer(cfg, aiChatService)
	return err
}
