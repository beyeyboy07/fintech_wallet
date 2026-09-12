package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrIdempotencyKeyRequired = errors.New("idempotency key is required")
	ErrIdempotencyConflict    = errors.New("idempotency key conflict")
)

func GenerateRequestHash(
	receiverID uint,
	amount int64,
	description string,
) string {

	data := fmt.Sprintf(
		"%d:%d:%s",
		receiverID,
		amount,
		description,
	)

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
