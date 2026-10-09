package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestPlanService_FreePlanCRUD(t *testing.T) {
	_, repo := setupTestDB(t)
	planService := service.NewPlanService(repo)
	ctx := context.Background()

	// 1. Criação de Plano Gratuito com Price = 0 e IsFree = true
	isActive := true
	createDTO := domain.CreatePlanDTO{
		Name:             "Plano Degustação Gratuito",
		Description:      "Acesso gratuito sem cobrança",
		Price:            999.00, // Mesmo enviando valor, se IsFree=true o serviço força 0.00
		BillingCycle:     domain.CycleMonthly,
		MaxProfessionals: 1,
		MaxServices:      5,
		Features:         `["1 Profissional", "Até 5 Serviços"]`,
		SortOrder:        0,
		IsActive:         &isActive,
		IsFree:           true,
	}

	plan, err := planService.Create(ctx, createDTO)
	if err != nil {
		t.Fatalf("Erro ao criar plano gratuito: %v", err)
	}

	if !plan.IsFree {
		t.Errorf("Esperava IsFree=true, obteve false")
	}
	if plan.Price != 0.00 {
		t.Errorf("Esperava Price=0.00 para plano gratuito, obteve %.2f", plan.Price)
	}

	// 2. Atualização mantendo IsFree
	updateDTO := domain.UpdatePlanDTO{
		Name:             "Plano Degustação Gratuito V2",
		Description:      "Descrição atualizada",
		Price:            0.00,
		BillingCycle:     domain.CycleMonthly,
		MaxProfessionals: 2,
		MaxServices:      10,
		IsActive:         &isActive,
		IsFree:           true,
	}

	updatedPlan, err := planService.Update(ctx, plan.ID, updateDTO)
	if err != nil {
		t.Fatalf("Erro ao atualizar plano gratuito: %v", err)
	}

	if updatedPlan.MaxProfessionals != 2 || updatedPlan.Price != 0.00 || !updatedPlan.IsFree {
		t.Errorf("Divergência ao atualizar plano gratuito: %+v", updatedPlan)
	}
}

func TestAuthService_RegisterWithFreePlan(t *testing.T) {
	_, repo := setupTestDB(t)
	mockAsaas := &MockAsaasProvider{ShouldFail: true} // Mesmo que Asaas falhasse, plano gratuito NÃO deve chamar Asaas
	subService := service.NewSubscriptionService(repo, mockAsaas)
	jwtService := jwt.NewJWTService("secret-test-key-2026", 24)
	authService := service.NewAuthService(repo, jwtService, subService)
	planService := service.NewPlanService(repo)
	ctx := context.Background()

	// 1. Cadastra plano gratuito
	isActive := true
	freePlan, err := planService.Create(ctx, domain.CreatePlanDTO{
		Name:             "Plano Comunidade Free",
		Price:            0.00,
		BillingCycle:     domain.CycleMonthly,
		MaxProfessionals: 1,
		MaxServices:      5,
		IsActive:         &isActive,
		IsFree:           true,
	})
	if err != nil {
		t.Fatalf("Erro ao criar plano gratuito: %v", err)
	}

	// 2. Onboarding público utilizando plano gratuito
	regDTO := domain.RegisterTenantWithPlanDTO{
		PlanID:     freePlan.ID,
		TenantName: "Barbearia do Povo",
		Slug:       "barbearia-do-povo",
		Document:   "99.888.777/0001-66",
		Phone:      "(11) 91111-2222",
		City:       "Guarulhos",
		State:      "SP",
		AdminName:  "Administrador Free",
		AdminEmail: "admin@barbeariadopovo.com",
		Password:   "senhaSegura123",
	}

	res, err := authService.RegisterTenant(ctx, regDTO)
	if err != nil {
		t.Fatalf("Falha no onboarding com plano gratuito: %v", err)
	}

	if res.Subscription == nil {
		t.Fatalf("Assinatura não foi gerada para plano gratuito")
	}

	// Status deve ser ACTIVE imediatamente
	if res.Subscription.Status != domain.SubscriptionStatusActive {
		t.Errorf("Esperava status ACTIVE para plano gratuito, obteve: %s", res.Subscription.Status)
	}

	// Origem deve ser FREE_PLAN
	if res.Subscription.Origin != domain.SubscriptionOriginFreePlan {
		t.Errorf("Esperava Origin FREE_PLAN, obteve: %s", res.Subscription.Origin)
	}

	// Sem IDs de cliente/assinatura do Asaas
	if res.Subscription.AsaasSubscriptionID != "" || res.Subscription.AsaasCustomerID != "" {
		t.Errorf("Plano gratuito não deveria registrar IDs no Asaas, obteve: sub=%s cus=%s", res.Subscription.AsaasSubscriptionID, res.Subscription.AsaasCustomerID)
	}

	// Verifica se registrou log de auditoria
	logs, err := subService.ListAuditLogs(ctx, res.Subscription.ID)
	if err != nil {
		t.Fatalf("Erro ao listar logs de auditoria: %v", err)
	}
	if len(logs) == 0 {
		t.Errorf("Esperava ao menos 1 log de auditoria para o onboarding gratuito")
	} else if logs[0].Action != domain.ActionFreeRegistration {
		t.Errorf("Esperava ação FREE_REGISTRATION no log de auditoria, obteve: %s", logs[0].Action)
	}
}

func TestSubscriptionService_GrantManualAndAuditLogs(t *testing.T) {
	_, repo := setupTestDB(t)
	subService := service.NewSubscriptionService(repo, nil)
	planService := service.NewPlanService(repo)
	ctx := context.Background()

	// 1. Cria Tenant diretamente
	tenantID := uuid.New()
	tenant := domain.Tenant{
		ID:        tenantID,
		Name:      "Barbearia VIP Manual",
		Slug:      "barbearia-vip-manual",
		Phone:     "(11) 97777-6666",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.CreateTenant(ctx, &tenant); err != nil {
		t.Fatalf("Erro ao criar tenant: %v", err)
	}

	// 2. Cria Usuário Admin Global (autor da liberação)
	adminUserID := uuid.New()
	adminUser := domain.User{
		ID:           adminUserID,
		Name:         "Super Admin Global",
		Email:        "superadmin@sistema.com",
		PasswordHash: "hash-fake",
		Role:         domain.RoleAdminGlobal,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := repo.CreateUser(ctx, &adminUser); err != nil {
		t.Fatalf("Erro ao criar admin: %v", err)
	}

	// 3. Cria Plano Pro
	isActive := true
	planPro, err := planService.Create(ctx, domain.CreatePlanDTO{
		Name:             "Plano Pro Parceiros",
		Price:            149.90,
		BillingCycle:     domain.CycleMonthly,
		MaxProfessionals: 10,
		MaxServices:      50,
		IsActive:         &isActive,
	})
	if err != nil {
		t.Fatalf("Erro ao criar plano: %v", err)
	}

	// 4. Executa Concessão Manual com Validade de 60 dias
	expiresAt := time.Now().AddDate(0, 2, 0)
	grantDTO := domain.GrantManualSubscriptionDTO{
		PlanID:    planPro.ID,
		ExpiresAt: &expiresAt,
		Reason:    "Parceria comercial autorizada pela diretoria - cortesia 60 dias",
	}

	sub, err := subService.GrantManualSubscription(ctx, tenant.ID, grantDTO, adminUserID)
	if err != nil {
		t.Fatalf("Erro ao conceder assinatura manual: %v", err)
	}

	if sub.Status != domain.SubscriptionStatusActive {
		t.Errorf("Esperava status ACTIVE após concessão manual, obteve: %s", sub.Status)
	}
	if sub.Origin != domain.SubscriptionOriginManual {
		t.Errorf("Esperava Origin MANUAL, obteve: %s", sub.Origin)
	}
	if sub.GrantedByUserID == nil || *sub.GrantedByUserID != adminUserID {
		t.Errorf("GrantedByUserID incorreto: %+v", sub.GrantedByUserID)
	}
	if sub.ManualGrantReason != grantDTO.Reason {
		t.Errorf("ManualGrantReason divergente: %s", sub.ManualGrantReason)
	}

	// 5. Consulta histórico de auditoria
	logs, err := subService.ListAuditLogs(ctx, sub.ID)
	if err != nil {
		t.Fatalf("Erro ao listar logs: %v", err)
	}
	if len(logs) == 0 {
		t.Fatalf("Esperava logs de auditoria gravados")
	}
	if logs[0].Action != domain.ActionManualGrant {
		t.Errorf("Esperava Action MANUAL_GRANT no log, obteve: %s", logs[0].Action)
	}
	if logs[0].PerformedByID == nil || *logs[0].PerformedByID != adminUserID {
		t.Errorf("PerformedByID incorreto no log: %+v", logs[0].PerformedByID)
	}
}

func TestSubscriptionService_CheckAndExpireSubscriptions(t *testing.T) {
	_, repo := setupTestDB(t)
	subService := service.NewSubscriptionService(repo, nil)
	planService := service.NewPlanService(repo)
	ctx := context.Background()

	// Cria plano e tenant
	isActive := true
	plan, err := planService.Create(ctx, domain.CreatePlanDTO{
		Name:         "Plano Curto",
		Price:        50.00,
		BillingCycle: domain.CycleMonthly,
		IsActive:     &isActive,
	})
	if err != nil {
		t.Fatalf("Erro ao criar plano: %v", err)
	}

	tenantID := uuid.New()
	tenant := domain.Tenant{
		ID:        tenantID,
		Name:      "Barbearia Expirada Teste",
		Slug:      "barbearia-expirada-teste",
		Phone:     "(11) 96666-5555",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.CreateTenant(ctx, &tenant); err != nil {
		t.Fatalf("Erro ao criar tenant: %v", err)
	}

	// Concede com validade já vencida no passado (ontem)
	pastDate := time.Now().AddDate(0, 0, -1)
	grantDTO := domain.GrantManualSubscriptionDTO{
		PlanID:    plan.ID,
		ExpiresAt: &pastDate,
		Reason:    "Período de testes de 1 dia",
	}

	adminUserID := uuid.New()
	sub, err := subService.GrantManualSubscription(ctx, tenant.ID, grantDTO, adminUserID)
	if err != nil {
		t.Fatalf("Erro ao conceder assinatura: %v", err)
	}

	// Executa rotina de expiração
	expiredCount, err := subService.CheckAndExpireSubscriptions(ctx)
	if err != nil {
		t.Fatalf("Erro ao checar expirações: %v", err)
	}
	if expiredCount != 1 {
		t.Errorf("Esperava 1 assinatura expirada, obteve: %d", expiredCount)
	}

	// Verifica se status mudou para EXPIRED
	subRefreshed, err := subService.GetByID(ctx, sub.ID)
	if err != nil {
		t.Fatalf("Erro ao buscar assinatura: %v", err)
	}
	if subRefreshed.Status != domain.SubscriptionStatusExpired {
		t.Errorf("Esperava status EXPIRED, obteve: %s", subRefreshed.Status)
	}
}
