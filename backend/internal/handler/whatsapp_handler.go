package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type WhatsAppHandler struct {
	waService *service.WhatsAppService
}

func NewWhatsAppHandler(waService *service.WhatsAppService) *WhatsAppHandler {
	return &WhatsAppHandler{
		waService: waService,
	}
}

// resolveTenantID extrai o tenantID autenticado do contexto Gin
func resolveTenantID(c *gin.Context) (uuid.UUID, bool) {
	tenantIDVal, exists := c.Get("tenant_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Tenant não identificado no contexto da sessão")
		return uuid.Nil, false
	}

	tenantID, ok := tenantIDVal.(uuid.UUID)
	if !ok {
		// Tenta string
		if tStr, okStr := tenantIDVal.(string); okStr {
			parsed, err := uuid.Parse(tStr)
			if err == nil {
				return parsed, true
			}
		}
		response.Error(c, http.StatusUnauthorized, "Identificador de tenant inválido")
		return uuid.Nil, false
	}
	return tenantID, true
}

// GetStatus retorna o status da conexão do WhatsApp
func (h *WhatsAppHandler) GetStatus(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	status, err := h.waService.GetStatus(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, status, "Status do WhatsApp consultado com sucesso")
}

// EnsureInstance garante a instância criada no WUZAPI
func (h *WhatsAppHandler) EnsureInstance(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	cfg, err := h.waService.EnsureInstance(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"instance_name": cfg.InstanceName,
		"status":        cfg.Status,
	}, "Instância inicializada com sucesso")
}

// Connect inicia a conexão da sessão e geração de QR code
func (h *WhatsAppHandler) Connect(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	qrResp, err := h.waService.Connect(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, qrResp, "Conexão solicitada com sucesso")
}

// GetQRCode retorna o QR Code atual para leitura
func (h *WhatsAppHandler) GetQRCode(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	qrResp, err := h.waService.GetQRCode(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, qrResp, "QR Code obtido com sucesso")
}

// Disconnect desconecta a sessão
func (h *WhatsAppHandler) Disconnect(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	if err := h.waService.Disconnect(c.Request.Context(), tenantID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil, "WhatsApp desconectado com sucesso")
}

// Logout desconecta e apaga a sessão no WUZAPI
func (h *WhatsAppHandler) Logout(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	if err := h.waService.Logout(c.Request.Context(), tenantID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil, "Sessão do WhatsApp encerrada com sucesso")
}

// GetAIConfig retorna a configuração de IA do estabelecimento
func (h *WhatsAppHandler) GetAIConfig(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	aiConfig, err := h.waService.GetAIConfig(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, aiConfig, "Configuração de IA recuperada com sucesso")
}

// UpdateAIConfig atualiza parâmetros e API keys de IA
func (h *WhatsAppHandler) UpdateAIConfig(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	var dto domain.UpdateAIConfigDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, "Dados de configuração inválidos: "+err.Error())
		return
	}

	updated, err := h.waService.UpdateAIConfig(c.Request.Context(), tenantID, dto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, updated, "Configurações de IA salvas com sucesso")
}
