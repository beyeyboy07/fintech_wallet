package handler

import (
	"net/http"
	"strconv"

	"fintech-wallet/backend/internal/response"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(
	transactionService *service.TransactionService,
) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")

	page := 1
	limit := 10

	if pageParam := c.Query("page"); pageParam != "" {
		value, err := strconv.Atoi(pageParam)

		if err != nil || value < 1 {
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid page",
				"INVALID_PAGE",
			)
			return
		}

		page = value
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		value, err := strconv.Atoi(limitParam)

		if err != nil || value < 1 {
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid limit",
				"INVALID_LIMIT",
			)
			return
		}

		if value > 100 {
			value = 100
		}

		limit = value
	}

	transactions, err :=
		h.transactionService.GetByUserID(
			userID,
			page,
			limit,
		)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to get transactions",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Transactions retrieved successfully",
		transactions,
	)
}
