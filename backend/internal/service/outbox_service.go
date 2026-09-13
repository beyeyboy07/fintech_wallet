package service

import (
	"encoding/json"
	"errors"
	"time"

	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/repository"

	"github.com/hibiken/asynq"
)

type OutboxService struct {
	outboxRepository *repository.OutboxRepository
	queueClient      *asynq.Client
}

func NewOutboxService(
	outboxRepository *repository.OutboxRepository,
	queueClient *asynq.Client,
) *OutboxService {
	return &OutboxService{
		outboxRepository: outboxRepository,
		queueClient:      queueClient,
	}
}

func (s *OutboxService) Dispatch(limit int) error {

	err := s.outboxRepository.ResetStaleProcessing(
		10 * time.Minute,
	)

	if err != nil {
		return err
	}

	events, err := s.outboxRepository.FindPending(limit)
	if err != nil {
		return err
	}

	for _, event := range events {

		claimedEvent, err :=
			s.outboxRepository.ClaimPending(event.ID)

		if err != nil {
			continue
		}

		if claimedEvent == nil {
			continue
		}

		err = s.dispatchEvent(*claimedEvent)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *OutboxService) dispatchEvent(
	event model.OutboxEvent,
) error {

	if event.EventType == "" {
		return errors.New("event type is empty")
	}

	task := asynq.NewTask(
		event.EventType,
		[]byte(event.Payload),
	)

	_, err := s.queueClient.Enqueue(task)
	if err != nil {
		return err
	}

	return s.outboxRepository.UpdateStatus(
		event.ID,
		"PROCESSED",
	)
}

func decodePayload(
	payload string,
	target interface{},
) error {
	return json.Unmarshal(
		[]byte(payload),
		target,
	)
}
