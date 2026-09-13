package worker

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

func NewTopUpTask(
	userID uint,
	transactionID uint,
	amount int64,
) (*asynq.Task, error) {

	payload := TopUpPayload{
		UserID:        userID,
		TransactionID: transactionID,
		Amount:        amount,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TaskTypeTopUp,
		data,
	), nil
}

func EnqueueTopUp(
	client *asynq.Client,
	userID uint,
	transactionID uint,
	amount int64,
) error {

	task, err := NewTopUpTask(
		userID,
		transactionID,
		amount,
	)

	if err != nil {
		return err
	}

	_, err = client.Enqueue(task)

	return err
}
