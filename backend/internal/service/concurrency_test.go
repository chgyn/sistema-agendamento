package service_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/internal/service"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, *postgres.Repository) {
	// Cria banco SQLite em memória com busy_timeout para serializar requisições concorrentes em teste
	db, err := gorm.Open(sqlite.Open("file:memtest?mode=memory&cache=shared&_busy_timeout=10000"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("falha ao abrir banco de teste: %v", err)
	}

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
	)
	if err != nil {
		t.Fatalf("falha ao migrar banco de teste: %v", err)
	}

	repo := postgres.NewRepository(db)
	return db, repo
}

// TestDoubleBookingConcurrencyValidation testa 20 clientes tentando agendar exatamente no mesmo profissional e mesmo horário simultaneamente
func TestDoubleBookingConcurrencyValidation(t *testing.T) {
	ctx := context.Background()
	db, repo := setupTestDB(t)

	tenantID := uuid.New()
	tenant := domain.Tenant{
		ID:       tenantID,
		Slug:     "barbearia-teste",
		Name:     "Barbearia Teste Concorrência",
		IsActive: true,
	}
	db.Create(&tenant)

	svcID := uuid.New()
	svc := domain.Service{
		ID:              svcID,
		TenantID:        tenantID,
		Name:            "Corte Teste",
		DurationMinutes: 30,
		Price:           50.0,
		IsActive:        true,
	}
	db.Create(&svc)

	proID := uuid.New()
	pro := domain.Professional{
		ID:       proID,
		TenantID: tenantID,
		Name:     "Barbeiro Concorrente",
		IsActive: true,
	}
	db.Create(&pro)

	aptSvc := service.NewAppointmentService(repo, nil)

	// Horário alvo
	targetStartTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)

	concurrencyCount := 20
	var wg sync.WaitGroup
	var successCount int64
	var conflictCount int64

	startBarrier := make(chan struct{})

	for i := 0; i < concurrencyCount; i++ {
		wg.Add(1)
		clientNum := i + 1
		go func(num int) {
			defer wg.Done()
			<-startBarrier // Aguarda todos estarem prontos para disparar no mesmo milissegundo

			dto := service.CreateAppointmentDTO{
				TenantID:       tenantID,
				ServiceID:      svcID,
				ProfessionalID: proID,
				StartAt:        targetStartTime,
				CustomerName:   fmt.Sprintf("Cliente Concorrente %d", num),
				CustomerPhone:  fmt.Sprintf("1199999%04d", num),
				CustomerEmail:  fmt.Sprintf("cliente%d@teste.com", num),
			}

			_, err := aptSvc.CreateAppointment(ctx, dto)
			if err == nil {
				atomic.AddInt64(&successCount, 1)
			} else if err == domain.ErrSlotAlreadyBooked {
				atomic.AddInt64(&conflictCount, 1)
			} else {
				t.Errorf("Erro inesperado do cliente %d: %v", num, err)
			}
		}(clientNum)
	}

	// Libera todas as 20 goroutines simultaneamente
	close(startBarrier)
	wg.Wait()

	t.Logf("Resultado da simulação de concorrência:")
	t.Logf("   -> Sucessos (Agendamentos confirmados): %d", successCount)
	t.Logf("   -> Conflitos bloqueados (ErrSlotAlreadyBooked): %d", conflictCount)

	if successCount != 1 {
		t.Fatalf("FALHA DE SEGURANÇA: Esperado exatamente 1 sucesso, obtido %d (Duplo agendamento detectado!)", successCount)
	}

	if conflictCount != int64(concurrencyCount-1) {
		t.Fatalf("Esperado %d conflitos bloqueados, obtido %d", concurrencyCount-1, conflictCount)
	}

	t.Log("✅ Teste de concorrência passou: Apenas 1 agendamento foi permitido no horário disputado!")
}

// TestMultiTenantDataIsolation valida que a Barbearia B não pode ver nem colidir com a Barbearia A
func TestMultiTenantDataIsolation(t *testing.T) {
	ctx := context.Background()
	db, repo := setupTestDB(t)

	// Tenant A
	tenantAID := uuid.New()
	db.Create(&domain.Tenant{ID: tenantAID, Slug: "tenant-a", Name: "Barbearia A", IsActive: true})
	svcA := domain.Service{ID: uuid.New(), TenantID: tenantAID, Name: "Corte A", DurationMinutes: 30, Price: 40}
	db.Create(&svcA)

	// Tenant B
	tenantBID := uuid.New()
	db.Create(&domain.Tenant{ID: tenantBID, Slug: "tenant-b", Name: "Salão B", IsActive: true})
	svcB := domain.Service{ID: uuid.New(), TenantID: tenantBID, Name: "Corte B", DurationMinutes: 30, Price: 60}
	db.Create(&svcB)

	// Tenant A listando seus serviços
	servicesA, err := repo.ListServices(ctx, tenantAID, false)
	if err != nil {
		t.Fatalf("erro ao listar serviços: %v", err)
	}
	if len(servicesA) != 1 || servicesA[0].Name != "Corte A" {
		t.Fatalf("Isolamento violado: Tenant A retornou %v", servicesA)
	}

	// Tenant B listando seus serviços
	servicesB, err := repo.ListServices(ctx, tenantBID, false)
	if err != nil {
		t.Fatalf("erro ao listar serviços: %v", err)
	}
	if len(servicesB) != 1 || servicesB[0].Name != "Corte B" {
		t.Fatalf("Isolamento violado: Tenant B retornou %v", servicesB)
	}

	// Tentativa do Tenant B buscar o serviço pertencente ao Tenant A
	_, err = repo.GetServiceByID(ctx, tenantBID, svcA.ID)
	if err != domain.ErrServiceNotFound {
		t.Fatalf("Isolamento violado: Tenant B conseguiu carregar serviço do Tenant A!")
	}

	t.Log("✅ Teste de isolamento multi-tenant passou: Dados de estabelecimentos completamente isolados!")
}
