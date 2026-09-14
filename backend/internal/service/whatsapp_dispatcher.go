package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/integrations/wuzapi"
)

type WhatsAppMessageDispatcher struct {
	wuzapiClient *wuzapi.Client
}

func NewWhatsAppMessageDispatcher(wuzapiClient *wuzapi.Client) *WhatsAppMessageDispatcher {
	return &WhatsAppMessageDispatcher{
		wuzapiClient: wuzapiClient,
	}
}

// DispatchParams parâmetros para o envio humanizado
type DispatchParams struct {
	InstanceToken       string
	CustomerPhone       string
	Message             string
	MinDelaySec         int
	MaxDelaySec         int
	TypingSpeedCharsSec int
}

// CalculateTypingDuration calcula o tempo de digitação simulada proporcional à mensagem
func CalculateTypingDuration(text string, charsPerSec int) time.Duration {
	if charsPerSec <= 0 {
		charsPerSec = 35
	}
	charCount := len([]rune(text))
	seconds := float64(charCount) / float64(charsPerSec)

	// Delimita o tempo de digitação entre 1.5s e 7.5s para manter a experiência natural
	if seconds < 1.5 {
		seconds = 1.5
	} else if seconds > 7.5 {
		seconds = 7.5
	}

	return time.Duration(seconds * float64(time.Second))
}

// CalculateInitialDelay calcula o tempo de espera aleatório antes de iniciar a resposta
func CalculateInitialDelay(minSec, maxSec int) time.Duration {
	if minSec <= 0 {
		minSec = 2
	}
	if maxSec < minSec {
		maxSec = minSec + 2
	}

	diff := maxSec - minSec
	delaySec := minSec
	if diff > 0 {
		delaySec = minSec + rand.Intn(diff+1)
	}

	// Adiciona jitter em milissegundos
	jitterMs := rand.Intn(800)
	return time.Duration(delaySec)*time.Second + time.Duration(jitterMs)*time.Millisecond
}

// SendHumanizedMessage executa o fluxo completo de envio humanizado
func (d *WhatsAppMessageDispatcher) SendHumanizedMessage(ctx context.Context, params DispatchParams) error {
	trimmed := strings.TrimSpace(params.Message)
	if trimmed == "" {
		return nil
	}

	// 1. Aplica tempo de espera inicial (variável para evitar padrões previsíveis)
	initialDelay := CalculateInitialDelay(params.MinDelaySec, params.MaxDelaySec)
	log.Printf("⏳ [Humanized Dispatcher] Aguardando pausa inicial de %v antes de responder %s...", initialDelay, params.CustomerPhone)
	select {
	case <-time.After(initialDelay):
	case <-ctx.Done():
		return ctx.Err()
	}

	// 2. Se a mensagem for muito longa ou contiver separação explícita de blocos ("---"), divide
	chunks := splitIntoNaturalChunks(trimmed)

	for i, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		// 3. Ativa status de "digitando..." no WhatsApp
		log.Printf("⌨️ [Humanized Dispatcher] Sinalizando 'digitando...' para %s...", params.CustomerPhone)
		_ = d.wuzapiClient.SetPresence(ctx, params.InstanceToken, params.CustomerPhone, "composing")

		// 4. Aguarda período adequado de digitação baseado no tamanho do texto
		typingDuration := CalculateTypingDuration(chunk, params.TypingSpeedCharsSec)
		select {
		case <-time.After(typingDuration):
		case <-ctx.Done():
			_ = d.wuzapiClient.SetPresence(ctx, params.InstanceToken, params.CustomerPhone, "paused")
			return ctx.Err()
		}

		// 5. Envia a mensagem através do WUZAPI
		_ = d.wuzapiClient.SetPresence(ctx, params.InstanceToken, params.CustomerPhone, "paused")
		if err := d.wuzapiClient.SendTextMessage(ctx, params.InstanceToken, params.CustomerPhone, chunk); err != nil {
			log.Printf("❌ [Humanized Dispatcher] Erro ao enviar mensagem para %s: %v", params.CustomerPhone, err)
			return fmt.Errorf("falha ao enviar mensagem via WUZAPI: %w", err)
		}

		log.Printf("✅ [Humanized Dispatcher] Mensagem enviada para %s (%d/%d)", params.CustomerPhone, i+1, len(chunks))

		// 6. Se houver próximo envio, aplica intervalo humano entre mensagens
		if i < len(chunks)-1 {
			interval := time.Duration(1500+rand.Intn(1200)) * time.Millisecond
			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

// splitIntoNaturalChunks divide textos excessivamente longos ou com separadores em mensagens separadas
func splitIntoNaturalChunks(text string) []string {
	if strings.Contains(text, "\n---\n") {
		return strings.Split(text, "\n---\n")
	}

	// Caso padrão: envia a mensagem em bloco único preservando formatação markdown do WhatsApp
	return []string{text}
}

// BuildDispatcherParams helper para criar parâmetros a partir da configuração do Tenant
func BuildDispatcherParams(cfg *domain.TenantWhatsAppConfig, phone, message string) DispatchParams {
	minDelay := cfg.HumanizedMinDelaySec
	if minDelay <= 0 {
		minDelay = 2
	}
	maxDelay := cfg.HumanizedMaxDelaySec
	if maxDelay <= 0 {
		maxDelay = 5
	}
	speed := cfg.TypingSpeedCharsSec
	if speed <= 0 {
		speed = 35
	}

	return DispatchParams{
		InstanceToken:       cfg.InstanceToken,
		CustomerPhone:       phone,
		Message:             message,
		MinDelaySec:         minDelay,
		MaxDelaySec:         maxDelay,
		TypingSpeedCharsSec: speed,
	}
}
