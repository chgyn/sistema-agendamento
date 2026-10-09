package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

// RequireActiveSubscription valida se o estabelecimento possui uma assinatura ativa/válida
func RequireActiveSubscription(repo *postgres.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Administradores Globais têm acesso irrestrito
		if IsAdminGlobal(c) {
			c.Next()
			return
		}

		tenantID, ok := GetTenantID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Estabelecimento não identificado",
				"code":    "TENANT_REQUIRED",
			})
			c.Abort()
			return
		}

		sub, err := repo.GetSubscriptionByTenantID(c.Request.Context(), tenantID)
		if err != nil || sub == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Nenhuma assinatura ativa vinculada ao estabelecimento",
				"code":    "SUBSCRIPTION_REQUIRED",
			})
			c.Abort()
			return
		}

		// Libera acesso se o status for ACTIVE ou TRIAL
		if sub.Status == domain.SubscriptionStatusActive || sub.Status == domain.SubscriptionStatusTrial {
			// Valida se o período de validade da assinatura manual/promocional expirou
			if sub.CurrentPeriodEnd != nil && time.Now().After(*sub.CurrentPeriodEnd) {
				_ = repo.UpdateSubscriptionStatus(c.Request.Context(), sub.ID, domain.SubscriptionStatusExpired)
				c.JSON(http.StatusForbidden, gin.H{
					"success":             false,
					"error":               "O período de validade da sua assinatura expirou. Entre em contato com a administração para renovação.",
					"code":                "SUBSCRIPTION_EXPIRED",
					"subscription_status": domain.SubscriptionStatusExpired,
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// Status PENDING, OVERDUE, CANCELLED ou EXPIRED: Bloqueia acesso operacional
		c.JSON(http.StatusForbidden, gin.H{
			"success":             false,
			"error":               "Assinatura pendente ou inativa. Regularize o pagamento do seu plano para liberar as funcionalidades operacionais.",
			"code":                "SUBSCRIPTION_REQUIRED",
			"subscription_status": sub.Status,
			"payment_url":         sub.PaymentURL,
			"plan_name":           sub.Plan.Name,
		})
		c.Abort()
	}
}
