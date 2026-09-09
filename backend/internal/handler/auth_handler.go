package handler

import (
	"errors"
	"net/http"

	"fintech-wallet/backend/internal/response"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (handler *AuthHandler) Register(data *gin.Context) {
	var request RegisterRequest

	if err := data.ShouldBindJSON(&request); err != nil {
		response.Error(
			data,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
		)
		return
	}

	user, err := handler.authService.Register(
		request.Name,
		request.Email,
		request.Password,
	)

	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.Error(
				data,
				http.StatusConflict,
				"Email already exists",
				"EMAIL_ALREADY_EXISTS",
			)
			return
		}

		response.Error(
			data,
			http.StatusInternalServerError,
			"Failed to register user",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		data,
		http.StatusCreated,
		"User registered successfully",
		user,
	)
}

func (handler *AuthHandler) Login(data *gin.Context) {
	var request LoginRequest

	if err := data.ShouldBindJSON(&request); err != nil {
		response.Error(
			data,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
		)
		return
	}

	token, err := handler.authService.Login(
		request.Email,
		request.Password,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(
				data,
				http.StatusUnauthorized,
				"Invalid email or password",
				"INVALID_CREDENTIALS",
			)
			return
		}

		response.Error(
			data,
			http.StatusInternalServerError,
			"Failed to login",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		data,
		http.StatusOK,
		"Login successful",
		gin.H{
			"token": token,
		},
	)
}

func (h *AuthHandler) Me(c *gin.Context) {
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

	userIDFloat, ok := userIDValue.(float64)

	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			"Invalid user ID",
			"UNAUTHORIZED",
		)
		return
	}

	userID := uint(userIDFloat)

	user, err := h.authService.GetCurrentUser(userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				"User not found",
				"USER_NOT_FOUND",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to get user",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"User retrieved successfully",
		user,
	)
}
