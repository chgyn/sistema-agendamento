package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/handler"
	"github.com/sistema-agendamento/backend/internal/queue"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/database"
	"github.com/sistema-agendamento/backend/pkg/jwt"
)

func main() {
	log.Println("🚀 Inicializando Sistema de Agendamento — API Backend...")

	// 1. Carrega configurações
	cfg := config.Load()
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. Conecta ao banco de dados e roda migrações
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Erro fatal ao conectar ao banco de dados: %v", err)
	}

	// 3. Popula dados de demonstração (se o banco estiver vazio)
	if err := database.SeedInitialData(db); err != nil {
		log.Printf("⚠️ Aviso durante seeding inicial: %v", err)
	}

	// 4. Inicializa Repositórios
	repo := postgres.NewRepository(db)

	// 5. Inicializa Cliente de Filas Asynq / Redis
	queueClient := queue.NewQueueClient(cfg)
	defer queueClient.Close()

	// 6. Inicializa Serviços
	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn)
	authService := service.NewAuthService(repo, jwtService)
	availService := service.NewAvailabilityService(repo)
	aptService := service.NewAppointmentService(repo, queueClient)
	tenantService := service.NewTenantService(repo)

	// 7. Inicializa Handlers
	authHandler := handler.NewAuthHandler(authService, repo)
	publicBookingHandler := handler.NewPublicBookingHandler(repo, availService, aptService)
	appointmentHandler := handler.NewAppointmentHandler(aptService, repo)
	serviceHandler := handler.NewServiceHandler(repo)
	proHandler := handler.NewProfessionalHandler(repo)
	customerHandler := handler.NewCustomerHandler(repo)
	dashboardHandler := handler.NewDashboardHandler(tenantService)
	tenantHandler := handler.NewTenantHandler(repo)

	// 8. Configura Router Gin
	r := gin.Default()
	handler.SetupRoutes(handler.RouterConfig{
		Engine:               r,
		JWTService:           jwtService,
		AuthHandler:          authHandler,
		PublicBookingHandler: publicBookingHandler,
		AppointmentHandler:   appointmentHandler,
		ServiceHandler:       serviceHandler,
		ProfessionalHandler:  proHandler,
		CustomerHandler:      customerHandler,
		DashboardHandler:     dashboardHandler,
		TenantHandler:        tenantHandler,
	})

	// 9. Inicia Servidor HTTP
	addr := ":" + cfg.Port
	log.Printf("⚡ Servidor HTTP ouvindo na porta %s (Ambiente: %s)", cfg.Port, cfg.Environment)
	log.Printf("🔗 Link público demo: http://localhost:%s/api/v1/public/dom-navalha", cfg.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Falha ao iniciar servidor HTTP: %v", err)
	}
}
