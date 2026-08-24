package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type TenantHandler struct {
	repo *postgres.Repository
}

func NewTenantHandler(repo *postgres.Repository) *TenantHandler {
	return &TenantHandler{repo: repo}
}

func (h *TenantHandler) GetSettings(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	tenant, err := h.repo.GetTenantByID(c.Request.Context(), tenantID)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}
	response.Success(c, tenant)
}

func (h *TenantHandler) UpdateSettings(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	tenant, err := h.repo.GetTenantByID(c.Request.Context(), tenantID)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	var req struct {
		Name                string `json:"name"`
		Document            string `json:"document"`
		Phone               string `json:"phone"`
		Email               string `json:"email"`
		Address             string `json:"address"`
		City                string `json:"city"`
		State               string `json:"state"`
		LogoURL             string `json:"logo_url"`
		PrimaryColor        string `json:"primary_color"`
		SlotIntervalMinutes int    `json:"slot_interval_minutes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos")
		return
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	tenant.Document = req.Document
	tenant.Phone = req.Phone
	tenant.Email = req.Email
	tenant.Address = req.Address
	tenant.City = req.City
	tenant.State = req.State
	if req.LogoURL != "" {
		tenant.LogoURL = req.LogoURL
	}
	if req.PrimaryColor != "" {
		tenant.PrimaryColor = req.PrimaryColor
	}
	if req.SlotIntervalMinutes > 0 {
		tenant.SlotIntervalMinutes = req.SlotIntervalMinutes
	}
	tenant.UpdatedAt = time.Now()

	if err := h.repo.UpdateTenant(c.Request.Context(), tenant); err != nil {
		response.InternalServerError(c, "Erro ao atualizar dados do estabelecimento")
		return
	}

	response.Success(c, tenant, "Configurações atualizadas com sucesso")
}
