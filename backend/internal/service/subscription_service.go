package service

import (
	"context"
	"fmt"
	"log"
	"strings"
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

// CreateTenantSubscription vincula o tenant a um plano e cria cliente + assinatura no Asaas (ou ativa diretamente se plano gratuito)
func (s *SubscriptionService) CreateTenantSubscription(ctx context.Context, tenant *domain.Tenant, plan *domain.Plan, adminUser *domain.User) (*domain.Subscription, error) {
	// 1. Caso seja plano gratuito, ativa imediatamente sem integração no Asaas
	if plan.IsFree {
		sub := domain.Subscription{
			ID:                  uuid.New(),
			TenantID:            tenant.ID,
			PlanID:              plan.ID,
			AsaasCustomerID:     "",
			AsaasSubscriptionID: "",
			Status:              domain.SubscriptionStatusActive,
			BillingCycle:        plan.BillingCycle,
			Price:               0.00,
			Origin:              domain.SubscriptionOriginFreePlan,
			NextDueDate:         nil,
			CurrentPeriodEnd:    nil,
			PaymentMethod:       "FREE",
			PaymentURL:          "",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		if err := s.repo.CreateSubscription(ctx, &sub); err != nil {
			return nil, err
		}

		auditLog := domain.SubscriptionAuditLog{
			ID:             uuid.New(),
			SubscriptionID: sub.ID,
			TenantID:       tenant.ID,
			PlanID:         plan.ID,
			Action:         domain.ActionFreeRegistration,
			PreviousStatus: "",
			NewStatus:      domain.SubscriptionStatusActive,
			PerformedByID:  nil,
			Reason:         "Cadastro público com plano gratuito ativado imediatamente",
			CreatedAt:      time.Now(),
		}
		_ = s.repo.CreateSubscriptionAuditLog(ctx, &auditLog)

		return &sub, nil
	}

	// 2. Plano pago: Cria ou obtém cliente no Asaas
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
		Origin:              domain.SubscriptionOriginAsaas,
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
		previousStatus := sub.Status
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

		auditLog := domain.SubscriptionAuditLog{
			ID:             uuid.New(),
			SubscriptionID: sub.ID,
			TenantID:       sub.TenantID,
			PlanID:         sub.PlanID,
			Action:         domain.ActionAsaasSync,
			PreviousStatus: previousStatus,
			NewStatus:      domain.SubscriptionStatusActive,
			Reason:         fmt.Sprintf("Pagamento confirmado via Asaas (%s)", event),
			CreatedAt:      now,
		}
		_ = s.repo.CreateSubscriptionAuditLog(ctx, &auditLog)

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
		previousStatus := sub.Status
		sub.Status = domain.SubscriptionStatusOverdue
		sub.UpdatedAt = time.Now()
		_ = s.repo.UpdateSubscription(ctx, sub)

		auditLog := domain.SubscriptionAuditLog{
			ID:             uuid.New(),
			SubscriptionID: sub.ID,
			TenantID:       sub.TenantID,
			PlanID:         sub.PlanID,
			Action:         domain.ActionAsaasSync,
			PreviousStatus: previousStatus,
			NewStatus:      domain.SubscriptionStatusOverdue,
			Reason:         "Cobrança vencida sem pagamento no Asaas",
			CreatedAt:      time.Now(),
		}
		_ = s.repo.CreateSubscriptionAuditLog(ctx, &auditLog)

	case "PAYMENT_DELETED", "PAYMENT_REFUNDED":
		log.Printf("⚠️ [Webhook Asaas] Pagamento estornado ou deletado para assinatura %s", sub.ID)

	case "SUBSCRIPTION_DELETED", "SUBSCRIPTION_CANCELLED":
		previousStatus := sub.Status
		sub.Status = domain.SubscriptionStatusCancelled
		sub.UpdatedAt = time.Now()
		_ = s.repo.UpdateSubscription(ctx, sub)

		auditLog := domain.SubscriptionAuditLog{
			ID:             uuid.New(),
			SubscriptionID: sub.ID,
			TenantID:       sub.TenantID,
			PlanID:         sub.PlanID,
			Action:         domain.ActionAsaasSync,
			PreviousStatus: previousStatus,
			NewStatus:      domain.SubscriptionStatusCancelled,
			Reason:         "Assinatura cancelada no Asaas",
			CreatedAt:      time.Now(),
		}
		_ = s.repo.CreateSubscriptionAuditLog(ctx, &auditLog)
	}

	return nil
}

func (s *SubscriptionService) GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetSubscriptionByTenantID(ctx, tenantID)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetSubscriptionByID(ctx, id)
}

func (s *SubscriptionService) ListAll(ctx context.Context, status, search string) ([]domain.Subscription, error) {
	return s.repo.ListAllSubscriptions(ctx, status, search)
}

// GrantManualSubscription concede ou atualiza manualmente a assinatura de um estabelecimento
func (s *SubscriptionService) GrantManualSubscription(ctx context.Context, tenantID uuid.UUID, dto domain.GrantManualSubscriptionDTO, adminUserID uuid.UUID) (*domain.Subscription, error) {
	tenant, err := s.repo.GetTenantByID(ctx, tenantID)
	if err != nil || tenant == nil {
		return nil, domain.ErrTenantNotFound
	}

	plan, err := s.repo.GetPlanByID(ctx, dto.PlanID)
	if err != nil || plan == nil {
		return nil, domain.ErrPlanNotFound
	}

	existingSub, _ := s.repo.GetSubscriptionByTenantID(ctx, tenantID)

	now := time.Now()
	if existingSub != nil {
		previousStatus := existingSub.Status
		existingSub.PlanID = plan.ID
		existingSub.Status = domain.SubscriptionStatusActive
		existingSub.Price = plan.Price
		existingSub.BillingCycle = plan.BillingCycle
		existingSub.Origin = domain.SubscriptionOriginManual
		existingSub.CurrentPeriodEnd = dto.ExpiresAt
		existingSub.NextDueDate = dto.ExpiresAt
		existingSub.PaymentMethod = "MANUAL"
		existingSub.ManualGrantReason = strings.TrimSpace(dto.Reason)
		existingSub.GrantedByUserID = &adminUserID
		existingSub.UpdatedAt = now

		if err := s.repo.UpdateSubscription(ctx, existingSub); err != nil {
			return nil, fmt.Errorf("falha ao atualizar assinatura manual: %w", err)
		}

		auditLog := domain.SubscriptionAuditLog{
			ID:             uuid.New(),
			SubscriptionID: existingSub.ID,
			TenantID:       tenant.ID,
			PlanID:         plan.ID,
			Action:         domain.ActionManualGrant,
			PreviousStatus: previousStatus,
			NewStatus:      domain.SubscriptionStatusActive,
			PerformedByID:  &adminUserID,
			Reason:         strings.TrimSpace(dto.Reason),
			ExpiresAt:      dto.ExpiresAt,
			CreatedAt:      now,
		}
		_ = s.repo.CreateSubscriptionAuditLog(ctx, &auditLog)

		return s.repo.GetSubscriptionByID(ctx, existingSub.ID)
	}

	// Criação de nova assinatura se não existir
	newSub := domain.Subscription{
		ID:                  uuid.New(),
		TenantID:            tenant.ID,
		PlanID:              plan.ID,
		AsaasCustomerID:     "",
		AsaasSubscriptionID: "",
		Status:              domain.SubscriptionStatusActive,
		BillingCycle:        plan.BillingCycle,
		Price:               plan.Price,
		Origin:              domain.SubscriptionOriginManual,
		NextDueDate:         dto.ExpiresAt,
		CurrentPeriodEnd:    dto.ExpiresAt,
		PaymentMethod:       "MANUAL",
		PaymentURL:          "",
		ManualGrantReason:   strings.TrimSpace(dto.Reason),
		GrantedByUserID:     &adminUserID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.repo.CreateSubscription(ctx, &newSub); err != nil {
		return nil, fmt.Errorf("falha ao criar assinatura manual: %w", err)
	}

	auditLog := domain.SubscriptionAuditLog{
		ID:             uuid.New(),
		SubscriptionID: newSub.ID,
		TenantID:       tenant.ID,
		PlanID:         plan.ID,
		Action:         domain.ActionManualGrant,
		PreviousStatus: "",
		NewStatus:      domain.SubscriptionStatusActive,
		PerformedByID:  &adminUserID,
		Reason:         strings.TrimSpace(dto.Reason),
		ExpiresAt:      dto.ExpiresAt,
		CreatedAt:      now,
	}
	_ = s.repo.CreateSubscriptionAuditLog(ctx, &auditLog)

	return s.repo.GetSubscriptionByID(ctx, newSub.ID)
}

func (s *SubscriptionService) OverrideStatus(ctx context.Context, subID uuid.UUID, status domain.SubscriptionStatus, reason string, adminUserID *uuid.UUID) error {
	sub, err := s.repo.GetSubscriptionByID(ctx, subID)
	if err != nil {
		return err
	}

	previousStatus := sub.Status
	if err := s.repo.UpdateSubscriptionStatus(ctx, subID, status); err != nil {
		return err
	}

	auditLog := domain.SubscriptionAuditLog{
		ID:             uuid.New(),
		SubscriptionID: sub.ID,
		TenantID:       sub.TenantID,
		PlanID:         sub.PlanID,
		Action:         domain.ActionStatusOverride,
		PreviousStatus: previousStatus,
		NewStatus:      status,
		PerformedByID:  adminUserID,
		Reason:         strings.TrimSpace(reason),
		CreatedAt:      time.Now(),
	}
	return s.repo.CreateSubscriptionAuditLog(ctx, &auditLog)
}

func (s *SubscriptionService) ListAuditLogs(ctx context.Context, subscriptionID uuid.UUID) ([]domain.SubscriptionAuditLog, error) {
	return s.repo.ListSubscriptionAuditLogs(ctx, subscriptionID)
}

func (s *SubscriptionService) CheckAndExpireSubscriptions(ctx context.Context) (int64, error) {
	return s.repo.CheckAndExpireSubscriptions(ctx)
}
