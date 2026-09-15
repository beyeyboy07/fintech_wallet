package handler

import (
	"net/http"

	"fintech-wallet/backend/internal/response"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary Get users
// @Description Get list of users for admin
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to get users",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Users retrieved successfully",
		users,
	)
}
