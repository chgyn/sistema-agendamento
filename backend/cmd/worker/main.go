package main

import (
	"log"

	"github.com/sistema-agendamento/backend/internal/config"
	"github.com/sistema-agendamento/backend/internal/queue"
)

func main() {
	log.Println("👷 Inicializando Sistema de Agendamento — Worker Daemon (Asynq + Redis)...")

	cfg := config.Load()

	if err := runWorker(cfg); err != nil {
		log.Fatalf("❌ Erro no Worker Asynq: %v", err)
	}
}

func runWorker(cfg *config.Config) error {
	_, err := queue.StartWorkerServer(cfg)
	return err
}
