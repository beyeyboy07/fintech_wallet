package router

import (
	"fintech-wallet/backend/internal/auth"
	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/handler"
	"fintech-wallet/backend/internal/middleware"
	"fintech-wallet/backend/internal/repository"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(
	db *gorm.DB,
	cfg config.Config,
) *gin.Engine {

	router := gin.Default()

	// JWT
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Repository
	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	// Service
	userService := service.NewUserService(
		userRepository,
	)

	authService := service.NewAuthService(
		userRepository,
		refreshTokenRepository,
		jwtManager,
		db,
	)
	// Handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	// Routes
	router.GET("/api/v1/health", handler.Health)

	router.GET(
		"/api/v1/users",
		middleware.AuthMiddleware(jwtManager),
		userHandler.GetUsers,
	)

	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)
	router.POST(
		"/api/v1/auth/refresh",
		authHandler.Refresh,
	)
	router.POST("/api/v1/auth/logout", authHandler.Logout)

	router.GET(
		"/api/v1/auth/me",
		middleware.AuthMiddleware(jwtManager),
		authHandler.Me,
	)

	return router
}
