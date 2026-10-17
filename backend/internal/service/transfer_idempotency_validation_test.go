package service

import (
	"errors"
	"testing"
)

func TestTransferWithIdempotency_EmptyKey(t *testing.T) {
	service := &TransferService{}

	_, err := service.TransferWithIdempotency(
		1,
		"2",
		1000,
		99,
		"Test transfer",
		"",
	)

	if !errors.Is(err, ErrIdempotencyKeyRequired) {
		t.Fatalf(
			"expected ErrIdempotencyKeyRequired, got %v",
			err,
		)
	}
}

func TestTransferWithIdempotency_InvalidAmount(t *testing.T) {
	service := &TransferService{}

	_, err := service.TransferWithIdempotency(
		1,
		"",
		2,
		0,
		"Test transfer",
		"test-key",
	)

	if !errors.Is(err, ErrInvalidTransferAmount) {
		t.Fatalf(
			"expected ErrInvalidTransferAmount, got %v",
			err,
		)
	}
}

func TestTransferWithIdempotency_SelfTransfer(t *testing.T) {
	service := &TransferService{}

	_, err := service.TransferWithIdempotency(
		1,
		"09",
		1,
		1000,
		"Test transfer",
		"test-key",
	)

	if !errors.Is(err, ErrSelfTransfer) {
		t.Fatalf(
			"expected ErrSelfTransfer, got %v",
			err,
		)
	}
}
