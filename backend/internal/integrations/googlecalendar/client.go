package googlecalendar

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sistema-agendamento/backend/internal/domain"
)

const (
	defaultAuthBaseURL        = "https://accounts.google.com/o/oauth2/v2/auth"
	defaultOAuthTokenURL      = "https://oauth2.googleapis.com/token"
	defaultOAuthRevokeURL     = "https://oauth2.googleapis.com/revoke"
	defaultUserInfoURL        = "https://www.googleapis.com/oauth2/v2/userinfo"
	defaultCalendarAPIBaseURL = "https://www.googleapis.com/calendar/v3"
)

var (
	ErrEventConflict     = errors.New("evento já existe no Google Calendar (409)")
	ErrEventNotFound     = errors.New("evento não encontrado no Google Calendar (404/410)")
	ErrTokenUnauthorized = errors.New("token de acesso Google Calendar recusado (401)")
)

// Endpoints permite sobrescrever URLs base em testes unitários via httptest
type Endpoints struct {
	AuthBaseURL        string
	OAuthTokenURL      string
	OAuthRevokeURL     string
	UserInfoURL        string
	CalendarAPIBaseURL string
}

type Client struct {
	authBaseURL        string
	oauthTokenURL      string
	oauthRevokeURL     string
	userInfoURL        string
	calendarAPIBaseURL string
	httpClient         *http.Client
}

var _ domain.CalendarProvider = (*Client)(nil)

func NewClient() *Client {
	return NewClientWithEndpoints(Endpoints{})
}

func NewClientWithEndpoints(ep Endpoints) *Client {
	authURL := strings.TrimRight(ep.AuthBaseURL, "/")
	if authURL == "" {
		authURL = defaultAuthBaseURL
	}
	tokenURL := strings.TrimRight(ep.OAuthTokenURL, "/")
	if tokenURL == "" {
		tokenURL = defaultOAuthTokenURL
	}
	revokeURL := strings.TrimRight(ep.OAuthRevokeURL, "/")
	if revokeURL == "" {
		revokeURL = defaultOAuthRevokeURL
	}
	userInfoURL := strings.TrimRight(ep.UserInfoURL, "/")
	if userInfoURL == "" {
		userInfoURL = defaultUserInfoURL
	}
	calBaseURL := strings.TrimRight(ep.CalendarAPIBaseURL, "/")
	if calBaseURL == "" {
		calBaseURL = defaultCalendarAPIBaseURL
	}

	return &Client{
		authBaseURL:        authURL,
		oauthTokenURL:      tokenURL,
		oauthRevokeURL:     revokeURL,
		userInfoURL:        userInfoURL,
		calendarAPIBaseURL: calBaseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type oauthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type googleAPIErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type googleEventDateTime struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone,omitempty"`
}

type googleExtendedProperties struct {
	Private map[string]string `json:"private,omitempty"`
}

type googleEventPayload struct {
	ID                 string                    `json:"id,omitempty"`
	Summary            string                    `json:"summary"`
	Description        string                    `json:"description,omitempty"`
	Start              googleEventDateTime       `json:"start"`
	End                googleEventDateTime       `json:"end"`
	ExtendedProperties *googleExtendedProperties `json:"extendedProperties,omitempty"`
}

type googleEventResponse struct {
	ID string `json:"id"`
}

// BuildAuthURL gera a URL de consentimento OAuth 2.0 com escopos mínimos do Google Calendar
func (c *Client) BuildAuthURL(clientID, redirectURI, state string) string {
	scopes := strings.Join([]string{
		"https://www.googleapis.com/auth/calendar.events",
		"https://www.googleapis.com/auth/calendar.calendarlist.readonly",
		"https://www.googleapis.com/auth/userinfo.email",
	}, " ")

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", scopes)
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")
	params.Set("include_granted_scopes", "true")
	params.Set("state", state)

	return c.authBaseURL + "?" + params.Encode()
}

// ExchangeCode troca o authorization code por access_token e refresh_token
func (c *Client) ExchangeCode(ctx context.Context, clientID, clientSecret, redirectURI, code string) (*domain.OAuthTokenResponse, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")

	return c.postTokenForm(ctx, form)
}

// RefreshToken renova o access_token utilizando o refresh_token salvo
func (c *Client) RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*domain.OAuthTokenResponse, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("grant_type", "refresh_token")

	return c.postTokenForm(ctx, form)
}

func (c *Client) postTokenForm(ctx context.Context, form url.Values) (*domain.OAuthTokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.oauthTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição OAuth token: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao comunicar com OAuth Google: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta OAuth Google: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var oauthErr oauthErrorResponse
		if jsonErr := json.Unmarshal(body, &oauthErr); jsonErr == nil && oauthErr.Error != "" {
			if oauthErr.Error == "invalid_grant" {
				return nil, fmt.Errorf("%w: %s", domain.ErrGoogleCalendarAuthExpired, oauthErr.ErrorDescription)
			}
			return nil, fmt.Errorf("erro OAuth Google (%d - %s): %s", resp.StatusCode, oauthErr.Error, oauthErr.ErrorDescription)
		}
		if strings.Contains(string(body), "invalid_grant") {
			return nil, fmt.Errorf("%w: %s", domain.ErrGoogleCalendarAuthExpired, string(body))
		}
		return nil, fmt.Errorf("erro OAuth Google (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp domain.OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta OAuth Google: %w", err)
	}
	return &tokenResp, nil
}

// RevokeToken revoga o token junto ao Google OAuth (usado no Disconnect)
func (c *Client) RevokeToken(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	form := url.Values{}
	form.Set("token", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.oauthRevokeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("falha ao criar requisição de revogação Google: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao comunicar com revogação Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao revogar token Google (HTTP %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// GetUserEmail retorna o e-mail da conta Google autenticada
func (c *Client) GetUserEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.userInfoURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha ao buscar perfil Google: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("%w: %s", ErrTokenUnauthorized, string(body))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("erro ao buscar e-mail Google (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var info struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", fmt.Errorf("falha ao decodificar userinfo Google: %w", err)
	}
	return info.Email, nil
}

// ListCalendars lista os calendários com permissão de escrita (owner ou writer)
func (c *Client) ListCalendars(ctx context.Context, accessToken string) ([]domain.GoogleCalendarItemDTO, error) {
	var result []domain.GoogleCalendarItemDTO
	pageToken := ""

	for {
		endpoint := c.calendarAPIBaseURL + "/users/me/calendarList"
		if pageToken != "" {
			endpoint += "?pageToken=" + url.QueryEscape(pageToken)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("falha ao listar calendários Google: %w", err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("%w: %s", ErrTokenUnauthorized, string(body))
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, parseGoogleAPIError("listar calendários", resp.StatusCode, body)
		}

		var listResp struct {
			NextPageToken string `json:"nextPageToken"`
			Items         []struct {
				ID         string `json:"id"`
				Summary    string `json:"summary"`
				Primary    bool   `json:"primary"`
				AccessRole string `json:"accessRole"`
			} `json:"items"`
		}
		if err := json.Unmarshal(body, &listResp); err != nil {
			return nil, fmt.Errorf("falha ao decodificar lista de calendários Google: %w", err)
		}

		for _, item := range listResp.Items {
			if item.AccessRole == "owner" || item.AccessRole == "writer" {
				result = append(result, domain.GoogleCalendarItemDTO{
					ID:         item.ID,
					Summary:    item.Summary,
					Primary:    item.Primary,
					AccessRole: item.AccessRole,
				})
			}
		}

		if listResp.NextPageToken == "" {
			break
		}
		pageToken = listResp.NextPageToken
	}

	return result, nil
}

// CreateEvent cria um evento no Google Calendar com sendUpdates=none
func (c *Client) CreateEvent(ctx context.Context, accessToken, calendarID string, event domain.CalendarEventInput) (string, error) {
	if calendarID == "" {
		calendarID = "primary"
	}

	payload := buildGoogleEventPayload(event, true)
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/calendars/%s/events?sendUpdates=none", c.calendarAPIBaseURL, url.PathEscape(calendarID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha ao criar evento no Google Calendar: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusConflict {
		return "", fmt.Errorf("%w: %s", ErrEventConflict, string(body))
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("%w: %s", ErrTokenUnauthorized, string(body))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", parseGoogleAPIError("criar evento", resp.StatusCode, body)
	}

	var created googleEventResponse
	if err := json.Unmarshal(body, &created); err != nil {
		return "", fmt.Errorf("falha ao decodificar evento criado: %w", err)
	}
	if created.ID == "" {
		created.ID = event.EventID
	}
	return created.ID, nil
}

// UpdateEvent atualiza um evento existente via PATCH com sendUpdates=none
func (c *Client) UpdateEvent(ctx context.Context, accessToken, calendarID, eventID string, event domain.CalendarEventInput) (string, error) {
	if calendarID == "" {
		calendarID = "primary"
	}

	payload := buildGoogleEventPayload(event, false)
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/calendars/%s/events/%s?sendUpdates=none",
		c.calendarAPIBaseURL, url.PathEscape(calendarID), url.PathEscape(eventID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha ao atualizar evento no Google Calendar: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return "", fmt.Errorf("%w (HTTP %d): %s", ErrEventNotFound, resp.StatusCode, string(body))
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("%w: %s", ErrTokenUnauthorized, string(body))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", parseGoogleAPIError("atualizar evento", resp.StatusCode, body)
	}

	var updated googleEventResponse
	if err := json.Unmarshal(body, &updated); err != nil {
		return "", fmt.Errorf("falha ao decodificar evento atualizado: %w", err)
	}
	if updated.ID == "" {
		updated.ID = eventID
	}
	return updated.ID, nil
}

// DeleteEvent remove o evento no Google Calendar (404/410 tratados como sucesso idempotente)
func (c *Client) DeleteEvent(ctx context.Context, accessToken, calendarID, eventID string) error {
	if calendarID == "" {
		calendarID = "primary"
	}
	if eventID == "" {
		return nil
	}

	endpoint := fmt.Sprintf("%s/calendars/%s/events/%s?sendUpdates=none",
		c.calendarAPIBaseURL, url.PathEscape(calendarID), url.PathEscape(eventID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao remover evento no Google Calendar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone ||
		resp.StatusCode == http.StatusNoContent || (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("%w: %s", ErrTokenUnauthorized, string(body))
	}
	return parseGoogleAPIError("remover evento", resp.StatusCode, body)
}

func buildGoogleEventPayload(event domain.CalendarEventInput, includeID bool) googleEventPayload {
	payload := googleEventPayload{
		Summary:     event.Summary,
		Description: event.Description,
		Start: googleEventDateTime{
			DateTime: event.StartAt.Format(time.RFC3339),
			TimeZone: event.TimeZone,
		},
		End: googleEventDateTime{
			DateTime: event.EndAt.Format(time.RFC3339),
			TimeZone: event.TimeZone,
		},
	}
	if includeID && event.EventID != "" {
		payload.ID = event.EventID
	}
	if len(event.PrivateProps) > 0 {
		payload.ExtendedProperties = &googleExtendedProperties{
			Private: event.PrivateProps,
		}
	}
	return payload
}

func parseGoogleAPIError(op string, statusCode int, body []byte) error {
	var apiErr googleAPIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error.Message != "" {
		return fmt.Errorf("erro Google Calendar ao %s (HTTP %d): %s", op, statusCode, apiErr.Error.Message)
	}
	return fmt.Errorf("erro Google Calendar ao %s (HTTP %d): %s", op, statusCode, string(body))
}
