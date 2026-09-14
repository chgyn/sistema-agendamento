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

type GeminiProvider struct {
	httpClient *http.Client
}

func NewGeminiProvider() *GeminiProvider {
	return &GeminiProvider{
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// Structs internas da API v1beta do Google Gemini
type geminiPart struct {
	Text             string                 `json:"text,omitempty"`
	FunctionCall     *geminiFunctionCall    `json:"functionCall,omitempty"`
	FunctionResponse *geminiFunctionResponse `json:"functionResponse,omitempty"`
}

type geminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type geminiFunctionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

type geminiContent struct {
	Role  string       `json:"role"` // "user" ou "model"
	Parts []geminiPart `json:"parts"`
}

type geminiSystemInstruction struct {
	Parts []geminiPart `json:"parts"`
}

type geminiFunctionDeclaration struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  FunctionParameters `json:"parameters"`
}

type geminiTool struct {
	FunctionDeclarations []geminiFunctionDeclaration `json:"functionDeclarations"`
}

type geminiGenerateContentRequest struct {
	SystemInstruction *geminiSystemInstruction `json:"systemInstruction,omitempty"`
	Contents          []geminiContent          `json:"contents"`
	Tools             []geminiTool             `json:"tools,omitempty"`
}

type geminiCandidate struct {
	Content struct {
		Parts []geminiPart `json:"parts"`
		Role  string       `json:"role"`
	} `json:"content"`
	FinishReason string `json:"finishReason"`
}

type geminiGenerateContentResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func (g *GeminiProvider) Chat(ctx context.Context, req AIRequest) (*AIResponse, error) {
	model := req.Model
	if model == "" {
		model = "gemini-2.5-flash"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, req.APIKey)

	geminiReq := geminiGenerateContentRequest{
		Contents: make([]geminiContent, 0, len(req.Messages)),
	}

	if req.SystemPrompt != "" {
		geminiReq.SystemInstruction = &geminiSystemInstruction{
			Parts: []geminiPart{{Text: req.SystemPrompt}},
		}
	}

	// Converte tools
	if len(req.Tools) > 0 {
		declarations := make([]geminiFunctionDeclaration, 0, len(req.Tools))
		for _, t := range req.Tools {
			declarations = append(declarations, geminiFunctionDeclaration{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			})
		}
		geminiReq.Tools = []geminiTool{{FunctionDeclarations: declarations}}
	}

	// Converte mensagens
	for _, m := range req.Messages {
		switch m.Role {
		case RoleUser:
			geminiReq.Contents = append(geminiReq.Contents, geminiContent{
				Role:  "user",
				Parts: []geminiPart{{Text: m.Content}},
			})
		case RoleAssistant:
			parts := make([]geminiPart, 0)
			if m.Content != "" {
				parts = append(parts, geminiPart{Text: m.Content})
			}
			for _, tc := range m.ToolCalls {
				parts = append(parts, geminiPart{
					FunctionCall: &geminiFunctionCall{
						Name: tc.Name,
						Args: tc.Arguments,
					},
				})
			}
			geminiReq.Contents = append(geminiReq.Contents, geminiContent{
				Role:  "model",
				Parts: parts,
			})
		case RoleTool:
			// Retorno de execução de tool para o Gemini
			var resObj map[string]interface{}
			if err := json.Unmarshal([]byte(m.Content), &resObj); err != nil {
				resObj = map[string]interface{}{"output": m.Content}
			}
			geminiReq.Contents = append(geminiReq.Contents, geminiContent{
				Role: "user",
				Parts: []geminiPart{
					{
						FunctionResponse: &geminiFunctionResponse{
							Name:     m.ToolName,
							Response: resObj,
						},
					},
				},
			})
		}
	}

	bodyBytes, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar payload do Gemini: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição HTTP para a API do Gemini: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler resposta do Gemini: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API retornou erro HTTP %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var geminiResp geminiGenerateContentResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("falha ao decodificar JSON do Gemini: %w", err)
	}

	if geminiResp.Error != nil {
		return nil, fmt.Errorf("erro retornado pelo Gemini: %s (status: %s)", geminiResp.Error.Message, geminiResp.Error.Status)
	}

	if len(geminiResp.Candidates) == 0 {
		return &AIResponse{Content: "Desculpe, não foi possível gerar uma resposta no momento."}, nil
	}

	aiResp := &AIResponse{}
	candidate := geminiResp.Candidates[0]

	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			if aiResp.Content != "" {
				aiResp.Content += "\n"
			}
			aiResp.Content += part.Text
		}
		if part.FunctionCall != nil {
			aiResp.ToolCalls = append(aiResp.ToolCalls, ToolCall{
				ID:        fmt.Sprintf("call_%s_%d", part.FunctionCall.Name, time.Now().UnixNano()),
				Name:      part.FunctionCall.Name,
				Arguments: part.FunctionCall.Args,
			})
		}
	}

	return aiResp, nil
}
