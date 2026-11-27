package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupGoogleCalendarTestDB(t *testing.T) (*gorm.DB, *postgres.Repository) {
	dbName := fmt.Sprintf("file:mem_gcal_repo_%d?mode=memory&cache=shared&_busy_timeout=10000", time.Now().UnixNano())
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
		&domain.Professional{},
		&domain.Service{},
		&domain.Customer{},
		&domain.Appointment{},
		&domain.TenantGoogleCalendarConfig{},
		&domain.AppointmentCalendarSync{},
	)
	if err != nil {
		t.Fatalf("falha no automigrate: %v", err)
	}

	return db, postgres.NewRepository(db)
}

func TestGoogleCalendarRepository_MultiTenantIsolationAndFailures(t *testing.T) {
	db, repo := setupGoogleCalendarTestDB(t)
	ctx := context.Background()

	t1 := uuid.New()
	t2 := uuid.New()

	cfg1 := &domain.TenantGoogleCalendarConfig{
		ID:                          uuid.New(),
		TenantID:                    t1,
		IsEnabled:                   true,
		Status:                      domain.GoogleCalendarStatusConnected,
		GoogleClientID:              "client-t1",
		GoogleClientSecretEncrypted: "enc-secret-t1",
		CalendarID:                  "primary",
	}
	cfg2 := &domain.TenantGoogleCalendarConfig{
		ID:                          uuid.New(),
		TenantID:                    t2,
		IsEnabled:                   false,
		Status:                      domain.GoogleCalendarStatusDisconnected,
		GoogleClientID:              "client-t2",
		GoogleClientSecretEncrypted: "enc-secret-t2",
		CalendarID:                  "cal-t2",
	}

	if err := repo.UpsertGoogleCalendarConfig(ctx, cfg1); err != nil {
		t.Fatalf("upsert cfg1 falhou: %v", err)
	}
	if err := repo.UpsertGoogleCalendarConfig(ctx, cfg2); err != nil {
		t.Fatalf("upsert cfg2 falhou: %v", err)
	}

	got1, err := repo.GetGoogleCalendarConfigByTenantID(ctx, t1)
	if err != nil || got1 == nil || got1.GoogleClientID != "client-t1" {
		t.Fatalf("config t1 inesperada: %+v err=%v", got1, err)
	}
	got2, err := repo.GetGoogleCalendarConfigByTenantID(ctx, t2)
	if err != nil || got2 == nil || got2.GoogleClientID != "client-t2" {
		t.Fatalf("config t2 inesperada: %+v err=%v", got2, err)
	}

	// Cria agendamentos e syncs com falha para t1 e t2
	pro1, svc1, cust1 := uuid.New(), uuid.New(), uuid.New()
	_ = db.Create(&domain.Professional{ID: pro1, TenantID: t1, Name: "João", IsActive: true}).Error
	_ = db.Create(&domain.Service{ID: svc1, TenantID: t1, Name: "Corte", DurationMinutes: 30, Price: 50, IsActive: true}).Error
	_ = db.Create(&domain.Customer{ID: cust1, TenantID: t1, Name: "Carlos", Phone: "11999990000"}).Error

	apt1 := uuid.New()
	apt2 := uuid.New()
	_ = db.Create(&domain.Appointment{
		ID: apt1, TenantID: t1, ProfessionalID: pro1, ServiceID: svc1, CustomerID: cust1,
		StartAt: time.Now(), EndAt: time.Now().Add(30 * time.Minute), DurationMinutes: 30, TotalPrice: 50, Status: domain.StatusConfirmed,
	}).Error
	_ = db.Create(&domain.Appointment{
		ID: apt2, TenantID: t2, ProfessionalID: uuid.New(), ServiceID: uuid.New(), CustomerID: uuid.New(),
		StartAt: time.Now(), EndAt: time.Now().Add(30 * time.Minute), DurationMinutes: 30, TotalPrice: 50, Status: domain.StatusConfirmed,
	}).Error

	_ = repo.UpsertAppointmentCalendarSync(ctx, &domain.AppointmentCalendarSync{
		ID:            uuid.New(),
		TenantID:      t1,
		AppointmentID: apt1,
		Provider:      domain.CalendarProviderGoogle,
		SyncStatus:    domain.CalendarSyncStatusFailed,
		LastAction:    domain.CalendarSyncActionCreate,
		AttemptCount:  2,
		LastError:     "timeout google",
	})
	_ = repo.UpsertAppointmentCalendarSync(ctx, &domain.AppointmentCalendarSync{
		ID:            uuid.New(),
		TenantID:      t2,
		AppointmentID: apt2,
		Provider:      domain.CalendarProviderGoogle,
		SyncStatus:    domain.CalendarSyncStatusFailed,
		LastAction:    domain.CalendarSyncActionCreate,
		AttemptCount:  1,
		LastError:     "erro t2",
	})

	// Busca cross-tenant não pode retornar registro de outro tenant
	cross, err := repo.GetAppointmentCalendarSync(ctx, t1, apt2, domain.CalendarProviderGoogle)
	if err != nil || cross != nil {
		t.Fatalf("esperava nil para busca cross-tenant, obteve %+v err=%v", cross, err)
	}

	countT1, err := repo.CountFailedCalendarSyncs(ctx, t1, domain.CalendarProviderGoogle)
	if err != nil || countT1 != 1 {
		t.Fatalf("esperava 1 falha para t1, obteve %d err=%v", countT1, err)
	}

	failuresT1, err := repo.ListFailedCalendarSyncs(ctx, t1, domain.CalendarProviderGoogle)
	if err != nil || len(failuresT1) != 1 {
		t.Fatalf("esperava 1 registro de falha para t1, obteve %d err=%v", len(failuresT1), err)
	}
	if failuresT1[0].Appointment == nil || failuresT1[0].Appointment.Customer == nil || failuresT1[0].Appointment.Customer.Name != "Carlos" {
		t.Errorf("esperava preload de Appointment/Customer no ListFailedCalendarSyncs, obteve %+v", failuresT1[0].Appointment)
	}
}
