package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/sistema-agendamento/backend/internal/crypto"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/integrations/googlecalendar"
	"github.com/sistema-agendamento/backend/internal/queue"
)

const (
	oauthStateTTL         = 10 * time.Minute
	oauthStateRedisPrefix = "gcal:oauth:state:"
	tokenRefreshBuffer    = 2 * time.Minute
)

// AppointmentReader abstrai leitura de agendamento para desacoplamento e testes unitários
type AppointmentReader interface {
	GetAppointmentByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Appointment, error)
}

// CalendarSyncEnqueuer abstrai o enfileiramento assíncrono de sincronizações do Google Calendar
type CalendarSyncEnqueuer interface {
	EnqueueGoogleCalendarSync(ctx context.Context, payload queue.GoogleCalendarSyncPayload) error
}

type oauthStatePayload struct {
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
	Nonce    string    `json:"nonce"`
	Exp      int64     `json:"exp"`
}

type GoogleCalendarService struct {
	repo          domain.GoogleCalendarRepository
	aptRepo       AppointmentReader
	provider      domain.CalendarProvider
	enqueuer      CalendarSyncEnqueuer
	redisClient   *redis.Client
	jwtSecret     string
	encryptionKey string
	redirectURI   string
}

func NewGoogleCalendarService(
	repo domain.GoogleCalendarRepository,
	aptRepo AppointmentReader,
	provider domain.CalendarProvider,
	enqueuer CalendarSyncEnqueuer,
	redisClient *redis.Client,
	jwtSecret string,
	encryptionKey string,
	redirectURI string,
) *GoogleCalendarService {
	return &GoogleCalendarService{
		repo:          repo,
		aptRepo:       aptRepo,
		provider:      provider,
		enqueuer:      enqueuer,
		redisClient:   redisClient,
		jwtSecret:     jwtSecret,
		encryptionKey: encryptionKey,
		redirectURI:   redirectURI,
	}
}

// GetStatus retorna configuração atual e contagem de falhas sem expor segredos ou tokens
func (s *GoogleCalendarService) GetStatus(ctx context.Context, tenantID uuid.UUID) (*domain.GoogleCalendarConfigResponseDTO, error) {
	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	failedCount, err := s.repo.CountFailedCalendarSyncs(ctx, tenantID, domain.CalendarProviderGoogle)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		return &domain.GoogleCalendarConfigResponseDTO{
			RedirectURI:      s.redirectURI,
			Status:           domain.GoogleCalendarStatusDisconnected,
			CalendarID:       "primary",
			FailedSyncsCount: failedCount,
		}, nil
	}

	calID := cfg.CalendarID
	if calID == "" {
		calID = "primary"
	}

	hasClientSecret := cfg.GoogleClientSecretEncrypted != ""
	hasCredentials := strings.TrimSpace(cfg.GoogleClientID) != "" && hasClientSecret

	return &domain.GoogleCalendarConfigResponseDTO{
		GoogleClientID:     cfg.GoogleClientID,
		HasClientSecret:    hasClientSecret,
		HasCredentials:     hasCredentials,
		RedirectURI:        s.redirectURI,
		IsEnabled:          cfg.IsEnabled,
		Status:             cfg.Status,
		GoogleAccountEmail: cfg.GoogleAccountEmail,
		CalendarID:         calID,
		CalendarSummary:    cfg.CalendarSummary,
		LastSyncAt:         cfg.LastSyncAt,
		LastError:          cfg.LastError,
		FailedSyncsCount:   failedCount,
	}, nil
}

// UpdateSettings atualiza credenciais BYOK, calendário alvo e flag de ativação do tenant
func (s *GoogleCalendarService) UpdateSettings(ctx context.Context, tenantID uuid.UUID, dto domain.UpdateGoogleCalendarSettingsDTO) (*domain.GoogleCalendarConfigResponseDTO, error) {
	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &domain.TenantGoogleCalendarConfig{
			ID:         uuid.New(),
			TenantID:   tenantID,
			Status:     domain.GoogleCalendarStatusDisconnected,
			CalendarID: "primary",
		}
	}

	credentialsChanged := false

	if dto.GoogleClientID != nil {
		newClientID := strings.TrimSpace(*dto.GoogleClientID)
		if newClientID != cfg.GoogleClientID {
			cfg.GoogleClientID = newClientID
			credentialsChanged = true
		}
	}

	if dto.GoogleClientSecret != nil {
		newSecret := strings.TrimSpace(*dto.GoogleClientSecret)
		if newSecret != "" {
			encSecret, err := crypto.Encrypt(newSecret, s.encryptionKey)
			if err != nil {
				return nil, fmt.Errorf("falha ao cifrar Google Client Secret: %w", err)
			}
			cfg.GoogleClientSecretEncrypted = encSecret
			credentialsChanged = true
		}
	}

	if credentialsChanged {
		cfg.AccessTokenEncrypted = ""
		cfg.RefreshTokenEncrypted = ""
		cfg.TokenExpiresAt = nil
		cfg.GoogleAccountEmail = ""
		cfg.Status = domain.GoogleCalendarStatusDisconnected
		cfg.IsEnabled = false
		cfg.LastError = ""
	}

	if dto.CalendarID != nil {
		calID := strings.TrimSpace(*dto.CalendarID)
		if calID == "" {
			calID = "primary"
		}
		cfg.CalendarID = calID
	}

	if dto.CalendarSummary != nil {
		cfg.CalendarSummary = strings.TrimSpace(*dto.CalendarSummary)
	}

	if dto.IsEnabled != nil && !credentialsChanged {
		if *dto.IsEnabled {
			if cfg.Status != domain.GoogleCalendarStatusConnected {
				return nil, domain.ErrGoogleCalendarNotConnected
			}
			cfg.IsEnabled = true
		} else {
			cfg.IsEnabled = false
		}
	}

	if err := s.repo.UpsertGoogleCalendarConfig(ctx, cfg); err != nil {
		return nil, err
	}

	return s.GetStatus(ctx, tenantID)
}

// GetAuthURL gera a URL de autorização OAuth 2.0 usando as credenciais BYOK do tenant
func (s *GoogleCalendarService) GetAuthURL(ctx context.Context, tenantID, userID uuid.UUID) (string, error) {
	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return "", err
	}
	if cfg == nil || strings.TrimSpace(cfg.GoogleClientID) == "" || cfg.GoogleClientSecretEncrypted == "" {
		return "", domain.ErrGoogleCalendarNotConfigured
	}
	if _, err := crypto.Decrypt(cfg.GoogleClientSecretEncrypted, s.encryptionKey); err != nil {
		return "", domain.ErrGoogleCalendarNotConfigured
	}

	state, err := s.generateOAuthState(ctx, tenantID, userID)
	if err != nil {
		return "", fmt.Errorf("falha ao gerar state OAuth: %w", err)
	}

	return s.provider.BuildAuthURL(cfg.GoogleClientID, s.redirectURI, state), nil
}

// HandleCallback valida o state assinado, troca o code por tokens e persiste criptografados para o tenant
func (s *GoogleCalendarService) HandleCallback(ctx context.Context, code, state string) (uuid.UUID, error) {
	payload, err := s.validateAndConsumeOAuthState(ctx, state)
	if err != nil {
		return uuid.Nil, err
	}
	tenantID := payload.TenantID

	if strings.TrimSpace(code) == "" {
		return tenantID, errors.New("código de autorização OAuth não informado")
	}

	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return tenantID, err
	}
	if cfg == nil || strings.TrimSpace(cfg.GoogleClientID) == "" || cfg.GoogleClientSecretEncrypted == "" {
		return tenantID, domain.ErrGoogleCalendarNotConfigured
	}

	clientSecret, err := crypto.Decrypt(cfg.GoogleClientSecretEncrypted, s.encryptionKey)
	if err != nil || clientSecret == "" {
		return tenantID, domain.ErrGoogleCalendarNotConfigured
	}

	tokenResp, err := s.provider.ExchangeCode(ctx, cfg.GoogleClientID, clientSecret, s.redirectURI, code)
	if err != nil {
		cfg.Status = domain.GoogleCalendarStatusError
		cfg.LastError = err.Error()
		_ = s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
		return tenantID, err
	}

	email, err := s.provider.GetUserEmail(ctx, tokenResp.AccessToken)
	if err != nil {
		cfg.Status = domain.GoogleCalendarStatusError
		cfg.LastError = err.Error()
		_ = s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
		return tenantID, err
	}

	encAccess, err := crypto.Encrypt(tokenResp.AccessToken, s.encryptionKey)
	if err != nil {
		return tenantID, fmt.Errorf("falha ao cifrar access token: %w", err)
	}

	if tokenResp.RefreshToken != "" {
		encRefresh, err := crypto.Encrypt(tokenResp.RefreshToken, s.encryptionKey)
		if err != nil {
			return tenantID, fmt.Errorf("falha ao cifrar refresh token: %w", err)
		}
		cfg.RefreshTokenEncrypted = encRefresh
	} else if cfg.RefreshTokenEncrypted == "" {
		return tenantID, errors.New("o Google não retornou refresh_token; revogue o acesso do aplicativo na sua conta Google e reconecte")
	}

	expiresIn := tokenResp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	cfg.AccessTokenEncrypted = encAccess
	cfg.TokenExpiresAt = &expiresAt
	cfg.GoogleAccountEmail = email
	if cfg.CalendarID == "" {
		cfg.CalendarID = "primary"
	}
	cfg.Status = domain.GoogleCalendarStatusConnected
	cfg.LastError = ""

	if err := s.repo.UpsertGoogleCalendarConfig(ctx, cfg); err != nil {
		return tenantID, err
	}

	return tenantID, nil
}

// Disconnect revoga o token (best-effort) e limpa a sessão mantendo credenciais BYOK do tenant
func (s *GoogleCalendarService) Disconnect(ctx context.Context, tenantID uuid.UUID) error {
	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}
	if cfg == nil {
		return nil
	}

	if cfg.RefreshTokenEncrypted != "" {
		if token, err := crypto.Decrypt(cfg.RefreshTokenEncrypted, s.encryptionKey); err == nil && token != "" {
			_ = s.provider.RevokeToken(ctx, token)
		}
	} else if cfg.AccessTokenEncrypted != "" {
		if token, err := crypto.Decrypt(cfg.AccessTokenEncrypted, s.encryptionKey); err == nil && token != "" {
			_ = s.provider.RevokeToken(ctx, token)
		}
	}

	cfg.AccessTokenEncrypted = ""
	cfg.RefreshTokenEncrypted = ""
	cfg.TokenExpiresAt = nil
	cfg.GoogleAccountEmail = ""
	cfg.IsEnabled = false
	cfg.Status = domain.GoogleCalendarStatusDisconnected
	cfg.LastError = ""

	return s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
}

// GetValidAccessToken retorna um access_token válido do tenant, renovando via refresh_token quando necessário
func (s *GoogleCalendarService) GetValidAccessToken(ctx context.Context, tenantID uuid.UUID) (string, error) {
	return s.getAccessToken(ctx, tenantID, false)
}

func (s *GoogleCalendarService) getAccessToken(ctx context.Context, tenantID uuid.UUID, forceRefresh bool) (string, error) {
	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return "", err
	}
	if cfg == nil {
		return "", domain.ErrGoogleCalendarNotConnected
	}
	if strings.TrimSpace(cfg.GoogleClientID) == "" || cfg.GoogleClientSecretEncrypted == "" {
		return "", domain.ErrGoogleCalendarNotConfigured
	}
	if cfg.Status == domain.GoogleCalendarStatusDisconnected {
		return "", domain.ErrGoogleCalendarNotConnected
	}
	if cfg.Status == domain.GoogleCalendarStatusExpired {
		return "", domain.ErrGoogleCalendarAuthExpired
	}

	if !forceRefresh && cfg.AccessTokenEncrypted != "" && cfg.TokenExpiresAt != nil && time.Until(*cfg.TokenExpiresAt) > tokenRefreshBuffer {
		accessToken, err := crypto.Decrypt(cfg.AccessTokenEncrypted, s.encryptionKey)
		if err == nil && accessToken != "" {
			return accessToken, nil
		}
	}

	if cfg.RefreshTokenEncrypted == "" {
		cfg.Status = domain.GoogleCalendarStatusExpired
		cfg.LastError = domain.ErrGoogleCalendarAuthExpired.Error()
		_ = s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
		return "", domain.ErrGoogleCalendarAuthExpired
	}

	clientSecret, err := crypto.Decrypt(cfg.GoogleClientSecretEncrypted, s.encryptionKey)
	if err != nil || clientSecret == "" {
		return "", domain.ErrGoogleCalendarNotConfigured
	}

	refreshToken, err := crypto.Decrypt(cfg.RefreshTokenEncrypted, s.encryptionKey)
	if err != nil || refreshToken == "" {
		cfg.Status = domain.GoogleCalendarStatusExpired
		cfg.LastError = domain.ErrGoogleCalendarAuthExpired.Error()
		_ = s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
		return "", domain.ErrGoogleCalendarAuthExpired
	}

	tokenResp, err := s.provider.RefreshToken(ctx, cfg.GoogleClientID, clientSecret, refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrGoogleCalendarAuthExpired) {
			cfg.Status = domain.GoogleCalendarStatusExpired
			cfg.LastError = err.Error()
			_ = s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
		}
		return "", err
	}

	encAccess, err := crypto.Encrypt(tokenResp.AccessToken, s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("falha ao cifrar novo access token: %w", err)
	}
	cfg.AccessTokenEncrypted = encAccess

	if tokenResp.RefreshToken != "" {
		if encRefresh, err := crypto.Encrypt(tokenResp.RefreshToken, s.encryptionKey); err == nil {
			cfg.RefreshTokenEncrypted = encRefresh
		}
	}

	expiresIn := tokenResp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	cfg.TokenExpiresAt = &expiresAt
	cfg.Status = domain.GoogleCalendarStatusConnected
	cfg.LastError = ""

	if err := s.repo.UpsertGoogleCalendarConfig(ctx, cfg); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// ListCalendars lista os calendários graváveis da conta Google conectada ao tenant
func (s *GoogleCalendarService) ListCalendars(ctx context.Context, tenantID uuid.UUID) ([]domain.GoogleCalendarItemDTO, error) {
	accessToken, err := s.GetValidAccessToken(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	calendars, err := s.provider.ListCalendars(ctx, accessToken)
	if errors.Is(err, googlecalendar.ErrTokenUnauthorized) {
		accessToken, refreshErr := s.getAccessToken(ctx, tenantID, true)
		if refreshErr != nil {
			return nil, refreshErr
		}
		return s.provider.ListCalendars(ctx, accessToken)
	}
	return calendars, err
}

// ListFailures retorna os agendamentos com falha de sincronização no tenant
func (s *GoogleCalendarService) ListFailures(ctx context.Context, tenantID uuid.UUID) ([]domain.AppointmentCalendarSyncDTO, error) {
	records, err := s.repo.ListFailedCalendarSyncs(ctx, tenantID, domain.CalendarProviderGoogle)
	if err != nil {
		return nil, err
	}

	result := make([]domain.AppointmentCalendarSyncDTO, 0, len(records))
	for _, rec := range records {
		item := domain.AppointmentCalendarSyncDTO{
			AppointmentID: rec.AppointmentID,
			SyncStatus:    rec.SyncStatus,
			LastAction:    rec.LastAction,
			AttemptCount:  rec.AttemptCount,
			LastError:     rec.LastError,
			UpdatedAt:     rec.UpdatedAt,
		}
		if rec.Appointment != nil {
			item.StartAt = rec.Appointment.StartAt
			if rec.Appointment.Customer != nil {
				item.CustomerName = rec.Appointment.Customer.Name
			}
			if rec.Appointment.Service != nil {
				item.ServiceName = rec.Appointment.Service.Name
			}
			if rec.Appointment.Professional != nil {
				item.ProfessionalName = rec.Appointment.Professional.Name
			}
		}
		result = append(result, item)
	}
	return result, nil
}

// RetryAppointmentSync valida que o agendamento pertence ao tenant e reenfileira a sincronização
func (s *GoogleCalendarService) RetryAppointmentSync(ctx context.Context, tenantID, appointmentID uuid.UUID) error {
	apt, err := s.aptRepo.GetAppointmentByID(ctx, tenantID, appointmentID)
	if err != nil {
		return err
	}

	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}
	if cfg == nil || cfg.Status != domain.GoogleCalendarStatusConnected {
		if cfg != nil && cfg.Status == domain.GoogleCalendarStatusExpired {
			return domain.ErrGoogleCalendarAuthExpired
		}
		return domain.ErrGoogleCalendarNotConnected
	}
	if !cfg.IsEnabled {
		return domain.ErrGoogleCalendarNotConnected
	}

	action := domain.CalendarSyncActionUpdate
	if apt.Status == domain.StatusCancelled {
		action = domain.CalendarSyncActionDelete
	} else {
		syncRec, err := s.repo.GetAppointmentCalendarSync(ctx, tenantID, appointmentID, domain.CalendarProviderGoogle)
		if err == nil && (syncRec == nil || syncRec.ExternalEventID == "" || syncRec.SyncStatus == domain.CalendarSyncStatusDeleted) {
			action = domain.CalendarSyncActionCreate
		}
	}

	if s.enqueuer == nil {
		return errors.New("fila de processamento indisponível")
	}
	return s.enqueuer.EnqueueGoogleCalendarSync(ctx, queue.GoogleCalendarSyncPayload{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		Action:        action,
	})
}

// SyncAppointment reconcilia o estado atual do agendamento no banco com o Google Calendar de forma idempotente
func (s *GoogleCalendarService) SyncAppointment(ctx context.Context, tenantID, appointmentID uuid.UUID, action domain.CalendarSyncAction) error {
	cfg, err := s.repo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}
	if cfg == nil || !cfg.IsEnabled || cfg.Status != domain.GoogleCalendarStatusConnected {
		return nil
	}

	apt, err := s.aptRepo.GetAppointmentByID(ctx, tenantID, appointmentID)
	if err != nil {
		if errors.Is(err, domain.ErrAppointmentNotFound) {
			return nil
		}
		return err
	}

	syncRec, err := s.repo.GetAppointmentCalendarSync(ctx, tenantID, appointmentID, domain.CalendarProviderGoogle)
	if err != nil {
		return err
	}
	if syncRec == nil {
		syncRec = &domain.AppointmentCalendarSync{
			ID:            uuid.New(),
			TenantID:      tenantID,
			AppointmentID: appointmentID,
			Provider:      domain.CalendarProviderGoogle,
			SyncStatus:    domain.CalendarSyncStatusPending,
		}
	}

	accessToken, err := s.GetValidAccessToken(ctx, tenantID)
	if err != nil {
		return s.recordSyncFailure(ctx, cfg, syncRec, action, err)
	}

	// Reconciliação pelo estado atual no DB (evita race condition em eventos fora de ordem)
	if apt.Status == domain.StatusCancelled || action == domain.CalendarSyncActionDelete {
		return s.syncDeletion(ctx, cfg, syncRec, accessToken)
	}

	if apt.Status != domain.StatusConfirmed && apt.Status != domain.StatusPending {
		return nil
	}

	return s.syncUpsert(ctx, cfg, syncRec, apt, accessToken, action)
}

func (s *GoogleCalendarService) syncDeletion(
	ctx context.Context,
	cfg *domain.TenantGoogleCalendarConfig,
	syncRec *domain.AppointmentCalendarSync,
	accessToken string,
) error {
	if syncRec.ExternalEventID == "" || syncRec.SyncStatus == domain.CalendarSyncStatusDeleted {
		if syncRec.SyncStatus == domain.CalendarSyncStatusFailed {
			now := time.Now()
			syncRec.SyncStatus = domain.CalendarSyncStatusDeleted
			syncRec.LastAction = domain.CalendarSyncActionDelete
			syncRec.LastError = ""
			syncRec.LastSyncedAt = &now
			_ = s.repo.UpsertAppointmentCalendarSync(ctx, syncRec)
		}
		return nil
	}

	calID := syncRec.ExternalCalendarID
	if calID == "" {
		calID = cfg.CalendarID
	}
	if calID == "" {
		calID = "primary"
	}

	err := s.provider.DeleteEvent(ctx, accessToken, calID, syncRec.ExternalEventID)
	if errors.Is(err, googlecalendar.ErrTokenUnauthorized) {
		newToken, refreshErr := s.getAccessToken(ctx, cfg.TenantID, true)
		if refreshErr != nil {
			return s.recordSyncFailure(ctx, cfg, syncRec, domain.CalendarSyncActionDelete, refreshErr)
		}
		err = s.provider.DeleteEvent(ctx, newToken, calID, syncRec.ExternalEventID)
	}
	if err != nil {
		return s.recordSyncFailure(ctx, cfg, syncRec, domain.CalendarSyncActionDelete, err)
	}

	now := time.Now()
	syncRec.SyncStatus = domain.CalendarSyncStatusDeleted
	syncRec.LastAction = domain.CalendarSyncActionDelete
	syncRec.AttemptCount++
	syncRec.LastError = ""
	syncRec.LastSyncedAt = &now
	if err := s.repo.UpsertAppointmentCalendarSync(ctx, syncRec); err != nil {
		return err
	}

	cfg.LastSyncAt = &now
	cfg.LastError = ""
	return s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
}

func (s *GoogleCalendarService) syncUpsert(
	ctx context.Context,
	cfg *domain.TenantGoogleCalendarConfig,
	syncRec *domain.AppointmentCalendarSync,
	apt *domain.Appointment,
	accessToken string,
	action domain.CalendarSyncAction,
) error {
	targetCalID := cfg.CalendarID
	if targetCalID == "" {
		targetCalID = "primary"
	}

	eventInput := buildCalendarEventInput(cfg.TenantID, apt)
	var finalEventID string
	var effectiveAction domain.CalendarSyncAction

	runUpsert := func(token string) (string, domain.CalendarSyncAction, error) {
		// Evento já existe e não foi deletado
		if syncRec.ExternalEventID != "" && syncRec.SyncStatus != domain.CalendarSyncStatusDeleted {
			// Se o calendário alvo mudou, remove do antigo e cria no novo
			if syncRec.ExternalCalendarID != "" && syncRec.ExternalCalendarID != targetCalID {
				if err := s.provider.DeleteEvent(ctx, token, syncRec.ExternalCalendarID, syncRec.ExternalEventID); err != nil {
					return "", domain.CalendarSyncActionUpdate, err
				}
				eventInput.EventID = newRandomGoogleEventID()
				createdID, err := s.createWithConflictRecovery(ctx, token, targetCalID, eventInput)
				return createdID, domain.CalendarSyncActionCreate, err
			}

			updatedID, err := s.provider.UpdateEvent(ctx, token, targetCalID, syncRec.ExternalEventID, eventInput)
			if errors.Is(err, googlecalendar.ErrEventNotFound) {
				// Evento apagado externamente (404/410) -> recria com novo ID válido
				eventInput.EventID = newRandomGoogleEventID()
				createdID, createErr := s.createWithConflictRecovery(ctx, token, targetCalID, eventInput)
				return createdID, domain.CalendarSyncActionCreate, createErr
			}
			return updatedID, domain.CalendarSyncActionUpdate, err
		}

		// Primeira criação usa ID determinístico em base32hex; se já foi deletado antes, usa novo UUID em base32hex
		if syncRec.SyncStatus == domain.CalendarSyncStatusDeleted {
			eventInput.EventID = newRandomGoogleEventID()
		} else {
			eventInput.EventID = deterministicGoogleEventID(apt.ID)
		}
		createdID, err := s.createWithConflictRecovery(ctx, token, targetCalID, eventInput)
		return createdID, domain.CalendarSyncActionCreate, err
	}

	finalEventID, effectiveAction, err := runUpsert(accessToken)
	if errors.Is(err, googlecalendar.ErrTokenUnauthorized) {
		newToken, refreshErr := s.getAccessToken(ctx, cfg.TenantID, true)
		if refreshErr != nil {
			return s.recordSyncFailure(ctx, cfg, syncRec, action, refreshErr)
		}
		finalEventID, effectiveAction, err = runUpsert(newToken)
	}
	if err != nil {
		return s.recordSyncFailure(ctx, cfg, syncRec, action, err)
	}

	now := time.Now()
	syncRec.ExternalCalendarID = targetCalID
	syncRec.ExternalEventID = finalEventID
	syncRec.SyncStatus = domain.CalendarSyncStatusSynced
	syncRec.LastAction = effectiveAction
	syncRec.AttemptCount++
	syncRec.LastError = ""
	syncRec.LastSyncedAt = &now
	if err := s.repo.UpsertAppointmentCalendarSync(ctx, syncRec); err != nil {
		return err
	}

	cfg.LastSyncAt = &now
	cfg.LastError = ""
	return s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
}

func (s *GoogleCalendarService) createWithConflictRecovery(
	ctx context.Context,
	accessToken string,
	calendarID string,
	eventInput domain.CalendarEventInput,
) (string, error) {
	createdID, err := s.provider.CreateEvent(ctx, accessToken, calendarID, eventInput)
	if errors.Is(err, googlecalendar.ErrEventConflict) && eventInput.EventID != "" {
		return s.provider.UpdateEvent(ctx, accessToken, calendarID, eventInput.EventID, eventInput)
	}
	return createdID, err
}

func (s *GoogleCalendarService) recordSyncFailure(
	ctx context.Context,
	cfg *domain.TenantGoogleCalendarConfig,
	syncRec *domain.AppointmentCalendarSync,
	action domain.CalendarSyncAction,
	syncErr error,
) error {
	syncRec.SyncStatus = domain.CalendarSyncStatusFailed
	syncRec.LastAction = action
	syncRec.AttemptCount++
	syncRec.LastError = syncErr.Error()
	_ = s.repo.UpsertAppointmentCalendarSync(ctx, syncRec)

	// Atualiza LastError do tenant apenas se não for reload de expiração já gravado
	if cfg != nil && !errors.Is(syncErr, domain.ErrGoogleCalendarAuthExpired) {
		cfg.LastError = syncErr.Error()
		_ = s.repo.UpsertGoogleCalendarConfig(ctx, cfg)
	}

	if errors.Is(syncErr, domain.ErrGoogleCalendarAuthExpired) {
		return fmt.Errorf("%w: %w", syncErr, asynq.SkipRetry)
	}
	return syncErr
}

func (s *GoogleCalendarService) generateOAuthState(ctx context.Context, tenantID, userID uuid.UUID) (string, error) {
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(nonceBytes)

	payload := oauthStatePayload{
		TenantID: tenantID,
		UserID:   userID,
		Nonce:    nonce,
		Exp:      time.Now().Add(oauthStateTTL).Unix(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)
	sigB64 := s.signStatePayload(payloadB64)

	if s.redisClient != nil {
		if err := s.redisClient.Set(ctx, oauthStateRedisPrefix+nonce, tenantID.String(), oauthStateTTL).Err(); err != nil {
			return "", fmt.Errorf("falha ao registrar nonce OAuth no Redis: %w", err)
		}
	}

	return payloadB64 + "." + sigB64, nil
}

func (s *GoogleCalendarService) validateAndConsumeOAuthState(ctx context.Context, state string) (*oauthStatePayload, error) {
	parts := strings.Split(state, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, domain.ErrInvalidOAuthState
	}

	expectedSig := s.signStatePayload(parts[0])
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSig)) {
		return nil, domain.ErrInvalidOAuthState
	}

	rawJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, domain.ErrInvalidOAuthState
	}

	var payload oauthStatePayload
	if err := json.Unmarshal(rawJSON, &payload); err != nil {
		return nil, domain.ErrInvalidOAuthState
	}

	if payload.TenantID == uuid.Nil || payload.Nonce == "" || time.Now().Unix() > payload.Exp {
		return nil, domain.ErrInvalidOAuthState
	}

	if s.redisClient != nil {
		storedTenantID, err := s.redisClient.GetDel(ctx, oauthStateRedisPrefix+payload.Nonce).Result()
		if err != nil || storedTenantID != payload.TenantID.String() {
			return nil, domain.ErrInvalidOAuthState
		}
	}

	return &payload, nil
}

func (s *GoogleCalendarService) signStatePayload(payloadB64 string) string {
	mac := hmac.New(sha256.New, []byte(s.jwtSecret))
	mac.Write([]byte(payloadB64))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func buildCalendarEventInput(tenantID uuid.UUID, apt *domain.Appointment) domain.CalendarEventInput {
	svcName := "Atendimento"
	if apt.Service != nil && apt.Service.Name != "" {
		svcName = apt.Service.Name
	}
	customerName := "Cliente"
	customerPhone := "-"
	if apt.Customer != nil {
		if apt.Customer.Name != "" {
			customerName = apt.Customer.Name
		}
		if apt.Customer.Phone != "" {
			customerPhone = apt.Customer.Phone
		}
	}
	proName := "-"
	if apt.Professional != nil && apt.Professional.Name != "" {
		proName = apt.Professional.Name
	}

	lines := []string{
		fmt.Sprintf("Cliente: %s", customerName),
		fmt.Sprintf("Telefone: %s", customerPhone),
		fmt.Sprintf("Profissional: %s", proName),
		fmt.Sprintf("Serviço: %s (%d min)", svcName, apt.DurationMinutes),
	}
	if strings.TrimSpace(apt.Notes) != "" {
		lines = append(lines, fmt.Sprintf("Observações: %s", strings.TrimSpace(apt.Notes)))
	}

	return domain.CalendarEventInput{
		Summary:     fmt.Sprintf("%s - %s", svcName, customerName),
		Description: strings.Join(lines, "\n"),
		StartAt:     apt.StartAt.In(time.Local),
		EndAt:       apt.EndAt.In(time.Local),
		TimeZone:    resolveLocalTimeZone(),
		PrivateProps: map[string]string{
			"agende_facil_appointment_id": apt.ID.String(),
			"agende_facil_tenant_id":      tenantID.String(),
		},
	}
}

func resolveLocalTimeZone() string {
	if tz := strings.TrimSpace(os.Getenv("TZ")); tz != "" {
		return tz
	}
	locName := time.Local.String()
	if locName != "" && locName != "Local" {
		return locName
	}
	return ""
}

// deterministicGoogleEventID gera um ID base32hex válido ([a-v0-9]) para o Google Calendar
func deterministicGoogleEventID(appointmentID uuid.UUID) string {
	return "af" + strings.ReplaceAll(appointmentID.String(), "-", "")
}

func newRandomGoogleEventID() string {
	return "af" + strings.ReplaceAll(uuid.NewString(), "-", "")
}
