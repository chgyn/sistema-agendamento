package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/crypto"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/integrations/wuzapi"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type WhatsAppService struct {
	repo          *postgres.Repository
	wuzapiClient  *wuzapi.Client
	encryptionKey string
	appBaseURL    string
}

func NewWhatsAppService(repo *postgres.Repository, wuzapiClient *wuzapi.Client, encryptionKey, appBaseURL string) *WhatsAppService {
	return &WhatsAppService{
		repo:          repo,
		wuzapiClient:  wuzapiClient,
		encryptionKey: encryptionKey,
		appBaseURL:    strings.TrimRight(appBaseURL, "/"),
	}
}

// EnsureInstance garante que o Tenant tenha registro no banco e usuário no WUZAPI
func (s *WhatsAppService) EnsureInstance(ctx context.Context, tenantID uuid.UUID) (*domain.TenantWhatsAppConfig, error) {
	cfg, err := s.repo.GetWhatsAppConfigByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar config do WhatsApp: %w", err)
	}

	tenant, err := s.repo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("estabelecimento não encontrado: %w", err)
	}

	instanceName := fmt.Sprintf("tenant_%s", tenant.Slug)

	if cfg == nil {
		instanceToken := uuid.New().String()
		cfg = &domain.TenantWhatsAppConfig{
			ID:                   uuid.New(),
			TenantID:             tenantID,
			InstanceName:         instanceName,
			InstanceToken:        instanceToken,
			Status:               domain.WhatsAppStatusDisconnected,
			IsAIEnabled:          false,
			AIProvider:           domain.AIProviderGemini,
			AIModel:              "gemini-2.5-flash",
			HumanizedMinDelaySec: 2,
			HumanizedMaxDelaySec: 5,
			TypingSpeedCharsSec:  35,
			DebounceWindowSec:    4,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		if err := s.repo.UpsertWhatsAppConfig(ctx, cfg); err != nil {
			return nil, fmt.Errorf("falha ao salvar config do WhatsApp: %w", err)
		}
	}

	// Garante que o usuário/instância exista no WUZAPI com webhook correto
	webhookURL := fmt.Sprintf("%s/api/v1/webhooks/whatsapp", s.appBaseURL)
	if err := s.wuzapiClient.EnsureUser(ctx, cfg.InstanceName, cfg.InstanceToken, webhookURL); err != nil {
		log.Printf("⚠️ [WUZAPI] Aviso ao registrar usuário no WUZAPI (%s): %v", cfg.InstanceName, err)
	}

	return cfg, nil
}

// GetStatus sincroniza e retorna o status atual da conexão
func (s *WhatsAppService) GetStatus(ctx context.Context, tenantID uuid.UUID) (*domain.WhatsAppStatusResponseDTO, error) {
	cfg, err := s.EnsureInstance(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Consulta status em tempo real no WUZAPI
	wuzapiStatus, err := s.wuzapiClient.GetStatus(ctx, cfg.InstanceToken)
	if err == nil && wuzapiStatus != nil {
		if wuzapiStatus.Data.LoggedIn && wuzapiStatus.Data.Connected {
			if cfg.Status != domain.WhatsAppStatusConnected {
				_ = s.repo.UpdateWhatsAppStatus(ctx, tenantID, domain.WhatsAppStatusConnected, "")
				cfg.Status = domain.WhatsAppStatusConnected
			}
		} else if !wuzapiStatus.Data.LoggedIn && cfg.Status == domain.WhatsAppStatusConnected {
			_ = s.repo.UpdateWhatsAppStatus(ctx, tenantID, domain.WhatsAppStatusDisconnected, "")
			cfg.Status = domain.WhatsAppStatusDisconnected
		}
	}

	hasQR := cfg.QRCodeBase64 != "" && (cfg.QRCodeExpiresAt == nil || cfg.QRCodeExpiresAt.After(time.Now()))

	return &domain.WhatsAppStatusResponseDTO{
		Status:          cfg.Status,
		PhoneNumber:     cfg.PhoneNumber,
		IsConnected:     cfg.Status == domain.WhatsAppStatusConnected,
		IsLoggedIn:      cfg.Status == domain.WhatsAppStatusConnected,
		HasQRCode:       hasQR,
		QRCodeBase64:    cfg.QRCodeBase64,
		QRCodeExpiresAt: cfg.QRCodeExpiresAt,
		InstanceName:    cfg.InstanceName,
	}, nil
}

// Connect inicia o pareamento e tenta gerar QR Code
func (s *WhatsAppService) Connect(ctx context.Context, tenantID uuid.UUID) (*domain.QRCodeResponseDTO, error) {
	cfg, err := s.EnsureInstance(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 1. Invoca connect no WUZAPI
	if err := s.wuzapiClient.Connect(ctx, cfg.InstanceToken); err != nil {
		return nil, fmt.Errorf("falha ao iniciar conexão no WUZAPI: %w", err)
	}

	// 2. Aguarda 1.5s e obtém o QR Code
	time.Sleep(1500 * time.Millisecond)
	qrBase64, err := s.wuzapiClient.GetQRCode(ctx, cfg.InstanceToken)
	if err != nil {
		// Verifica se já conectou diretamente
		status, sErr := s.wuzapiClient.GetStatus(ctx, cfg.InstanceToken)
		if sErr == nil && status.Data.LoggedIn {
			_ = s.repo.UpdateWhatsAppStatus(ctx, tenantID, domain.WhatsAppStatusConnected, "")
			return &domain.QRCodeResponseDTO{
				Status: domain.WhatsAppStatusConnected,
			}, nil
		}
		return nil, fmt.Errorf("falha ao recuperar QR code gerado: %w", err)
	}

	expiresAt := time.Now().Add(60 * time.Second)
	if err := s.repo.UpdateQRCode(ctx, tenantID, qrBase64, &expiresAt); err != nil {
		return nil, err
	}

	return &domain.QRCodeResponseDTO{
		QRCodeBase64: qrBase64,
		Status:       domain.WhatsAppStatusQRCode,
	}, nil
}

// GetQRCode retorna o QR Code atual ou gera um novo caso esteja expirado
func (s *WhatsAppService) GetQRCode(ctx context.Context, tenantID uuid.UUID) (*domain.QRCodeResponseDTO, error) {
	cfg, err := s.EnsureInstance(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if cfg.Status == domain.WhatsAppStatusConnected {
		return &domain.QRCodeResponseDTO{
			Status: domain.WhatsAppStatusConnected,
		}, nil
	}

	qrBase64, err := s.wuzapiClient.GetQRCode(ctx, cfg.InstanceToken)
	if err != nil {
		return nil, fmt.Errorf("QR Code indisponível no momento: %w", err)
	}

	expiresAt := time.Now().Add(60 * time.Second)
	_ = s.repo.UpdateQRCode(ctx, tenantID, qrBase64, &expiresAt)

	return &domain.QRCodeResponseDTO{
		QRCodeBase64: qrBase64,
		Status:       domain.WhatsAppStatusQRCode,
	}, nil
}

// Disconnect desconecta a sessão
func (s *WhatsAppService) Disconnect(ctx context.Context, tenantID uuid.UUID) error {
	cfg, err := s.repo.GetWhatsAppConfigByTenantID(ctx, tenantID)
	if err != nil || cfg == nil {
		return domain.ErrTenantNotFound
	}

	_ = s.wuzapiClient.Disconnect(ctx, cfg.InstanceToken)
	return s.repo.UpdateWhatsAppStatus(ctx, tenantID, domain.WhatsAppStatusDisconnected, "")
}

// Logout efetua logout completo e remove a sessão
func (s *WhatsAppService) Logout(ctx context.Context, tenantID uuid.UUID) error {
	cfg, err := s.repo.GetWhatsAppConfigByTenantID(ctx, tenantID)
	if err != nil || cfg == nil {
		return domain.ErrTenantNotFound
	}

	_ = s.wuzapiClient.Logout(ctx, cfg.InstanceToken)
	return s.repo.UpdateWhatsAppStatus(ctx, tenantID, domain.WhatsAppStatusLoggedOut, "")
}

// GetAIConfig retorna a configuração atual de IA com chaves mascaradas
func (s *WhatsAppService) GetAIConfig(ctx context.Context, tenantID uuid.UUID) (*domain.AIConfigResponseDTO, error) {
	cfg, err := s.EnsureInstance(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return &domain.AIConfigResponseDTO{
		IsAIEnabled:          cfg.IsAIEnabled,
		AIProvider:           cfg.AIProvider,
		AIModel:              cfg.AIModel,
		HasGeminiAPIKey:      cfg.GeminiAPIKeyEncrypted != "",
		HasOpenAIAPIKey:      cfg.OpenAIAPIKeyEncrypted != "",
		SystemPromptCustom:   cfg.SystemPromptCustom,
		HumanizedMinDelaySec: cfg.HumanizedMinDelaySec,
		HumanizedMaxDelaySec: cfg.HumanizedMaxDelaySec,
		TypingSpeedCharsSec:  cfg.TypingSpeedCharsSec,
		DebounceWindowSec:    cfg.DebounceWindowSec,
	}, nil
}

// UpdateAIConfig atualiza configurações e criptografa as API keys fornecidas
func (s *WhatsAppService) UpdateAIConfig(ctx context.Context, tenantID uuid.UUID, dto domain.UpdateAIConfigDTO) (*domain.AIConfigResponseDTO, error) {
	cfg, err := s.EnsureInstance(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	cfg.IsAIEnabled = dto.IsAIEnabled
	cfg.AIProvider = dto.AIProvider
	cfg.AIModel = dto.AIModel
	cfg.SystemPromptCustom = dto.SystemPromptCustom

	if dto.HumanizedMinDelaySec > 0 {
		cfg.HumanizedMinDelaySec = dto.HumanizedMinDelaySec
	}
	if dto.HumanizedMaxDelaySec > 0 {
		cfg.HumanizedMaxDelaySec = dto.HumanizedMaxDelaySec
	}
	if dto.TypingSpeedCharsSec > 0 {
		cfg.TypingSpeedCharsSec = dto.TypingSpeedCharsSec
	}
	if dto.DebounceWindowSec > 0 {
		cfg.DebounceWindowSec = dto.DebounceWindowSec
	}

	// Se nova Gemini API Key enviada, criptografa com AES-256-GCM
	if dto.GeminiAPIKey != nil && *dto.GeminiAPIKey != "" {
		encKey, err := crypto.Encrypt(*dto.GeminiAPIKey, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("falha ao criptografar chave do Gemini: %w", err)
		}
		cfg.GeminiAPIKeyEncrypted = encKey
	}

	// Se nova OpenAI API Key enviada, criptografa com AES-256-GCM
	if dto.OpenAIAPIKey != nil && *dto.OpenAIAPIKey != "" {
		encKey, err := crypto.Encrypt(*dto.OpenAIAPIKey, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("falha ao criptografar chave da OpenAI: %w", err)
		}
		cfg.OpenAIAPIKeyEncrypted = encKey
	}

	cfg.UpdatedAt = time.Now()
	if err := s.repo.UpsertWhatsAppConfig(ctx, cfg); err != nil {
		return nil, fmt.Errorf("falha ao salvar configurações de IA: %w", err)
	}

	return s.GetAIConfig(ctx, tenantID)
}
