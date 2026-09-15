package handler

import (
	"errors"
	"net/http"
	"strings"

	"fintech-wallet/backend/internal/response"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *service.AuthService
}

type RefreshTokenRequest struct {
	Token string `json:"refresh_token" binding:"required"`
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account and wallet
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/register [post]
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

	request.Email = strings.ToLower(strings.TrimSpace(request.Email))

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

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token and refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/login [post]
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

// Me godoc
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/me [get]
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

// Refresh godoc
// @Summary Refresh access token
// @Description Generate a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/refresh [post]
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

// Logout godoc
// @Summary Logout user
// @Description Revoke the current refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/logout [post]
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
