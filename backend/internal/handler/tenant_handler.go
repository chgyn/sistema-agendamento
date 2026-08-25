package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/hash"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type TenantHandler struct {
	repo *postgres.Repository
}

func NewTenantHandler(repo *postgres.Repository) *TenantHandler {
	return &TenantHandler{repo: repo}
}

// -------------------------------------------------------------
// GESTÃO DE ESTABELECIMENTOS (ADMIN_GLOBAL)
// -------------------------------------------------------------

// ListTenants lista todos os estabelecimentos da plataforma
func (h *TenantHandler) ListTenants(c *gin.Context) {
	search := c.Query("search")
	var onlyActive *bool
	if activeParam := c.Query("active"); activeParam != "" {
		val := activeParam == "true" || activeParam == "1"
		onlyActive = &val
	}

	tenants, err := h.repo.ListTenants(c.Request.Context(), search, onlyActive)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar estabelecimentos: "+err.Error())
		return
	}

	response.Success(c, tenants)
}

// GetTenantByID busca detalhes de um estabelecimento
func (h *TenantHandler) GetTenantByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	tenant, err := h.repo.GetTenantByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	response.Success(c, tenant)
}

type CreateTenantWithAdminDTO struct {
	Name                string `json:"name" binding:"required"`
	Slug                string `json:"slug" binding:"required"`
	Document            string `json:"document"`
	Phone               string `json:"phone" binding:"required"`
	Email               string `json:"email"`
	Address             string `json:"address"`
	City                string `json:"city"`
	State               string `json:"state"`
	LogoURL             string `json:"logo_url"`
	PrimaryColor        string `json:"primary_color"`
	SlotIntervalMinutes int    `json:"slot_interval_minutes"`
	
	// Administrador inicial do Tenant
	AdminName     string `json:"admin_name" binding:"required"`
	AdminEmail    string `json:"admin_email" binding:"required,email"`
	AdminPassword string `json:"admin_password" binding:"required,min=6"`
}

// CreateTenantWithAdmin cadastra um novo estabelecimento e seu administrador inicial em transação atômica
func (h *TenantHandler) CreateTenantWithAdmin(c *gin.Context) {
	var dto CreateTenantWithAdminDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	passHash, err := hash.HashPassword(dto.AdminPassword)
	if err != nil {
		response.InternalServerError(c, "Erro ao criptografar senha do administrador")
		return
	}

	slug := strings.ToLower(strings.TrimSpace(dto.Slug))
	interval := dto.SlotIntervalMinutes
	if interval <= 0 {
		interval = 30
	}
	color := dto.PrimaryColor
	if color == "" {
		color = "#10b981"
	}

	tenantID := uuid.New()
	tenant := domain.Tenant{
		ID:                  tenantID,
		Slug:                slug,
		Name:                dto.Name,
		Document:            dto.Document,
		Phone:               dto.Phone,
		Email:               dto.Email,
		Address:             dto.Address,
		City:                dto.City,
		State:               dto.State,
		LogoURL:             dto.LogoURL,
		PrimaryColor:        color,
		SlotIntervalMinutes: interval,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	adminUser := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenantID,
		Name:         dto.AdminName,
		Email:        dto.AdminEmail,
		PasswordHash: passHash,
		Role:         domain.RoleAdminTenant,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.CreateTenantWithAdmin(c.Request.Context(), &tenant, &adminUser); err != nil {
		if errors.Is(err, domain.ErrSlugAlreadyExists) {
			response.Conflict(c, "O identificador (slug) já está em uso")
			return
		}
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(c, "O e-mail do administrador já está cadastrado")
			return
		}
		response.InternalServerError(c, "Erro ao criar estabelecimento: "+err.Error())
		return
	}

	response.Created(c, gin.H{
		"tenant": tenant,
		"admin":  adminUser,
	}, "Estabelecimento e Administrador criados com sucesso")
}

// UpdateTenantGlobal atualiza qualquer estabelecimento pelo ADMIN_GLOBAL
func (h *TenantHandler) UpdateTenantGlobal(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	tenant, err := h.repo.GetTenantByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	var req struct {
		Name                string `json:"name"`
		Slug                string `json:"slug"`
		Document            string `json:"document"`
		Phone               string `json:"phone"`
		Email               string `json:"email"`
		Address             string `json:"address"`
		City                string `json:"city"`
		State               string `json:"state"`
		LogoURL             string `json:"logo_url"`
		PrimaryColor        string `json:"primary_color"`
		SlotIntervalMinutes int    `json:"slot_interval_minutes"`
		IsActive            *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Slug != "" {
		tenant.Slug = strings.ToLower(strings.TrimSpace(req.Slug))
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
	if req.IsActive != nil {
		tenant.IsActive = *req.IsActive
	}
	tenant.UpdatedAt = time.Now()

	if err := h.repo.UpdateTenant(c.Request.Context(), tenant); err != nil {
		response.InternalServerError(c, "Erro ao atualizar estabelecimento")
		return
	}

	response.Success(c, tenant, "Estabelecimento atualizado com sucesso")
}

// ToggleTenantStatus ativa ou inativa um estabelecimento
func (h *TenantHandler) ToggleTenantStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	tenant, err := h.repo.GetTenantByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.IsActive = !tenant.IsActive
	}

	if err := h.repo.UpdateTenantStatus(c.Request.Context(), id, req.IsActive); err != nil {
		response.InternalServerError(c, "Erro ao alterar status do estabelecimento")
		return
	}

	tenant.IsActive = req.IsActive
	statusMsg := "inativado"
	if req.IsActive {
		statusMsg = "ativado"
	}
	response.Success(c, tenant, "Estabelecimento "+statusMsg+" com sucesso")
}

// GetGlobalDashboard retorna estatísticas consolidadas da plataforma
func (h *TenantHandler) GetGlobalDashboard(c *gin.Context) {
	stats, err := h.repo.GetGlobalDashboardStats(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Erro ao obter dados do dashboard global")
		return
	}
	response.Success(c, stats)
}

// -------------------------------------------------------------
// CONFIGURAÇÕES DO TENANT (ADMIN_TENANT)
// -------------------------------------------------------------

func (h *TenantHandler) GetSettings(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok || tenantID == uuid.Nil {
		response.NotFound(c, "Estabelecimento não identificado")
		return
	}

	tenant, err := h.repo.GetTenantByID(c.Request.Context(), tenantID)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}
	response.Success(c, tenant)
}

func (h *TenantHandler) UpdateSettings(c *gin.Context) {
	tenantID, ok := middleware.GetTenantID(c)
	if !ok || tenantID == uuid.Nil {
		response.Forbidden(c, "Estabelecimento não identificado")
		return
	}

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
