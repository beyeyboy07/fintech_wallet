package main

import (
	"context"
	"encoding/json"
	"log"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/redis"
	"fintech-wallet/backend/internal/repository"
	"fintech-wallet/backend/internal/service"
	"fintech-wallet/backend/internal/worker"

	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(cfg)

	redisService := service.NewRedisService(redisClient)

	walletRepository := repository.NewWalletRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	outboxRepository := repository.NewOutboxRepository(db)

	walletService := service.NewWalletService(
		walletRepository,
		transactionRepository,
		outboxRepository,
		redisService,
		nil,
		db,
	)

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: cfg.RedisHost + ":" + cfg.RedisPort,
		},
		asynq.Config{
			Concurrency: 5,
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(worker.TaskTypeTopUp, func(
		ctx context.Context,
		task *asynq.Task,
	) error {

		var payload worker.TopUpPayload

		err := json.Unmarshal(
			task.Payload(),
			&payload,
		)

		if err != nil {
			return err
		}

		err = walletService.ProcessTopUp(
			payload.TransactionID,
			payload.UserID,
			payload.Amount,
		)

		if err != nil {
			log.Printf(
				"Failed to process top up: transaction_id=%d error=%v",
				payload.TransactionID,
				err,
			)

			return err
		}

		log.Printf(
			"Top up processed successfully: transaction_id=%d",
			payload.TransactionID,
		)

		return nil
	})

	log.Println("Worker started")

	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}
}
