package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type CustomerHandler struct {
	repo *postgres.Repository
}

func NewCustomerHandler(repo *postgres.Repository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

func (h *CustomerHandler) List(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	search := c.Query("search")

	customers, err := h.repo.ListCustomers(c.Request.Context(), tenantID, search)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar clientes")
		return
	}
	response.Success(c, customers)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	cust, err := h.repo.GetCustomerByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "Cliente não encontrado")
		return
	}

	// Busca histórico de agendamentos deste cliente
	apts, _ := h.repo.ListAppointments(c.Request.Context(), tenantID, postgres.AppointmentFilter{
		CustomerID: &id,
	})

	response.Success(c, gin.H{
		"customer":     cust,
		"appointments": apts,
	})
}
