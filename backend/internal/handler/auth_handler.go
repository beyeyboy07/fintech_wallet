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

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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

	result, err := handler.authService.Login(
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
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
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

func (h *AuthHandler) Refresh(c *gin.Context) {
	var request RefreshRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
		)
		return
	}

	result, err := h.authService.Refresh(
		request.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(
				c,
				http.StatusUnauthorized,
				"Invalid or expired refresh token",
				"INVALID_REFRESH_TOKEN",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to refresh token",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Token refreshed successfully",
		gin.H{
			"access_token":  result.AccessToken,
			"refresh_token": result.RefreshToken,
		},
	)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var request RefreshRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
		)
		return
	}

	err := h.authService.Logout(
		request.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(
				c,
				http.StatusUnauthorized,
				"Invalid or expired refresh token",
				"INVALID_REFRESH_TOKEN",
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"Failed to logout",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"Logout successful",
		nil,
	)
}
