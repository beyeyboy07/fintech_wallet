package service

import (
	"fintech-wallet/backend/internal/worker"

	"github.com/hibiken/asynq"
)

type QueueService struct {
	client *asynq.Client
}

func NewQueueService(client *asynq.Client) *QueueService {
	return &QueueService{
		client: client,
	}
}

func (s *QueueService) EnqueueTopUp(
	userID uint,
	transactionID uint,
	amount int64,
) error {
	return worker.EnqueueTopUp(
		s.client,
		userID,
		transactionID,
		amount,
	)
}
