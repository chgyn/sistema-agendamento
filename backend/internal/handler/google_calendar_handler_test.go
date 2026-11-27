package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/handler"
	"github.com/sistema-agendamento/backend/internal/queue"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/hash"
	"github.com/sistema-agendamento/backend/pkg/jwt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type stubCalendarProvider struct{}

func (s *stubCalendarProvider) BuildAuthURL(clientID, redirectURI, state string) string {
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + clientID + "&state=" + state
}
func (s *stubCalendarProvider) ExchangeCode(ctx context.Context, clientID, clientSecret, redirectURI, code string) (*domain.OAuthTokenResponse, error) {
	return &domain.OAuthTokenResponse{AccessToken: "acc", RefreshToken: "ref", ExpiresIn: 3600}, nil
}
func (s *stubCalendarProvider) RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*domain.OAuthTokenResponse, error) {
	return &domain.OAuthTokenResponse{AccessToken: "acc2", ExpiresIn: 3600}, nil
}
func (s *stubCalendarProvider) RevokeToken(ctx context.Context, token string) error { return nil }
func (s *stubCalendarProvider) GetUserEmail(ctx context.Context, accessToken string) (string, error) {
	return "contato@salao.com", nil
}
func (s *stubCalendarProvider) ListCalendars(ctx context.Context, accessToken string) ([]domain.GoogleCalendarItemDTO, error) {
	return []domain.GoogleCalendarItemDTO{{ID: "primary", Summary: "Principal", Primary: true, AccessRole: "owner"}}, nil
}
func (s *stubCalendarProvider) CreateEvent(ctx context.Context, accessToken, calendarID string, event domain.CalendarEventInput) (string, error) {
	return "evt-1", nil
}
func (s *stubCalendarProvider) UpdateEvent(ctx context.Context, accessToken, calendarID, eventID string, event domain.CalendarEventInput) (string, error) {
	return eventID, nil
}
func (s *stubCalendarProvider) DeleteEvent(ctx context.Context, accessToken, calendarID, eventID string) error {
	return nil
}

type noopEnqueuer struct{}

func (n *noopEnqueuer) EnqueueGoogleCalendarSync(ctx context.Context, payload queue.GoogleCalendarSyncPayload) error {
	return nil
}

type stubOAuthStateStore struct {
	mu   sync.Mutex
	data map[string]string
}

func (s *stubOAuthStateStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = make(map[string]string)
	}
	s.data[key] = value
	return nil
}

func (s *stubOAuthStateStore) GetDel(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	delete(s.data, key)
	return val, nil
}

func setupGoogleCalendarHandlerTestApp(t *testing.T) (*gin.Engine, *postgres.Repository, *jwt.JWTService) {
	gin.SetMode(gin.TestMode)
	dbName := fmt.Sprintf("file:mem_gcal_handler_%d?mode=memory&cache=shared&_busy_timeout=10000", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("falha ao conectar sqlite teste: %v", err)
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	_ = db.AutoMigrate(
		&domain.Tenant{},
		&domain.User{},
		&domain.Professional{},
		&domain.Service{},
		&domain.Customer{},
		&domain.Appointment{},
		&domain.Plan{},
		&domain.Subscription{},
		&domain.SubscriptionInvoice{},
		&domain.SubscriptionAuditLog{},
		&domain.TenantGoogleCalendarConfig{},
		&domain.AppointmentCalendarSync{},
	)

	repo := postgres.NewRepository(db)
	jwtSecret := "super-secret-test-key-32-chars-long"
	encKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	jwtSvc := jwt.NewJWTService(jwtSecret, 24)
	gcalSvc := service.NewGoogleCalendarService(repo, repo, &stubCalendarProvider{}, &noopEnqueuer{}, &stubOAuthStateStore{}, jwtSecret, encKey, "https://app.test.com/api/v1/integrations/google-calendar/callback")
	gcalHandler := handler.NewGoogleCalendarHandler(gcalSvc, "https://app.test.com")
	subService := service.NewSubscriptionService(repo, nil)
	planService := service.NewPlanService(repo)
	cfg := &config.Config{AsaasWebhookSecret: "test"}

	r := gin.New()
	handler.SetupRoutes(handler.RouterConfig{
		Engine:                r,
		JWTService:            jwtSvc,
		Repo:                  repo,
		UserHandler:           handler.NewUserHandler(repo),
		TenantHandler:         handler.NewTenantHandler(repo),
		PlanHandler:           handler.NewPlanHandler(planService),
		SubscriptionHandler:   handler.NewSubscriptionHandler(subService),
		WebhookHandler:        handler.NewWebhookHandler(cfg, nil, subService),
		DashboardHandler:      handler.NewDashboardHandler(service.NewTenantService(repo)),
		GoogleCalendarHandler: gcalHandler,
	})

	return r, repo, jwtSvc
}

func TestGoogleCalendarHandler_RBAC_Callback_And_Isolation(t *testing.T) {
	engine, repo, jwtSvc := setupGoogleCalendarHandlerTestApp(t)
	passHash, _ := hash.HashPassword("123456")

	planID := uuid.New()
	_ = repo.CreatePlan(context.Background(), &domain.Plan{ID: planID, Name: "Pro", Price: 99, IsActive: true})

	t1ID := uuid.New()
	t2ID := uuid.New()
	_ = repo.CreateTenant(context.Background(), &domain.Tenant{ID: t1ID, Slug: "barber-1", Name: "Barber 1", IsActive: true})
	_ = repo.CreateTenant(context.Background(), &domain.Tenant{ID: t2ID, Slug: "barber-2", Name: "Barber 2", IsActive: true})
	_ = repo.CreateSubscription(context.Background(), &domain.Subscription{ID: uuid.New(), TenantID: t1ID, PlanID: planID, Status: domain.SubscriptionStatusActive, BillingCycle: domain.CycleMonthly, Price: 99})
	_ = repo.CreateSubscription(context.Background(), &domain.Subscription{ID: uuid.New(), TenantID: t2ID, PlanID: planID, Status: domain.SubscriptionStatusActive, BillingCycle: domain.CycleMonthly, Price: 99})

	adminT1 := domain.User{ID: uuid.New(), TenantID: &t1ID, Name: "Admin T1", Email: "a1@gcal.com", PasswordHash: passHash, Role: domain.RoleAdminTenant, IsActive: true}
	opT1 := domain.User{ID: uuid.New(), TenantID: &t1ID, Name: "Op T1", Email: "op1@gcal.com", PasswordHash: passHash, Role: domain.RoleOperator, IsActive: true}
	adminT2 := domain.User{ID: uuid.New(), TenantID: &t2ID, Name: "Admin T2", Email: "a2@gcal.com", PasswordHash: passHash, Role: domain.RoleAdminTenant, IsActive: true}
	_ = repo.CreateUser(context.Background(), &adminT1)
	_ = repo.CreateUser(context.Background(), &opT1)
	_ = repo.CreateUser(context.Background(), &adminT2)

	tokAdminT1, _ := jwtSvc.GenerateToken(&adminT1)
	tokOpT1, _ := jwtSvc.GenerateToken(&opT1)
	tokAdminT2, _ := jwtSvc.GenerateToken(&adminT2)

	// 1. OPERATOR não pode acessar Google Calendar
	reqOp, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/integrations/google-calendar/status", nil)
	reqOp.Header.Set("Authorization", "Bearer "+tokOpT1)
	wOp := httptest.NewRecorder()
	engine.ServeHTTP(wOp, reqOp)
	if wOp.Code != http.StatusForbidden {
		t.Fatalf("esperava 403 para OPERATOR, obteve %d", wOp.Code)
	}

	// 2. ADMIN_TENANT salva credenciais e gera auth-url
	bodyBytes, _ := json.Marshal(map[string]string{
		"google_client_id":     "cid-t1",
		"google_client_secret": "csec-t1",
	})
	reqSave, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/integrations/google-calendar/settings", bytes.NewReader(bodyBytes))
	reqSave.Header.Set("Authorization", "Bearer "+tokAdminT1)
	reqSave.Header.Set("Content-Type", "application/json")
	wSave := httptest.NewRecorder()
	engine.ServeHTTP(wSave, reqSave)
	if wSave.Code != http.StatusOK {
		t.Fatalf("esperava 200 ao salvar credenciais, obteve %d: %s", wSave.Code, wSave.Body.String())
	}
	if strings.Contains(wSave.Body.String(), "csec-t1") {
		t.Fatalf("resposta da API vazou o client secret!")
	}

	reqAuth, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/integrations/google-calendar/auth-url", nil)
	reqAuth.Header.Set("Authorization", "Bearer "+tokAdminT1)
	wAuth := httptest.NewRecorder()
	engine.ServeHTTP(wAuth, reqAuth)
	if wAuth.Code != http.StatusOK {
		t.Fatalf("esperava 200 ao gerar auth-url, obteve %d: %s", wAuth.Code, wAuth.Body.String())
	}
	var authResp struct {
		Data struct {
			AuthURL string `json:"auth_url"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wAuth.Body.Bytes(), &authResp)
	state := strings.Split(authResp.Data.AuthURL, "state=")[1]

	// 3. Callback público redireciona para frontend com status=connected
	reqCb, _ := http.NewRequest(http.MethodGet, "/api/v1/integrations/google-calendar/callback?code=ok&state="+state, nil)
	wCb := httptest.NewRecorder()
	engine.ServeHTTP(wCb, reqCb)
	if wCb.Code != http.StatusFound || wCb.Header().Get("Location") != "https://app.test.com/admin/google-calendar?status=connected" {
		t.Fatalf("redirect callback inesperado: code=%d location=%s", wCb.Code, wCb.Header().Get("Location"))
	}

	// 4. Callback com erro redireciona com reason sanitizado
	reqCbErr, _ := http.NewRequest(http.MethodGet, "/api/v1/integrations/google-calendar/callback?error=access_denied", nil)
	wCbErr := httptest.NewRecorder()
	engine.ServeHTTP(wCbErr, reqCbErr)
	if wCbErr.Code != http.StatusFound || !strings.Contains(wCbErr.Header().Get("Location"), "status=error") || !strings.Contains(wCbErr.Header().Get("Location"), "reason=access_denied") {
		t.Fatalf("redirect de erro inesperado: code=%d location=%s", wCbErr.Code, wCbErr.Header().Get("Location"))
	}

	// 5. Retry cross-tenant bloqueado (404)
	aptT1 := uuid.New()
	_ = repo.DB().Create(&domain.Appointment{
		ID: aptT1, TenantID: t1ID, ProfessionalID: uuid.New(), ServiceID: uuid.New(), CustomerID: uuid.New(),
		StartAt: time.Now(), EndAt: time.Now().Add(time.Hour), DurationMinutes: 60, TotalPrice: 50, Status: domain.StatusConfirmed,
	}).Error

	reqCross, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/integrations/google-calendar/syncs/"+aptT1.String()+"/retry", nil)
	reqCross.Header.Set("Authorization", "Bearer "+tokAdminT2)
	wCross := httptest.NewRecorder()
	engine.ServeHTTP(wCross, reqCross)
	if wCross.Code != http.StatusNotFound {
		t.Fatalf("esperava 404 para retry cross-tenant, obteve %d", wCross.Code)
	}
}
