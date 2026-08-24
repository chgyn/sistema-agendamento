package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type PublicBookingHandler struct {
	repo                *postgres.Repository
	availabilityService *service.AvailabilityService
	appointmentService  *service.AppointmentService
}

func NewPublicBookingHandler(repo *postgres.Repository, availSvc *service.AvailabilityService, aptSvc *service.AppointmentService) *PublicBookingHandler {
	return &PublicBookingHandler{
		repo:                repo,
		availabilityService: availSvc,
		appointmentService:  aptSvc,
	}
}

// GetTenantInfo retorna informações públicas do estabelecimento
func (h *PublicBookingHandler) GetTenantInfo(c *gin.Context) {
	slug := c.Param("slug")
	tenant, err := h.repo.GetTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrTenantNotFound) {
			response.NotFound(c, "Estabelecimento não encontrado")
			return
		}
		response.InternalServerError(c, "Erro ao buscar dados do estabelecimento: "+err.Error())
		return
	}

	response.Success(c, tenant)
}

// ListServices retorna serviços ativos do estabelecimento
func (h *PublicBookingHandler) ListServices(c *gin.Context) {
	slug := c.Param("slug")
	tenant, err := h.repo.GetTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	services, err := h.repo.ListServices(c.Request.Context(), tenant.ID, true)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar serviços")
		return
	}

	response.Success(c, services)
}

// ListProfessionalsByService retorna profissionais habilitados para o serviço
func (h *PublicBookingHandler) ListProfessionalsByService(c *gin.Context) {
	slug := c.Param("slug")
	serviceIDStr := c.Param("service_id")

	tenant, err := h.repo.GetTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		response.BadRequest(c, "ID do serviço inválido")
		return
	}

	pros, err := h.repo.ListProfessionalsByService(c.Request.Context(), tenant.ID, serviceID)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar profissionais")
		return
	}

	response.Success(c, pros)
}

// GetAvailability consulta horários disponíveis
func (h *PublicBookingHandler) GetAvailability(c *gin.Context) {
	slug := c.Param("slug")
	serviceIDStr := c.Query("service_id")
	professionalIDStr := c.Query("professional_id")
	dateStr := c.Query("date") // "YYYY-MM-DD"

	if dateStr == "" || serviceIDStr == "" {
		response.BadRequest(c, "Parâmetros 'service_id' e 'date' são obrigatórios")
		return
	}

	tenant, err := h.repo.GetTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		response.BadRequest(c, "ID do serviço inválido")
		return
	}

	var professionalID uuid.UUID
	if professionalIDStr != "" && professionalIDStr != "any" {
		professionalID, err = uuid.Parse(professionalIDStr)
		if err != nil {
			response.BadRequest(c, "ID do profissional inválido")
			return
		}
	}

	avail, err := h.availabilityService.GetAvailableSlotsForDate(c.Request.Context(), tenant.ID, serviceID, professionalID, dateStr)
	if err != nil {
		response.BadRequest(c, "Erro ao calcular disponibilidade: "+err.Error())
		return
	}

	response.Success(c, avail)
}

type PublicBookingRequest struct {
	ServiceID      uuid.UUID `json:"service_id" binding:"required"`
	ProfessionalID uuid.UUID `json:"professional_id" binding:"required"`
	StartAt        string    `json:"start_at" binding:"required"` // ISO string / RFC3339
	CustomerName   string    `json:"customer_name" binding:"required"`
	CustomerPhone  string    `json:"customer_phone" binding:"required"`
	CustomerEmail  string    `json:"customer_email"`
	Notes          string    `json:"notes"`
}

// CreateAppointment cria o agendamento público com proteção concorrente
func (h *PublicBookingHandler) CreateAppointment(c *gin.Context) {
	slug := c.Param("slug")
	tenant, err := h.repo.GetTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	var req PublicBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		// Tenta formato "2006-01-02T15:04:05" ou "2006-01-02 15:04"
		startAt, err = time.Parse("2006-01-02T15:04:00", req.StartAt)
		if err != nil {
			startAt, err = time.Parse("2006-01-02 15:04", req.StartAt)
			if err != nil {
				response.BadRequest(c, "Formato de data e hora inválido (use ISO8601 ou RFC3339)")
				return
			}
		}
	}

	dto := service.CreateAppointmentDTO{
		TenantID:        tenant.ID,
		ServiceID:       req.ServiceID,
		ProfessionalID:  req.ProfessionalID,
		StartAt:         startAt,
		CustomerName:    req.CustomerName,
		CustomerPhone:   req.CustomerPhone,
		CustomerEmail:   req.CustomerEmail,
		Notes:           req.Notes,
	}

	apt, err := h.appointmentService.CreateAppointment(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domain.ErrSlotAlreadyBooked) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidAppointmentTime) || errors.Is(err, domain.ErrServiceNotFound) || errors.Is(err, domain.ErrProfessionalNotFound) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Erro ao registrar agendamento: "+err.Error())
		return
	}

	response.Created(c, apt, "Agendamento realizado com sucesso!")
}

// GetAppointmentDetails consulta os detalhes do agendamento
func (h *PublicBookingHandler) GetAppointmentDetails(c *gin.Context) {
	slug := c.Param("slug")
	tenant, err := h.repo.GetTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID de agendamento inválido")
		return
	}

	apt, err := h.appointmentService.GetByID(c.Request.Context(), tenant.ID, id)
	if err != nil {
		response.NotFound(c, "Agendamento não encontrado")
		return
	}

	response.Success(c, apt)
}

// CancelAppointment cancela um agendamento público
func (h *PublicBookingHandler) CancelAppointment(c *gin.Context) {
	slug := c.Param("slug")
	tenant, err := h.repo.GetTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID de agendamento inválido")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := h.appointmentService.Cancel(c.Request.Context(), tenant.ID, id, req.Reason, "Cliente (Página Pública)"); err != nil {
		response.BadRequest(c, "Erro ao cancelar agendamento: "+err.Error())
		return
	}

	response.Success(c, nil, "Agendamento cancelado com sucesso")
}
