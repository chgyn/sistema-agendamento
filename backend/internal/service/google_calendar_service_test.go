package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/integrations/googlecalendar"
	"github.com/sistema-agendamento/backend/internal/queue"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testEncKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
const testJWTSecret = "super-secret-jwt-key-for-google-calendar-tests"

func setupServiceTestDB(t *testing.T) (*gorm.DB, *postgres.Repository) {
	dbName := fmt.Sprintf("file:mem_gcal_svc_%d?mode=memory&cache=shared&_busy_timeout=10000", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("falha ao abrir sqlite em memória: %v", err)
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	err = db.AutoMigrate(
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
		&domain.TenantGoogleCalendarConfig{},
		&domain.AppointmentCalendarSync{},
	)
	if err != nil {
		t.Fatalf("falha no automigrate: %v", err)
	}

	return db, postgres.NewRepository(db)
}

type mockCalendarProvider struct {
	authURL          string
	exchangeFn       func(ctx context.Context, clientID, clientSecret, redirectURI, code string) (*domain.OAuthTokenResponse, error)
	refreshFn        func(ctx context.Context, clientID, clientSecret, refreshToken string) (*domain.OAuthTokenResponse, error)
	revokeFn         func(ctx context.Context, token string) error
	userEmailFn      func(ctx context.Context, accessToken string) (string, error)
	listCalendarsFn  func(ctx context.Context, accessToken string) ([]domain.GoogleCalendarItemDTO, error)
	createEventFn    func(ctx context.Context, accessToken, calendarID string, event domain.CalendarEventInput) (string, error)
	updateEventFn    func(ctx context.Context, accessToken, calendarID, eventID string, event domain.CalendarEventInput) (string, error)
	deleteEventFn    func(ctx context.Context, accessToken, calendarID, eventID string) error
	createCalls      []domain.CalendarEventInput
	updateCalls      []string
	deleteCalls      []string
}

func (m *mockCalendarProvider) BuildAuthURL(clientID, redirectURI, state string) string {
	if m.authURL != "" {
		return m.authURL + "?state=" + state + "&client_id=" + clientID
	}
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + clientID + "&redirect_uri=" + redirectURI + "&state=" + state
}

func (m *mockCalendarProvider) ExchangeCode(ctx context.Context, clientID, clientSecret, redirectURI, code string) (*domain.OAuthTokenResponse, error) {
	if m.exchangeFn != nil {
		return m.exchangeFn(ctx, clientID, clientSecret, redirectURI, code)
	}
	return &domain.OAuthTokenResponse{
		AccessToken:  "access-token-plain",
		RefreshToken: "refresh-token-plain",
		ExpiresIn:    3600,
	}, nil
}

func (m *mockCalendarProvider) RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*domain.OAuthTokenResponse, error) {
	if m.refreshFn != nil {
		return m.refreshFn(ctx, clientID, clientSecret, refreshToken)
	}
	return &domain.OAuthTokenResponse{
		AccessToken: "refreshed-access-token",
		ExpiresIn:   3600,
	}, nil
}

func (m *mockCalendarProvider) RevokeToken(ctx context.Context, token string) error {
	if m.revokeFn != nil {
		return m.revokeFn(ctx, token)
	}
	return nil
}

func (m *mockCalendarProvider) GetUserEmail(ctx context.Context, accessToken string) (string, error) {
	if m.userEmailFn != nil {
		return m.userEmailFn(ctx, accessToken)
	}
	return "barbearia@gmail.com", nil
}

func (m *mockCalendarProvider) ListCalendars(ctx context.Context, accessToken string) ([]domain.GoogleCalendarItemDTO, error) {
	if m.listCalendarsFn != nil {
		return m.listCalendarsFn(ctx, accessToken)
	}
	return []domain.GoogleCalendarItemDTO{{ID: "primary", Summary: "Principal", Primary: true, AccessRole: "owner"}}, nil
}

func (m *mockCalendarProvider) CreateEvent(ctx context.Context, accessToken, calendarID string, event domain.CalendarEventInput) (string, error) {
	m.createCalls = append(m.createCalls, event)
	if m.createEventFn != nil {
		return m.createEventFn(ctx, accessToken, calendarID, event)
	}
	if event.EventID != "" {
		return event.EventID, nil
	}
	return "evt-created", nil
}

func (m *mockCalendarProvider) UpdateEvent(ctx context.Context, accessToken, calendarID, eventID string, event domain.CalendarEventInput) (string, error) {
	m.updateCalls = append(m.updateCalls, calendarID+":"+eventID)
	if m.updateEventFn != nil {
		return m.updateEventFn(ctx, accessToken, calendarID, eventID, event)
	}
	return eventID, nil
}

func (m *mockCalendarProvider) DeleteEvent(ctx context.Context, accessToken, calendarID, eventID string) error {
	m.deleteCalls = append(m.deleteCalls, calendarID+":"+eventID)
	if m.deleteEventFn != nil {
		return m.deleteEventFn(ctx, accessToken, calendarID, eventID)
	}
	return nil
}

type mockEnqueuer struct {
	payloads []queue.GoogleCalendarSyncPayload
	err      error
}

func (m *mockEnqueuer) EnqueueGoogleCalendarSync(ctx context.Context, payload queue.GoogleCalendarSyncPayload) error {
	if m.err != nil {
		return m.err
	}
	m.payloads = append(m.payloads, payload)
	return nil
}

type mockOAuthStateStore struct {
	mu   sync.Mutex
	data map[string]string
}

func newMockOAuthStateStore() *mockOAuthStateStore {
	return &mockOAuthStateStore{data: make(map[string]string)}
}

func (s *mockOAuthStateStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *mockOAuthStateStore) GetDel(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if !ok {
		return "", errors.New("redis: nil")
	}
	delete(s.data, key)
	return val, nil
}

func seedAppointmentFixtures(t *testing.T, db *gorm.DB, tenantID uuid.UUID) (proID, svcID, custID uuid.UUID) {
	proID, svcID, custID = uuid.New(), uuid.New(), uuid.New()
	if err := db.Create(&domain.Tenant{ID: tenantID, Slug: "t-" + tenantID.String()[:8], Name: "Barbearia Teste", IsActive: true}).Error; err != nil {
		t.Fatalf("erro ao criar tenant: %v", err)
	}
	if err := db.Create(&domain.Professional{ID: proID, TenantID: tenantID, Name: "Marcos Barber", IsActive: true}).Error; err != nil {
		t.Fatalf("erro ao criar profissional: %v", err)
	}
	if err := db.Create(&domain.Service{ID: svcID, TenantID: tenantID, Name: "Corte Degradê", DurationMinutes: 45, Price: 60, IsActive: true}).Error; err != nil {
		t.Fatalf("erro ao criar serviço: %v", err)
	}
	if err := db.Create(&domain.Customer{ID: custID, TenantID: tenantID, Name: "João Pedro", Phone: "11999887766"}).Error; err != nil {
		t.Fatalf("erro ao criar cliente: %v", err)
	}
	return proID, svcID, custID
}

func TestGoogleCalendarService_OAuthStateSecurityAndEncryption(t *testing.T) {
	_, repo := setupServiceTestDB(t)
	ctx := context.Background()
	provider := &mockCalendarProvider{}
	stateStore := newMockOAuthStateStore()
	svc := NewGoogleCalendarService(repo, repo, provider, nil, stateStore, testJWTSecret, testEncKey, "https://app.test.com/api/v1/integrations/google-calendar/callback")

	tenant1 := uuid.New()
	tenant2 := uuid.New()
	user1 := uuid.New()

	// Fail-closed quando stateStore for nil
	nilStoreSvc := NewGoogleCalendarService(repo, repo, provider, nil, nil, testJWTSecret, testEncKey, "https://app.test.com/api/v1/integrations/google-calendar/callback")
	cidNil, csecNil := "cid-nil", "csec-nil"
	tenantNil := uuid.New()
	_, _ = nilStoreSvc.UpdateSettings(ctx, tenantNil, domain.UpdateGoogleCalendarSettingsDTO{
		GoogleClientID:     &cidNil,
		GoogleClientSecret: &csecNil,
	})
	if _, err := nilStoreSvc.GetAuthURL(ctx, tenantNil, user1); err == nil {
		t.Fatalf("esperava erro fail-closed em GetAuthURL com stateStore=nil")
	}
	if _, err := nilStoreSvc.HandleCallback(ctx, "code-1", "qualquer.state"); !errors.Is(err, domain.ErrInvalidOAuthState) {
		t.Fatalf("esperava ErrInvalidOAuthState fail-closed em HandleCallback com stateStore=nil, obteve %v", err)
	}

	// Sem credenciais -> GetAuthURL falha com ErrGoogleCalendarNotConfigured
	if _, err := svc.GetAuthURL(ctx, tenant1, user1); !errors.Is(err, domain.ErrGoogleCalendarNotConfigured) {
		t.Fatalf("esperava ErrGoogleCalendarNotConfigured, obteve %v", err)
	}

	cid := "client-id-t1"
	csec := "super-secret-oauth-client-secret"
	_, err := svc.UpdateSettings(ctx, tenant1, domain.UpdateGoogleCalendarSettingsDTO{
		GoogleClientID:     &cid,
		GoogleClientSecret: &csec,
	})
	if err != nil {
		t.Fatalf("UpdateSettings falhou: %v", err)
	}

	// Verifica que o ClientSecret está criptografado no banco e ausente no DTO/JSON
	rawCfg, err := repo.GetGoogleCalendarConfigByTenantID(ctx, tenant1)
	if err != nil || rawCfg == nil {
		t.Fatalf("config não encontrada: %v", err)
	}
	if rawCfg.GoogleClientSecretEncrypted == "" || rawCfg.GoogleClientSecretEncrypted == csec {
		t.Fatalf("client secret não foi cifrado no banco!")
	}
	rawJSON, _ := json.Marshal(rawCfg)
	if strings.Contains(string(rawJSON), csec) || strings.Contains(string(rawJSON), rawCfg.GoogleClientSecretEncrypted) {
		t.Fatalf("client secret vazou no JSON do model: %s", string(rawJSON))
	}

	// Gera AuthURL e extrai state válido
	authURL, err := svc.GetAuthURL(ctx, tenant1, user1)
	if err != nil {
		t.Fatalf("GetAuthURL falhou: %v", err)
	}
	parts := strings.Split(authURL, "state=")
	if len(parts) != 2 {
		t.Fatalf("state não encontrado na URL: %s", authURL)
	}
	validState := strings.Split(parts[1], "&")[0]

	// 1) Adulteração de assinatura HMAC
	tamperedState := validState + "x"
	if _, err := svc.HandleCallback(ctx, "code-1", tamperedState); !errors.Is(err, domain.ErrInvalidOAuthState) {
		t.Fatalf("esperava ErrInvalidOAuthState para assinatura adulterada, obteve %v", err)
	}

	// 2) Adulteração cross-tenant no payload (substitui tenant1 por tenant2 mantendo assinatura de tenant1)
	stateParts := strings.Split(validState, ".")
	tamperedCrossPayload, _ := json.Marshal(oauthStatePayload{
		TenantID: tenant2,
		UserID:   user1,
		Nonce:    "cross-tenant-nonce",
		Exp:      time.Now().Add(5 * time.Minute).Unix(),
	})
	fakePayload := base64URLEncode(tamperedCrossPayload)
	if _, err := svc.HandleCallback(ctx, "code-1", fakePayload+"."+stateParts[1]); !errors.Is(err, domain.ErrInvalidOAuthState) {
		t.Fatalf("esperava ErrInvalidOAuthState para payload adulterado, obteve %v", err)
	}

	// 3) State expirado manualmente
	expiredState, _ := svc.generateOAuthState(ctx, tenant1, user1)
	expParts := strings.Split(expiredState, ".")
	// Cria payload expirado assinado com a mesma chave para testar checagem de exp
	expiredJSON, _ := json.Marshal(oauthStatePayload{
		TenantID: tenant1,
		UserID:   user1,
		Nonce:    "nonce123",
		Exp:      time.Now().Add(-1 * time.Minute).Unix(),
	})
	_ = expParts
	expiredB64 := base64URLEncode(expiredJSON)
	expiredSigned := expiredB64 + "." + svc.signStatePayload(expiredB64)
	if _, err := svc.HandleCallback(ctx, "code-1", expiredSigned); !errors.Is(err, domain.ErrInvalidOAuthState) {
		t.Fatalf("esperava ErrInvalidOAuthState para state expirado, obteve %v", err)
	}

	// 4) Callback válido conecta tenant1 e cifra access/refresh tokens
	gotTenant, err := svc.HandleCallback(ctx, "code-1", validState)
	if err != nil || gotTenant != tenant1 {
		t.Fatalf("HandleCallback falhou: tenant=%s err=%v", gotTenant, err)
	}

	// 5) Tentativa de replay do mesmo state já consumido deve ser rejeitada
	if _, err := svc.HandleCallback(ctx, "code-1", validState); !errors.Is(err, domain.ErrInvalidOAuthState) {
		t.Fatalf("esperava ErrInvalidOAuthState em tentativa de replay do mesmo state, obteve %v", err)
	}

	// Tenant2 continua desconectado (isolamento multi-tenant)
	statusT2, _ := svc.GetStatus(ctx, tenant2)
	if statusT2.Status != domain.GoogleCalendarStatusDisconnected || statusT2.HasCredentials {
		t.Fatalf("vazamento cross-tenant detectado: %+v", statusT2)
	}

	connectedCfg, _ := repo.GetGoogleCalendarConfigByTenantID(ctx, tenant1)
	if connectedCfg.Status != domain.GoogleCalendarStatusConnected {
		t.Fatalf("esperava CONNECTED, obteve %s", connectedCfg.Status)
	}
	if connectedCfg.AccessTokenEncrypted == "access-token-plain" || connectedCfg.RefreshTokenEncrypted == "refresh-token-plain" {
		t.Fatalf("tokens não foram cifrados com AES-256-GCM!")
	}

	dto, err := svc.GetStatus(ctx, tenant1)
	if err != nil {
		t.Fatalf("GetStatus falhou: %v", err)
	}
	dtoJSON, _ := json.Marshal(dto)
	if strings.Contains(string(dtoJSON), csec) || strings.Contains(string(dtoJSON), "access-token-plain") || strings.Contains(string(dtoJSON), "refresh-token-plain") {
		t.Fatalf("DTO expôs segredos no JSON: %s", string(dtoJSON))
	}

	// Alterar ClientID deve limpar tokens e desconectar
	newCID := "client-id-t1-novo"
	updatedDTO, err := svc.UpdateSettings(ctx, tenant1, domain.UpdateGoogleCalendarSettingsDTO{
		GoogleClientID: &newCID,
	})
	if err != nil {
		t.Fatalf("UpdateSettings ao alterar client_id falhou: %v", err)
	}
	if updatedDTO.Status != domain.GoogleCalendarStatusDisconnected || updatedDTO.IsEnabled {
		t.Fatalf("esperava DISCONNECTED e IsEnabled=false ao alterar credenciais, obteve %+v", updatedDTO)
	}
}

func TestGoogleCalendarService_SyncIdempotencyAndReconciliation(t *testing.T) {
	db, repo := setupServiceTestDB(t)
	ctx := context.Background()
	tenantID := uuid.New()
	proID, svcID, custID := seedAppointmentFixtures(t, db, tenantID)

	apt := &domain.Appointment{
		ID:              uuid.New(),
		TenantID:        tenantID,
		ProfessionalID:  proID,
		ServiceID:       svcID,
		CustomerID:      custID,
		StartAt:         time.Now().Add(2 * time.Hour),
		EndAt:           time.Now().Add(2*time.Hour + 45*time.Minute),
		DurationMinutes: 45,
		TotalPrice:      60,
		Status:          domain.StatusConfirmed,
		Notes:           "Preferência tesoura",
	}
	if err := db.Create(apt).Error; err != nil {
		t.Fatalf("erro ao criar agendamento: %v", err)
	}

	provider := &mockCalendarProvider{}
	stateStore := newMockOAuthStateStore()
	svc := NewGoogleCalendarService(repo, repo, provider, &mockEnqueuer{}, stateStore, testJWTSecret, testEncKey, "https://cb.test")

	// Configura e conecta tenant
	cid, csec := "cid-sync", "csec-sync"
	_, _ = svc.UpdateSettings(ctx, tenantID, domain.UpdateGoogleCalendarSettingsDTO{
		GoogleClientID:     &cid,
		GoogleClientSecret: &csec,
	})
	state, _ := svc.generateOAuthState(ctx, tenantID, uuid.New())
	_, _ = svc.HandleCallback(ctx, "code-ok", state)

	// Com is_enabled=false -> no-op silencioso
	if err := svc.SyncAppointment(ctx, tenantID, apt.ID, domain.CalendarSyncActionCreate); err != nil {
		t.Fatalf("esperava no-op quando desativado, obteve %v", err)
	}
	if len(provider.createCalls) != 0 {
		t.Fatalf("não deveria chamar CreateEvent com integração desativada")
	}

	// Habilita integração
	enabled := true
	_, err := svc.UpdateSettings(ctx, tenantID, domain.UpdateGoogleCalendarSettingsDTO{IsEnabled: &enabled})
	if err != nil {
		t.Fatalf("erro ao ativar integração: %v", err)
	}

	// 1) Primeira criação -> usa ID determinístico "af..." e PrivateProps
	if err := svc.SyncAppointment(ctx, tenantID, apt.ID, domain.CalendarSyncActionCreate); err != nil {
		t.Fatalf("SyncAppointment create falhou: %v", err)
	}
	if len(provider.createCalls) != 1 {
		t.Fatalf("esperava 1 CreateEvent, obteve %d", len(provider.createCalls))
	}
	detID := deterministicGoogleEventID(apt.ID)
	if provider.createCalls[0].EventID != detID {
		t.Errorf("esperava ID determinístico %s, obteve %s", detID, provider.createCalls[0].EventID)
	}
	if provider.createCalls[0].PrivateProps["agende_facil_appointment_id"] != apt.ID.String() {
		t.Errorf("PrivateProps ausente")
	}

	// 2) Reexecução (ou reagendamento) -> chama UpdateEvent (zero duplicidade)
	if err := svc.SyncAppointment(ctx, tenantID, apt.ID, domain.CalendarSyncActionUpdate); err != nil {
		t.Fatalf("SyncAppointment update falhou: %v", err)
	}
	if len(provider.createCalls) != 1 || len(provider.updateCalls) != 1 {
		t.Fatalf("esperava 1 create e 1 update, obteve create=%d update=%d", len(provider.createCalls), len(provider.updateCalls))
	}

	// 3) Task DELETE atrasada chegando quando agendamento está CONFIRMED no banco (ex: reagendado/reconfirmado) NÃO deve apagar o evento
	if err := svc.SyncAppointment(ctx, tenantID, apt.ID, domain.CalendarSyncActionDelete); err != nil {
		t.Fatalf("SyncAppointment com action DELETE atrasada falhou: %v", err)
	}
	if len(provider.deleteCalls) != 0 {
		t.Fatalf("task DELETE atrasada não deveria apagar evento de agendamento CONFIRMED, obteve deletes=%v", provider.deleteCalls)
	}
	if len(provider.updateCalls) != 2 {
		t.Fatalf("esperava reconciliação via UpdateEvent para agendamento CONFIRMED (total 2 updates), obteve %d", len(provider.updateCalls))
	}

	// 4) Mudança de calendário -> deleta no antigo ("primary") e cria no novo ("cal-vip")
	newCal := "cal-vip"
	_, _ = svc.UpdateSettings(ctx, tenantID, domain.UpdateGoogleCalendarSettingsDTO{CalendarID: &newCal})
	if err := svc.SyncAppointment(ctx, tenantID, apt.ID, domain.CalendarSyncActionUpdate); err != nil {
		t.Fatalf("SyncAppointment mudança de calendário falhou: %v", err)
	}
	if len(provider.deleteCalls) != 1 || provider.deleteCalls[0] != "primary:"+detID {
		t.Fatalf("esperava delete no calendário antigo primary, obteve %v", provider.deleteCalls)
	}
	if len(provider.createCalls) != 2 {
		t.Fatalf("esperava criação no novo calendário cal-vip, obteve %d creates", len(provider.createCalls))
	}

	// 5) Cancelamento no banco -> remove evento do calendário atual e marca DELETED
	_ = repo.UpdateAppointmentStatus(ctx, tenantID, apt.ID, domain.StatusCancelled, "Admin", "Cliente desistiu")
	if err := svc.SyncAppointment(ctx, tenantID, apt.ID, domain.CalendarSyncActionDelete); err != nil {
		t.Fatalf("SyncAppointment delete falhou: %v", err)
	}
	if len(provider.deleteCalls) != 2 || !strings.HasPrefix(provider.deleteCalls[1], "cal-vip:") {
		t.Fatalf("esperava delete em cal-vip, obteve %v", provider.deleteCalls)
	}
	syncRec, _ := repo.GetAppointmentCalendarSync(ctx, tenantID, apt.ID, domain.CalendarProviderGoogle)
	if syncRec.SyncStatus != domain.CalendarSyncStatusDeleted {
		t.Fatalf("esperava SyncStatus=DELETED, obteve %s", syncRec.SyncStatus)
	}

	// 5) 409 Conflict na criação -> converte automaticamente em UpdateEvent
	apt2 := &domain.Appointment{
		ID: uuid.New(), TenantID: tenantID, ProfessionalID: proID, ServiceID: svcID, CustomerID: custID,
		StartAt: time.Now().Add(4 * time.Hour), EndAt: time.Now().Add(5 * time.Hour),
		DurationMinutes: 60, TotalPrice: 60, Status: domain.StatusConfirmed,
	}
	_ = db.Create(apt2).Error
	provider.createEventFn = func(ctx context.Context, accessToken, calendarID string, event domain.CalendarEventInput) (string, error) {
		return "", googlecalendar.ErrEventConflict
	}
	if err := svc.SyncAppointment(ctx, tenantID, apt2.ID, domain.CalendarSyncActionCreate); err != nil {
		t.Fatalf("esperava recuperação automática de 409 Conflict, obteve: %v", err)
	}

	// 6) Token expirado/revogado (ErrGoogleCalendarAuthExpired) -> marca EXPIRED e retorna asynq.SkipRetry
	apt3 := &domain.Appointment{
		ID: uuid.New(), TenantID: tenantID, ProfessionalID: proID, ServiceID: svcID, CustomerID: custID,
		StartAt: time.Now().Add(6 * time.Hour), EndAt: time.Now().Add(7 * time.Hour),
		DurationMinutes: 60, TotalPrice: 60, Status: domain.StatusConfirmed,
	}
	_ = db.Create(apt3).Error
	// Força expiração imediata do access token para disparar RefreshToken
	expPast := time.Now().Add(-5 * time.Minute)
	cfgNow, _ := repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	cfgNow.TokenExpiresAt = &expPast
	_ = repo.UpsertGoogleCalendarConfig(ctx, cfgNow)
	provider.refreshFn = func(ctx context.Context, clientID, clientSecret, refreshToken string) (*domain.OAuthTokenResponse, error) {
		return nil, domain.ErrGoogleCalendarAuthExpired
	}
	err = svc.SyncAppointment(ctx, tenantID, apt3.ID, domain.CalendarSyncActionCreate)
	if !errors.Is(err, domain.ErrGoogleCalendarAuthExpired) || !errors.Is(err, asynq.SkipRetry) {
		t.Fatalf("esperava ErrGoogleCalendarAuthExpired + asynq.SkipRetry, obteve %v", err)
	}
	cfgAfter, _ := repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if cfgAfter.Status != domain.GoogleCalendarStatusExpired {
		t.Fatalf("esperava status EXPIRED após invalid_grant, obteve %s", cfgAfter.Status)
	}
}

func TestCentralization_AppointmentServiceAndAIChatService(t *testing.T) {
	db, repo := setupServiceTestDB(t)
	ctx := context.Background()
	tenantID := uuid.New()
	proID, svcID, _ := seedAppointmentFixtures(t, db, tenantID)

	enqueuer := &mockEnqueuer{}
	aptSvc := NewAppointmentService(repo, nil)
	aptSvc.SetCalendarSyncDependencies(enqueuer, repo)

	// 0. Sem integração Google Calendar ativa, NÃO deve enfileirar nenhuma task no Asynq
	startAt := time.Now().Add(24 * time.Hour)
	aptNoGcal, err := aptSvc.CreateAppointment(ctx, CreateAppointmentDTO{
		TenantID:       tenantID,
		ServiceID:      svcID,
		ProfessionalID: proID,
		StartAt:        startAt,
		CustomerName:   "Sem Agenda",
		CustomerPhone:  "11900001111",
	})
	if err != nil {
		t.Fatalf("CreateAppointment sem Google Calendar falhou: %v", err)
	}
	if len(enqueuer.payloads) != 0 {
		t.Fatalf("não deveria enfileirar task sem integração Google Calendar ativa, obteve %+v", enqueuer.payloads)
	}
	_ = aptSvc.Cancel(ctx, tenantID, aptNoGcal.ID, "Teste", "Admin")
	if len(enqueuer.payloads) != 0 {
		t.Fatalf("não deveria enfileirar DELETE sem integração Google Calendar ativa, obteve %+v", enqueuer.payloads)
	}

	// Configura Google Calendar conectado porém ainda desativado (IsEnabled=false) -> também não enfileira
	_ = repo.UpsertGoogleCalendarConfig(ctx, &domain.TenantGoogleCalendarConfig{
		ID:                          uuid.New(),
		TenantID:                    tenantID,
		IsEnabled:                   false,
		Status:                      domain.GoogleCalendarStatusConnected,
		GoogleClientID:              "cid",
		GoogleClientSecretEncrypted: "enc",
	})
	aptDisabled, err := aptSvc.CreateAppointment(ctx, CreateAppointmentDTO{
		TenantID:       tenantID,
		ServiceID:      svcID,
		ProfessionalID: proID,
		StartAt:        startAt.Add(1 * time.Hour),
		CustomerName:   "Agenda Desativada",
		CustomerPhone:  "11900002222",
	})
	if err != nil || aptDisabled == nil {
		t.Fatalf("CreateAppointment com Google Calendar desativado falhou: %v", err)
	}
	if len(enqueuer.payloads) != 0 {
		t.Fatalf("não deveria enfileirar task com IsEnabled=false, obteve %+v", enqueuer.payloads)
	}

	// Habilita integração Google Calendar (IsEnabled=true e CONNECTED)
	_ = repo.UpsertGoogleCalendarConfig(ctx, &domain.TenantGoogleCalendarConfig{
		ID:                          uuid.New(),
		TenantID:                    tenantID,
		IsEnabled:                   true,
		Status:                      domain.GoogleCalendarStatusConnected,
		GoogleClientID:              "cid",
		GoogleClientSecretEncrypted: "enc",
	})

	// 1. CreateAppointment dispara CREATE
	apt, err := aptSvc.CreateAppointment(ctx, CreateAppointmentDTO{
		TenantID:       tenantID,
		ServiceID:      svcID,
		ProfessionalID: proID,
		StartAt:        startAt.Add(2 * time.Hour),
		CustomerName:   "Lucas Mendes",
		CustomerPhone:  "11988887777",
	})
	if err != nil {
		t.Fatalf("CreateAppointment falhou: %v", err)
	}
	if len(enqueuer.payloads) != 1 || enqueuer.payloads[0].Action != domain.CalendarSyncActionCreate || enqueuer.payloads[0].AppointmentID != apt.ID {
		t.Fatalf("esperava enqueue CREATE para %s, obteve %+v", apt.ID, enqueuer.payloads)
	}

	// 2. Reschedule dispara UPDATE
	if err := aptSvc.Reschedule(ctx, tenantID, apt.ID, startAt.Add(4*time.Hour), proID); err != nil {
		t.Fatalf("Reschedule falhou: %v", err)
	}
	if len(enqueuer.payloads) != 2 || enqueuer.payloads[1].Action != domain.CalendarSyncActionUpdate {
		t.Fatalf("esperava enqueue UPDATE, obteve %+v", enqueuer.payloads)
	}

	// 3. AIChatService cancelar_agendamento usa AppointmentService.Cancel e dispara DELETE
	aiSvc := NewAIChatService(repo, NewAvailabilityService(repo), aptSvc, nil, nil, nil, testEncKey)
	toolOut := aiSvc.executeTool(ctx, tenantID, "cancelar_agendamento", map[string]interface{}{
		"agendamento_id": apt.ID.String(),
		"motivo":         "Imprevisto",
	}, "11988887777", "Lucas Mendes")
	if !strings.Contains(toolOut, `"sucesso": true`) {
		t.Fatalf("cancelar_agendamento via IA falhou: %s", toolOut)
	}
	if len(enqueuer.payloads) != 3 || enqueuer.payloads[2].Action != domain.CalendarSyncActionDelete {
		t.Fatalf("esperava enqueue DELETE via AIChatService, obteve %+v", enqueuer.payloads)
	}

	// 4. Falha no enqueuer NÃO bloqueia agendamento e registra FAILED quando Google Calendar está conectado+habilitado
	aptSvc.SetCalendarSyncDependencies(&mockEnqueuer{err: errors.New("redis connection refused")}, repo)

	aptFail, err := aptSvc.CreateAppointment(ctx, CreateAppointmentDTO{
		TenantID:       tenantID,
		ServiceID:      svcID,
		ProfessionalID: proID,
		StartAt:        startAt.Add(6 * time.Hour),
		CustomerName:   "Maria Souza",
		CustomerPhone:  "11977776666",
	})
	if err != nil {
		t.Fatalf("falha no Redis/enqueuer NÃO deveria bloquear agendamento: %v", err)
	}
	failedRec, _ := repo.GetAppointmentCalendarSync(ctx, tenantID, aptFail.ID, domain.CalendarProviderGoogle)
	if failedRec == nil || failedRec.SyncStatus != domain.CalendarSyncStatusFailed {
		t.Fatalf("esperava registro FAILED visível quando enqueue falha com integração ativa, obteve %+v", failedRec)
	}
}

func base64URLEncode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}
