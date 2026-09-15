package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/response"
	"fintech-wallet/backend/internal/service"
	"fintech-wallet/backend/internal/webhook"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
	cfg            config.Config
}

func NewWebhookHandler(
	webhookService *service.WebhookService,
	cfg config.Config,
) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
		cfg:            cfg,
	}
}

// Payment godoc
// @Summary Payment webhook
// @Description Receive and process payment webhook with HMAC signature verification
// @Tags Webhook
// @Accept json
// @Produce json
// @Param X-Webhook-Signature header string true "HMAC-SHA256 webhook signature"
// @Param X-Webhook-Timestamp header string true "Unix timestamp of the webhook request"
// @Param request body model.PaymentWebhookRequest true "Payment webhook request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/webhooks/payment [post]
func (h *WebhookHandler) Payment(c *gin.Context) {

	var request model.PaymentWebhookRequest

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Failed to read webhook body",
			"INVALID_WEBHOOK_BODY",
		)
		return
	}

	signature := c.GetHeader("X-Webhook-Signature")
	timestamp := c.GetHeader("X-Webhook-Timestamp")

	if !webhook.VerifyTimestamp(timestamp) {
		response.Error(
			c,
			http.StatusUnauthorized,
			"Invalid webhook timestamp",
			"INVALID_WEBHOOK_TIMESTAMP",
		)
		return
	}

	if !webhook.VerifyTimestampedSignature(
		body,
		timestamp,
		signature,
		h.cfg.WebhookSecret,
	) {
		response.Error(
			c,
			http.StatusUnauthorized,
			"Invalid webhook signature",
			"INVALID_WEBHOOK_SIGNATURE",
		)
		return
	}

	c.Request.Body = io.NopCloser(
		bytes.NewReader(body),
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid webhook payload",
			"INVALID_WEBHOOK_PAYLOAD",
		)
		return
	}

	err = h.webhookService.Receive(request)

	if err != nil {

		switch {
		case errors.Is(err, service.ErrDuplicateWebhook):
			response.Error(
				c,
				http.StatusConflict,
				"Duplicate webhook event",
				"DUPLICATE_WEBHOOK",
			)

		case errors.Is(err, service.ErrTransactionNotFound):
			response.Error(
				c,
				http.StatusNotFound,
				"Transaction not found",
				"TRANSACTION_NOT_FOUND",
			)

		case errors.Is(err, service.ErrInvalidTransactionType):
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid transaction type",
				"INVALID_TRANSACTION_TYPE",
			)

		case errors.Is(err, service.ErrInvalidWebhookAmount):
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid webhook amount",
				"INVALID_WEBHOOK_AMOUNT",
			)

		case errors.Is(err, service.ErrTransactionProcessed):
			response.Error(
				c,
				http.StatusConflict,
				"Transaction already processed",
				"TRANSACTION_ALREADY_PROCESSED",
			)

		case errors.Is(err, service.ErrInvalidWebhookStatus):
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid webhook status",
				"INVALID_WEBHOOK_STATUS",
			)

		default:
			response.Error(
				c,
				http.StatusInternalServerError,
				"Failed to process webhook",
				"WEBHOOK_PROCESSING_FAILED",
			)
		}

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Webhook processed successfully",
		nil,
	)
}
