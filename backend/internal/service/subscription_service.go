package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type SubscriptionService struct {
	repo          *postgres.Repository
	asaasProvider domain.AsaasProvider
}

func NewSubscriptionService(repo *postgres.Repository, asaasProvider domain.AsaasProvider) *SubscriptionService {
	return &SubscriptionService{
		repo:          repo,
		asaasProvider: asaasProvider,
	}
}

// CreateTenantSubscription vincula o tenant a um plano e cria cliente + assinatura no Asaas
func (s *SubscriptionService) CreateTenantSubscription(ctx context.Context, tenant *domain.Tenant, plan *domain.Plan, adminUser *domain.User) (*domain.Subscription, error) {
	// 1. Cria ou obtém cliente no Asaas
	var asaasCustID string
	var asaasSubID string
	var paymentURL string

	if s.asaasProvider != nil {
		custReq := domain.AsaasCustomerRequest{
			Name:              tenant.Name,
			Email:             adminUser.Email,
			Phone:             tenant.Phone,
			MobilePhone:       tenant.Phone,
			CpfCnpj:           tenant.Document,
			ExternalReference: tenant.ID.String(),
		}
		custResp, err := s.asaasProvider.CreateCustomer(ctx, custReq)
		if err != nil {
			log.Printf("⚠️ [Asaas] Aviso ao criar cliente Asaas (%v). Prosseguindo com fallback simulado...", err)
			asaasCustID = fmt.Sprintf("cus_mock_%s", tenant.ID.String()[:8])
		} else {
			asaasCustID = custResp.ID
		}

		// Próximo vencimento: 1 dia ou imediato
		nextDueDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		subReq := domain.AsaasSubscriptionRequest{
			Customer:          asaasCustID,
			BillingType:       "UNDEFINED",
			Value:             plan.Price,
			NextDueDate:       nextDueDate,
			Cycle:             string(plan.BillingCycle),
			Description:       fmt.Sprintf("Assinatura %s - %s", plan.Name, tenant.Name),
			ExternalReference: tenant.ID.String(),
		}
		subResp, err := s.asaasProvider.CreateSubscription(ctx, subReq)
		if err != nil {
			log.Printf("⚠️ [Asaas] Aviso ao criar assinatura no Asaas (%v). Prosseguindo com mock...", err)
			asaasSubID = fmt.Sprintf("sub_mock_%s", tenant.ID.String()[:8])
			paymentURL = "https://sandbox.asaas.com/payment/mock"
		} else {
			asaasSubID = subResp.ID
			paymentURL = subResp.InvoiceURL
			if paymentURL == "" {
				paymentURL = subResp.BankSlipURL
			}
		}
	} else {
		asaasCustID = fmt.Sprintf("cus_mock_%s", tenant.ID.String()[:8])
		asaasSubID = fmt.Sprintf("sub_mock_%s", tenant.ID.String()[:8])
	}

	dueDate := time.Now().AddDate(0, 0, 1)
	sub := domain.Subscription{
		ID:                  uuid.New(),
		TenantID:            tenant.ID,
		PlanID:              plan.ID,
		AsaasCustomerID:     asaasCustID,
		AsaasSubscriptionID: asaasSubID,
		Status:              domain.SubscriptionStatusPending,
		BillingCycle:        plan.BillingCycle,
		Price:               plan.Price,
		NextDueDate:         &dueDate,
		PaymentMethod:       "UNDEFINED",
		PaymentURL:          paymentURL,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.repo.CreateSubscription(ctx, &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

// ProcessAsaasWebhookEvent processa notificações de eventos do Asaas (pagamentos, alterações de assinatura)
func (s *SubscriptionService) ProcessAsaasWebhookEvent(ctx context.Context, event string, paymentMap, subMap map[string]interface{}) error {
	log.Printf("🔔 [Webhook Asaas] Processando evento: %s", event)

	var asaasSubID string
	if sID, ok := paymentMap["subscription"].(string); ok && sID != "" {
		asaasSubID = sID
	} else if sID, ok := subMap["id"].(string); ok && sID != "" {
		asaasSubID = sID
	}

	if asaasSubID == "" {
		log.Printf("ℹ️ [Webhook Asaas] Evento %s sem ID de assinatura vinculado. Ignorando.", event)
		return nil
	}

	sub, err := s.repo.GetSubscriptionByAsaasID(ctx, asaasSubID)
	if err != nil {
		log.Printf("⚠️ [Webhook Asaas] Assinatura não localizada para asaasSubID %s: %v", asaasSubID, err)
		return nil // Não falha para evitar retentativas desnecessárias de webhook não relacionado
	}

	switch event {
	case "PAYMENT_RECEIVED", "PAYMENT_CONFIRMED":
		sub.Status = domain.SubscriptionStatusActive
		now := time.Now()
		// Adiciona período conforme o ciclo
		var periodEnd time.Time
		switch sub.BillingCycle {
		case domain.CycleQuarterly:
			periodEnd = now.AddDate(0, 3, 0)
		case domain.CycleSemiannual:
			periodEnd = now.AddDate(0, 6, 0)
		case domain.CycleYearly:
			periodEnd = now.AddDate(1, 0, 0)
		default: // Monthly
			periodEnd = now.AddDate(0, 1, 0)
		}
		sub.CurrentPeriodEnd = &periodEnd
		sub.NextDueDate = &periodEnd
		sub.UpdatedAt = now

		if err := s.repo.UpdateSubscription(ctx, sub); err != nil {
			return fmt.Errorf("falha ao atualizar status da assinatura: %w", err)
		}

		// Registra fatura
		if pID, ok := paymentMap["id"].(string); ok && pID != "" {
			value, _ := paymentMap["value"].(float64)
			netValue, _ := paymentMap["netValue"].(float64)
			billingType, _ := paymentMap["billingType"].(string)
			invoiceURL, _ := paymentMap["invoiceUrl"].(string)
			bankSlipURL, _ := paymentMap["bankSlipUrl"].(string)

			invoice := domain.SubscriptionInvoice{
				ID:             uuid.New(),
				TenantID:       sub.TenantID,
				SubscriptionID: sub.ID,
				AsaasPaymentID: pID,
				Status:         "CONFIRMED",
				Value:          value,
				NetValue:       netValue,
				BillingType:    billingType,
				DueDate:        now,
				PaymentDate:    &now,
				InvoiceURL:     invoiceURL,
				BankSlipURL:    bankSlipURL,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			_ = s.repo.SaveSubscriptionInvoice(ctx, &invoice)
		}

	case "PAYMENT_OVERDUE":
		sub.Status = domain.SubscriptionStatusOverdue
		sub.UpdatedAt = time.Now()
		_ = s.repo.UpdateSubscription(ctx, sub)

	case "PAYMENT_DELETED", "PAYMENT_REFUNDED":
		log.Printf("⚠️ [Webhook Asaas] Pagamento estornado ou deletado para assinatura %s", sub.ID)

	case "SUBSCRIPTION_DELETED", "SUBSCRIPTION_CANCELLED":
		sub.Status = domain.SubscriptionStatusCancelled
		sub.UpdatedAt = time.Now()
		_ = s.repo.UpdateSubscription(ctx, sub)
	}

	return nil
}

func (s *SubscriptionService) GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetSubscriptionByTenantID(ctx, tenantID)
}

func (s *SubscriptionService) ListAll(ctx context.Context, status, search string) ([]domain.Subscription, error) {
	return s.repo.ListAllSubscriptions(ctx, status, search)
}

func (s *SubscriptionService) OverrideStatus(ctx context.Context, subID uuid.UUID, status domain.SubscriptionStatus) error {
	return s.repo.UpdateSubscriptionStatus(ctx, subID, status)
}
