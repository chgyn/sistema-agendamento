package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/pkg/jwt"
	"github.com/sistema-agendamento/backend/pkg/response"
)

const (
	ContextUserID   = "user_id"
	ContextTenantID = "tenant_id"
	ContextUserRole = "user_role"
	ContextUserEmail= "user_email"
	ContextUserName = "user_name"
)

// AuthMiddleware valida o cabeçalho Authorization: Bearer <token> e injeta tenant_id e user_id no contexto
func AuthMiddleware(jwtSvc *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Cabeçalho de autenticação não fornecido")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Formato do token inválido, utilize Bearer <token>")
			c.Abort()
			return
		}

		claims, err := jwtSvc.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Token inválido ou expirado")
			c.Abort()
			return
		}

		// Injeta no contexto do Gin para uso em todos os handlers e repositórios
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextTenantID, claims.TenantID)
		c.Set(ContextUserRole, claims.Role)
		c.Set(ContextUserEmail, claims.Email)
		c.Set(ContextUserName, claims.Name)

		c.Next()
	}
}

// RequireRole restringe o endpoint apenas para certos perfis (ex: ADMIN)
func RequireRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(ContextUserRole)
		if !exists {
			response.Forbidden(c, "Acesso negado")
			c.Abort()
			return
		}

		role, ok := val.(domain.Role)
		if !ok {
			response.Forbidden(c, "Acesso negado")
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if role == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Você não tem permissão para acessar este recurso")
		c.Abort()
	}
}

// Helper para obter TenantID com segurança do contexto do request
func GetTenantID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextTenantID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// Helper para obter UserID
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextUserID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
