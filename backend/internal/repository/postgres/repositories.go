package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository gerencia todas as operações de persistência com isolamento multi-tenant
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

// -------------------------------------------------------------
// TENANT REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) GetTenantBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).
		Preload("Services", "is_active = ?", true).
		Preload("Professionals", "is_active = ?", true).
		First(&tenant, "slug = ? AND is_active = ?", slug, true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) CreateTenant(ctx context.Context, tenant *domain.Tenant) error {
	var count int64
	r.db.WithContext(ctx).Model(&domain.Tenant{}).Where("slug = ?", tenant.Slug).Count(&count)
	if count > 0 {
		return domain.ErrSlugAlreadyExists
	}
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *Repository) UpdateTenant(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

// -------------------------------------------------------------
// USER REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Preload("Tenant").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Preload("Tenant").First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	var count int64
	r.db.WithContext(ctx).Model(&domain.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return domain.ErrEmailAlreadyExists
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&users).Error
	return users, err
}

// -------------------------------------------------------------
// SERVICE REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetServiceByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Service, error) {
	var s domain.Service
	if err := r.db.WithContext(ctx).Preload("Professionals").First(&s, "tenant_id = ? AND id = ?", tenantID, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrServiceNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListServices(ctx context.Context, tenantID uuid.UUID, onlyActive bool) ([]domain.Service, error) {
	var services []domain.Service
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Preload("Professionals").Order("name asc").Find(&services).Error
	return services, err
}

func (r *Repository) CreateService(ctx context.Context, s *domain.Service) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *Repository) UpdateService(ctx context.Context, s *domain.Service) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", s.TenantID, s.ID).Save(s).Error
}

func (r *Repository) DeleteService(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Service{}).Error
}

// -------------------------------------------------------------
// PROFESSIONAL REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetProfessionalByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Professional, error) {
	var p domain.Professional
	if err := r.db.WithContext(ctx).
		Preload("Services").
		Preload("WorkingHours").
		Preload("Exceptions").
		First(&p, "tenant_id = ? AND id = ?", tenantID, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProfessionalNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListProfessionals(ctx context.Context, tenantID uuid.UUID, onlyActive bool) ([]domain.Professional, error) {
	var pros []domain.Professional
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Preload("Services").Preload("WorkingHours").Order("name asc").Find(&pros).Error
	return pros, err
}

func (r *Repository) ListProfessionalsByService(ctx context.Context, tenantID, serviceID uuid.UUID) ([]domain.Professional, error) {
	var pros []domain.Professional
	err := r.db.WithContext(ctx).
		Joins("JOIN professional_services ps ON ps.professional_id = professionals.id").
		Where("professionals.tenant_id = ? AND professionals.is_active = ? AND ps.service_id = ?", tenantID, true, serviceID).
		Preload("WorkingHours").
		Order("professionals.name asc").
		Find(&pros).Error
	return pros, err
}

func (r *Repository) CreateProfessional(ctx context.Context, p *domain.Professional, serviceIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		for _, sID := range serviceIDs {
			ps := domain.ProfessionalService{
				TenantID:       p.TenantID,
				ProfessionalID: p.ID,
				ServiceID:      sID,
			}
			if err := tx.Create(&ps).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) UpdateProfessional(ctx context.Context, p *domain.Professional, serviceIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tenant_id = ? AND id = ?", p.TenantID, p.ID).Save(p).Error; err != nil {
			return err
		}
		if serviceIDs != nil {
			// Atualiza associação
			if err := tx.Where("tenant_id = ? AND professional_id = ?", p.TenantID, p.ID).Delete(&domain.ProfessionalService{}).Error; err != nil {
				return err
			}
			for _, sID := range serviceIDs {
				ps := domain.ProfessionalService{
					TenantID:       p.TenantID,
					ProfessionalID: p.ID,
					ServiceID:      sID,
				}
				if err := tx.Create(&ps).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *Repository) DeleteProfessional(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_ = tx.Where("tenant_id = ? AND professional_id = ?", tenantID, id).Delete(&domain.ProfessionalService{})
		_ = tx.Where("tenant_id = ? AND professional_id = ?", tenantID, id).Delete(&domain.WorkingHour{})
		_ = tx.Where("tenant_id = ? AND professional_id = ?", tenantID, id).Delete(&domain.AvailabilityException{})
		return tx.Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.Professional{}).Error
	})
}

// -------------------------------------------------------------
// WORKING HOURS & EXCEPTIONS REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetWorkingHoursByProfessional(ctx context.Context, tenantID, professionalID uuid.UUID) ([]domain.WorkingHour, error) {
	var whs []domain.WorkingHour
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND professional_id = ?", tenantID, professionalID).
		Order("day_of_week asc").
		Find(&whs).Error
	return whs, err
}

func (r *Repository) SaveWorkingHours(ctx context.Context, tenantID, professionalID uuid.UUID, whs []domain.WorkingHour) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tenant_id = ? AND professional_id = ?", tenantID, professionalID).Delete(&domain.WorkingHour{}).Error; err != nil {
			return err
		}
		for i := range whs {
			whs[i].ID = uuid.New()
			whs[i].TenantID = tenantID
			whs[i].ProfessionalID = professionalID
			if err := tx.Create(&whs[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) CreateException(ctx context.Context, exp *domain.AvailabilityException) error {
	return r.db.WithContext(ctx).Create(exp).Error
}

func (r *Repository) ListExceptions(ctx context.Context, tenantID, professionalID uuid.UUID, dateFrom, dateTo string) ([]domain.AvailabilityException, error) {
	var exps []domain.AvailabilityException
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if professionalID != uuid.Nil {
		query = query.Where("professional_id = ?", professionalID)
	}
	if dateFrom != "" {
		query = query.Where("date >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("date <= ?", dateTo)
	}
	err := query.Order("date asc, start_time asc").Find(&exps).Error
	return exps, err
}

func (r *Repository) DeleteException(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.AvailabilityException{}).Error
}

// -------------------------------------------------------------
// CUSTOMER REPOSITORY
// -------------------------------------------------------------

func (r *Repository) GetCustomerByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Customer, error) {
	var c domain.Customer
	if err := r.db.WithContext(ctx).First(&c, "tenant_id = ? AND id = ?", tenantID, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCustomerNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) FindOrCreateCustomer(ctx context.Context, tenantID uuid.UUID, name, phone, email string) (*domain.Customer, error) {
	var customer domain.Customer
	// Procura por telefone ou email dentro do tenant
	query := r.db.WithContext(ctx).Where("tenant_id = ? AND (phone = ? OR (email != '' AND email = ?))", tenantID, phone, email)
	err := query.First(&customer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		customer = domain.Customer{
			ID:          uuid.New(),
			TenantID:    tenantID,
			Name:        name,
			Phone:       phone,
			Email:       email,
			TotalVisits: 0,
			TotalSpent:  0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := r.db.WithContext(ctx).Create(&customer).Error; err != nil {
			return nil, err
		}
		return &customer, nil
	} else if err != nil {
		return nil, err
	}

	// Atualiza o nome se foi fornecido um mais completo
	if name != "" && customer.Name != name {
		customer.Name = name
		_ = r.db.WithContext(ctx).Save(&customer)
	}

	return &customer, nil
}

func (r *Repository) ListCustomers(ctx context.Context, tenantID uuid.UUID, search string) ([]domain.Customer, error) {
	var customers []domain.Customer
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like)
	}
	err := query.Order("name asc").Find(&customers).Error
	return customers, err
}

// -------------------------------------------------------------
// APPOINTMENT REPOSITORY (CONCURRENCY SAFE ENGINE)
// -------------------------------------------------------------

type AppointmentFilter struct {
	ProfessionalID *uuid.UUID
	CustomerID     *uuid.UUID
	Status         *domain.AppointmentStatus
	DateFrom       *time.Time
	DateTo         *time.Time
}

func (r *Repository) GetAppointmentByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Appointment, error) {
	var apt domain.Appointment
	if err := r.db.WithContext(ctx).
		Preload("Professional").
		Preload("Service").
		Preload("Customer").
		First(&apt, "tenant_id = ? AND id = ?", tenantID, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAppointmentNotFound
		}
		return nil, err
	}
	return &apt, nil
}

func (r *Repository) ListAppointments(ctx context.Context, tenantID uuid.UUID, filter AppointmentFilter) ([]domain.Appointment, error) {
	var appointments []domain.Appointment
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if filter.ProfessionalID != nil && *filter.ProfessionalID != uuid.Nil {
		query = query.Where("professional_id = ?", *filter.ProfessionalID)
	}
	if filter.CustomerID != nil && *filter.CustomerID != uuid.Nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.DateFrom != nil {
		query = query.Where("start_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("start_at <= ?", *filter.DateTo)
	}

	err := query.
		Preload("Professional").
		Preload("Service").
		Preload("Customer").
		Order("start_at asc").
		Find(&appointments).Error
	return appointments, err
}

// CreateAppointmentWithLock realiza a criação segura sob concorrência e transação atômica
func (r *Repository) CreateAppointmentWithLock(ctx context.Context, apt *domain.Appointment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Bloqueia linha do profissional para serializar requisições concorrentes no mesmo profissional
		var pro domain.Professional
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&pro, "tenant_id = ? AND id = ?", apt.TenantID, apt.ProfessionalID).Error; err != nil {
			return domain.ErrProfessionalNotFound
		}

		// 2. Verifica se já existe agendamento ativo com sobreposição de horário
		// Intervalo [start_at, end_at] sobrepõe [Apt.StartAt, Apt.EndAt] se:
		// existing.StartAt < apt.EndAt AND existing.EndAt > apt.StartAt
		var conflictCount int64
		err := tx.Model(&domain.Appointment{}).
			Where("tenant_id = ? AND professional_id = ? AND status IN (?, ?) AND start_at < ? AND end_at > ?",
				apt.TenantID, apt.ProfessionalID, domain.StatusPending, domain.StatusConfirmed, apt.EndAt, apt.StartAt).
			Count(&conflictCount).Error
		if err != nil {
			return err
		}

		if conflictCount > 0 {
			return domain.ErrSlotAlreadyBooked
		}

		// 3. Insere o agendamento sem re-inserir associações aninhadas
		if err := tx.Omit(clause.Associations).Create(apt).Error; err != nil {
			return fmt.Errorf("falha ao salvar agendamento: %w", err)
		}

		// 4. Registra histórico de status
		history := domain.AppointmentStatusHistory{
			ID:            uuid.New(),
			TenantID:      apt.TenantID,
			AppointmentID: apt.ID,
			OldStatus:     "",
			NewStatus:     apt.Status,
			ChangedBy:     "Sistema de Agendamento",
			Reason:        "Criação inicial do agendamento",
			CreatedAt:     time.Now(),
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		// 5. Atualiza métricas do cliente (visitas e gasto)
		_ = tx.Model(&domain.Customer{}).
			Where("tenant_id = ? AND id = ?", apt.TenantID, apt.CustomerID).
			Updates(map[string]interface{}{
				"total_visits": gorm.Expr("total_visits + 1"),
				"total_spent":  gorm.Expr("total_spent + ?", apt.TotalPrice),
				"updated_at":   time.Now(),
			})

		return nil
	})
}

// UpdateAppointmentStatus altera o status com log de auditoria
func (r *Repository) UpdateAppointmentStatus(ctx context.Context, tenantID, id uuid.UUID, newStatus domain.AppointmentStatus, changedBy, reason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var apt domain.Appointment
		if err := tx.First(&apt, "tenant_id = ? AND id = ?", tenantID, id).Error; err != nil {
			return domain.ErrAppointmentNotFound
		}

		oldStatus := apt.Status
		apt.Status = newStatus
		if newStatus == domain.StatusCancelled && reason != "" {
			apt.CancellationReason = reason
		}

		if err := tx.Save(&apt).Error; err != nil {
			return err
		}

		// Registra no histórico
		history := domain.AppointmentStatusHistory{
			ID:            uuid.New(),
			TenantID:      tenantID,
			AppointmentID: id,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			ChangedBy:     changedBy,
			Reason:        reason,
			CreatedAt:     time.Now(),
		}
		return tx.Create(&history).Error
	})
}

// RescheduleAppointment altera horário com validação de colisão e lock
func (r *Repository) RescheduleAppointment(ctx context.Context, tenantID, id uuid.UUID, newStartAt, newEndAt time.Time, newProfessionalID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var apt domain.Appointment
		if err := tx.First(&apt, "tenant_id = ? AND id = ?", tenantID, id).Error; err != nil {
			return domain.ErrAppointmentNotFound
		}

		targetProID := apt.ProfessionalID
		if newProfessionalID != uuid.Nil {
			targetProID = newProfessionalID
		}

		// Bloqueia linha do profissional
		var pro domain.Professional
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&pro, "tenant_id = ? AND id = ?", tenantID, targetProID).Error; err != nil {
			return domain.ErrProfessionalNotFound
		}

		// Checa sobreposição com outros agendamentos exceto ele mesmo
		var conflictCount int64
		err := tx.Model(&domain.Appointment{}).
			Where("tenant_id = ? AND professional_id = ? AND id != ? AND status IN (?, ?) AND start_at < ? AND end_at > ?",
				tenantID, targetProID, id, domain.StatusPending, domain.StatusConfirmed, newEndAt, newStartAt).
			Count(&conflictCount).Error
		if err != nil {
			return err
		}
		if conflictCount > 0 {
			return domain.ErrSlotAlreadyBooked
		}

		apt.ProfessionalID = targetProID
		apt.StartAt = newStartAt
		apt.EndAt = newEndAt
		apt.Status = domain.StatusConfirmed
		apt.UpdatedAt = time.Now()

		if err := tx.Save(&apt).Error; err != nil {
			return err
		}

		history := domain.AppointmentStatusHistory{
			ID:            uuid.New(),
			TenantID:      tenantID,
			AppointmentID: id,
			OldStatus:     apt.Status,
			NewStatus:     domain.StatusConfirmed,
			ChangedBy:     "Reagendamento",
			Reason:        fmt.Sprintf("Reagendado para %s", newStartAt.Format("02/01/2006 15:04")),
			CreatedAt:     time.Now(),
		}
		return tx.Create(&history).Error
	})
}
