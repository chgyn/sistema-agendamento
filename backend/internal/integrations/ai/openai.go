package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIProvider struct {
	httpClient *http.Client
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// Structs internas da API OpenAI Chat Completions
type openAIMessage struct {
	Role       string            `json:"role"`
	Content    string            `json:"content,omitempty"`
	ToolCalls  []openAIToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string            `json:"tool_call_id,omitempty"`
	Name       string            `json:"name,omitempty"`
}

type openAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"` // "function"
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Tools       []Tool          `json:"tools,omitempty"`
	Temperature float64         `json:"temperature"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Role      string           `json:"role"`
			Content   string           `json:"content"`
			ToolCalls []openAIToolCall `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

func (o *OpenAIProvider) Chat(ctx context.Context, req AIRequest) (*AIResponse, error) {
	model := req.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	messages := make([]openAIMessage, 0, len(req.Messages)+1)

	if req.SystemPrompt != "" {
		messages = append(messages, openAIMessage{
			Role:    "system",
			Content: req.SystemPrompt,
		})
	}

	for _, m := range req.Messages {
		switch m.Role {
		case RoleUser:
			messages = append(messages, openAIMessage{
				Role:    "user",
				Content: m.Content,
			})
		case RoleAssistant:
			msg := openAIMessage{
				Role:    "assistant",
				Content: m.Content,
			}
			if len(m.ToolCalls) > 0 {
				msg.ToolCalls = make([]openAIToolCall, 0, len(m.ToolCalls))
				for _, tc := range m.ToolCalls {
					argsJSON, _ := json.Marshal(tc.Arguments)
					tcOpenAI := openAIToolCall{
						ID:   tc.ID,
						Type: "function",
					}
					tcOpenAI.Function.Name = tc.Name
					tcOpenAI.Function.Arguments = string(argsJSON)
					msg.ToolCalls = append(msg.ToolCalls, tcOpenAI)
				}
			}
			messages = append(messages, msg)
		case RoleTool:
			messages = append(messages, openAIMessage{
				Role:       "tool",
				Content:    m.Content,
				ToolCallID: m.ToolCallID,
			})
		}
	}

	openAIReq := openAIChatRequest{
		Model:       model,
		Messages:    messages,
		Tools:       req.Tools,
		Temperature: 0.3,
	}

	bodyBytes, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar payload OpenAI: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))

	httpResp, err := o.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição HTTP para a API da OpenAI: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta da OpenAI: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API retornou erro HTTP %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var openAIResp openAIChatResponse
	if err := json.Unmarshal(respBytes, &openAIResp); err != nil {
		return nil, fmt.Errorf("falha ao decodificar JSON da OpenAI: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("erro retornado pela OpenAI: %s (código: %s)", openAIResp.Error.Message, openAIResp.Error.Code)
	}

	if len(openAIResp.Choices) == 0 {
		return &AIResponse{Content: "Desculpe, não foi possível gerar uma resposta no momento."}, nil
	}

	choice := openAIResp.Choices[0]
	aiResp := &AIResponse{
		Content: choice.Message.Content,
	}

	if len(choice.Message.ToolCalls) > 0 {
		aiResp.ToolCalls = make([]ToolCall, 0, len(choice.Message.ToolCalls))
		for _, tc := range choice.Message.ToolCalls {
			var args map[string]interface{}
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
			aiResp.ToolCalls = append(aiResp.ToolCalls, ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: args,
			})
		}
	}

	return aiResp, nil
}
