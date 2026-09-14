package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/queue"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type WhatsAppWebhookHandler struct {
	repo        *postgres.Repository
	queueClient *queue.QueueClient
}

func NewWhatsAppWebhookHandler(repo *postgres.Repository, queueClient *queue.QueueClient) *WhatsAppWebhookHandler {
	return &WhatsAppWebhookHandler{
		repo:        repo,
		queueClient: queueClient,
	}
}

// ReceiveWebhook recebe os eventos do WUZAPI de forma não-bloqueante
func (h *WhatsAppWebhookHandler) ReceiveWebhook(c *gin.Context) {
	var payload domain.WuzapiWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("⚠️ [WhatsApp Webhook] Erro ao decodificar JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}

	// Responde 200 OK imediatamente para liberar a conexão do WUZAPI
	c.JSON(http.StatusOK, gin.H{"status": "received"})

	// Processa apenas eventos do tipo "Message"
	if payload.Type != "Message" {
		return
	}

	event := payload.Event
	// Ignora mensagens enviadas pelo próprio bot (IsFromMe)
	if event.Info.Source.IsFromMe {
		return
	}

	// Ignora mensagens de grupos de WhatsApp (atendimento apenas individual 1:1)
	if event.Info.Source.IsGroup {
		return
	}

	// Extrai o número do telefone do remetente
	rawSender := event.Info.Source.Sender
	if rawSender == "" {
		rawSender = event.Info.Source.Chat
	}
	customerPhone := strings.Split(rawSender, "@")[0]
	customerPhone = strings.TrimSpace(customerPhone)

	if customerPhone == "" {
		return
	}

	// Extrai o conteúdo textual da mensagem
	messageText := strings.TrimSpace(event.Message.Conversation)
	if messageText == "" {
		messageText = strings.TrimSpace(event.Message.ExtendedTextMessage.Text)
	}

	if messageText == "" {
		// Mensagens de mídia pura sem legenda ignoradas
		return
	}

	customerName := strings.TrimSpace(event.Info.PushName)
	if customerName == "" {
		customerName = "Cliente"
	}

	ctx := c.Request.Context()

	// Identifica o Tenant através do Token da instância enviado pelo WUZAPI
	cfg, err := h.repo.GetWhatsAppConfigByInstanceToken(ctx, payload.Token)
	if err != nil || cfg == nil {
		log.Printf("⚠️ [WhatsApp Webhook] Token de instância não reconhecido: %s", payload.Token)
		return
	}

	if !cfg.IsAIEnabled {
		log.Printf("ℹ️ [WhatsApp Webhook] Atendimento com IA desativado para o tenant %s", cfg.TenantID)
		return
	}

	// Enfileira a tarefa no Asynq para processamento assíncrono pela IA
	taskPayload := queue.WhatsAppIncomingPayload{
		TenantID:      cfg.TenantID,
		InstanceToken: cfg.InstanceToken,
		CustomerPhone: customerPhone,
		CustomerName:  customerName,
		MessageID:     event.Info.ID,
		MessageText:   messageText,
		Timestamp:     time.Now(),
	}

	if h.queueClient != nil {
		if err := h.queueClient.EnqueueWhatsAppIncoming(ctx, taskPayload); err != nil {
			log.Printf("❌ [WhatsApp Webhook] Falha ao enfileirar task de atendimento: %v", err)
		}
	}
}
