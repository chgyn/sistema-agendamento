package service_test

import (
	"context"
	"testing"

	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/service"
	"github.com/sistema-agendamento/backend/pkg/jwt"
)

// MockAsaasProvider mock para testes unitários sem chamadas de rede
type MockAsaasProvider struct {
	ShouldFail bool
}

func (m *MockAsaasProvider) CreateCustomer(ctx context.Context, req domain.AsaasCustomerRequest) (*domain.AsaasCustomerResponse, error) {
	if m.ShouldFail {
		return nil, domain.ErrInvalidCredentials
	}
	return &domain.AsaasCustomerResponse{
		ID:    "cus_mock_test_123",
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

func (m *MockAsaasProvider) CreateSubscription(ctx context.Context, req domain.AsaasSubscriptionRequest) (*domain.AsaasSubscriptionResponse, error) {
	if m.ShouldFail {
		return nil, domain.ErrInvalidCredentials
	}
	return &domain.AsaasSubscriptionResponse{
		ID:          "sub_mock_test_123",
		Customer:    req.Customer,
		Value:       req.Value,
		Cycle:       req.Cycle,
		Status:      "ACTIVE",
		BankSlipURL: "https://sandbox.asaas.com/b/mock",
		InvoiceURL:  "https://sandbox.asaas.com/i/mock",
	}, nil
}

func (m *MockAsaasProvider) GetSubscription(ctx context.Context, asaasSubID string) (*domain.AsaasSubscriptionResponse, error) {
	return &domain.AsaasSubscriptionResponse{
		ID:     asaasSubID,
		Status: "ACTIVE",
	}, nil
}

func (m *MockAsaasProvider) CancelSubscription(ctx context.Context, asaasSubID string) error {
	return nil
}

func (m *MockAsaasProvider) GetSubscriptionPayments(ctx context.Context, asaasSubID string) ([]domain.AsaasPaymentResponse, error) {
	return []domain.AsaasPaymentResponse{}, nil
}

func TestPlanService_CRUD(t *testing.T) {
	_, repo := setupTestDB(t)
	planService := service.NewPlanService(repo)
	ctx := context.Background()

	// 1. Criação de Plano
	isActive := true
	createDTO := domain.CreatePlanDTO{
		Name:             "Plano Teste Pro",
		Description:      "Descrição do plano teste",
		Price:            99.90,
		BillingCycle:     domain.CycleMonthly,
		MaxProfessionals: 5,
		MaxServices:      20,
		Features:         `["Recurso 1", "Recurso 2"]`,
		SortOrder:        1,
		IsActive:         &isActive,
	}

	plan, err := planService.Create(ctx, createDTO)
	if err != nil {
		t.Fatalf("Erro ao criar plano: %v", err)
	}
	if plan.Name != "Plano Teste Pro" || plan.Price != 99.90 {
		t.Errorf("Dados do plano divergentes: %+v", plan)
	}

	// 2. Listagem de Planos
	plans, err := planService.List(ctx, "", nil)
	if err != nil {
		t.Fatalf("Erro ao listar planos: %v", err)
	}
	if len(plans) == 0 {
		t.Errorf("Esperava pelo menos 1 plano, encontrou 0")
	}

	// 3. Atualização de Plano
	updateDTO := domain.UpdatePlanDTO{
		Name:             "Plano Teste Pro Atualizado",
		Description:      "Nova descrição",
		Price:            109.90,
		BillingCycle:     domain.CycleMonthly,
		MaxProfessionals: 6,
		MaxServices:      25,
		IsActive:         &isActive,
	}
	updatedPlan, err := planService.Update(ctx, plan.ID, updateDTO)
	if err != nil {
		t.Fatalf("Erro ao atualizar plano: %v", err)
	}
	if updatedPlan.Name != "Plano Teste Pro Atualizado" || updatedPlan.Price != 109.90 {
		t.Errorf("Dados atualizados inválidos: %+v", updatedPlan)
	}

	// 4. Inativação
	err = planService.ToggleStatus(ctx, plan.ID, false)
	if err != nil {
		t.Fatalf("Erro ao desativar plano: %v", err)
	}
	p, _ := planService.GetByID(ctx, plan.ID)
	if p.IsActive {
		t.Errorf("Esperava plano inativo")
	}

	// 5. Exclusão (sem assinantes deve permitir)
	err = planService.Delete(ctx, plan.ID)
	if err != nil {
		t.Fatalf("Erro ao excluir plano sem assinantes: %v", err)
	}
}

func TestAuthService_RegisterWithPlanAndSubscription(t *testing.T) {
	_, repo := setupTestDB(t)
	mockAsaas := &MockAsaasProvider{}
	subService := service.NewSubscriptionService(repo, mockAsaas)
	jwtService := jwt.NewJWTService("secret-test-key-2026", 24)
	authService := service.NewAuthService(repo, jwtService, subService)
	planService := service.NewPlanService(repo)
	ctx := context.Background()

	// Cria um plano para o teste
	isActive := true
	plan, err := planService.Create(ctx, domain.CreatePlanDTO{
		Name:         "Plano Onboarding",
		Price:        79.90,
		BillingCycle: domain.CycleMonthly,
		IsActive:     &isActive,
	})
	if err != nil {
		t.Fatalf("Erro ao preparar plano: %v", err)
	}

	// Executa Onboarding do Tenant
	regDTO := domain.RegisterTenantWithPlanDTO{
		PlanID:     plan.ID,
		TenantName: "Barbearia Imperial Teste",
		Slug:       "barbearia-imperial-teste",
		Document:   "12.345.678/0001-99",
		Phone:      "(11) 98888-7777",
		City:       "São Paulo",
		State:      "SP",
		AdminName:  "Dono do Estabelecimento",
		AdminEmail: "dono@imperialteste.com",
		Password:   "senhaSegura123",
	}

	res, err := authService.RegisterTenant(ctx, regDTO)
	if err != nil {
		t.Fatalf("Erro ao registrar tenant com plano: %v", err)
	}

	if res.Tenant == nil || res.User == nil || res.Subscription == nil {
		t.Fatalf("Dados de resposta incompletos: %+v", res)
	}

	if res.Subscription.Status != domain.SubscriptionStatusPending {
		t.Errorf("Esperava status inicial PENDING, recebeu %s", res.Subscription.Status)
	}

	if res.Subscription.AsaasCustomerID != "cus_mock_test_123" {
		t.Errorf("AsaasCustomerID incorreto: %s", res.Subscription.AsaasCustomerID)
	}

	// Tentar deletar o plano com assinante vinculado deve retornar erro
	err = planService.Delete(ctx, plan.ID)
	if err != domain.ErrPlanHasSubscribers {
		t.Errorf("Esperava ErrPlanHasSubscribers, recebeu: %v", err)
	}

	// Processa Webhook do Asaas simulando confirmação de pagamento
	paymentMap := map[string]interface{}{
		"id":           "pay_12345",
		"subscription": res.Subscription.AsaasSubscriptionID,
		"value":        79.90,
		"netValue":     78.00,
		"billingType":  "PIX",
		"invoiceUrl":   "https://asaas.com/i/12345",
	}
	err = subService.ProcessAsaasWebhookEvent(ctx, "PAYMENT_CONFIRMED", paymentMap, nil)
	if err != nil {
		t.Fatalf("Erro ao processar webhook: %v", err)
	}

	// Verifica se a assinatura foi ativada
	subUpdated, err := subService.GetByTenantID(ctx, res.Tenant.ID)
	if err != nil {
		t.Fatalf("Erro ao buscar assinatura atualizada: %v", err)
	}
	if subUpdated.Status != domain.SubscriptionStatusActive {
		t.Errorf("Esperava assinatura ACTIVE após webhook, recebeu %s", subUpdated.Status)
	}
}
