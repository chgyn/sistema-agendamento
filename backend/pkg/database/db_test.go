package database_test

import (
	"testing"

	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/pkg/database"
)

func TestDatabase_ConnectAndApplySeed_SQLiteMemory(t *testing.T) {
	cfg := &config.Config{
		DBDriver:             "sqlite",
		DBURL:                "file:test_apply_seed?mode=memory&cache=shared",
		SeedDemo:             false,
		InitialAdminEmail:    "admin@teste.com",
		InitialAdminPassword: "password123",
	}

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("falha ao conectar e migrar: %v", err)
	}

	// 1. Aplica o seed mínimo de produção
	if err := database.ApplySeed(db, cfg); err != nil {
		t.Fatalf("falha ao aplicar seed: %v", err)
	}

	// Verifica se planos foram criados
	var planCount int64
	if err := db.Model(&domain.Plan{}).Count(&planCount).Error; err != nil {
		t.Fatalf("falha ao contar planos: %v", err)
	}
	if planCount != 3 {
		t.Errorf("esperado 3 planos cadastrados, obteve %d", planCount)
	}

	// Verifica se admin inicial foi criado
	var adminCount int64
	if err := db.Model(&domain.User{}).Where("role = ? AND email = ?", domain.RoleAdminGlobal, "admin@teste.com").Count(&adminCount).Error; err != nil {
		t.Fatalf("falha ao contar admin: %v", err)
	}
	if adminCount != 1 {
		t.Errorf("esperado 1 admin global cadastrado, obteve %d", adminCount)
	}

	// 2. Idempotência: aplicar seed novamente não deve duplicar planos ou admins
	if err := database.ApplySeed(db, cfg); err != nil {
		t.Fatalf("falha na segunda execução de seed (idempotência): %v", err)
	}

	var planCount2 int64
	db.Model(&domain.Plan{}).Count(&planCount2)
	if planCount2 != 3 {
		t.Errorf("esperado que contagem de planos se mantivesse em 3, obteve %d", planCount2)
	}

	var adminCount2 int64
	db.Model(&domain.User{}).Where("role = ?", domain.RoleAdminGlobal).Count(&adminCount2)
	if adminCount2 != 1 {
		t.Errorf("esperado que contagem de admins se mantivesse em 1, obteve %d", adminCount2)
	}
}

func TestDatabase_ConnectWorker_SQLiteMemory(t *testing.T) {
	cfg := &config.Config{
		DBDriver: "sqlite",
		DBURL:    "file:test_worker?mode=memory&cache=shared",
	}

	db, err := database.ConnectWorker(cfg)
	if err != nil {
		t.Fatalf("falha ao conectar via ConnectWorker: %v", err)
	}
	if db == nil {
		t.Fatal("esperado db não nulo")
	}
}
