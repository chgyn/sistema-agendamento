package database

import (
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/pkg/hash"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if cfg.DBDriver == "sqlite" {
		dialector = sqlite.Open("agendamento.db")
	} else {
		var dsn string
		if cfg.DBURL != "" {
			dsn = cfg.DBURL
		} else {
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
		}
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		// Se falhar a conexão com Postgres e estivermos em dev, tenta fallback para sqlite local
		if cfg.DBDriver == "postgres" && cfg.Environment == "development" {
			log.Printf("⚠️ Falha ao conectar ao Postgres (%v). Inicializando SQLite local de desenvolvimento...", err)
			dialector = sqlite.Open("agendamento_dev.db")
			db, err = gorm.Open(dialector, &gorm.Config{
				Logger: logger.Default.LogMode(logger.Warn),
			})
			if err != nil {
				return nil, fmt.Errorf("falha ao abrir banco de fallback: %w", err)
			}
		} else {
			return nil, err
		}
	}

	// Auto Migration de todas as tabelas
	err = db.AutoMigrate(
		&domain.Tenant{},
		&domain.User{},
		&domain.Professional{},
		&domain.Service{},
		&domain.ProfessionalService{},
		&domain.WorkingHour{},
		&domain.AvailabilityException{},
		&domain.Customer{},
		&domain.Appointment{},
		&domain.AppointmentStatusHistory{},
		&domain.AuditLog{},
		&domain.Plan{},
		&domain.Subscription{},
		&domain.SubscriptionInvoice{},
	)
	if err != nil {
		return nil, fmt.Errorf("falha na migração do banco: %w", err)
	}

	log.Println("✅ Banco de dados conectado e tabelas migradas com sucesso.")
	return db, nil
}

// SeedInitialData popula dados de teste ricos para demonstração imediata
func SeedInitialData(db *gorm.DB) error {
	// Senha padrão demo: admin123
	passHash, _ := hash.HashPassword("admin123")

	// 0. Garante existência do Administrador Geral da Plataforma
	var globalAdminCount int64
	db.Model(&domain.User{}).Where("email = ?", "admin@plataforma.com").Count(&globalAdminCount)
	if globalAdminCount == 0 {
		globalAdmin := domain.User{
			ID:           uuid.New(),
			TenantID:     nil,
			Name:         "Administrador Geral da Plataforma",
			Email:        "admin@plataforma.com",
			PasswordHash: passHash,
			Role:         domain.RoleAdminGlobal,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := db.Create(&globalAdmin).Error; err != nil {
			log.Printf("Aviso ao criar admin global: %v", err)
		} else {
			log.Println("👑 Administrador Geral inicial criado: admin@plataforma.com")
		}
	}

	// 0.1 Garante existência de Planos Comerciais Padrão
	var planCount int64
	db.Model(&domain.Plan{}).Count(&planCount)
	var defaultPlanID uuid.UUID
	if planCount == 0 {
		planBasicoID := uuid.New()
		defaultPlanID = planBasicoID
		planProID := uuid.New()
		planPremiumID := uuid.New()

		plans := []domain.Plan{
			{
				ID:               planBasicoID,
				Name:             "Plano Starter",
				Description:      "Ideal para profissionais autônomos e barbearias individuais que buscam praticidade.",
				Price:            49.90,
				BillingCycle:     domain.CycleMonthly,
				MaxProfessionals: 1,
				MaxServices:      10,
				Features:         `["1 Profissional", "Até 10 Serviços", "Agendamento Online 24/7", "Lembretes no WhatsApp", "Painel Básico"]`,
				IsActive:         true,
				SortOrder:        1,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			{
				ID:               planProID,
				Name:             "Plano Profissional",
				Description:      "O plano mais popular para estabelecimentos em crescimento com equipe.",
				Price:            99.90,
				BillingCycle:     domain.CycleMonthly,
				MaxProfessionals: 5,
				MaxServices:      30,
				Features:         `["Até 5 Profissionais", "Até 30 Serviços", "Agendamento Online 24/7", "Lembretes Automáticos", "Relatórios Financeiros", "Suporte Prioritário"]`,
				IsActive:         true,
				SortOrder:        2,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			{
				ID:               planPremiumID,
				Name:             "Plano Scale / VIP",
				Description:      "Para salões de grande porte e redes com múltiplos profissionais e alta demanda.",
				Price:            189.90,
				BillingCycle:     domain.CycleMonthly,
				MaxProfessionals: 0, // Ilimitado
				MaxServices:      0, // Ilimitado
				Features:         `["Profissionais Ilimitados", "Serviços Ilimitados", "Personalização Completa", "Taxa Zero por Agendamento", "API & Webhooks", "Gerente de Conta"]`,
				IsActive:         true,
				SortOrder:        3,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}

		for _, p := range plans {
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Aviso ao criar plano padrão: %v", err)
			}
		}
		log.Println("💎 Planos de assinatura padrão cadastrados com sucesso!")
	} else {
		var firstPlan domain.Plan
		db.First(&firstPlan)
		defaultPlanID = firstPlan.ID
	}

	var count int64
	db.Model(&domain.Tenant{}).Count(&count)
	if count > 0 {
		return nil // Já existem tenants
	}

	log.Println("🌱 Populando banco com dados de demonstração multi-tenant...")

	// 1. TENANT 1: Barbearia Dom Navalha
	tenant1ID := uuid.New()
	tenant1 := domain.Tenant{
		ID:                  tenant1ID,
		Slug:                "dom-navalha",
		Name:                "Barbearia Dom Navalha",
		Document:            "12.345.678/0001-90",
		Phone:               "(11) 98765-4321",
		Email:               "contato@domnavalha.com.br",
		Address:             "Av. Paulista, 1578 - Bela Vista",
		City:                "São Paulo",
		State:               "SP",
		LogoURL:             "https://images.unsplash.com/photo-1503951914875-452162b0f3f1?w=300&h=300&fit=crop",
		PrimaryColor:        "#10b981",
		SlotIntervalMinutes: 30,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	if err := db.Create(&tenant1).Error; err != nil {
		return err
	}

	// Assinatura Demo Tenant 1
	sub1NextDue := time.Now().AddDate(0, 1, 0)
	sub1 := domain.Subscription{
		ID:                  uuid.New(),
		TenantID:            tenant1ID,
		PlanID:              defaultPlanID,
		AsaasCustomerID:     "cus_demo_domnavalha",
		AsaasSubscriptionID: "sub_demo_domnavalha",
		Status:              domain.SubscriptionStatusActive,
		BillingCycle:        domain.CycleMonthly,
		Price:               99.90,
		NextDueDate:         &sub1NextDue,
		PaymentMethod:       "PIX",
		PaymentURL:          "https://sandbox.asaas.com/payment/mock",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	db.Create(&sub1)

	// Usuários Tenant 1
	admin1 := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenant1ID,
		Name:         "Carlos Administrador",
		Email:        "admin@domnavalha.com",
		PasswordHash: passHash,
		Role:         domain.RoleAdminTenant,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	op1 := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenant1ID,
		Name:         "Marcos Operador",
		Email:        "operador@domnavalha.com",
		PasswordHash: passHash,
		Role:         domain.RoleOperator,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	db.Create(&admin1)
	db.Create(&op1)

	// Serviços Tenant 1
	s1 := domain.Service{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		Name:            "Corte Degradê / Moderno",
		Description:     "Corte com máquina, tesoura, alinhamento e finalização com pomada modeladora.",
		DurationMinutes: 30,
		Price:           45.00,
		Color:           "#10b981",
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	s2 := domain.Service{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		Name:            "Barba Terapia com Toalha Quente",
		Description:     "Modelagem de barba com navalha, toalha quente, óleos essenciais e pós-barba refrescante.",
		DurationMinutes: 30,
		Price:           40.00,
		Color:           "#059669",
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	s3 := domain.Service{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		Name:            "Combo VIP: Corte + Barba Terapia",
		Description:     "Experiência completa de cuidado capilar e barba com atendimento de luxo e café/cerveja cortesia.",
		DurationMinutes: 60,
		Price:           75.00,
		Color:           "#f59e0b",
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	s4 := domain.Service{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		Name:            "Acabamento & Pezinho",
		Description:     "Alinhamento dos contornos do cabelo e nuca com navalha.",
		DurationMinutes: 15,
		Price:           20.00,
		Color:           "#3b82f6",
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	s5 := domain.Service{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		Name:            "Pigmentação de Barba",
		Description:     "Disfarce de falhas e realce da cor da barba com produtos hipoalergênicos.",
		DurationMinutes: 30,
		Price:           35.00,
		Color:           "#8b5cf6",
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	db.Create(&s1)
	db.Create(&s2)
	db.Create(&s3)
	db.Create(&s4)
	db.Create(&s5)

	// Profissionais Tenant 1
	p1ID := uuid.New()
	p1 := domain.Professional{
		ID:        p1ID,
		TenantID:  tenant1ID,
		Name:      "Carlos Navalha",
		Email:     "carlos@domnavalha.com",
		Phone:     "(11) 99111-2222",
		Title:     "Barbeiro Master",
		Specialty: "Especialista em Degradê Navalhado, Pompadour e Barboterapia",
		Bio:       "Mais de 10 anos de experiência transformando o visual masculino com técnicas clássicas e contemporâneas.",
		AvatarURL: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=200&h=200&fit=crop",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	p2ID := uuid.New()
	p2 := domain.Professional{
		ID:        p2ID,
		TenantID:  tenant1ID,
		Name:      "Lucas Tesoura",
		Email:     "lucas@domnavalha.com",
		Phone:     "(11) 99222-3333",
		Title:     "Barbeiro Especialista",
		Specialty: "Cortes Clássicos, Texturização na Tesoura e Barba Alinhada",
		Bio:       "Formado pelas melhores academias de barbearia do Brasil, focado em precisão e detalhes.",
		AvatarURL: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=200&h=200&fit=crop",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	p3ID := uuid.New()
	p3 := domain.Professional{
		ID:        p3ID,
		TenantID:  tenant1ID,
		Name:      "Diego Visagista",
		Email:     "diego@domnavalha.com",
		Phone:     "(11) 99333-4444",
		Title:     "Visagista & Estilista",
		Specialty: "Harmonização de Barba e Cabelo de acordo com o formato do rosto",
		Bio:       "Consultoria personalizada para você encontrar o estilo ideal para seu dia a dia profissional.",
		AvatarURL: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&h=200&fit=crop",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(&p1)
	db.Create(&p2)
	db.Create(&p3)

	// Associar Serviços aos Profissionais
	proServices := []domain.ProfessionalService{
		{TenantID: tenant1ID, ProfessionalID: p1ID, ServiceID: s1.ID},
		{TenantID: tenant1ID, ProfessionalID: p1ID, ServiceID: s2.ID},
		{TenantID: tenant1ID, ProfessionalID: p1ID, ServiceID: s3.ID},
		{TenantID: tenant1ID, ProfessionalID: p1ID, ServiceID: s4.ID},
		{TenantID: tenant1ID, ProfessionalID: p1ID, ServiceID: s5.ID},

		{TenantID: tenant1ID, ProfessionalID: p2ID, ServiceID: s1.ID},
		{TenantID: tenant1ID, ProfessionalID: p2ID, ServiceID: s2.ID},
		{TenantID: tenant1ID, ProfessionalID: p2ID, ServiceID: s3.ID},
		{TenantID: tenant1ID, ProfessionalID: p2ID, ServiceID: s4.ID},

		{TenantID: tenant1ID, ProfessionalID: p3ID, ServiceID: s1.ID},
		{TenantID: tenant1ID, ProfessionalID: p3ID, ServiceID: s2.ID},
		{TenantID: tenant1ID, ProfessionalID: p3ID, ServiceID: s3.ID},
		{TenantID: tenant1ID, ProfessionalID: p3ID, ServiceID: s5.ID},
	}
	for _, ps := range proServices {
		db.Create(&ps)
	}

	// Criar Horários de Trabalho (Segunda a Sábado, 08:00 às 19:00 com almoço 12:00 às 13:00)
	professionals := []uuid.UUID{p1ID, p2ID, p3ID}
	for _, pid := range professionals {
		// Segunda (1) a Sexta (5)
		for day := 1; day <= 5; day++ {
			wh := domain.WorkingHour{
				ID:             uuid.New(),
				TenantID:       tenant1ID,
				ProfessionalID: pid,
				DayOfWeek:      day,
				StartTime:      "08:00",
				EndTime:        "19:00",
				BreakStart:     "12:00",
				BreakEnd:       "13:00",
				IsActive:       true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			db.Create(&wh)
		}
		// Sábado (6) - 08:00 às 17:00
		whSat := domain.WorkingHour{
			ID:             uuid.New(),
			TenantID:       tenant1ID,
			ProfessionalID: pid,
			DayOfWeek:      6,
			StartTime:      "08:00",
			EndTime:        "17:00",
			BreakStart:     "12:00",
			BreakEnd:       "13:00",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		db.Create(&whSat)
	}

	// Clientes de Teste
	c1 := domain.Customer{
		ID:          uuid.New(),
		TenantID:    tenant1ID,
		Name:        "Guilherme Santos",
		Phone:       "(11) 98888-7777",
		Email:       "guilherme.santos@email.com",
		TotalVisits: 5,
		TotalSpent:  320.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	c2 := domain.Customer{
		ID:          uuid.New(),
		TenantID:    tenant1ID,
		Name:        "Rodrigo Oliveira",
		Phone:       "(11) 97777-6666",
		Email:       "rodrigo.oliveira@email.com",
		TotalVisits: 3,
		TotalSpent:  185.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	c3 := domain.Customer{
		ID:          uuid.New(),
		TenantID:    tenant1ID,
		Name:        "Felipe Silveira",
		Phone:       "(11) 96666-5555",
		Email:       "felipe.silveira@email.com",
		TotalVisits: 8,
		TotalSpent:  560.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(&c1)
	db.Create(&c2)
	db.Create(&c3)

	// Agendamentos para hoje e próximos dias
	today := time.Now().Truncate(24 * time.Hour)
	apt1 := domain.Appointment{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		ProfessionalID:  p1ID,
		ServiceID:       s3.ID,
		CustomerID:      c1.ID,
		StartAt:         today.Add(9 * time.Hour),
		EndAt:           today.Add(10 * time.Hour),
		DurationMinutes: 60,
		TotalPrice:      75.00,
		Status:          domain.StatusConfirmed,
		Notes:           "Cliente prefere toalha bem quente na barba",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	apt2 := domain.Appointment{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		ProfessionalID:  p2ID,
		ServiceID:       s1.ID,
		CustomerID:      c2.ID,
		StartAt:         today.Add(10 * time.Hour),
		EndAt:           today.Add(10*time.Hour + 30*time.Minute),
		DurationMinutes: 30,
		TotalPrice:      45.00,
		Status:          domain.StatusConfirmed,
		Notes:           "Corte degradê navalhado nas laterais",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	apt3 := domain.Appointment{
		ID:              uuid.New(),
		TenantID:        tenant1ID,
		ProfessionalID:  p1ID,
		ServiceID:       s1.ID,
		CustomerID:      c3.ID,
		StartAt:         today.Add(14 * time.Hour),
		EndAt:           today.Add(14*time.Hour + 30*time.Minute),
		DurationMinutes: 30,
		TotalPrice:      45.00,
		Status:          domain.StatusConfirmed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	db.Create(&apt1)
	db.Create(&apt2)
	db.Create(&apt3)

	// 2. TENANT 2: Salão Bella Vista
	tenant2ID := uuid.New()
	tenant2 := domain.Tenant{
		ID:                  tenant2ID,
		Slug:                "bella-vista",
		Name:                "Salão Bella Vista Studio",
		Document:            "98.765.432/0001-10",
		Phone:               "(21) 97654-3210",
		Email:               "contato@bellavista.com.br",
		Address:             "Rua Visconde de Pirajá, 300 - Ipanema",
		City:                "Rio de Janeiro",
		State:               "RJ",
		LogoURL:             "https://images.unsplash.com/photo-1560066984-138dadb4c035?w=300&h=300&fit=crop",
		PrimaryColor:        "#ec4899",
		SlotIntervalMinutes: 30,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	if err := db.Create(&tenant2).Error; err != nil {
		return err
	}

	// Assinatura Demo Tenant 2
	sub2NextDue := time.Now().AddDate(0, 1, 0)
	sub2 := domain.Subscription{
		ID:                  uuid.New(),
		TenantID:            tenant2ID,
		PlanID:              defaultPlanID,
		AsaasCustomerID:     "cus_demo_bellavista",
		AsaasSubscriptionID: "sub_demo_bellavista",
		Status:              domain.SubscriptionStatusActive,
		BillingCycle:        domain.CycleMonthly,
		Price:               99.90,
		NextDueDate:         &sub2NextDue,
		PaymentMethod:       "CREDIT_CARD",
		PaymentURL:          "https://sandbox.asaas.com/payment/mock",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	db.Create(&sub2)

	admin2 := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenant2ID,
		Name:         "Juliana Administradora",
		Email:        "admin@bellavista.com",
		PasswordHash: passHash,
		Role:         domain.RoleAdminTenant,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	db.Create(&admin2)

	log.Println("✨ Dados de demonstração criados com sucesso!")
	return nil
}
