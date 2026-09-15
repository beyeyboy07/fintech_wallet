package handler

import (
	"errors"
	"net/http"

	"fintech-wallet/backend/internal/response"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WalletHandler struct {
	walletService *service.WalletService
}

// type TransferRequest struct {
// 	ReceiverID  uint   `json:"receiver_id" binding:"required"`
// 	Amount      int64  `json:"amount" binding:"required"`
// 	Description string `json:"description"`
// }

func NewWalletHandler(
	walletService *service.WalletService,
) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// GetWallet godoc
// @Summary Get wallet
// @Description Get wallet information for the currently authenticated user
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/wallet [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {

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

	userID, ok := userIDValue.(uint)

	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"Invalid user ID",
			"UNAUTHORIZED",
		)
		return
	}

	wallet, err := h.walletService.GetWallet(userID)

	if err != nil {
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
			"Failed to get wallet",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Wallet retrieved successfully",
		wallet,
	)
}

type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

// TopUp godoc
// @Summary Top up wallet
// @Description Add balance to the authenticated user's wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TopUpRequest true "Top up request"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/wallet/topup [post]
func (h *WalletHandler) TopUp(c *gin.Context) {

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

	userID, ok := userIDValue.(uint)

	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"Invalid user ID",
			"UNAUTHORIZED",
		)
		return
	}

	var request TopUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
		)
		return
	}

	transaction, err := h.walletService.TopUpAsync(
		userID,
		request.Amount,
	)

	if err != nil {

		if errors.Is(err, service.ErrInvalidTopUpAmount) {
			response.Error(
				c,
				http.StatusBadRequest,
				"Invalid top up amount",
				"INVALID_TOPUP_AMOUNT",
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
			"Failed to top up wallet",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusAccepted,
		"Top up request accepted",
		transaction,
	)
}
