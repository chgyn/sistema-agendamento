package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/pkg/jwt"
)

type RouterConfig struct {
	Engine               *gin.Engine
	JWTService           *jwt.JWTService
	AuthHandler          *AuthHandler
	PublicBookingHandler *PublicBookingHandler
	AppointmentHandler   *AppointmentHandler
	ServiceHandler       *ServiceHandler
	ProfessionalHandler  *ProfessionalHandler
	CustomerHandler      *CustomerHandler
	DashboardHandler     *DashboardHandler
	TenantHandler        *TenantHandler
}

func SetupRoutes(cfg RouterConfig) {
	r := cfg.Engine

	// Middlewares globais
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	api := r.Group("/api/v1")
	{
		// Healthcheck
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "timestamp": "online"})
		})

		// -------------------------------------------------------------
		// ROTAS PÚBLICAS DE AGENDAMENTO (NÃO REQUEREM LOGIN)
		// Exemplo: /api/v1/public/dom-navalha/...
		// -------------------------------------------------------------
		public := api.Group("/public/:slug")
		{
			public.GET("", cfg.PublicBookingHandler.GetTenantInfo)
			public.GET("/services", cfg.PublicBookingHandler.ListServices)
			public.GET("/services/:service_id/professionals", cfg.PublicBookingHandler.ListProfessionalsByService)
			public.GET("/availability", cfg.PublicBookingHandler.GetAvailability)
			public.POST("/appointments", cfg.PublicBookingHandler.CreateAppointment)
			public.GET("/appointments/:id", cfg.PublicBookingHandler.GetAppointmentDetails)
			public.POST("/appointments/:id/cancel", cfg.PublicBookingHandler.CancelAppointment)
		}

		// -------------------------------------------------------------
		// ROTAS DE AUTENTICAÇÃO
		// -------------------------------------------------------------
		auth := api.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.RegisterTenant)
			auth.POST("/login", cfg.AuthHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(cfg.JWTService), cfg.AuthHandler.Me)
		}

		// -------------------------------------------------------------
		// ROTAS ADMINISTRATIVAS PROTEGIDAS (JWT + TENANT_ID SCOPER)
		// -------------------------------------------------------------
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTService))
		{
			// Dashboard & Métricas
			admin.GET("/dashboard", cfg.DashboardHandler.GetKPIs)

			// Agendamentos / Agenda
			admin.GET("/appointments", cfg.AppointmentHandler.List)
			admin.GET("/appointments/:id", cfg.AppointmentHandler.GetByID)
			admin.POST("/appointments", cfg.AppointmentHandler.CreateAdmin)
			admin.PATCH("/appointments/:id/complete", cfg.AppointmentHandler.Complete)
			admin.PATCH("/appointments/:id/cancel", cfg.AppointmentHandler.Cancel)
			admin.PATCH("/appointments/:id/reschedule", cfg.AppointmentHandler.Reschedule)

			// Serviços
			admin.GET("/services", cfg.ServiceHandler.List)
			admin.GET("/services/:id", cfg.ServiceHandler.GetByID)
			admin.POST("/services", middleware.RequireRole(domain.RoleAdmin), cfg.ServiceHandler.Create)
			admin.PUT("/services/:id", middleware.RequireRole(domain.RoleAdmin), cfg.ServiceHandler.Update)
			admin.DELETE("/services/:id", middleware.RequireRole(domain.RoleAdmin), cfg.ServiceHandler.Delete)

			// Profissionais, Horários e Bloqueios
			admin.GET("/professionals", cfg.ProfessionalHandler.List)
			admin.GET("/professionals/:id", cfg.ProfessionalHandler.GetByID)
			admin.POST("/professionals", middleware.RequireRole(domain.RoleAdmin), cfg.ProfessionalHandler.Create)
			admin.PUT("/professionals/:id", middleware.RequireRole(domain.RoleAdmin), cfg.ProfessionalHandler.Update)
			admin.DELETE("/professionals/:id", middleware.RequireRole(domain.RoleAdmin), cfg.ProfessionalHandler.Delete)

			admin.GET("/professionals/:id/working-hours", cfg.ProfessionalHandler.GetWorkingHours)
			admin.PUT("/professionals/:id/working-hours", middleware.RequireRole(domain.RoleAdmin), cfg.ProfessionalHandler.SaveWorkingHours)

			admin.GET("/professionals/:id/exceptions", cfg.ProfessionalHandler.ListExceptions)
			admin.POST("/professionals/:id/exceptions", cfg.ProfessionalHandler.CreateException)
			admin.DELETE("/professionals/:id/exceptions/:exception_id", cfg.ProfessionalHandler.DeleteException)

			// Clientes e Histórico
			admin.GET("/customers", cfg.CustomerHandler.List)
			admin.GET("/customers/:id", cfg.CustomerHandler.GetByID)

			// Configurações do Estabelecimento
			admin.GET("/settings", cfg.TenantHandler.GetSettings)
			admin.PUT("/settings", middleware.RequireRole(domain.RoleAdmin), cfg.TenantHandler.UpdateSettings)
		}
	}
}
