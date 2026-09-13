package main

import (
	"log"
	"time"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/repository"
	"fintech-wallet/backend/internal/service"
	"fintech-wallet/backend/internal/worker"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	outboxRepository := repository.NewOutboxRepository(db)

	queueClient := worker.NewClient(
		cfg.RedisHost + ":" + cfg.RedisPort,
	)
	defer queueClient.Close()

	outboxService := service.NewOutboxService(
		outboxRepository,
		queueClient,
	)

	log.Println("Outbox dispatcher started")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		err := outboxService.Dispatch(100)

		if err != nil {
			log.Printf(
				"Outbox dispatch failed: %v",
				err,
			)
		}

		<-ticker.C
	}
}
