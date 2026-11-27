package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/queue"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type CreateAppointmentDTO struct {
	TenantID        uuid.UUID `json:"tenant_id"`
	ServiceID       uuid.UUID `json:"service_id" binding:"required"`
	ProfessionalID  uuid.UUID `json:"professional_id" binding:"required"`
	StartAt         time.Time `json:"start_at" binding:"required"`
	CustomerName    string    `json:"customer_name" binding:"required"`
	CustomerPhone   string    `json:"customer_phone" binding:"required"`
	CustomerEmail   string    `json:"customer_email"`
	Notes           string    `json:"notes"`
}

type AppointmentService struct {
	repo             *postgres.Repository
	queueClient      *queue.QueueClient
	calendarEnqueuer CalendarSyncEnqueuer
	googleCalRepo    domain.GoogleCalendarRepository
}

func NewAppointmentService(repo *postgres.Repository, qClient *queue.QueueClient) *AppointmentService {
	svc := &AppointmentService{
		repo:        repo,
		queueClient: qClient,
	}
	if qClient != nil {
		svc.calendarEnqueuer = qClient
	}
	if repo != nil {
		svc.googleCalRepo = repo
	}
	return svc
}

// SetCalendarSyncDependencies permite injetar mocks ou enqueuer customizado em testes
func (s *AppointmentService) SetCalendarSyncDependencies(enqueuer CalendarSyncEnqueuer, calRepo domain.GoogleCalendarRepository) {
	s.calendarEnqueuer = enqueuer
	s.googleCalRepo = calRepo
}

// CreateAppointment executa a reserva do atendimento com proteção contra dupla reserva
func (s *AppointmentService) CreateAppointment(ctx context.Context, dto CreateAppointmentDTO) (*domain.Appointment, error) {
	// 1. Carrega o serviço
	svc, err := s.repo.GetServiceByID(ctx, dto.TenantID, dto.ServiceID)
	if err != nil {
		return nil, domain.ErrServiceNotFound
	}

	// 2. Carrega o profissional
	pro, err := s.repo.GetProfessionalByID(ctx, dto.TenantID, dto.ProfessionalID)
	if err != nil {
		return nil, domain.ErrProfessionalNotFound
	}

	// 3. Valida horário não no passado
	now := time.Now()
	if dto.StartAt.Before(now.Add(-2 * time.Minute)) {
		return nil, domain.ErrInvalidAppointmentTime
	}

	// 4. Calcula data final com base na duração do serviço
	endAt := dto.StartAt.Add(time.Duration(svc.DurationMinutes) * time.Minute)

	// 5. Busca ou cadastra o cliente no tenant
	customer, err := s.repo.FindOrCreateCustomer(ctx, dto.TenantID, dto.CustomerName, dto.CustomerPhone, dto.CustomerEmail)
	if err != nil {
		return nil, err
	}

	// 6. Monta o agendamento
	apt := &domain.Appointment{
		ID:              uuid.New(),
		TenantID:        dto.TenantID,
		ProfessionalID:  dto.ProfessionalID,
		ServiceID:       dto.ServiceID,
		CustomerID:      customer.ID,
		StartAt:         dto.StartAt,
		EndAt:           endAt,
		DurationMinutes: svc.DurationMinutes,
		TotalPrice:      svc.Price,
		Status:          domain.StatusConfirmed,
		Notes:           dto.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Professional:    pro,
		Service:         svc,
		Customer:        customer,
	}

	// 7. Salva com lock concorrente no Postgres
	if err := s.repo.CreateAppointmentWithLock(ctx, apt); err != nil {
		return nil, err
	}

	// 8. Enfileira task assíncrona de notificação no Redis/Asynq
	if s.queueClient != nil {
		_ = s.queueClient.EnqueueBookingConfirmation(ctx, queue.BookingConfirmationPayload{
			TenantID:         dto.TenantID,
			AppointmentID:    apt.ID,
			CustomerName:     customer.Name,
			CustomerPhone:    customer.Phone,
			CustomerEmail:    customer.Email,
			ServiceName:      svc.Name,
			ProfessionalName: pro.Name,
			StartAt:          apt.StartAt,
			TotalPrice:       apt.TotalPrice,
		})
	}

	// 9. Enfileira sincronização com Google Calendar (não-bloqueante)
	s.enqueueGoogleCalendarSync(ctx, dto.TenantID, apt.ID, domain.CalendarSyncActionCreate)

	return apt, nil
}

func (s *AppointmentService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Appointment, error) {
	return s.repo.GetAppointmentByID(ctx, tenantID, id)
}

func (s *AppointmentService) List(ctx context.Context, tenantID uuid.UUID, filter postgres.AppointmentFilter) ([]domain.Appointment, error) {
	return s.repo.ListAppointments(ctx, tenantID, filter)
}

func (s *AppointmentService) Cancel(ctx context.Context, tenantID, id uuid.UUID, reason, changedBy string) error {
	if reason == "" {
		reason = "Cancelado pelo usuário/cliente"
	}
	if err := s.repo.UpdateAppointmentStatus(ctx, tenantID, id, domain.StatusCancelled, changedBy, reason); err != nil {
		return err
	}
	s.enqueueGoogleCalendarSync(ctx, tenantID, id, domain.CalendarSyncActionDelete)
	return nil
}

func (s *AppointmentService) Complete(ctx context.Context, tenantID, id uuid.UUID, changedBy string) error {
	return s.repo.UpdateAppointmentStatus(ctx, tenantID, id, domain.StatusCompleted, changedBy, "Atendimento concluído com sucesso")
}

func (s *AppointmentService) Reschedule(ctx context.Context, tenantID, id uuid.UUID, newStartAt time.Time, newProID uuid.UUID) error {
	apt, err := s.repo.GetAppointmentByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	newEndAt := newStartAt.Add(time.Duration(apt.DurationMinutes) * time.Minute)
	if err := s.repo.RescheduleAppointment(ctx, tenantID, id, newStartAt, newEndAt, newProID); err != nil {
		return err
	}
	s.enqueueGoogleCalendarSync(ctx, tenantID, id, domain.CalendarSyncActionUpdate)
	return nil
}

func (s *AppointmentService) enqueueGoogleCalendarSync(ctx context.Context, tenantID, appointmentID uuid.UUID, action domain.CalendarSyncAction) {
	if s.calendarEnqueuer == nil {
		return
	}

	err := s.calendarEnqueuer.EnqueueGoogleCalendarSync(ctx, queue.GoogleCalendarSyncPayload{
		TenantID:      tenantID,
		AppointmentID: appointmentID,
		Action:        action,
	})
	if err == nil {
		return
	}

	log.Printf("⚠️ [AppointmentService] Falha ao enfileirar sync Google Calendar (tenant=%s apt=%s action=%s): %v",
		tenantID, appointmentID, action, err)

	if s.googleCalRepo == nil {
		return
	}
	cfg, cfgErr := s.googleCalRepo.GetGoogleCalendarConfigByTenantID(ctx, tenantID)
	if cfgErr != nil || cfg == nil || !cfg.IsEnabled || cfg.Status != domain.GoogleCalendarStatusConnected {
		return
	}

	syncRec, recErr := s.googleCalRepo.GetAppointmentCalendarSync(ctx, tenantID, appointmentID, domain.CalendarProviderGoogle)
	if recErr != nil || syncRec == nil {
		syncRec = &domain.AppointmentCalendarSync{
			ID:            uuid.New(),
			TenantID:      tenantID,
			AppointmentID: appointmentID,
			Provider:      domain.CalendarProviderGoogle,
		}
	}
	syncRec.SyncStatus = domain.CalendarSyncStatusFailed
	syncRec.LastAction = action
	syncRec.LastError = "Falha ao enfileirar sincronização: " + err.Error()
	_ = s.googleCalRepo.UpsertAppointmentCalendarSync(ctx, syncRec)
}
