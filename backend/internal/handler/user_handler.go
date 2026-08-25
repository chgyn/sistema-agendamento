package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/hash"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type UserHandler struct {
	repo *postgres.Repository
}

func NewUserHandler(repo *postgres.Repository) *UserHandler {
	return &UserHandler{repo: repo}
}

// -------------------------------------------------------------
// GESTÃO DE USUÁRIOS DE ESTABELECIMENTOS (TENANT USERS)
// -------------------------------------------------------------

// List lista usuários considerando o perfil e isolamento de tenant
func (h *UserHandler) List(c *gin.Context) {
	isGlobal := middleware.IsAdminGlobal(c)

	if isGlobal {
		var tenantID *uuid.UUID
		if param := c.Query("tenant_id"); param != "" {
			if parsed, err := uuid.Parse(param); err == nil {
				tenantID = &parsed
			}
		}
		roleFilter := c.Query("role")

		users, err := h.repo.ListAllUsers(c.Request.Context(), tenantID, roleFilter)
		if err != nil {
			response.InternalServerError(c, "Erro ao listar usuários: "+err.Error())
			return
		}
		response.Success(c, users)
		return
	}

	// Administrador de Tenant: estritamente o seu tenant
	tenantID, ok := middleware.GetTenantID(c)
	if !ok || tenantID == uuid.Nil {
		response.Forbidden(c, "Tenant não identificado")
		return
	}

	users, err := h.repo.ListUsersByTenant(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar usuários: "+err.Error())
		return
	}

	response.Success(c, users)
}

// GetByID busca um usuário com validação rigorosa de escopo de tenant
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "Usuário não encontrado")
			return
		}
		response.InternalServerError(c, "Erro ao buscar usuário")
		return
	}

	// Validação de isolamento: se não for admin global, o usuário alvo DEVE ser do mesmo tenant
	if !middleware.IsAdminGlobal(c) {
		callerTenantID, ok := middleware.GetTenantID(c)
		if !ok || user.TenantID == nil || *user.TenantID != callerTenantID {
			response.Forbidden(c, domain.ErrCannotManageOtherTenantUser.Error())
			return
		}
	}

	response.Success(c, user)
}

type CreateUserDTO struct {
	TenantID *uuid.UUID  `json:"tenant_id"`
	Name     string      `json:"name" binding:"required"`
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required,min=6"`
	Role     domain.Role `json:"role" binding:"required"`
}

// Create cria novo usuário respeitando as fronteiras de autorização
func (h *UserHandler) Create(c *gin.Context) {
	var dto CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	isGlobal := middleware.IsAdminGlobal(c)

	var targetTenantID *uuid.UUID

	if isGlobal {
		// Administrador Geral pode definir o tenant do usuário ou criar usuário sem tenant
		if dto.Role != domain.RoleAdminGlobal {
			if dto.TenantID == nil || *dto.TenantID == uuid.Nil {
				response.BadRequest(c, "O estabelecimento (tenant_id) é obrigatório para este perfil")
				return
			}
			targetTenantID = dto.TenantID
		}
	} else {
		// Administrador de Tenant:
		// 1. NÃO pode criar ADMIN_GLOBAL
		if dto.Role == domain.RoleAdminGlobal {
			response.Forbidden(c, domain.ErrCannotCreateGlobalAdmin.Error())
			return
		}

		// 2. Sempre vincula estritamente ao seu próprio Tenant
		callerTenantID, ok := middleware.GetTenantID(c)
		if !ok || callerTenantID == uuid.Nil {
			response.Forbidden(c, "Tenant não identificado no contexto")
			return
		}
		targetTenantID = &callerTenantID
	}

	passHash, err := hash.HashPassword(dto.Password)
	if err != nil {
		response.InternalServerError(c, "Erro ao processar senha")
		return
	}

	// Normaliza Role
	role := dto.Role
	if role == domain.RoleAdmin {
		role = domain.RoleAdminTenant
	}

	newUser := domain.User{
		ID:           uuid.New(),
		TenantID:     targetTenantID,
		Name:         dto.Name,
		Email:        dto.Email,
		PasswordHash: passHash,
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.CreateUser(c.Request.Context(), &newUser); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(c, "O e-mail já está cadastrado")
			return
		}
		response.InternalServerError(c, "Erro ao cadastrar usuário: "+err.Error())
		return
	}

	response.Created(c, newUser, "Usuário cadastrado com sucesso")
}

type UpdateUserDTO struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
	IsActive *bool       `json:"is_active"`
}

// Update atualiza os dados de um usuário
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Usuário não encontrado")
		return
	}

	isGlobal := middleware.IsAdminGlobal(c)

	if !isGlobal {
		// Valida se o usuário pertence ao mesmo tenant
		callerTenantID, ok := middleware.GetTenantID(c)
		if !ok || user.TenantID == nil || *user.TenantID != callerTenantID {
			response.Forbidden(c, domain.ErrCannotManageOtherTenantUser.Error())
			return
		}
	}

	var dto UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Role != "" {
		if !isGlobal && dto.Role == domain.RoleAdminGlobal {
			response.Forbidden(c, domain.ErrCannotCreateGlobalAdmin.Error())
			return
		}
		role := dto.Role
		if role == domain.RoleAdmin {
			role = domain.RoleAdminTenant
		}
		user.Role = role
	}
	if dto.IsActive != nil {
		user.IsActive = *dto.IsActive
	}
	if dto.Password != "" {
		passHash, err := hash.HashPassword(dto.Password)
		if err != nil {
			response.InternalServerError(c, "Erro ao criptografar nova senha")
			return
		}
		user.PasswordHash = passHash
	}
	user.UpdatedAt = time.Now()

	if err := h.repo.UpdateUser(c.Request.Context(), user); err != nil {
		response.InternalServerError(c, "Erro ao salvar alterações do usuário: "+err.Error())
		return
	}

	response.Success(c, user, "Usuário atualizado com sucesso")
}

// ToggleStatus ativa ou inativa um usuário
func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Usuário não encontrado")
		return
	}

	if !middleware.IsAdminGlobal(c) {
		callerTenantID, ok := middleware.GetTenantID(c)
		if !ok || user.TenantID == nil || *user.TenantID != callerTenantID {
			response.Forbidden(c, domain.ErrCannotManageOtherTenantUser.Error())
			return
		}
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Se não passar body, simplesmente inverte
		req.IsActive = !user.IsActive
	}

	if err := h.repo.UpdateUserStatus(c.Request.Context(), id, req.IsActive); err != nil {
		response.InternalServerError(c, "Erro ao atualizar status do usuário")
		return
	}

	user.IsActive = req.IsActive
	statusMsg := "inativado"
	if req.IsActive {
		statusMsg = "ativado"
	}
	response.Success(c, user, "Usuário "+statusMsg+" com sucesso")
}

// -------------------------------------------------------------
// GESTÃO DE ADMINISTRADORES GERAIS (GLOBAL ADMINS ONLY)
// -------------------------------------------------------------

// ListGlobalAdmins lista todos os administradores gerais da plataforma
func (h *UserHandler) ListGlobalAdmins(c *gin.Context) {
	admins, err := h.repo.ListGlobalAdmins(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Erro ao listar administradores gerais")
		return
	}
	response.Success(c, admins)
}

type CreateGlobalAdminDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateGlobalAdmin cadastra um novo Administrador Geral
func (h *UserHandler) CreateGlobalAdmin(c *gin.Context) {
	var dto CreateGlobalAdminDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	passHash, err := hash.HashPassword(dto.Password)
	if err != nil {
		response.InternalServerError(c, "Erro ao processar senha")
		return
	}

	adminUser := domain.User{
		ID:           uuid.New(),
		TenantID:     nil, // Global Admin não possui tenant_id
		Name:         dto.Name,
		Email:        dto.Email,
		PasswordHash: passHash,
		Role:         domain.RoleAdminGlobal,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.CreateUser(c.Request.Context(), &adminUser); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(c, "O e-mail já está cadastrado")
			return
		}
		response.InternalServerError(c, "Erro ao cadastrar administrador geral: "+err.Error())
		return
	}

	response.Created(c, adminUser, "Administrador Geral cadastrado com sucesso")
}
