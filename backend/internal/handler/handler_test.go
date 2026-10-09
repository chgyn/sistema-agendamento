package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/handler"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/hash"
	"github.com/sistema-agendamento/backend/pkg/jwt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestApp(t *testing.T) (*gin.Engine, *postgres.Repository, *jwt.JWTService) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("falha ao conectar sqlite teste: %v", err)
	}

	_ = db.AutoMigrate(
		&domain.Tenant{},
		&domain.User{},
		&domain.Professional{},
		&domain.Service{},
		&domain.ProfessionalService{},
		&domain.WorkingHour{},
		&domain.AvailabilityException{},
		&domain.Customer{},
		&domain.Appointment{},
		&domain.AppointmentStatusHistory{},
		&domain.AuditLog{},
		&domain.Plan{},
		&domain.Subscription{},
		&domain.SubscriptionInvoice{},
		&domain.SubscriptionAuditLog{},
	)

	repo := postgres.NewRepository(db)
	jwtSvc := jwt.NewJWTService("super-secret-test-key-32-chars-long", 24)
	userHandler := handler.NewUserHandler(repo)
	tenantHandler := handler.NewTenantHandler(repo)
	planService := service.NewPlanService(repo)
	planHandler := handler.NewPlanHandler(planService)
	subService := service.NewSubscriptionService(repo, nil)
	subHandler := handler.NewSubscriptionHandler(subService)
	cfg := &config.Config{AsaasWebhookSecret: "test-secret"}
	webhookHandler := handler.NewWebhookHandler(cfg, nil, subService)
	tenantService := service.NewTenantService(repo)
	dashboardHandler := handler.NewDashboardHandler(tenantService)

	r := gin.New()
	handler.SetupRoutes(handler.RouterConfig{
		Engine:              r,
		JWTService:          jwtSvc,
		Repo:                repo,
		UserHandler:         userHandler,
		TenantHandler:       tenantHandler,
		PlanHandler:         planHandler,
		SubscriptionHandler: subHandler,
		WebhookHandler:      webhookHandler,
		DashboardHandler:     dashboardHandler,
	})

	return r, repo, jwtSvc
}

func TestRouteAuthorizationAndTenantIsolation(t *testing.T) {
	engine, repo, jwtSvc := setupTestApp(t)
	passHash, _ := hash.HashPassword("123456")

	// Setup Tenant 1 & Tenant 2
	t1ID := uuid.New()
	t1 := domain.Tenant{ID: t1ID, Slug: "t1", Name: "Tenant 1", IsActive: true}
	_ = repo.CreateTenant(nil, &t1)

	t2ID := uuid.New()
	t2 := domain.Tenant{ID: t2ID, Slug: "t2", Name: "Tenant 2", IsActive: true}
	_ = repo.CreateTenant(nil, &t2)

	// Users
	globalAdmin := domain.User{ID: uuid.New(), Name: "Global", Email: "g@p.com", PasswordHash: passHash, Role: domain.RoleAdminGlobal, IsActive: true}
	adminT1 := domain.User{ID: uuid.New(), TenantID: &t1ID, Name: "Admin 1", Email: "a1@t1.com", PasswordHash: passHash, Role: domain.RoleAdminTenant, IsActive: true}
	opT1 := domain.User{ID: uuid.New(), TenantID: &t1ID, Name: "Op 1", Email: "op1@t1.com", PasswordHash: passHash, Role: domain.RoleOperator, IsActive: true}
	adminT2 := domain.User{ID: uuid.New(), TenantID: &t2ID, Name: "Admin 2", Email: "a2@t2.com", PasswordHash: passHash, Role: domain.RoleAdminTenant, IsActive: true}

	_ = repo.CreateUser(nil, &globalAdmin)
	_ = repo.CreateUser(nil, &adminT1)
	_ = repo.CreateUser(nil, &opT1)
	_ = repo.CreateUser(nil, &adminT2)

	tokenGlobal, _ := jwtSvc.GenerateToken(&globalAdmin)
	tokenAdminT1, _ := jwtSvc.GenerateToken(&adminT1)
	tokenOpT1, _ := jwtSvc.GenerateToken(&opT1)

	// 1. Operador tentando acessar rota de usuários (deve retornar 403 Forbidden)
	t.Run("Operador não pode acessar /admin/users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+tokenOpT1)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("esperado status 403 Forbidden para Operador, obtido %d", w.Code)
		}
	})

	// 2. Admin Tenant 1 listando usuários (só deve ver usuários do Tenant 1)
	t.Run("Admin Tenant 1 só lista usuários do Tenant 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+tokenAdminT1)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("esperado 200 OK, obtido %d", w.Code)
		}

		var res struct {
			Success bool          `json:"success"`
			Data    []domain.User `json:"data"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &res)

		if len(res.Data) != 2 {
			t.Errorf("esperado 2 usuários para Tenant 1, obtido %d", len(res.Data))
		}
		for _, u := range res.Data {
			if u.TenantID == nil || *u.TenantID != t1ID {
				t.Errorf("usuário de outro tenant listado indevidamente: %s", u.Email)
			}
		}
	})

	// 3. Admin Tenant 1 tentando acessar usuário do Tenant 2 via ID (deve retornar 403 Forbidden)
	t.Run("Admin Tenant 1 não pode ver/gerenciar usuário do Tenant 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/"+adminT2.ID.String(), nil)
		req.Header.Set("Authorization", "Bearer "+tokenAdminT1)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("esperado 403 Forbidden ao tentar acessar usuário de outro tenant, obtido %d", w.Code)
		}
	})

	// 4. Admin Tenant tentando criar ADMIN_GLOBAL (deve retornar 403 Forbidden)
	t.Run("Admin Tenant não pode criar ADMIN_GLOBAL", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"name":     "Hacker Global",
			"email":    "hacker@plataforma.com",
			"password": "password123",
			"role":     "ADMIN_GLOBAL",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+tokenAdminT1)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("esperado 403 Forbidden ao tentar criar ADMIN_GLOBAL como Tenant Admin, obtido %d", w.Code)
		}
	})

	// 5. ADMIN_GLOBAL acessando listagem de tenants e dashboard global
	t.Run("ADMIN_GLOBAL pode acessar /admin/tenants e /admin/global-dashboard", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/tenants", nil)
		req.Header.Set("Authorization", "Bearer "+tokenGlobal)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("esperado 200 OK para ADMIN_GLOBAL em /admin/tenants, obtido %d", w.Code)
		}

		reqDash := httptest.NewRequest(http.MethodGet, "/api/v1/admin/global-dashboard", nil)
		reqDash.Header.Set("Authorization", "Bearer "+tokenGlobal)
		wDash := httptest.NewRecorder()
		engine.ServeHTTP(wDash, reqDash)

		if wDash.Code != http.StatusOK {
			t.Errorf("esperado 200 OK para ADMIN_GLOBAL em /admin/global-dashboard, obtido %d", wDash.Code)
		}
	})

	// 6. Tenant Admin pode acessar /admin/dashboard/analytics (Padrão Figma SAAS)
	t.Run("Tenant Admin pode acessar /admin/dashboard/analytics", func(t *testing.T) {
		plan := domain.Plan{ID: uuid.New(), Name: "Plano Pro", IsActive: true}
		_ = repo.DB().Create(&plan)
		subT1 := domain.Subscription{
			ID:       uuid.New(),
			TenantID: t1ID,
			PlanID:   plan.ID,
			Status:   domain.SubscriptionStatusActive,
		}
		if err := repo.DB().Create(&subT1).Error; err != nil {
			t.Fatalf("falha ao criar sub: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard/analytics", nil)
		req.Header.Set("Authorization", "Bearer "+tokenAdminT1)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("esperado 200 OK para Tenant Admin em /admin/dashboard/analytics, obtido %d, body: %s", w.Code, w.Body.String())
		}
	})
}
