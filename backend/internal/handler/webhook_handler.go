package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/queue"
	"github.com/sistema-agendamento/backend/internal/service"
)

type WebhookHandler struct {
	cfg         *config.Config
	queueClient *queue.QueueClient
	subService  *service.SubscriptionService
}

func NewWebhookHandler(cfg *config.Config, qClient *queue.QueueClient, subService *service.SubscriptionService) *WebhookHandler {
	return &WebhookHandler{
		cfg:         cfg,
		queueClient: qClient,
		subService:  subService,
	}
}

// HandleAsaas recebe eventos de webhook do Asaas e despacha para a fila assíncrona
func (h *WebhookHandler) HandleAsaas(c *gin.Context) {
	// Validação de token de acesso do webhook (se configurado)
	if h.cfg.AsaasWebhookSecret != "" {
		token := c.GetHeader("asaas-access-token")
		if token == "" {
			token = c.GetHeader("access_token")
		}
		if token != "" && token != h.cfg.AsaasWebhookSecret {
			log.Printf("⚠️ [Webhook Asaas] Token inválido recebido: %s", token)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de webhook inválido"})
			return
		}
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("⚠️ [Webhook Asaas] Payload JSON inválido: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	event, _ := payload["event"].(string)
	paymentMap, _ := payload["payment"].(map[string]interface{})
	subMap, _ := payload["subscription"].(map[string]interface{})

	log.Printf("📩 [Webhook Asaas Recebido] Evento: %s", event)

	// 1. Tenta enfileirar via Asynq / Redis para processamento desacoplado
	if h.queueClient != nil {
		err := h.queueClient.EnqueueAsaasWebhook(c.Request.Context(), queue.AsaasWebhookTaskPayload{
			Event:        event,
			Payment:      paymentMap,
			Subscription: subMap,
		})
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "queued", "event": event})
			return
		}
		log.Printf("⚠️ [Webhook Asaas] Fila indisponível (%v). Executando processamento síncrono de fallback...", err)
	}

	// 2. Fallback síncrono caso o Redis esteja indisponível
	if h.subService != nil {
		_ = h.subService.ProcessAsaasWebhookEvent(c.Request.Context(), event, paymentMap, subMap)
	}

	c.JSON(http.StatusOK, gin.H{"status": "received", "event": event})
}
