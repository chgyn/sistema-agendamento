package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type PlanHandler struct {
	planService *service.PlanService
}

func NewPlanHandler(planService *service.PlanService) *PlanHandler {
	return &PlanHandler{planService: planService}
}

// ListActivePublic lista os planos ativos para exibição na página de cadastro
func (h *PlanHandler) ListActivePublic(c *gin.Context) {
	plans, err := h.planService.ListActivePublic(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Falha ao listar planos disponíveis: "+err.Error())
		return
	}
	response.Success(c, plans)
}

// List lista todos os planos com filtro de busca e status para ADMIN_GLOBAL
func (h *PlanHandler) List(c *gin.Context) {
	search := c.Query("search")
	var onlyActive *bool
	if status := c.Query("is_active"); status != "" {
		active := status == "true" || status == "1"
		onlyActive = &active
	}

	plans, err := h.planService.List(c.Request.Context(), search, onlyActive)
	if err != nil {
		response.InternalServerError(c, "Falha ao listar planos: "+err.Error())
		return
	}
	response.Success(c, plans)
}

// GetByID detalha um plano específico
func (h *PlanHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID do plano inválido")
		return
	}

	plan, err := h.planService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPlanNotFound) {
			response.NotFound(c, "Plano não encontrado")
			return
		}
		response.InternalServerError(c, "Falha ao buscar plano: "+err.Error())
		return
	}

	response.Success(c, plan)
}

// Create cria um novo plano comercial
func (h *PlanHandler) Create(c *gin.Context) {
	var dto domain.CreatePlanDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	plan, err := h.planService.Create(c.Request.Context(), dto)
	if err != nil {
		response.InternalServerError(c, "Falha ao criar plano: "+err.Error())
		return
	}

	response.Created(c, plan, "Plano criado com sucesso!")
}

// Update atualiza um plano existente
func (h *PlanHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID do plano inválido")
		return
	}

	var dto domain.UpdatePlanDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	plan, err := h.planService.Update(c.Request.Context(), id, dto)
	if err != nil {
		if errors.Is(err, domain.ErrPlanNotFound) {
			response.NotFound(c, "Plano não encontrado")
			return
		}
		response.InternalServerError(c, "Falha ao atualizar plano: "+err.Error())
		return
	}

	response.Success(c, plan, "Plano atualizado com sucesso!")
}

// ToggleStatus ativa ou inativa um plano comercial
func (h *PlanHandler) ToggleStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID do plano inválido")
		return
	}

	var dto domain.TogglePlanStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Status inválido: "+err.Error())
		return
	}

	if err := h.planService.ToggleStatus(c.Request.Context(), id, dto.IsActive); err != nil {
		if errors.Is(err, domain.ErrPlanNotFound) {
			response.NotFound(c, "Plano não encontrado")
			return
		}
		response.InternalServerError(c, "Falha ao alterar status do plano: "+err.Error())
		return
	}

	msg := "Plano ativado com sucesso!"
	if !dto.IsActive {
		msg = "Plano desativado com sucesso!"
	}
	response.Success(c, gin.H{"is_active": dto.IsActive}, msg)
}

// Delete remove um plano se não houver assinaturas vinculadas
func (h *PlanHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID do plano inválido")
		return
	}

	if err := h.planService.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrPlanHasSubscribers) {
			response.Conflict(c, err.Error())
			return
		}
		if errors.Is(err, domain.ErrPlanNotFound) {
			response.NotFound(c, "Plano não encontrado")
			return
		}
		response.InternalServerError(c, "Falha ao excluir plano: "+err.Error())
		return
	}

	response.Success(c, nil, "Plano excluído com sucesso!")
}
