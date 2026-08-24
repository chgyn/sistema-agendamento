package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
	repo               *postgres.Repository
}

func NewAppointmentHandler(aptSvc *service.AppointmentService, repo *postgres.Repository) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: aptSvc,
		repo:               repo,
	}
}

// List lista os agendamentos do tenant com filtros
func (h *AppointmentHandler) List(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		response.Unauthorized(c, "Não autorizado")
		return
	}

	filter := postgres.AppointmentFilter{}

	if proIDStr := c.Query("professional_id"); proIDStr != "" {
		if proID, err := uuid.Parse(proIDStr); err == nil {
			filter.ProfessionalID = &proID
		}
	}

	if custIDStr := c.Query("customer_id"); custIDStr != "" {
		if custID, err := uuid.Parse(custIDStr); err == nil {
			filter.CustomerID = &custID
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.AppointmentStatus(statusStr)
		filter.Status = &status
	}

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if d, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			filter.DateFrom = &d
		}
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if d, err := time.Parse("2006-01-02", dateToStr); err == nil {
			end := d.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.DateTo = &end
		}
	}

	apts, err := h.appointmentService.List(c.Request.Context(), tenantID, filter)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar agendamentos: "+err.Error())
		return
	}

	response.Success(c, apts)
}

// GetByID busca um agendamento por ID
func (h *AppointmentHandler) GetByID(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	apt, err := h.appointmentService.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "Agendamento não encontrado")
		return
	}

	response.Success(c, apt)
}

// CreateAdmin cria um agendamento direto pelo painel de controle
func (h *AppointmentHandler) CreateAdmin(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)

	var req struct {
		ServiceID      uuid.UUID `json:"service_id" binding:"required"`
		ProfessionalID uuid.UUID `json:"professional_id" binding:"required"`
		StartAt        string    `json:"start_at" binding:"required"`
		CustomerName   string    `json:"customer_name" binding:"required"`
		CustomerPhone  string    `json:"customer_phone" binding:"required"`
		CustomerEmail  string    `json:"customer_email"`
		Notes          string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		startAt, err = time.Parse("2006-01-02T15:04:00", req.StartAt)
		if err != nil {
			startAt, err = time.Parse("2006-01-02 15:04", req.StartAt)
			if err != nil {
				response.BadRequest(c, "Formato de data inválido")
				return
			}
		}
	}

	dto := service.CreateAppointmentDTO{
		TenantID:        tenantID,
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
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, apt, "Agendamento criado com sucesso!")
}

// Complete conclui o atendimento
func (h *AppointmentHandler) Complete(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.appointmentService.Complete(c.Request.Context(), tenantID, id, "Administrador"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil, "Agendamento marcado como concluído")
}

// Cancel cancela o agendamento
func (h *AppointmentHandler) Cancel(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := h.appointmentService.Cancel(c.Request.Context(), tenantID, id, req.Reason, "Administrador"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil, "Agendamento cancelado")
}

// Reschedule reagenda o horário do atendimento
func (h *AppointmentHandler) Reschedule(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		StartAt        string    `json:"start_at" binding:"required"`
		ProfessionalID uuid.UUID `json:"professional_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos")
		return
	}

	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		startAt, err = time.Parse("2006-01-02T15:04:00", req.StartAt)
		if err != nil {
			response.BadRequest(c, "Formato de data inválido")
			return
		}
	}

	if err := h.appointmentService.Reschedule(c.Request.Context(), tenantID, id, startAt, req.ProfessionalID); err != nil {
		if errors.Is(err, domain.ErrSlotAlreadyBooked) {
			response.Conflict(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil, "Agendamento reagendado com sucesso")
}
