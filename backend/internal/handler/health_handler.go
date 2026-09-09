package handler

import (
	"net/http"

	"fintech-wallet/backend/internal/response"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	response.Success(
		c,
		http.StatusOK,
		"Fintech Wallet API is healthy",
		nil,
	)
}
