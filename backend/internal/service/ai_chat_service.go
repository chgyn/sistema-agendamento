package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/crypto"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/integrations/ai"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type AIChatService struct {
	repo           *postgres.Repository
	availService   *AvailabilityService
	aptService     *AppointmentService
	geminiProvider *ai.GeminiProvider
	openaiProvider *ai.OpenAIProvider
	dispatcher     *WhatsAppMessageDispatcher
	encryptionKey  string
}

func NewAIChatService(
	repo *postgres.Repository,
	availService *AvailabilityService,
	aptService *AppointmentService,
	geminiProvider *ai.GeminiProvider,
	openaiProvider *ai.OpenAIProvider,
	dispatcher *WhatsAppMessageDispatcher,
	encryptionKey string,
) *AIChatService {
	return &AIChatService{
		repo:           repo,
		availService:   availService,
		aptService:     aptService,
		geminiProvider: geminiProvider,
		openaiProvider: openaiProvider,
		dispatcher:     dispatcher,
		encryptionKey:  encryptionKey,
	}
}

// ProcessCustomerMessage processa a mensagem recebida pelo WhatsApp utilizando o agente de IA com Tools da Agenda
func (s *AIChatService) ProcessCustomerMessage(ctx context.Context, tenantID uuid.UUID, customerPhone, customerName, incomingText, msgID string) error {
	// 1. Carrega configuração do WhatsApp e IA do Tenant
	cfg, err := s.repo.GetWhatsAppConfigByTenantID(ctx, tenantID)
	if err != nil || cfg == nil {
		log.Printf("⚠️ [AIChat] Tenant %s sem configuração de WhatsApp ativa", tenantID)
		return nil
	}

	if !cfg.IsAIEnabled {
		log.Printf("ℹ️ [AIChat] Atendimento por IA desativado para o tenant %s", tenantID)
		return nil
	}

	// 2. Decifra a API Key do provedor configurado
	var apiKey string
	if cfg.AIProvider == domain.AIProviderGemini {
		apiKey, err = crypto.Decrypt(cfg.GeminiAPIKeyEncrypted, s.encryptionKey)
	} else {
		apiKey, err = crypto.Decrypt(cfg.OpenAIAPIKeyEncrypted, s.encryptionKey)
	}

	if err != nil || apiKey == "" {
		log.Printf("⚠️ [AIChat] Tenant %s com IA ativada mas sem API Key válida decifrada: %v", tenantID, err)
		return nil
	}

	// 3. Obtém ou cria a sessão de conversa com o cliente
	conv, err := s.repo.GetOrCreateConversation(ctx, tenantID, customerPhone, customerName)
	if err != nil {
		return fmt.Errorf("falha ao obter/criar conversa: %w", err)
	}

	// 4. Salva a mensagem recebida do cliente no banco
	userMsg := &domain.WhatsAppMessage{
		ID:             uuid.New(),
		ConversationID: conv.ID,
		TenantID:       tenantID,
		Sender:         "customer",
		Role:           string(ai.RoleUser),
		Body:           incomingText,
		WuzapiMsgID:    msgID,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.SaveWhatsAppMessage(ctx, userMsg); err != nil {
		log.Printf("⚠️ [AIChat] Erro ao salvar mensagem do cliente no banco: %v", err)
	}

	// 5. Carrega dados do estabelecimento para contextualização
	tenant, err := s.repo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("estabelecimento não encontrado: %w", err)
	}

	// 6. Monta o histórico de mensagens recentes (memória multi-turn)
	recentMsgs, err := s.repo.GetRecentMessages(ctx, conv.ID, 12)
	if err != nil {
		log.Printf("⚠️ [AIChat] Erro ao consultar histórico recente: %v", err)
	}

	chatMessages := make([]ai.ChatMessage, 0, len(recentMsgs))
	for _, m := range recentMsgs {
		chatMessages = append(chatMessages, ai.ChatMessage{
			Role:    ai.ChatRole(m.Role),
			Content: m.Body,
		})
	}

	// 7. Constrói o System Prompt com regras rigorosas anti-alucinação
	now := time.Now()
	nowStr := now.Format("02/01/2006 às 15:04 (segunda-feira = Monday)")
	switch now.Weekday() {
	case time.Sunday:
		nowStr = now.Format("02/01/2006 às 15:04") + " (Domingo)"
	case time.Monday:
		nowStr = now.Format("02/01/2006 às 15:04") + " (Segunda-feira)"
	case time.Tuesday:
		nowStr = now.Format("02/01/2006 às 15:04") + " (Terça-feira)"
	case time.Wednesday:
		nowStr = now.Format("02/01/2006 às 15:04") + " (Quarta-feira)"
	case time.Thursday:
		nowStr = now.Format("02/01/2006 às 15:04") + " (Quinta-feira)"
	case time.Friday:
		nowStr = now.Format("02/01/2006 às 15:04") + " (Sexta-feira)"
	case time.Saturday:
		nowStr = now.Format("02/01/2006 às 15:04") + " (Sábado)"
	}

	systemPrompt := fmt.Sprintf(`Você é o assistente virtual oficial de atendimento da barbearia/salão "%s".
Seu objetivo é atender o cliente cordialmente, tirar dúvidas sobre serviços, verificar horários livres e realizar agendamentos diretamente pelo WhatsApp.

### DIRETRIZES FUNDAMENTAIS E REGRAS DE OURO:
1. MOMENTO ATUAL DE REFERÊNCIA: Agora é %s.
2. NUNCA INVENTE horários livres, profissionais ou serviços que não existam no sistema.
3. ANTES DE CONFIRMAR OU SUGERIR QUALQUER HORÁRIO AO CLIENTE, VOCÊ DEVE OBRIGATORIAMENTE EXECUTAR A FERRAMENTA "consultar_horarios_disponiveis" PARA A DATA DESEJADA.
4. Quando o cliente pedir para agendar (ex: "quero marcar para amanhã às 14h"), consulte os horários disponíveis primeiro. Se o horário estiver livre, confirme os dados com o cliente ou chame "agendar_horario".
5. Use linguagem natural, acolhedora, concisa e brasileira, ideal para mensagens de WhatsApp.
6. Nome do cliente registrado: %s (telefone: %s).
7. Endereço do estabelecimento: %s, %s - %s. Telefone de contato: %s.

### INSTRUÇÕES ESPECÍFICAS DO ESTABELECIMENTO:
%s
`,
		tenant.Name,
		nowStr,
		customerName,
		customerPhone,
		tenant.Address,
		tenant.City,
		tenant.State,
		tenant.Phone,
		cfg.SystemPromptCustom,
	)

	// 8. Declara as Ferramentas (Tools) disponíveis para a IA
	tools := s.getAITools()

	// 9. Loop de Execução de Tools (máximo de 5 iterações para evitar loops infinitos)
	var finalReply string
	maxTurns := 5

	for turn := 0; turn < maxTurns; turn++ {
		req := ai.AIRequest{
			Model:        cfg.AIModel,
			APIKey:       apiKey,
			SystemPrompt: systemPrompt,
			Messages:     chatMessages,
			Tools:        tools,
		}

		var resp *ai.AIResponse
		if cfg.AIProvider == domain.AIProviderGemini {
			resp, err = s.geminiProvider.Chat(ctx, req)
		} else {
			resp, err = s.openaiProvider.Chat(ctx, req)
		}

		if err != nil {
			log.Printf("❌ [AIChat] Erro na chamada à IA (%s): %v", cfg.AIProvider, err)
			return fmt.Errorf("erro na API da IA: %w", err)
		}

		if resp == nil {
			break
		}

		// Se o modelo retornou chamadas de ferramentas
		if len(resp.ToolCalls) > 0 {
			// Adiciona a resposta do assistente (com tool_calls) ao contexto
			chatMessages = append(chatMessages, ai.ChatMessage{
				Role:      ai.RoleAssistant,
				Content:   resp.Content,
				ToolCalls: resp.ToolCalls,
			})

			// Executa cada Tool chamada
			for _, tc := range resp.ToolCalls {
				log.Printf("🔧 [AIChat Tool Call] Executando %s com args: %v", tc.Name, tc.Arguments)
				toolResult := s.executeTool(ctx, tenantID, tc.Name, tc.Arguments, customerPhone, customerName)
				log.Printf("🔧 [AIChat Tool Result] %s: %s", tc.Name, toolResult)

				chatMessages = append(chatMessages, ai.ChatMessage{
					Role:       ai.RoleTool,
					Content:    toolResult,
					ToolCallID: tc.ID,
					ToolName:   tc.Name,
				})
			}
			continue
		}

		// Se não há tool calls, o modelo gerou a resposta textual final
		finalReply = resp.Content
		break
	}

	if finalReply == "" {
		finalReply = "Olá! Como posso te ajudar hoje com seus agendamentos?"
	}

	// 10. Salva a resposta do assistente no histórico
	assistantMsg := &domain.WhatsAppMessage{
		ID:             uuid.New(),
		ConversationID: conv.ID,
		TenantID:       tenantID,
		Sender:         "assistant",
		Role:           string(ai.RoleAssistant),
		Body:           finalReply,
		CreatedAt:      time.Now(),
	}
	_ = s.repo.SaveWhatsAppMessage(ctx, assistantMsg)

	// 11. Despacha através da Camada de Controle de Ritmo Humanizado
	dispatchParams := BuildDispatcherParams(cfg, customerPhone, finalReply)
	return s.dispatcher.SendHumanizedMessage(ctx, dispatchParams)
}

// executeTool executa a função solicitada pela IA
func (s *AIChatService) executeTool(ctx context.Context, tenantID uuid.UUID, toolName string, args map[string]interface{}, defaultPhone, defaultName string) string {
	switch toolName {
	case "listar_servicos":
		svcs, err := s.repo.ListServices(ctx, tenantID, true)
		if err != nil {
			return `{"erro": "Falha ao consultar serviços"}`
		}
		type SvcDTO struct {
			ID       string  `json:"id"`
			Nome     string  `json:"nome"`
			Preco    float64 `json:"preco"`
			Duracao  int     `json:"duracao_minutos"`
			Descricao string `json:"descricao"`
		}
		list := make([]SvcDTO, 0, len(svcs))
		for _, svc := range svcs {
			list = append(list, SvcDTO{
				ID:       svc.ID.String(),
				Nome:     svc.Name,
				Preco:    svc.Price,
				Duracao:  svc.DurationMinutes,
				Descricao: svc.Description,
			})
		}
		b, _ := json.Marshal(list)
		return string(b)

	case "listar_profissionais":
		var pros []domain.Professional
		var err error
		if svcIDStr, ok := args["servico_id"].(string); ok && svcIDStr != "" {
			if svcID, parseErr := uuid.Parse(svcIDStr); parseErr == nil {
				pros, err = s.repo.ListProfessionalsByService(ctx, tenantID, svcID)
			}
		}
		if len(pros) == 0 {
			pros, err = s.repo.ListProfessionals(ctx, tenantID, true)
		}
		if err != nil {
			return `{"erro": "Falha ao listar profissionais"}`
		}
		type ProDTO struct {
			ID          string `json:"id"`
			Nome        string `json:"nome"`
			Titulo      string `json:"titulo"`
			Especialidade string `json:"especialidade"`
		}
		list := make([]ProDTO, 0, len(pros))
		for _, p := range pros {
			list = append(list, ProDTO{
				ID:          p.ID.String(),
				Nome:        p.Name,
				Titulo:      p.Title,
				Especialidade: p.Specialty,
			})
		}
		b, _ := json.Marshal(list)
		return string(b)

	case "consultar_horarios_disponiveis":
		dateStr, _ := args["data"].(string)
		if dateStr == "" {
			return `{"erro": "A data (YYYY-MM-DD) é obrigatória"}`
		}
		var serviceID uuid.UUID
		if sid, ok := args["servico_id"].(string); ok {
			serviceID, _ = uuid.Parse(sid)
		}
		if serviceID == uuid.Nil {
			// Se o modelo não enviou o ID, tenta pegar o primeiro serviço ativo como padrão
			svcs, _ := s.repo.ListServices(ctx, tenantID, true)
			if len(svcs) > 0 {
				serviceID = svcs[0].ID
			}
		}

		var professionalID uuid.UUID
		if pid, ok := args["profissional_id"].(string); ok && pid != "" {
			professionalID, _ = uuid.Parse(pid)
		}

		avail, err := s.availService.GetAvailableSlotsForDate(ctx, tenantID, serviceID, professionalID, dateStr)
		if err != nil {
			return fmt.Sprintf(`{"erro": "Falha ao consultar disponibilidade: %s"}`, err.Error())
		}

		type SlotDTO struct {
			HorarioInicio string `json:"horario_inicio"`
			HorarioFim    string `json:"horario_fim"`
			Profissional  string `json:"profissional"`
		}
		freeSlots := make([]SlotDTO, 0)
		for _, s := range avail.Slots {
			if s.IsAvailable {
				freeSlots = append(freeSlots, SlotDTO{
					HorarioInicio: s.StartTime,
					HorarioFim:    s.EndTime,
					Profissional:  s.ProfessionalName,
				})
			}
		}

		res := map[string]interface{}{
			"data":              avail.Date,
			"tem_horarios":      len(freeSlots) > 0,
			"total_disponiveis": len(freeSlots),
			"horarios_livres":   freeSlots,
		}
		b, _ := json.Marshal(res)
		return string(b)

	case "agendar_horario":
		svcIDStr, _ := args["servico_id"].(string)
		proIDStr, _ := args["profissional_id"].(string)
		dataHoraStr, _ := args["data_hora"].(string)
		nomeCliente, _ := args["nome_cliente"].(string)
		telefoneCliente, _ := args["telefone_cliente"].(string)
		notes, _ := args["observacoes"].(string)

		if nomeCliente == "" {
			nomeCliente = defaultName
		}
		if telefoneCliente == "" {
			telefoneCliente = defaultPhone
		}

		serviceID, err := uuid.Parse(svcIDStr)
		if err != nil {
			return `{"erro": "ID de serviço inválido"}`
		}
		professionalID, err := uuid.Parse(proIDStr)
		if err != nil {
			return `{"erro": "ID de profissional inválido"}`
		}

		// Tenta múltiplos formatos de data/hora
		var startAt time.Time
		formats := []string{
			"2006-01-02 15:04",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04",
			"2006-01-02 15:04:05",
		}
		var parseErr error
		for _, f := range formats {
			startAt, parseErr = time.ParseInLocation(f, dataHoraStr, time.Local)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			return fmt.Sprintf(`{"erro": "Formato de data e hora inválido: %s. Use YYYY-MM-DD HH:MM"}`, dataHoraStr)
		}

		apt, err := s.aptService.CreateAppointment(ctx, CreateAppointmentDTO{
			TenantID:       tenantID,
			ServiceID:      serviceID,
			ProfessionalID: professionalID,
			StartAt:        startAt,
			CustomerName:   nomeCliente,
			CustomerPhone:  telefoneCliente,
			Notes:          notes,
		})

		if err != nil {
			return fmt.Sprintf(`{"sucesso": false, "erro": "%s"}`, err.Error())
		}

		res := map[string]interface{}{
			"sucesso":           true,
			"agendamento_id":    apt.ID.String(),
			"cliente":           nomeCliente,
			"servico":           apt.Service.Name,
			"profissional":      apt.Professional.Name,
			"data_hora_inicio":  apt.StartAt.Format("02/01/2006 às 15:04"),
			"valor_total":       apt.TotalPrice,
			"mensagem":          "Agendamento confirmado com sucesso no sistema!",
		}
		b, _ := json.Marshal(res)
		return string(b)

	case "consultar_meus_agendamentos":
		phone, _ := args["telefone_cliente"].(string)
		if phone == "" {
			phone = defaultPhone
		}
		// Consulta os agendamentos do cliente no tenant
		var customer domain.Customer
		if err := s.repo.DB().WithContext(ctx).Where("tenant_id = ? AND phone = ?", tenantID, phone).First(&customer).Error; err != nil {
			return `{"agendamentos": [], "mensagem": "Nenhum agendamento encontrado para este número."}`
		}

		var apts []domain.Appointment
		_ = s.repo.DB().WithContext(ctx).
			Preload("Service").
			Preload("Professional").
			Where("tenant_id = ? AND customer_id = ? AND status = ?", tenantID, customer.ID, domain.StatusConfirmed).
			Order("start_at asc").
			Find(&apts).Error

		type AptDTO struct {
			ID           string  `json:"id"`
			Servico      string  `json:"servico"`
			Profissional string  `json:"profissional"`
			DataHora     string  `json:"data_hora"`
			Preco        float64 `json:"preco"`
			Status       string  `json:"status"`
		}
		list := make([]AptDTO, 0, len(apts))
		for _, a := range apts {
			svcName := ""
			if a.Service != nil {
				svcName = a.Service.Name
			}
			proName := ""
			if a.Professional != nil {
				proName = a.Professional.Name
			}
			list = append(list, AptDTO{
				ID:           a.ID.String(),
				Servico:      svcName,
				Profissional: proName,
				DataHora:     a.StartAt.Format("02/01/2006 às 15:04"),
				Preco:        a.TotalPrice,
				Status:       string(a.Status),
			})
		}
		b, _ := json.Marshal(map[string]interface{}{"agendamentos": list})
		return string(b)

	case "cancelar_agendamento":
		aptIDStr, _ := args["agendamento_id"].(string)
		motivo, _ := args["motivo"].(string)
		if motivo == "" {
			motivo = "Cancelado a pedido do cliente via WhatsApp"
		}

		aptID, err := uuid.Parse(aptIDStr)
		if err != nil {
			return `{"erro": "ID de agendamento inválido"}`
		}

		err = s.repo.UpdateAppointmentStatus(ctx, tenantID, aptID, domain.StatusCancelled, "WhatsApp AI", motivo)
		if err != nil {
			return fmt.Sprintf(`{"sucesso": false, "erro": "%s"}`, err.Error())
		}
		return `{"sucesso": true, "mensagem": "Agendamento cancelado com sucesso."}`

	default:
		return fmt.Sprintf(`{"erro": "Ferramenta desconhecida: %s"}`, toolName)
	}
}

// getAITools retorna as definições de funções compatíveis com Gemini e OpenAI
func (s *AIChatService) getAITools() []ai.Tool {
	return []ai.Tool{
		{
			Type: "function",
			Function: ai.FunctionDefinition{
				Name:        "listar_servicos",
				Description: "Lista todos os serviços oferecidos pela barbearia/salão (nome, preço, duração em minutos e descrição).",
				Parameters: ai.FunctionParameters{
					Type:       "object",
					Properties: map[string]ai.ParameterProperty{},
				},
			},
		},
		{
			Type: "function",
			Function: ai.FunctionDefinition{
				Name:        "listar_profissionais",
				Description: "Lista os barbeiros/profissionais que atendem no estabelecimento. Pode opcionalmente filtrar por ID do serviço.",
				Parameters: ai.FunctionParameters{
					Type: "object",
					Properties: map[string]ai.ParameterProperty{
						"servico_id": {
							Type:        "string",
							Description: "ID UUID do serviço para filtrar profissionais qualificados (opcional).",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: ai.FunctionDefinition{
				Name:        "consultar_horarios_disponiveis",
				Description: "Consulta os horários livres e disponíveis para agendamento em uma data específica. OBRIGATÓRIO chamar antes de sugerir ou agendar.",
				Parameters: ai.FunctionParameters{
					Type: "object",
					Properties: map[string]ai.ParameterProperty{
						"data": {
							Type:        "string",
							Description: "Data a ser consultada no formato YYYY-MM-DD (ex: 2026-09-15).",
						},
						"servico_id": {
							Type:        "string",
							Description: "UUID do serviço desejado.",
						},
						"profissional_id": {
							Type:        "string",
							Description: "UUID do profissional específico, ou omitir/vazio para verificar disponibilidade de qualquer profissional.",
						},
					},
					Required: []string{"data", "servico_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai.FunctionDefinition{
				Name:        "agendar_horario",
				Description: "Realiza a reserva definitiva do atendimento no sistema de agenda. Só invoque após o cliente escolher um horário livre retornado por consultar_horarios_disponiveis.",
				Parameters: ai.FunctionParameters{
					Type: "object",
					Properties: map[string]ai.ParameterProperty{
						"servico_id": {
							Type:        "string",
							Description: "UUID do serviço escolhido.",
						},
						"profissional_id": {
							Type:        "string",
							Description: "UUID do profissional escolhido.",
						},
						"data_hora": {
							Type:        "string",
							Description: "Data e hora de início no formato 'YYYY-MM-DD HH:MM' (ex: 2026-09-15 14:00).",
						},
						"nome_cliente": {
							Type:        "string",
							Description: "Nome do cliente.",
						},
						"telefone_cliente": {
							Type:        "string",
							Description: "Telefone do cliente com DDD.",
						},
						"observacoes": {
							Type:        "string",
							Description: "Observações adicionais ou pedidos especiais do cliente (opcional).",
						},
					},
					Required: []string{"servico_id", "profissional_id", "data_hora", "nome_cliente", "telefone_cliente"},
				},
			},
		},
		{
			Type: "function",
			Function: ai.FunctionDefinition{
				Name:        "consultar_meus_agendamentos",
				Description: "Consulta agendamentos futuros marcados para este cliente pelo telefone.",
				Parameters: ai.FunctionParameters{
					Type: "object",
					Properties: map[string]ai.ParameterProperty{
						"telefone_cliente": {
							Type:        "string",
							Description: "Telefone com DDD do cliente.",
						},
					},
					Required: []string{"telefone_cliente"},
				},
			},
		},
		{
			Type: "function",
			Function: ai.FunctionDefinition{
				Name:        "cancelar_agendamento",
				Description: "Cancela um agendamento do cliente a pedido dele.",
				Parameters: ai.FunctionParameters{
					Type: "object",
					Properties: map[string]ai.ParameterProperty{
						"agendamento_id": {
							Type:        "string",
							Description: "UUID do agendamento que será cancelado.",
						},
						"motivo": {
							Type:        "string",
							Description: "Motivo do cancelamento informado pelo cliente.",
						},
					},
					Required: []string{"agendamento_id"},
				},
			},
		},
	}
}
