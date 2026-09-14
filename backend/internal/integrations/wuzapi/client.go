package wuzapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInstanceNotFound = errors.New("instância do WUZAPI não encontrada")
	ErrWuzapiOffline    = errors.New("serviço WUZAPI indisponível ou inacessível")
)

type Client struct {
	baseURL    string
	adminToken string
	httpClient *http.Client
}

func NewClient(baseURL, adminToken string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL:    baseURL,
		adminToken: adminToken,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// User representa uma instância/usuário cadastrada no WUZAPI
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Token     string `json:"token"`
	Webhook   string `json:"webhook"`
	JID       string `json:"jid"`
	QRCode    string `json:"qrcode"`
	Connected bool   `json:"connected"`
	Events    string `json:"events"`
}

// SessionStatusResponse retorno de /session/status
type SessionStatusResponse struct {
	Code    int `json:"code"`
	Success bool `json:"success"`
	Data    struct {
		Connected bool `json:"Connected"`
		LoggedIn  bool `json:"LoggedIn"`
	} `json:"data"`
}

// QRCodeResponse retorno de /session/qr
type QRCodeResponse struct {
	Code    int `json:"code"`
	Success bool `json:"success"`
	Data    struct {
		QRCode string `json:"QRCode"`
	} `json:"data"`
}

// CreateUserRequest payload para cadastrar instância no WUZAPI
type CreateUserRequest struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	Webhook string `json:"webhook"`
	Events  string `json:"events"`
}

// PresenceRequest payload para definir status digitando/gravando
type PresenceRequest struct {
	Phone string `json:"Phone"`
	State string `json:"State"` // "composing" ou "paused"
	Media string `json:"Media"` // vazio para texto, "audio" para gravação de voz
}

// SendTextMessageRequest payload para envio de texto
type SendTextMessageRequest struct {
	Phone string `json:"Phone"`
	Body  string `json:"Body"`
}

// EnsureUser verifica se a instância já existe no WUZAPI; se não, cadastra com webhook configurado
func (c *Client) EnsureUser(ctx context.Context, name, token, webhookURL string) error {
	users, err := c.ListUsers(ctx)
	if err == nil {
		for _, u := range users {
			if u.Name == name || u.Token == token {
				return nil // Instância já existe
			}
		}
	}

	reqBody := CreateUserRequest{
		Name:    name,
		Token:   token,
		Webhook: webhookURL,
		Events:  "Message,ReadReceipt,ChatPresence",
	}
	bodyData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/admin/users", c.baseURL), bytes.NewBuffer(bodyData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.adminToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao conectar ao WUZAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("falha ao criar usuário no WUZAPI (status %d): %s", resp.StatusCode, string(respBytes))
	}

	return nil
}

// ListUsers lista todos os usuários cadastrados no WUZAPI
func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/admin/users", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.adminToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao WUZAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WUZAPI retornou status %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// Connect inicia a conexão da sessão WhatsApp para a instância
func (c *Client) Connect(ctx context.Context, token string) error {
	payload := map[string]interface{}{
		"Subscribe": []string{"Message", "ReadReceipt", "ChatPresence"},
		"Immediate": false,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/session/connect", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao iniciar conexão da sessão no WUZAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WUZAPI /session/connect falhou (status %d): %s", resp.StatusCode, string(respBytes))
	}

	return nil
}

// GetStatus consulta o estado da sessão WhatsApp
func (c *Client) GetStatus(ctx context.Context, token string) (*SessionStatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/session/status", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter status da sessão: %w", err)
	}
	defer resp.Body.Close()

	var result SessionStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetQRCode recupera o QR Code base64 gerado pelo WUZAPI
func (c *Client) GetQRCode(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/session/qr", c.baseURL), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha ao obter QR code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("WUZAPI /session/qr falhou (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var result QRCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.QRCode, nil
}

// Disconnect pausa a conexão do WhatsApp sem invalidar a sessão salva
func (c *Client) Disconnect(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/session/disconnect", c.baseURL), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao desconectar: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// Logout encerra a sessão WhatsApp no aparelho e remove o pareamento
func (c *Client) Logout(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/session/logout", c.baseURL), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao realizar logout: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// SetPresence envia sinalização de "digitando..." (composing) ou "pausado" (paused)
func (c *Client) SetPresence(ctx context.Context, token, phone, state string) error {
	body := PresenceRequest{
		Phone: phone,
		State: state,
		Media: "",
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/chat/presence", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao definir presence: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// SendTextMessage envia uma mensagem de texto pelo WhatsApp
func (c *Client) SendTextMessage(ctx context.Context, token, phone, body string) error {
	reqPayload := SendTextMessageRequest{
		Phone: phone,
		Body:  body,
	}
	data, _ := json.Marshal(reqPayload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/chat/send/text", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("falha ao enviar mensagem de texto: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WUZAPI /chat/send/text falhou (status %d): %s", resp.StatusCode, string(respBytes))
	}

	return nil
}
