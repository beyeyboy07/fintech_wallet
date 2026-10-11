package service

import (
	"errors"
	"testing"
)

func TestTransferWithIdempotencyByAccountNumber_InvalidAccountNumber(t *testing.T) {
	service := &TransferService{}

	tests := []struct {
		name          string
		accountNumber string
	}{
		{
			name:          "account number too short",
			accountNumber: "12345",
		},
		{
			name:          "account number contains letters",
			accountNumber: "88010000000000AB",
		},
		{
			name:          "account number contains spaces",
			accountNumber: "88010000000000 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.TransferWithIdempotencyByAccountNumber(
				1,
				tt.accountNumber,
				1000,
				"Test transfer",
				"test-idempotency-key",
			)

			if !errors.Is(err, ErrInvalidAccountNumber) {
				t.Fatalf(
					"expected ErrInvalidAccountNumber, got %v",
					err,
				)
			}
		})
	}
}
