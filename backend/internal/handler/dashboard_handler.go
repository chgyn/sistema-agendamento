package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type DashboardHandler struct {
	tenantService *service.TenantService
}

func NewDashboardHandler(tenantService *service.TenantService) *DashboardHandler {
	return &DashboardHandler{tenantService: tenantService}
}

func (h *DashboardHandler) GetKPIs(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		response.Unauthorized(c, "Não autorizado")
		return
	}

	kpis, err := h.tenantService.GetDashboardKPIs(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalServerError(c, "Erro ao obter KPIs do dashboard: "+err.Error())
		return
	}

	response.Success(c, kpis)
}

func (h *DashboardHandler) GetAnalytics(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok {
		response.Unauthorized(c, "Não autorizado")
		return
	}

	analytics, err := h.tenantService.GetDashboardAnalytics(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalServerError(c, "Erro ao obter análises e relatórios do dashboard: "+err.Error())
		return
	}

	response.Success(c, analytics)
}

