package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type SubscriptionHandler struct {
	subService *service.SubscriptionService
}

func NewSubscriptionHandler(subService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: subService}
}

// GetMySubscription retorna a assinatura ativa e histórico de cobranças do tenant autenticado
func (h *SubscriptionHandler) GetMySubscription(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		response.BadRequest(c, "Estabelecimento não identificado no token")
		return
	}

	sub, err := h.subService.GetByTenantID(c.Request.Context(), tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			response.NotFound(c, "Nenhuma assinatura ativa encontrada para este estabelecimento")
			return
		}
		response.InternalServerError(c, "Falha ao buscar assinatura: "+err.Error())
		return
	}

	response.Success(c, sub)
}

// ListGlobal lista todas as assinaturas cadastradas na plataforma para o ADMIN_GLOBAL
func (h *SubscriptionHandler) ListGlobal(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")

	subs, err := h.subService.ListAll(c.Request.Context(), status, search)
	if err != nil {
		response.InternalServerError(c, "Falha ao listar assinaturas: "+err.Error())
		return
	}

	response.Success(c, subs)
}

// OverrideStatus permite ao Administrador Geral alterar manualmente o status de uma assinatura
func (h *SubscriptionHandler) OverrideStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID da assinatura inválido")
		return
	}

	var dto domain.OverrideSubscriptionStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Status inválido: "+err.Error())
		return
	}

	if err := h.subService.OverrideStatus(c.Request.Context(), id, dto.Status); err != nil {
		response.InternalServerError(c, "Falha ao alterar status da assinatura: "+err.Error())
		return
	}

	response.Success(c, gin.H{"status": dto.Status}, "Status da assinatura atualizado com sucesso!")
}
