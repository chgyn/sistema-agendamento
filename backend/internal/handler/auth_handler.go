package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
	repo        *postgres.Repository
}

func NewAuthHandler(authService *service.AuthService, repo *postgres.Repository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		repo:        repo,
	}
}

// RegisterTenant cria um novo estabelecimento e usuário administrador (Onboarding)
func (h *AuthHandler) RegisterTenant(c *gin.Context) {
	var dto service.RegisterTenantDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	res, err := h.authService.RegisterTenant(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domain.ErrSlugAlreadyExists) || errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, "Falha ao registrar estabelecimento: "+err.Error())
		return
	}

	response.Created(c, res, "Estabelecimento registrado com sucesso!")
}

// Login autentica um usuário administrador ou operador
func (h *AuthHandler) Login(c *gin.Context) {
	var dto service.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	res, err := h.authService.Login(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrUnauthorized) {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalServerError(c, "Falha ao autenticar: "+err.Error())
		return
	}

	response.Success(c, res, "Login efetuado com sucesso")
}

// Me retorna os dados do usuário autenticado e seu tenant
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "Não autenticado")
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "Usuário não encontrado")
		return
	}

	tenant, err := h.repo.GetTenantByID(c.Request.Context(), user.TenantID)
	if err != nil {
		response.NotFound(c, "Estabelecimento não encontrado")
		return
	}

	response.Success(c, gin.H{
		"user":   user,
		"tenant": tenant,
	})
}
