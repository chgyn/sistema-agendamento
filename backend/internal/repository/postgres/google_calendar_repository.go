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
// GOOGLE CALENDAR CONFIG & SYNC REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetGoogleCalendarConfigByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.TenantGoogleCalendarConfig, error) {
	var cfg domain.TenantGoogleCalendarConfig
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

func (r *Repository) UpsertGoogleCalendarConfig(ctx context.Context, cfg *domain.TenantGoogleCalendarConfig) error {
	now := time.Now()
	if cfg.ID == uuid.Nil {
		cfg.ID = uuid.New()
	}
	if cfg.CreatedAt.IsZero() {
		cfg.CreatedAt = now
	}
	cfg.UpdatedAt = now
	if cfg.CalendarID == "" {
		cfg.CalendarID = "primary"
	}
	if cfg.Status == "" {
		cfg.Status = domain.GoogleCalendarStatusDisconnected
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "tenant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"is_enabled",
			"status",
			"google_client_id",
			"google_client_secret_encrypted",
			"google_account_email",
			"access_token_encrypted",
			"refresh_token_encrypted",
			"token_expires_at",
			"calendar_id",
			"calendar_summary",
			"last_sync_at",
			"last_error",
			"updated_at",
		}),
	}).Create(cfg).Error
}

func (r *Repository) GetAppointmentCalendarSync(ctx context.Context, tenantID, appointmentID uuid.UUID, provider domain.CalendarSyncProvider) (*domain.AppointmentCalendarSync, error) {
	var syncRec domain.AppointmentCalendarSync
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND appointment_id = ? AND provider = ?", tenantID, appointmentID, provider).
		First(&syncRec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &syncRec, nil
}

func (r *Repository) UpsertAppointmentCalendarSync(ctx context.Context, syncRec *domain.AppointmentCalendarSync) error {
	now := time.Now()
	if syncRec.ID == uuid.Nil {
		syncRec.ID = uuid.New()
	}
	if syncRec.CreatedAt.IsZero() {
		syncRec.CreatedAt = now
	}
	syncRec.UpdatedAt = now
	if syncRec.Provider == "" {
		syncRec.Provider = domain.CalendarProviderGoogle
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "tenant_id"},
			{Name: "appointment_id"},
			{Name: "provider"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"external_calendar_id",
			"external_event_id",
			"sync_status",
			"last_action",
			"attempt_count",
			"last_error",
			"last_synced_at",
			"updated_at",
		}),
	}).Create(syncRec).Error
}

func (r *Repository) ListFailedCalendarSyncs(ctx context.Context, tenantID uuid.UUID, provider domain.CalendarSyncProvider) ([]domain.AppointmentCalendarSync, error) {
	var syncs []domain.AppointmentCalendarSync
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND provider = ? AND sync_status = ?", tenantID, provider, domain.CalendarSyncStatusFailed).
		Preload("Appointment").
		Preload("Appointment.Customer").
		Preload("Appointment.Service").
		Preload("Appointment.Professional").
		Order("updated_at desc").
		Find(&syncs).Error
	return syncs, err
}

func (r *Repository) CountFailedCalendarSyncs(ctx context.Context, tenantID uuid.UUID, provider domain.CalendarSyncProvider) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.AppointmentCalendarSync{}).
		Where("tenant_id = ? AND provider = ? AND sync_status = ?", tenantID, provider, domain.CalendarSyncStatusFailed).
		Count(&count).Error
	return count, err
}
