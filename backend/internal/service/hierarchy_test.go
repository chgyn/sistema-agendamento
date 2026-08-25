package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/pkg/hash"
	"github.com/sistema-agendamento/backend/pkg/jwt"
)

func TestHierarchyAndMultiTenantIsolation(t *testing.T) {
	ctx := context.Background()
	_, repo := setupTestDB(t)
	jwtSvc := jwt.NewJWTService("secret-key-test-1234567890123456", 24)

	passHash, _ := hash.HashPassword("senha123")

	// 1. Cria Tenant A e Tenant B
	tenantAID := uuid.New()
	tenantA := domain.Tenant{
		ID:        tenantAID,
		Slug:      "barbearia-a",
		Name:      "Barbearia A",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.CreateTenant(ctx, &tenantA); err != nil {
		t.Fatalf("falha ao criar Tenant A: %v", err)
	}

	tenantBID := uuid.New()
	tenantB := domain.Tenant{
		ID:        tenantBID,
		Slug:      "salao-b",
		Name:      "Salão B",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.CreateTenant(ctx, &tenantB); err != nil {
		t.Fatalf("falha ao criar Tenant B: %v", err)
	}

	// 2. Cria Usuários:
	// - Global Admin (sem TenantID)
	// - Admin Tenant A
	// - Operador Tenant A
	// - Admin Tenant B
	globalAdmin := domain.User{
		ID:           uuid.New(),
		TenantID:     nil,
		Name:         "Admin Geral",
		Email:        "global@plataforma.com",
		PasswordHash: passHash,
		Role:         domain.RoleAdminGlobal,
		IsActive:     true,
	}
	if err := repo.CreateUser(ctx, &globalAdmin); err != nil {
		t.Fatalf("falha ao criar Global Admin: %v", err)
	}

	adminA := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenantAID,
		Name:         "Admin A",
		Email:        "admin@tenant-a.com",
		PasswordHash: passHash,
		Role:         domain.RoleAdminTenant,
		IsActive:     true,
	}
	if err := repo.CreateUser(ctx, &adminA); err != nil {
		t.Fatalf("falha ao criar Admin A: %v", err)
	}

	opA := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenantAID,
		Name:         "Operador A",
		Email:        "operador@tenant-a.com",
		PasswordHash: passHash,
		Role:         domain.RoleOperator,
		IsActive:     true,
	}
	if err := repo.CreateUser(ctx, &opA); err != nil {
		t.Fatalf("falha ao criar Operador A: %v", err)
	}

	adminB := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenantBID,
		Name:         "Admin B",
		Email:        "admin@tenant-b.com",
		PasswordHash: passHash,
		Role:         domain.RoleAdminTenant,
		IsActive:     true,
	}
	if err := repo.CreateUser(ctx, &adminB); err != nil {
		t.Fatalf("falha ao criar Admin B: %v", err)
	}

	// 3. Teste de Validação de Tokens JWT
	t.Run("Validar Claims JWT para ADMIN_GLOBAL e ADMIN_TENANT", func(t *testing.T) {
		tokenGlobal, err := jwtSvc.GenerateToken(&globalAdmin)
		if err != nil {
			t.Fatalf("erro ao gerar token global: %v", err)
		}
		claimsGlobal, err := jwtSvc.ValidateToken(tokenGlobal)
		if err != nil {
			t.Fatalf("erro ao validar token global: %v", err)
		}
		if claimsGlobal.Role != domain.RoleAdminGlobal || claimsGlobal.TenantID != nil {
			t.Errorf("claims global inválidas: %+v", claimsGlobal)
		}

		tokenA, err := jwtSvc.GenerateToken(&adminA)
		if err != nil {
			t.Fatalf("erro ao gerar token adminA: %v", err)
		}
		claimsA, err := jwtSvc.ValidateToken(tokenA)
		if err != nil {
			t.Fatalf("erro ao validar token adminA: %v", err)
		}
		if claimsA.Role != domain.RoleAdminTenant || claimsA.TenantID == nil || *claimsA.TenantID != tenantAID {
			t.Errorf("claims adminA inválidas: %+v", claimsA)
		}
	})

	// 4. Teste de Isolamento de Usuários por Tenant
	t.Run("ListUsersByTenant deve retornar somente usuários do próprio Tenant", func(t *testing.T) {
		usersA, err := repo.ListUsersByTenant(ctx, tenantAID)
		if err != nil {
			t.Fatalf("erro ao listar usuários do Tenant A: %v", err)
		}
		if len(usersA) != 2 {
			t.Errorf("esperado 2 usuários no Tenant A, obtido %d", len(usersA))
		}
		for _, u := range usersA {
			if u.TenantID == nil || *u.TenantID != tenantAID {
				t.Errorf("vazamento de dados: usuário %s com TenantID diferente encontrado no Tenant A", u.Email)
			}
		}

		usersB, err := repo.ListUsersByTenant(ctx, tenantBID)
		if err != nil {
			t.Fatalf("erro ao listar usuários do Tenant B: %v", err)
		}
		if len(usersB) != 1 {
			t.Errorf("esperado 1 usuário no Tenant B, obtido %d", len(usersB))
		}
		if usersB[0].Email != "admin@tenant-b.com" {
			t.Errorf("esperado admin@tenant-b.com, obtido %s", usersB[0].Email)
		}
	})

	// 5. Teste de ADMIN_GLOBAL listando todos os usuários e tenants
	t.Run("ADMIN_GLOBAL pode listar todos os usuários e filtrar por Tenant", func(t *testing.T) {
		allUsers, err := repo.ListAllUsers(ctx, nil, "")
		if err != nil {
			t.Fatalf("erro ao listar todos os usuários: %v", err)
		}
		// 1 Global + 2 Tenant A + 1 Tenant B = 4
		if len(allUsers) < 4 {
			t.Errorf("esperado pelo menos 4 usuários globais, obtido %d", len(allUsers))
		}

		// Filtrar por Tenant B
		filteredB, err := repo.ListAllUsers(ctx, &tenantBID, "")
		if err != nil {
			t.Fatalf("erro ao listar filtrado: %v", err)
		}
		if len(filteredB) != 1 || filteredB[0].Email != "admin@tenant-b.com" {
			t.Errorf("filtro de tenant falhou: %+v", filteredB)
		}
	})

	// 6. Teste de Criação de Tenant com Admin Inicial Atômico
	t.Run("CreateTenantWithAdmin deve criar Tenant e Admin associado", func(t *testing.T) {
		newTenant := domain.Tenant{
			ID:        uuid.New(),
			Slug:      "nova-barbearia-vip",
			Name:      "Barbearia VIP",
			Phone:     "(11) 91234-5678",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		newAdmin := domain.User{
			ID:           uuid.New(),
			Name:         "Admin VIP",
			Email:        "admin@barbeariavip.com",
			PasswordHash: passHash,
			Role:         domain.RoleAdminTenant,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := repo.CreateTenantWithAdmin(ctx, &newTenant, &newAdmin)
		if err != nil {
			t.Fatalf("erro ao criar tenant com admin: %v", err)
		}

		createdTenant, err := repo.GetTenantByID(ctx, newTenant.ID)
		if err != nil {
			t.Fatalf("tenant criado não encontrado: %v", err)
		}
		if createdTenant.Slug != "nova-barbearia-vip" {
			t.Errorf("slug incorreto: %s", createdTenant.Slug)
		}

		createdAdmin, err := repo.GetUserByEmail(ctx, "admin@barbeariavip.com")
		if err != nil {
			t.Fatalf("admin criado não encontrado: %v", err)
		}
		if createdAdmin.TenantID == nil || *createdAdmin.TenantID != newTenant.ID {
			t.Errorf("admin não foi vinculado ao Tenant correto")
		}
	})

	// 7. Teste de ADMIN_GLOBAL listando outros Global Admins
	t.Run("ListGlobalAdmins deve retornar apenas Administradores Gerais", func(t *testing.T) {
		globalAdmins, err := repo.ListGlobalAdmins(ctx)
		if err != nil {
			t.Fatalf("erro ao listar global admins: %v", err)
		}
		if len(globalAdmins) < 1 {
			t.Fatalf("nenhum global admin retornado")
		}
		for _, ga := range globalAdmins {
			if ga.Role != domain.RoleAdminGlobal {
				t.Errorf("usuário com role %s retornado na lista de global admins", ga.Role)
			}
		}
	})
}
