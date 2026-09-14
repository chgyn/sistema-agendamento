package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// -------------------------------------------------------------
// WHATSAPP CONFIG & CONVERSATIONS REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetWhatsAppConfigByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.TenantWhatsAppConfig, error) {
	var cfg domain.TenantWhatsAppConfig
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *Repository) GetWhatsAppConfigByInstanceToken(ctx context.Context, token string) (*domain.TenantWhatsAppConfig, error) {
	var cfg domain.TenantWhatsAppConfig
	err := r.db.WithContext(ctx).
		Where("instance_token = ?", token).
		First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *Repository) UpsertWhatsAppConfig(ctx context.Context, cfg *domain.TenantWhatsAppConfig) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "tenant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"instance_name",
			"instance_token",
			"phone_number",
			"status",
			"qr_code_base64",
			"qr_code_expires_at",
			"is_ai_enabled",
			"ai_provider",
			"ai_model",
			"gemini_api_key_encrypted",
			"open_ai_api_key_encrypted",
			"system_prompt_custom",
			"humanized_min_delay_sec",
			"humanized_max_delay_sec",
			"typing_speed_chars_sec",
			"debounce_window_sec",
			"updated_at",
		}),
	}).Create(cfg).Error
}

func (r *Repository) UpdateWhatsAppStatus(ctx context.Context, tenantID uuid.UUID, status domain.WhatsAppConnectionStatus, phone string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if phone != "" {
		updates["phone_number"] = phone
	}
	if status == domain.WhatsAppStatusConnected {
		// Limpa QR Code quando conectado com sucesso
		updates["qr_code_base64"] = ""
		updates["qr_code_expires_at"] = nil
	}
	return r.db.WithContext(ctx).
		Model(&domain.TenantWhatsAppConfig{}).
		Where("tenant_id = ?", tenantID).
		Updates(updates).Error
}

func (r *Repository) UpdateQRCode(ctx context.Context, tenantID uuid.UUID, qrCode string, expiresAt *time.Time) error {
	updates := map[string]interface{}{
		"status":             domain.WhatsAppStatusQRCode,
		"qr_code_base64":     qrCode,
		"qr_code_expires_at": expiresAt,
		"updated_at":         time.Now(),
	}
	return r.db.WithContext(ctx).
		Model(&domain.TenantWhatsAppConfig{}).
		Where("tenant_id = ?", tenantID).
		Updates(updates).Error
}

func (r *Repository) GetOrCreateConversation(ctx context.Context, tenantID uuid.UUID, customerPhone, customerName string) (*domain.WhatsAppConversation, error) {
	var conv domain.WhatsAppConversation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND customer_phone = ?", tenantID, customerPhone).
		First(&conv).Error

	if err == nil {
		if customerName != "" && conv.CustomerName != customerName {
			_ = r.db.WithContext(ctx).Model(&conv).Update("customer_name", customerName).Error
			conv.CustomerName = customerName
		}
		return &conv, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Busca vínculo com Customer existente no tenant
	var customer domain.Customer
	var customerID *uuid.UUID
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND phone = ?", tenantID, customerPhone).First(&customer).Error; err == nil {
		customerID = &customer.ID
	}

	newConv := domain.WhatsAppConversation{
		ID:            uuid.New(),
		TenantID:      tenantID,
		CustomerPhone: customerPhone,
		CustomerName:  customerName,
		CustomerID:    customerID,
		LastMessageAt: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&newConv).Error; err != nil {
		return nil, err
	}

	return &newConv, nil
}

func (r *Repository) SaveWhatsAppMessage(ctx context.Context, msg *domain.WhatsAppMessage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
		}
		return tx.Model(&domain.WhatsAppConversation{}).
			Where("id = ?", msg.ConversationID).
			Update("last_message_at", msg.CreatedAt).Error
	})
}

func (r *Repository) GetRecentMessages(ctx context.Context, conversationID uuid.UUID, limit int) ([]domain.WhatsAppMessage, error) {
	if limit <= 0 {
		limit = 20
	}
	var msgs []domain.WhatsAppMessage
	// Busca as N mensagens mais recentes (ordem decrescente)
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at desc").
		Limit(limit).
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}

	// Inverte a fatia em memória para entregar em ordem cronológica (antiga -> recente)
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nil
}
