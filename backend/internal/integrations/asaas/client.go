package asaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sistema-agendamento/backend/internal/domain"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// AsaasErrorResponse estrutura padrão de erros da API Asaas
type AsaasErrorResponse struct {
	Errors []struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"errors"`
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("falha ao serializar payload Asaas: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição HTTP para Asaas: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("access_token", c.apiKey)
	req.Header.Set("User-Agent", "SistemaAgendamento/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na comunicação com Asaas (%s %s): %w", method, url, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta do Asaas: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var asaasErr AsaasErrorResponse
		if jsonErr := json.Unmarshal(respBytes, &asaasErr); jsonErr == nil && len(asaasErr.Errors) > 0 {
			return nil, fmt.Errorf("erro da API Asaas (%d): %s (código: %s)", resp.StatusCode, asaasErr.Errors[0].Description, asaasErr.Errors[0].Code)
		}
		return nil, fmt.Errorf("erro da API Asaas (HTTP %d): %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

// CreateCustomer cria um cliente no Asaas
func (c *Client) CreateCustomer(ctx context.Context, req domain.AsaasCustomerRequest) (*domain.AsaasCustomerResponse, error) {
	respBytes, err := c.doRequest(ctx, http.MethodPost, "/customers", req)
	if err != nil {
		return nil, err
	}

	var res domain.AsaasCustomerResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("falha ao desserializar resposta de criação de cliente Asaas: %w", err)
	}

	return &res, nil
}

// CreateSubscription cria uma assinatura recorrente no Asaas
func (c *Client) CreateSubscription(ctx context.Context, req domain.AsaasSubscriptionRequest) (*domain.AsaasSubscriptionResponse, error) {
	respBytes, err := c.doRequest(ctx, http.MethodPost, "/subscriptions", req)
	if err != nil {
		return nil, err
	}

	var res domain.AsaasSubscriptionResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("falha ao desserializar resposta de criação de assinatura Asaas: %w", err)
	}

	return &res, nil
}

// GetSubscription obtém dados de uma assinatura no Asaas
func (c *Client) GetSubscription(ctx context.Context, asaasSubID string) (*domain.AsaasSubscriptionResponse, error) {
	path := fmt.Sprintf("/subscriptions/%s", asaasSubID)
	respBytes, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var res domain.AsaasSubscriptionResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("falha ao desserializar assinatura Asaas: %w", err)
	}

	return &res, nil
}

// CancelSubscription cancela uma assinatura no Asaas
func (c *Client) CancelSubscription(ctx context.Context, asaasSubID string) error {
	path := fmt.Sprintf("/subscriptions/%s", asaasSubID)
	_, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	return err
}

// GetSubscriptionPayments lista as faturas/cobranças geradas para uma assinatura
func (c *Client) GetSubscriptionPayments(ctx context.Context, asaasSubID string) ([]domain.AsaasPaymentResponse, error) {
	path := fmt.Sprintf("/subscriptions/%s/payments", asaasSubID)
	respBytes, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var listRes domain.AsaasPaymentListResponse
	if err := json.Unmarshal(respBytes, &listRes); err != nil {
		return nil, fmt.Errorf("falha ao desserializar faturas da assinatura Asaas: %w", err)
	}

	return listRes.Data, nil
}
