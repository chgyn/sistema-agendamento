package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/jwt"
)

type RouterConfig struct {
	Engine               *gin.Engine
	JWTService           *jwt.JWTService
	Repo                 *postgres.Repository
	AuthHandler          *AuthHandler
	PublicBookingHandler *PublicBookingHandler
	AppointmentHandler   *AppointmentHandler
	ServiceHandler       *ServiceHandler
	ProfessionalHandler  *ProfessionalHandler
	CustomerHandler      *CustomerHandler
	DashboardHandler     *DashboardHandler
	TenantHandler        *TenantHandler
	UserHandler          *UserHandler
	PlanHandler            *PlanHandler
	SubscriptionHandler    *SubscriptionHandler
	WebhookHandler         *WebhookHandler
	WhatsAppHandler        *WhatsAppHandler
	WhatsAppWebhookHandler *WhatsAppWebhookHandler
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
		// ROTAS PÚBLICAS DE PLANOS (CATÁLOGO PARA CADASTRO)
		// -------------------------------------------------------------
		api.GET("/public/plans", cfg.PlanHandler.ListActivePublic)

		// -------------------------------------------------------------
		// ROTAS DE WEBHOOKS (INTEGRAÇÕES EXTERNAS: ASAAS & WUZAPI)
		// -------------------------------------------------------------
		api.POST("/webhooks/asaas", cfg.WebhookHandler.HandleAsaas)
		api.POST("/webhooks/whatsapp", cfg.WhatsAppWebhookHandler.ReceiveWebhook)

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
			// ==========================================
			// 1. RECURSOS EXCLUSIVOS DO ADMINISTRADOR GERAL
			// ==========================================
			globalOnly := admin.Group("")
			globalOnly.Use(middleware.RequireRole(domain.RoleAdminGlobal))
			{
				// Dashboard Global da Plataforma
				globalOnly.GET("/global-dashboard", cfg.TenantHandler.GetGlobalDashboard)

				// Gestão de Estabelecimentos (Tenants)
				globalOnly.GET("/tenants", cfg.TenantHandler.ListTenants)
				globalOnly.POST("/tenants", cfg.TenantHandler.CreateTenantWithAdmin)
				globalOnly.GET("/tenants/:id", cfg.TenantHandler.GetTenantByID)
				globalOnly.PUT("/tenants/:id", cfg.TenantHandler.UpdateTenantGlobal)
				globalOnly.PATCH("/tenants/:id/status", cfg.TenantHandler.ToggleTenantStatus)
				globalOnly.POST("/tenants/:id/subscriptions/grant-manual", cfg.SubscriptionHandler.GrantManual)

				// Gestão de Administradores Gerais
				globalOnly.GET("/global-admins", cfg.UserHandler.ListGlobalAdmins)
				globalOnly.POST("/global-admins", cfg.UserHandler.CreateGlobalAdmin)

				// Gestão Global de Planos Comerciais (CRUD Exclusivo)
				globalOnly.GET("/plans", cfg.PlanHandler.List)
				globalOnly.POST("/plans", cfg.PlanHandler.Create)
				globalOnly.GET("/plans/:id", cfg.PlanHandler.GetByID)
				globalOnly.PUT("/plans/:id", cfg.PlanHandler.Update)
				globalOnly.PATCH("/plans/:id/status", cfg.PlanHandler.ToggleStatus)
				globalOnly.DELETE("/plans/:id", cfg.PlanHandler.Delete)

				// Gestão Global de Assinaturas
				globalOnly.GET("/subscriptions", cfg.SubscriptionHandler.ListGlobal)
				globalOnly.PATCH("/subscriptions/:id/status", cfg.SubscriptionHandler.OverrideStatus)
				globalOnly.GET("/subscriptions/:id/audit-logs", cfg.SubscriptionHandler.ListAuditLogs)
			}

			// ==========================================
			// 2. GESTÃO DE USUÁRIOS (ADMIN_GLOBAL e ADMIN_TENANT)
			// ==========================================
			userMgmt := admin.Group("/users")
			userMgmt.Use(middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin))
			{
				userMgmt.GET("", cfg.UserHandler.List)
				userMgmt.GET("/:id", cfg.UserHandler.GetByID)
				userMgmt.POST("", cfg.UserHandler.Create)
				userMgmt.PUT("/:id", cfg.UserHandler.Update)
				userMgmt.PATCH("/:id/status", cfg.UserHandler.ToggleStatus)
			}

			// ==========================================
			// 3. CONSULTA DE ASSINATURA DO TENANT (LIVRE DE GATEKEEPER)
			// ==========================================
			admin.GET("/subscription", cfg.SubscriptionHandler.GetMySubscription)

			// ==========================================
			// 4. RECURSOS OPERACIONAIS DO TENANT (PROTEGIDOS POR ASSINATURA ATIVA)
			// ==========================================
			tenantOps := admin.Group("")
			tenantOps.Use(middleware.RequireActiveSubscription(cfg.Repo))
			{
				// Dashboard & Métricas do Tenant
				tenantOps.GET("/dashboard", cfg.DashboardHandler.GetKPIs)
				tenantOps.GET("/dashboard/analytics", cfg.DashboardHandler.GetAnalytics)

				// Agendamentos / Agenda (Operadores e Admins)
				tenantOps.GET("/appointments", cfg.AppointmentHandler.List)
				tenantOps.GET("/appointments/:id", cfg.AppointmentHandler.GetByID)
				tenantOps.POST("/appointments", cfg.AppointmentHandler.CreateAdmin)
				tenantOps.PATCH("/appointments/:id/complete", cfg.AppointmentHandler.Complete)
				tenantOps.PATCH("/appointments/:id/cancel", cfg.AppointmentHandler.Cancel)
				tenantOps.PATCH("/appointments/:id/reschedule", cfg.AppointmentHandler.Reschedule)

				// Serviços
				tenantOps.GET("/services", cfg.ServiceHandler.List)
				tenantOps.GET("/services/:id", cfg.ServiceHandler.GetByID)
				tenantOps.POST("/services", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ServiceHandler.Create)
				tenantOps.PUT("/services/:id", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ServiceHandler.Update)
				tenantOps.DELETE("/services/:id", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ServiceHandler.Delete)

				// Profissionais, Horários e Bloqueios
				tenantOps.GET("/professionals", cfg.ProfessionalHandler.List)
				tenantOps.GET("/professionals/:id", cfg.ProfessionalHandler.GetByID)
				tenantOps.POST("/professionals", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ProfessionalHandler.Create)
				tenantOps.PUT("/professionals/:id", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ProfessionalHandler.Update)
				tenantOps.DELETE("/professionals/:id", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ProfessionalHandler.Delete)

				tenantOps.GET("/professionals/:id/working-hours", cfg.ProfessionalHandler.GetWorkingHours)
				tenantOps.PUT("/professionals/:id/working-hours", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ProfessionalHandler.SaveWorkingHours)

				tenantOps.GET("/professionals/:id/exceptions", cfg.ProfessionalHandler.ListExceptions)
				tenantOps.POST("/professionals/:id/exceptions", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ProfessionalHandler.CreateException)
				tenantOps.DELETE("/professionals/:id/exceptions/:exception_id", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.ProfessionalHandler.DeleteException)

				// Clientes e Histórico
				tenantOps.GET("/customers", cfg.CustomerHandler.List)
				tenantOps.GET("/customers/:id", cfg.CustomerHandler.GetByID)

				// Configurações do Estabelecimento (Restrito a ADMIN_GLOBAL e ADMIN_TENANT)
				tenantOps.GET("/settings", cfg.TenantHandler.GetSettings)
				tenantOps.PUT("/settings", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.TenantHandler.UpdateSettings)

				// Integração WhatsApp & Atendimento IA (WUZAPI)
				whatsapp := tenantOps.Group("/whatsapp")
				{
					whatsapp.GET("/status", cfg.WhatsAppHandler.GetStatus)
					whatsapp.POST("/instance", cfg.WhatsAppHandler.EnsureInstance)
					whatsapp.POST("/connect", cfg.WhatsAppHandler.Connect)
					whatsapp.GET("/qr", cfg.WhatsAppHandler.GetQRCode)
					whatsapp.POST("/disconnect", cfg.WhatsAppHandler.Disconnect)
					whatsapp.POST("/logout", cfg.WhatsAppHandler.Logout)
					whatsapp.GET("/ai-config", cfg.WhatsAppHandler.GetAIConfig)
					whatsapp.PUT("/ai-config", middleware.RequireRole(domain.RoleAdminGlobal, domain.RoleAdminTenant, domain.RoleAdmin), cfg.WhatsAppHandler.UpdateAIConfig)
				}
			}
		}
	}
}
