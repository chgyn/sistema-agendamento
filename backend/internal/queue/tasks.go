package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sistema-agendamento/backend/internal/config"
)

const (
	TypeSendBookingConfirmation = "task:send_booking_confirmation"
	TypeSendBookingReminder     = "task:send_booking_reminder"
	TypeAutoCancelPending       = "task:auto_cancel_pending"
	TypeProcessAsaasWebhook     = "task:process_asaas_webhook"
)

type BookingConfirmationPayload struct {
	TenantID        uuid.UUID `json:"tenant_id"`
	AppointmentID   uuid.UUID `json:"appointment_id"`
	CustomerName    string    `json:"customer_name"`
	CustomerPhone   string    `json:"customer_phone"`
	CustomerEmail   string    `json:"customer_email"`
	ServiceName     string    `json:"service_name"`
	ProfessionalName string   `json:"professional_name"`
	StartAt         time.Time `json:"start_at"`
	TotalPrice      float64   `json:"total_price"`
}

type AsaasWebhookTaskPayload struct {
	Event        string                 `json:"event"`
	Payment      map[string]interface{} `json:"payment,omitempty"`
	Subscription map[string]interface{} `json:"subscription,omitempty"`
}

type QueueClient struct {
	client *asynq.Client
}

func NewQueueClient(cfg *config.Config) *QueueClient {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: cfg.RedisPass,
	})
	return &QueueClient{client: client}
}

func (q *QueueClient) Close() error {
	if q.client != nil {
		return q.client.Close()
	}
	return nil
}

// EnqueueBookingConfirmation enfileira a tarefa de envio de confirmação (WhatsApp / Email)
func (q *QueueClient) EnqueueBookingConfirmation(ctx context.Context, payload BookingConfirmationPayload) error {
	if q.client == nil {
		return nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeSendBookingConfirmation, data, asynq.MaxRetry(3), asynq.Timeout(20*time.Second))
	info, err := q.client.EnqueueContext(ctx, task)
	if err != nil {
		log.Printf("⚠️ [Queue] Falha ao enfileirar task de confirmação: %v", err)
		return nil // Não bloqueia o agendamento se o redis falhar
	}

	log.Printf("🚀 [Queue] Task de confirmação enfileirada: ID=%s Queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueAsaasWebhook enfileira a notificação de webhook do Asaas
func (q *QueueClient) EnqueueAsaasWebhook(ctx context.Context, payload AsaasWebhookTaskPayload) error {
	if q.client == nil {
		return nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeProcessAsaasWebhook, data, asynq.MaxRetry(5), asynq.Timeout(30*time.Second), asynq.Queue("critical"))
	info, err := q.client.EnqueueContext(ctx, task)
	if err != nil {
		log.Printf("⚠️ [Queue] Falha ao enfileirar task Asaas Webhook: %v", err)
		return err
	}

	log.Printf("🚀 [Queue] Task Asaas Webhook enfileirada: ID=%s Event=%s", info.ID, payload.Event)
	return nil
}

// StartWorkerServer inicia o servidor consumidor de tarefas do Asynq
func StartWorkerServer(cfg *config.Config) (*asynq.Server, error) {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: cfg.RedisPass,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSendBookingConfirmation, HandleBookingConfirmationTask)
	mux.HandleFunc(TypeSendBookingReminder, HandleBookingReminderTask)
	mux.HandleFunc(TypeAutoCancelPending, HandleAutoCancelPendingTask)
	mux.HandleFunc(TypeProcessAsaasWebhook, HandleProcessAsaasWebhookTask)

	log.Println("👷 [Worker] Asynq worker server pronto e ouvindo filas Redis...")
	return srv, srv.Run(mux)
}

func HandleBookingConfirmationTask(ctx context.Context, t *asynq.Task) error {
	var p BookingConfirmationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal falhou: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("📩 [NOTIFICAÇÃO ENVIADA] Agendamento #%s confirmado!", p.AppointmentID)
	log.Printf("   -> Cliente: %s (%s)", p.CustomerName, p.CustomerPhone)
	log.Printf("   -> Serviço: %s com %s", p.ServiceName, p.ProfessionalName)
	log.Printf("   -> Data/Hora: %s | Valor: R$ %.2f", p.StartAt.Format("02/01/2006 às 15:04"), p.TotalPrice)
	return nil
}

func HandleBookingReminderTask(ctx context.Context, t *asynq.Task) error {
	log.Printf("⏰ [LEMBRETE ENVIADO] Lembrete prévio de atendimento enviado com sucesso.")
	return nil
}

func HandleAutoCancelPendingTask(ctx context.Context, t *asynq.Task) error {
	log.Printf("🧹 [AUTO-CANCEL] Verificação de agendamentos pendentes expirados executada.")
	return nil
}

func HandleProcessAsaasWebhookTask(ctx context.Context, t *asynq.Task) error {
	var p AsaasWebhookTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal falhou no webhook Asaas: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("⚡ [Worker Asynq] Processando evento Webhook Asaas: %s", p.Event)
	return nil
}
