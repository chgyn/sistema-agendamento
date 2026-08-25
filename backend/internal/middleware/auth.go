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
		if claims.TenantID != nil {
			c.Set(ContextTenantID, *claims.TenantID)
		}
		
		// Normaliza role legada ADMIN para ADMIN_TENANT
		role := claims.Role
		if role == domain.RoleAdmin {
			role = domain.RoleAdminTenant
		}
		c.Set(ContextUserRole, role)
		c.Set(ContextUserEmail, claims.Email)
		c.Set(ContextUserName, claims.Name)

		c.Next()
	}
}

// RequireRole restringe o endpoint apenas para certos perfis (ex: ADMIN_GLOBAL, ADMIN_TENANT)
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

		// Normaliza ADMIN para ADMIN_TENANT
		if role == domain.RoleAdmin {
			role = domain.RoleAdminTenant
		}

		for _, r := range allowedRoles {
			if r == domain.RoleAdmin {
				r = domain.RoleAdminTenant
			}
			if role == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Você não tem permissão para acessar este recurso")
		c.Abort()
	}
}

// Helper para obter TenantID com segurança do contexto do request (autenticado)
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

// Helper para obter UserRole
func GetUserRole(c *gin.Context) (domain.Role, bool) {
	val, exists := c.Get(ContextUserRole)
	if !exists {
		return "", false
	}
	role, ok := val.(domain.Role)
	if !ok {
		return "", false
	}
	if role == domain.RoleAdmin {
		role = domain.RoleAdminTenant
	}
	return role, true
}

// IsAdminGlobal verifica se o usuário autenticado é Administrador Geral
func IsAdminGlobal(c *gin.Context) bool {
	role, ok := GetUserRole(c)
	return ok && role == domain.RoleAdminGlobal
}

// IsAdminTenant verifica se o usuário autenticado é Administrador de Tenant
func IsAdminTenant(c *gin.Context) bool {
	role, ok := GetUserRole(c)
	return ok && (role == domain.RoleAdminTenant || role == domain.RoleAdmin)
}

// IsOperator verifica se o usuário autenticado é Operador
func IsOperator(c *gin.Context) bool {
	role, ok := GetUserRole(c)
	return ok && role == domain.RoleOperator
}

// ResolveTenantID resolve o tenant_id de forma estrita contra adulteração:
// - Para ADMIN_TENANT e OPERATOR: sempre retorna estritamente o TenantID do token JWT.
// - Para ADMIN_GLOBAL: permite especificar ?tenant_id=... ou :tenant_id caso queira operar sobre um tenant específico.
func ResolveTenantID(c *gin.Context) (uuid.UUID, error) {
	if IsAdminGlobal(c) {
		// Tenta query param ou URL param
		param := c.Query("tenant_id")
		if param == "" {
			param = c.Param("tenant_id")
		}
		if param != "" {
			return uuid.Parse(param)
		}
		// Se o ADMIN_GLOBAL tiver um tenant_id próprio no token (opcional), retorna-o
		if tID, ok := GetTenantID(c); ok && tID != uuid.Nil {
			return tID, nil
		}
		return uuid.Nil, nil
	}

	// Usuário não é global: O tenant_id DO JWT é a única fonte da verdade
	tID, ok := GetTenantID(c)
	if !ok || tID == uuid.Nil {
		return uuid.Nil, domain.ErrTenantNotFound
	}
	return tID, nil
}
