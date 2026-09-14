package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupWhatsAppTestDB(t *testing.T) (*gorm.DB, *postgres.Repository) {
	dbName := fmt.Sprintf("file:mem_wa_test_%d?mode=memory&cache=shared&_busy_timeout=10000", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("falha ao abrir banco sqlite em memória: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	err = db.AutoMigrate(
		&domain.Tenant{},
		&domain.Customer{},
		&domain.TenantWhatsAppConfig{},
		&domain.WhatsAppConversation{},
		&domain.WhatsAppMessage{},
	)
	if err != nil {
		t.Fatalf("falha ao executar automigrate: %v", err)
	}

	repo := postgres.NewRepository(db)
	return db, repo
}

func TestWhatsAppRepository_GetRecentMessages(t *testing.T) {
	_, repo := setupWhatsAppTestDB(t)
	ctx := context.Background()

	tenantID := uuid.New()
	customerPhone := "5511999998888"
	customerName := "Carlos Silva"

	// 1. Cria conversa
	conv, err := repo.GetOrCreateConversation(ctx, tenantID, customerPhone, customerName)
	if err != nil {
		t.Fatalf("erro ao criar conversa: %v", err)
	}
	if conv == nil || conv.ID == uuid.Nil {
		t.Fatalf("conversa inválida retornada")
	}

	// 2. Insere 15 mensagens sequenciais (msg_1 até msg_15)
	baseTime := time.Date(2026, 9, 14, 10, 0, 0, 0, time.UTC)
	for i := 1; i <= 15; i++ {
		msg := &domain.WhatsAppMessage{
			ID:             uuid.New(),
			ConversationID: conv.ID,
			TenantID:       tenantID,
			Sender:         "customer",
			Role:           "user",
			Body:           fmt.Sprintf("Mensagem %d", i),
			CreatedAt:      baseTime.Add(time.Duration(i) * time.Minute),
		}
		if err := repo.SaveWhatsAppMessage(ctx, msg); err != nil {
			t.Fatalf("erro ao salvar mensagem %d: %v", i, err)
		}
	}

	// 3. Busca com limite de 5 mensagens
	// Espera-se obter as últimas 5 mensagens (11, 12, 13, 14, 15) em ordem cronológica
	limit := 5
	recentMsgs, err := repo.GetRecentMessages(ctx, conv.ID, limit)
	if err != nil {
		t.Fatalf("erro ao buscar mensagens recentes: %v", err)
	}

	if len(recentMsgs) != 5 {
		t.Fatalf("esperava 5 mensagens, obteve %d", len(recentMsgs))
	}

	expectedTexts := []string{
		"Mensagem 11",
		"Mensagem 12",
		"Mensagem 13",
		"Mensagem 14",
		"Mensagem 15",
	}

	for i, msg := range recentMsgs {
		if msg.Body != expectedTexts[i] {
			t.Errorf("posição %d: esperava %q, mas obteve %q", i, expectedTexts[i], msg.Body)
		}
	}
}

func TestWhatsAppRepository_UpsertAndStatus(t *testing.T) {
	_, repo := setupWhatsAppTestDB(t)
	ctx := context.Background()

	tenantID := uuid.New()
	cfg := &domain.TenantWhatsAppConfig{
		ID:            uuid.New(),
		TenantID:      tenantID,
		InstanceName:  "instancia-teste",
		InstanceToken: "token-secreto-123",
		Status:        domain.WhatsAppStatusDisconnected,
	}

	if err := repo.UpsertWhatsAppConfig(ctx, cfg); err != nil {
		t.Fatalf("falha ao salvar config: %v", err)
	}

	// Atualiza QR Code
	exp := time.Now().Add(60 * time.Second)
	if err := repo.UpdateQRCode(ctx, tenantID, "data:image/png;base64,mockqr", &exp); err != nil {
		t.Fatalf("falha ao atualizar qrcode: %v", err)
	}

	loaded, err := repo.GetWhatsAppConfigByTenantID(ctx, tenantID)
	if err != nil || loaded == nil {
		t.Fatalf("falha ao recarregar config: %v", err)
	}

	if loaded.Status != domain.WhatsAppStatusQRCode {
		t.Errorf("esperava status QRCODE, obteve %s", loaded.Status)
	}
	if loaded.QRCodeBase64 != "data:image/png;base64,mockqr" {
		t.Errorf("esperava qrcode base64 salvo, obteve %s", loaded.QRCodeBase64)
	}

	// Conecta e limpa QR
	if err := repo.UpdateWhatsAppStatus(ctx, tenantID, domain.WhatsAppStatusConnected, "5511999990000"); err != nil {
		t.Fatalf("falha ao atualizar status para CONNECTED: %v", err)
	}

	loadedConnected, err := repo.GetWhatsAppConfigByInstanceToken(ctx, "token-secreto-123")
	if err != nil || loadedConnected == nil {
		t.Fatalf("falha ao buscar por instance token: %v", err)
	}
	if loadedConnected.Status != domain.WhatsAppStatusConnected {
		t.Errorf("esperava status CONNECTED, obteve %s", loadedConnected.Status)
	}
	if loadedConnected.QRCodeBase64 != "" {
		t.Errorf("esperava qrcode limpo após conexão, obteve %s", loadedConnected.QRCodeBase64)
	}
	if loadedConnected.PhoneNumber != "5511999990000" {
		t.Errorf("esperava telefone atualizado, obteve %s", loadedConnected.PhoneNumber)
	}
}
