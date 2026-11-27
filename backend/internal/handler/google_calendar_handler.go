package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type GoogleCalendarHandler struct {
	svc         *service.GoogleCalendarService
	frontendURL string
}

func NewGoogleCalendarHandler(svc *service.GoogleCalendarService, frontendURL string) *GoogleCalendarHandler {
	return &GoogleCalendarHandler{
		svc:         svc,
		frontendURL: strings.TrimRight(frontendURL, "/"),
	}
}

func (h *GoogleCalendarHandler) buildFrontendRedirect(status, reason string) string {
	base := h.frontendURL + "/admin/google-calendar"
	q := url.Values{}
	q.Set("status", status)
	if reason != "" {
		q.Set("reason", reason)
	}
	return base + "?" + q.Encode()
}

// Callback recebe o redirecionamento OAuth 2.0 do Google e redireciona para o painel frontend
func (h *GoogleCalendarHandler) Callback(c *gin.Context) {
	if oauthErr := strings.TrimSpace(c.Query("error")); oauthErr != "" {
		reason := "oauth_denied"
		if oauthErr == "access_denied" {
			reason = "access_denied"
		}
		c.Redirect(http.StatusFound, h.buildFrontendRedirect("error", reason))
		return
	}

	code := c.Query("code")
	state := c.Query("state")
	if _, err := h.svc.HandleCallback(c.Request.Context(), code, state); err != nil {
		reason := "oauth_exchange_failed"
		switch {
		case errors.Is(err, domain.ErrInvalidOAuthState):
			reason = "invalid_state"
		case errors.Is(err, domain.ErrGoogleCalendarNotConfigured):
			reason = "not_configured"
		}
		c.Redirect(http.StatusFound, h.buildFrontendRedirect("error", reason))
		return
	}

	c.Redirect(http.StatusFound, h.buildFrontendRedirect("connected", ""))
}

// GetStatus retorna o status e configurações (sem segredos) do Google Calendar para o tenant
func (h *GoogleCalendarHandler) GetStatus(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	status, err := h.svc.GetStatus(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, status, "Status do Google Calendar recuperado com sucesso")
}

// GetAuthURL gera a URL de consentimento OAuth 2.0 para conectar a conta Google do tenant
func (h *GoogleCalendarHandler) GetAuthURL(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	userID, _ := middleware.GetUserID(c)
	authURL, err := h.svc.GetAuthURL(c.Request.Context(), tenantID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrGoogleCalendarNotConfigured) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"auth_url": authURL}, "URL de autorização gerada com sucesso")
}

// ListCalendars lista os calendários graváveis da conta Google conectada
func (h *GoogleCalendarHandler) ListCalendars(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	calendars, err := h.svc.ListCalendars(c.Request.Context(), tenantID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrGoogleCalendarNotConfigured),
			errors.Is(err, domain.ErrGoogleCalendarNotConnected),
			errors.Is(err, domain.ErrGoogleCalendarAuthExpired):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(c, calendars, "Calendários recuperados com sucesso")
}

// UpdateSettings atualiza credenciais BYOK, calendário selecionado e toggle de sincronização
func (h *GoogleCalendarHandler) UpdateSettings(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	var dto domain.UpdateGoogleCalendarSettingsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, "Dados de configuração inválidos: "+err.Error())
		return
	}

	updated, err := h.svc.UpdateSettings(c.Request.Context(), tenantID, dto)
	if err != nil {
		if errors.Is(err, domain.ErrGoogleCalendarNotConnected) || errors.Is(err, domain.ErrGoogleCalendarNotConfigured) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, updated, "Configurações do Google Calendar salvas com sucesso")
}

// Disconnect desconecta a conta Google mantendo as credenciais OAuth BYOK configuradas
func (h *GoogleCalendarHandler) Disconnect(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	if err := h.svc.Disconnect(c.Request.Context(), tenantID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil, "Conta Google Calendar desconectada com sucesso")
}

// ListFailures lista agendamentos com falha de sincronização no tenant
func (h *GoogleCalendarHandler) ListFailures(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	failures, err := h.svc.ListFailures(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, failures, "Falhas de sincronização recuperadas com sucesso")
}

// RetrySync reenfileira manualmente a sincronização de um agendamento do tenant
func (h *GoogleCalendarHandler) RetrySync(c *gin.Context) {
	tenantID, ok := resolveTenantID(c)
	if !ok {
		return
	}

	aptID, err := uuid.Parse(c.Param("appointment_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID de agendamento inválido")
		return
	}

	if err := h.svc.RetryAppointmentSync(c.Request.Context(), tenantID, aptID); err != nil {
		switch {
		case errors.Is(err, domain.ErrAppointmentNotFound):
			response.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrGoogleCalendarNotConnected),
			errors.Is(err, domain.ErrGoogleCalendarAuthExpired):
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(c, nil, "Sincronização reenfileirada com sucesso")
}
