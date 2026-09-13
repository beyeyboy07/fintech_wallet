package main

import (
	"log"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/worker"
)

func main() {
	cfg := config.Load()

	client := worker.NewClient(
		cfg.RedisHost + ":" + cfg.RedisPort,
	)

	defer client.Close()

	err := worker.EnqueueTopUp(
		client,
		2,
		1,
		50000,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Top up task enqueued successfully")
}
