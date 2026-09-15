package handler

import (
	"net/http"

	"fintech-wallet/backend/internal/response"

	"github.com/gin-gonic/gin"
)

// Health godoc
// @Summary Health check
// @Description Check whether the Fintech Wallet API is running
// @Tags Health
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /api/v1/health [get]
func Health(c *gin.Context) {
	response.Success(
		c,
		http.StatusOK,
		"Fintech Wallet API is healthy",
		nil,
	)
}
