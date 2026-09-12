package handler

import (
	"errors"
	"net/http"

	"fintech-wallet/backend/internal/response"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransferHandler struct {
	transferService *service.TransferService
}

func NewTransferHandler(
	transferService *service.TransferService,
) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

type TransferRequest struct {
	ReceiverID  uint   `json:"receiver_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	Description string `json:"description"`
}

func (h *TransferHandler) Transfer(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"User is not authenticated",
			"UNAUTHORIZED",
		)
		return
	}

	senderID, ok := userIDValue.(uint)

	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"Invalid user ID",
			"UNAUTHORIZED",
		)
		return
	}

	var request TransferRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
		)
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")

	transaction, err := h.transferService.TransferWithIdempotency(
		senderID,
		request.ReceiverID,
		request.Amount,
		request.Description,
		idempotencyKey,
	)

	if err != nil {

		if errors.Is(err, service.ErrIdempotencyKeyRequired) {
			response.Error(
				c,
				http.StatusBadRequest,
				"Idempotency-Key is required",
				"IDEMPOTENCY_KEY_REQUIRED",
			)
			return
		}

		if errors.Is(err, service.ErrIdempotencyConflict) {
			response.Error(
				c,
				http.StatusConflict,
				"Idempotency-Key has already been used with a different request",
				"IDEMPOTENCY_KEY_CONFLICT",
			)
			return
		}

		if errors.Is(err, service.ErrInvalidTransferAmount) {
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid transfer amount",
				"INVALID_TRANSFER_AMOUNT",
			)
			return
		}

		if errors.Is(err, service.ErrSelfTransfer) {
			response.Error(
				c,
				http.StatusBadRequest,
				"Cannot transfer to yourself",
				"SELF_TRANSFER",
			)
			return
		}

		if errors.Is(err, service.ErrSenderNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				"Sender wallet not found",
				"SENDER_WALLET_NOT_FOUND",
			)
			return
		}

		if errors.Is(err, service.ErrReceiverNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				"Receiver wallet not found",
				"RECEIVER_WALLET_NOT_FOUND",
			)
			return
		}

		if errors.Is(err, service.ErrInsufficientBalance) {
			response.Error(
				c,
				http.StatusBadRequest,
				"Insufficient balance",
				"INSUFFICIENT_BALANCE",
			)
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				"Wallet not found",
				"WALLET_NOT_FOUND",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to transfer",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"Transfer successful",
		transaction,
	)
}
