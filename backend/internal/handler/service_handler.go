package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type ServiceHandler struct {
	repo *postgres.Repository
}

func NewServiceHandler(repo *postgres.Repository) *ServiceHandler {
	return &ServiceHandler{repo: repo}
}

func (h *ServiceHandler) List(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	services, err := h.repo.ListServices(c.Request.Context(), tenantID, false)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar serviços")
		return
	}
	response.Success(c, services)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	s, err := h.repo.GetServiceByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "Serviço não encontrado")
		return
	}
	response.Success(c, s)
}

func (h *ServiceHandler) Create(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)

	var req struct {
		Name            string    `json:"name" binding:"required"`
		Description     string    `json:"description"`
		DurationMinutes int       `json:"duration_minutes" binding:"required,min=5"`
		Price           float64   `json:"price" binding:"required,min=0"`
		Color           string    `json:"color"`
		IsActive        *bool     `json:"is_active"`
		ProfessionalIDs []uuid.UUID `json:"professional_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	// Validação de limite do plano para serviços
	if sub, err := h.repo.GetSubscriptionByTenantID(c.Request.Context(), tenantID); err == nil && sub != nil && sub.Plan != nil {
		if sub.Plan.MaxServices > 0 {
			services, _ := h.repo.ListServices(c.Request.Context(), tenantID, true)
			if len(services) >= sub.Plan.MaxServices {
				response.Forbidden(c, fmt.Sprintf("Limite do plano atingido: seu plano atual permite no máximo %d serviço(s). Faça upgrade para adicionar mais serviços.", sub.Plan.MaxServices))
				return
			}
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	if req.Color == "" {
		req.Color = "#10b981"
	}

	svc := domain.Service{
		ID:              uuid.New(),
		TenantID:        tenantID,
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Price:           req.Price,
		Color:           req.Color,
		IsActive:        isActive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.repo.CreateService(c.Request.Context(), &svc); err != nil {
		response.InternalServerError(c, "Erro ao salvar serviço")
		return
	}

	// Associa aos profissionais indicados
	if len(req.ProfessionalIDs) > 0 {
		for _, pid := range req.ProfessionalIDs {
			_ = h.repo.DB().Create(&domain.ProfessionalService{
				TenantID:       tenantID,
				ProfessionalID: pid,
				ServiceID:      svc.ID,
			})
		}
	}

	response.Created(c, svc, "Serviço cadastrado com sucesso")
}

func (h *ServiceHandler) Update(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	svc, err := h.repo.GetServiceByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "Serviço não encontrado")
		return
	}

	var req struct {
		Name            string    `json:"name"`
		Description     string    `json:"description"`
		DurationMinutes int       `json:"duration_minutes"`
		Price           float64   `json:"price"`
		Color           string    `json:"color"`
		IsActive        *bool     `json:"is_active"`
		ProfessionalIDs []uuid.UUID `json:"professional_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos")
		return
	}

	if req.Name != "" {
		svc.Name = req.Name
	}
	svc.Description = req.Description
	if req.DurationMinutes > 0 {
		svc.DurationMinutes = req.DurationMinutes
	}
	if req.Price >= 0 {
		svc.Price = req.Price
	}
	if req.Color != "" {
		svc.Color = req.Color
	}
	if req.IsActive != nil {
		svc.IsActive = *req.IsActive
	}
	svc.UpdatedAt = time.Now()

	if err := h.repo.UpdateService(c.Request.Context(), svc); err != nil {
		response.InternalServerError(c, "Erro ao atualizar serviço")
		return
	}

	if req.ProfessionalIDs != nil {
		_ = h.repo.DB().Where("tenant_id = ? AND service_id = ?", tenantID, svc.ID).Delete(&domain.ProfessionalService{})
		for _, pid := range req.ProfessionalIDs {
			_ = h.repo.DB().Create(&domain.ProfessionalService{
				TenantID:       tenantID,
				ProfessionalID: pid,
				ServiceID:      svc.ID,
			})
		}
	}

	response.Success(c, svc, "Serviço atualizado com sucesso")
}

func (h *ServiceHandler) Delete(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.repo.DeleteService(c.Request.Context(), tenantID, id); err != nil {
		response.InternalServerError(c, "Erro ao remover serviço")
		return
	}

	response.Success(c, nil, "Serviço removido com sucesso")
}
