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

func TestConfigLoad_GoogleCalendarURLs(t *testing.T) {
	os.Unsetenv("APP_PUBLIC_URL")
	os.Unsetenv("FRONTEND_URL")
	os.Unsetenv("GOOGLE_OAUTH_REDIRECT_URI")
	os.Setenv("PORT", "8090")
	os.Setenv("APP_ENV", "development")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("APP_ENV")
	}()

	cfg := config.Load()
	if cfg.AppPublicURL != "http://localhost:8090" {
		t.Errorf("esperado AppPublicURL=http://localhost:8090, obteve %s", cfg.AppPublicURL)
	}
	if cfg.FrontendURL != "http://localhost:5173" {
		t.Errorf("esperado FrontendURL=http://localhost:5173 em development, obteve %s", cfg.FrontendURL)
	}
	expectedRedirect := "http://localhost:8090/api/v1/integrations/google-calendar/callback"
	if cfg.GoogleOAuthRedirectURI != expectedRedirect {
		t.Errorf("esperado GoogleOAuthRedirectURI=%s, obteve %s", expectedRedirect, cfg.GoogleOAuthRedirectURI)
	}

	// Em produção com APP_PUBLIC_URL customizada
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_PUBLIC_URL", "https://agenda.exemplo.com.br/")
	defer func() {
		os.Unsetenv("APP_PUBLIC_URL")
	}()

	prodCfg := config.Load()
	if prodCfg.FrontendURL != "https://agenda.exemplo.com.br/" {
		t.Errorf("esperado FrontendURL herdando AppPublicURL em produção, obteve %s", prodCfg.FrontendURL)
	}
	if prodCfg.GoogleOAuthRedirectURI != "https://agenda.exemplo.com.br/api/v1/integrations/google-calendar/callback" {
		t.Errorf("redirect URI inesperada em produção: %s", prodCfg.GoogleOAuthRedirectURI)
	}
}
