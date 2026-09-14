package ai

import (
	"context"
)

// ChatRole papéis suportados nas conversas
type ChatRole string

const (
	RoleSystem    ChatRole = "system"
	RoleUser      ChatRole = "user"
	RoleAssistant ChatRole = "assistant"
	RoleTool      ChatRole = "tool"
)

// ParameterProperty definição de propriedade JSON Schema para Tools
type ParameterProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// FunctionParameters esquema de parâmetros de uma ferramenta
type FunctionParameters struct {
	Type       string                       `json:"type"`
	Properties map[string]ParameterProperty `json:"properties"`
	Required   []string                     `json:"required,omitempty"`
}

// FunctionDefinition especificação de uma função invocável pela IA
type FunctionDefinition struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  FunctionParameters `json:"parameters"`
}

// Tool container para chamada de função
type Tool struct {
	Type     string             `json:"type"` // "function"
	Function FunctionDefinition `json:"function"`
}

// ToolCall chamada gerada pelo modelo de IA
type ToolCall struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ChatMessage mensagem no histórico da conversa
type ChatMessage struct {
	Role       ChatRole   `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolName   string     `json:"tool_name,omitempty"`
}

// AIRequest requisição padronizada para qualquer provedor de IA
type AIRequest struct {
	Model        string        `json:"model"`
	APIKey       string        `json:"api_key"`
	SystemPrompt string        `json:"system_prompt"`
	Messages     []ChatMessage `json:"messages"`
	Tools        []Tool        `json:"tools,omitempty"`
}

// AIResponse resposta gerada pelo provedor de IA
type AIResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// AIProvider contrato comum para Google Gemini e OpenAI
type AIProvider interface {
	Chat(ctx context.Context, req AIRequest) (*AIResponse, error)
}
