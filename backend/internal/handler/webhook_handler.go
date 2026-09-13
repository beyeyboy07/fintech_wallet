package handler

import (
	"net/http"

	"fintech-wallet/backend/internal/model"
	"fintech-wallet/backend/internal/response"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

func (h *WebhookHandler) Payment(c *gin.Context) {

	var request model.PaymentWebhookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid webhook payload",
			"INVALID_WEBHOOK_PAYLOAD",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Webhook received successfully",
		request,
	)
}
