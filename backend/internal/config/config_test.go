package config_test

import (
	"os"
	"testing"

	"github.com/sistema-agendamento/backend/internal/config"
)

func TestConfigLoad_Defaults(t *testing.T) {
	// Limpa variáveis relevantes
	os.Unsetenv("SEED_DEMO")
	os.Unsetenv("INITIAL_ADMIN_EMAIL")
	os.Unsetenv("INITIAL_ADMIN_PASSWORD")

	cfg := config.Load()
	if cfg == nil {
		t.Fatal("esperado cfg não nulo")
	}

	if cfg.InitialAdminEmail == "" {
		t.Error("esperado InitialAdminEmail padrão preenchido")
	}
	if cfg.InitialAdminPassword == "" {
		t.Error("esperado InitialAdminPassword padrão preenchido")
	}
}

func TestConfigLoad_CustomEnv(t *testing.T) {
	os.Setenv("SEED_DEMO", "true")
	os.Setenv("INITIAL_ADMIN_EMAIL", "custom-admin@teste.com")
	os.Setenv("INITIAL_ADMIN_PASSWORD", "minhasenhasupersegura123")
	defer func() {
		os.Unsetenv("SEED_DEMO")
		os.Unsetenv("INITIAL_ADMIN_EMAIL")
		os.Unsetenv("INITIAL_ADMIN_PASSWORD")
	}()

	cfg := config.Load()
	if !cfg.SeedDemo {
		t.Errorf("esperado SeedDemo=true, obteve %v", cfg.SeedDemo)
	}
	if cfg.InitialAdminEmail != "custom-admin@teste.com" {
		t.Errorf("esperado custom-admin@teste.com, obteve %s", cfg.InitialAdminEmail)
	}
	if cfg.InitialAdminPassword != "minhasenhasupersegura123" {
		t.Errorf("esperado minhasenhasupersegura123, obteve %s", cfg.InitialAdminPassword)
	}
}
