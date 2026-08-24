package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Role types
type Role string

const (
	RoleAdmin    Role = "ADMIN"
	RoleOperator Role = "OPERATOR"
)

// AppointmentStatus types
type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "PENDING"
	StatusConfirmed AppointmentStatus = "CONFIRMED"
	StatusCompleted AppointmentStatus = "COMPLETED"
	StatusCancelled AppointmentStatus = "CANCELLED"
	StatusNoShow    AppointmentStatus = "NO_SHOW"
)

// Domain Errors
var (
	ErrTenantNotFound          = errors.New("estabelecimento não encontrado")
	ErrUserNotFound            = errors.New("usuário não encontrado")
	ErrInvalidCredentials      = errors.New("credenciais inválidas")
	ErrEmailAlreadyExists      = errors.New("o e-mail já está cadastrado")
	ErrSlugAlreadyExists       = errors.New("o identificador (slug) da barbearia/salão já está em uso")
	ErrUnauthorized            = errors.New("acesso não autorizado")
	ErrForbidden               = errors.New("você não tem permissão para realizar esta ação")
	ErrServiceNotFound         = errors.New("serviço não encontrado")
	ErrProfessionalNotFound    = errors.New("profissional não encontrado")
	ErrCustomerNotFound        = errors.New("cliente não encontrado")
	ErrAppointmentNotFound     = errors.New("agendamento não encontrado")
	ErrSlotAlreadyBooked       = errors.New("este horário já foi reservado por outro cliente. Por favor, escolha outro horário disponível")
	ErrOutsideWorkingHours     = errors.New("o horário selecionado está fora do expediente de atendimento do profissional")
	ErrSlotInBreak             = errors.New("o horário selecionado coincide com o intervalo do profissional")
	ErrSlotInException         = errors.New("o profissional possui um bloqueio ou folga neste horário")
	ErrInvalidAppointmentTime  = errors.New("horário de agendamento inválido ou no passado")
	ErrProfessionalNotAssigned = errors.New("o profissional selecionado não executa este serviço")
)

// Tenant representa a barbearia ou salão de beleza (conta isolada)
type Tenant struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Slug                string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Name                string    `gorm:"type:varchar(150);not null" json:"name"`
	Document            string    `gorm:"type:varchar(30)" json:"document"`
	Phone               string    `gorm:"type:varchar(30)" json:"phone"`
	Email               string    `gorm:"type:varchar(150)" json:"email"`
	Address             string    `gorm:"type:varchar(255)" json:"address"`
	City                string    `gorm:"type:varchar(100)" json:"city"`
	State               string    `gorm:"type:varchar(2)" json:"state"`
	LogoURL             string    `gorm:"type:text" json:"logo_url"`
	PrimaryColor        string    `gorm:"type:varchar(20);default:'#10b981'" json:"primary_color"`
	SlotIntervalMinutes int       `gorm:"default:30" json:"slot_interval_minutes"`
	IsActive            bool      `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	Users         []User         `gorm:"foreignKey:TenantID" json:"users,omitempty"`
	Professionals []Professional `gorm:"foreignKey:TenantID" json:"professionals,omitempty"`
	Services      []Service      `gorm:"foreignKey:TenantID" json:"services,omitempty"`
}

// User representa os operadores ou administradores vinculados a um Tenant
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID     uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name"`
	Email        string    `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         Role      `gorm:"type:varchar(30);default:'ADMIN'" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// Professional representa o barbeiro, cabeleireiro, manicure, etc.
type Professional struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Email       string    `gorm:"type:varchar(150)" json:"email"`
	Phone       string    `gorm:"type:varchar(30)" json:"phone"`
	Title       string    `gorm:"type:varchar(100)" json:"title"`       // Ex: Barbeiro Master, Colorista
	Specialty   string    `gorm:"type:varchar(150)" json:"specialty"`   // Ex: Cortes degradê e barba navalhada
	Bio         string    `gorm:"type:text" json:"bio"`
	AvatarURL   string    `gorm:"type:text" json:"avatar_url"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Services     []Service               `gorm:"many2many:professional_services;" json:"services,omitempty"`
	WorkingHours []WorkingHour           `gorm:"foreignKey:ProfessionalID" json:"working_hours,omitempty"`
	Exceptions   []AvailabilityException `gorm:"foreignKey:ProfessionalID" json:"exceptions,omitempty"`
}

// Service representa o serviço prestado (Corte, Barba, Tratamento)
type Service struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID        uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name            string    `gorm:"type:varchar(150);not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	DurationMinutes int       `gorm:"not null" json:"duration_minutes"`
	Price           float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Color           string    `gorm:"type:varchar(20);default:'#3b82f6'" json:"color"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Professionals []Professional `gorm:"many2many:professional_services;" json:"professionals,omitempty"`
}

// ProfessionalService tabela de junção muitos-para-muitos
type ProfessionalService struct {
	TenantID       uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	ProfessionalID uuid.UUID `gorm:"type:uuid;primaryKey" json:"professional_id"`
	ServiceID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"service_id"`
}

// WorkingHour horários regulares de trabalho de um profissional por dia da semana
// DayOfWeek: 0 = Domingo, 1 = Segunda, 2 = Terça, ..., 6 = Sábado
type WorkingHour struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID       uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	ProfessionalID uuid.UUID `gorm:"type:uuid;index;not null" json:"professional_id"`
	DayOfWeek      int       `gorm:"not null" json:"day_of_week"` // 0 a 6
	StartTime      string    `gorm:"type:varchar(5);not null" json:"start_time"` // "08:00"
	EndTime        string    `gorm:"type:varchar(5);not null" json:"end_time"`   // "18:00"
	BreakStart     string    `gorm:"type:varchar(5)" json:"break_start"`         // "12:00"
	BreakEnd       string    `gorm:"type:varchar(5)" json:"break_end"`           // "13:00"
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AvailabilityException bloqueios pontuais, feriados, atestados, folgas
type AvailabilityException struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID       uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	ProfessionalID uuid.UUID `gorm:"type:uuid;index;not null" json:"professional_id"`
	Date           string    `gorm:"type:varchar(10);not null" json:"date"` // "YYYY-MM-DD"
	StartTime      string    `gorm:"type:varchar(5)" json:"start_time"`     // "14:00"
	EndTime        string    `gorm:"type:varchar(5)" json:"end_time"`       // "16:00"
	Reason         string    `gorm:"type:varchar(200)" json:"reason"`
	IsFullDay      bool      `gorm:"default:false" json:"is_full_day"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Customer representa o cliente que agenda serviços
type Customer struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID    uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Phone       string    `gorm:"type:varchar(30);index;not null" json:"phone"`
	Email       string    `gorm:"type:varchar(150);index" json:"email"`
	Notes       string    `gorm:"type:text" json:"notes"`
	TotalVisits int       `gorm:"default:0" json:"total_visits"`
	TotalSpent  float64   `gorm:"type:decimal(10,2);default:0" json:"total_spent"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Appointment representa o agendamento de um atendimento
type Appointment struct {
	ID                 uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID           uuid.UUID         `gorm:"type:uuid;index;not null" json:"tenant_id"`
	ProfessionalID     uuid.UUID         `gorm:"type:uuid;index;not null" json:"professional_id"`
	ServiceID          uuid.UUID         `gorm:"type:uuid;index;not null" json:"service_id"`
	CustomerID         uuid.UUID         `gorm:"type:uuid;index;not null" json:"customer_id"`
	StartAt            time.Time         `gorm:"index;not null" json:"start_at"`
	EndAt              time.Time         `gorm:"index;not null" json:"end_at"`
	DurationMinutes    int               `gorm:"not null" json:"duration_minutes"`
	TotalPrice         float64           `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status             AppointmentStatus `gorm:"type:varchar(30);default:'CONFIRMED';index;not null" json:"status"`
	Notes              string            `gorm:"type:text" json:"notes"`
	CancellationReason string            `gorm:"type:text" json:"cancellation_reason"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`

	Professional *Professional `gorm:"foreignKey:ProfessionalID" json:"professional,omitempty"`
	Service      *Service      `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	Customer     *Customer     `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

// AppointmentStatusHistory histórico e auditoria de mudanças de status
type AppointmentStatusHistory struct {
	ID            uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID      uuid.UUID         `gorm:"type:uuid;index;not null" json:"tenant_id"`
	AppointmentID uuid.UUID         `gorm:"type:uuid;index;not null" json:"appointment_id"`
	OldStatus     AppointmentStatus `gorm:"type:varchar(30)" json:"old_status"`
	NewStatus     AppointmentStatus `gorm:"type:varchar(30);not null" json:"new_status"`
	ChangedBy     string            `gorm:"type:varchar(150)" json:"changed_by"`
	Reason        string            `gorm:"type:text" json:"reason"`
	CreatedAt     time.Time         `json:"created_at"`
}

// AuditLog log geral de auditoria para ações no sistema
type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID    *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	Action    string    `gorm:"type:varchar(100);not null" json:"action"`
	Entity    string    `gorm:"type:varchar(100);not null" json:"entity"`
	EntityID  string    `gorm:"type:varchar(100)" json:"entity_id"`
	Details   string    `gorm:"type:text" json:"details"`
	IPAddress string    `gorm:"type:varchar(50)" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
