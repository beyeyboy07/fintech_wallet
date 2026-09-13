package router

import (
	"fintech-wallet/backend/internal/auth"
	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/handler"
	"fintech-wallet/backend/internal/middleware"
	"fintech-wallet/backend/internal/repository"
	"fintech-wallet/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func SetupRouter(
	db *gorm.DB,
	cfg config.Config,
	redisService *service.RedisService,
	queueClient *asynq.Client,
) *gin.Engine {

	router := gin.Default()

	// JWT
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Repository
	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	walletRepository := repository.NewWalletRepository(db)
	idempotencyRepository := repository.NewIdempotencyRepository(db)
	outboxRepository := repository.NewOutboxRepository(db)
	// Service
	userService := service.NewUserService(
		userRepository,
	)
	queueService := service.NewQueueService(queueClient)

	transactionRepository := repository.NewTransactionRepository(db)

	transactionEntryRepository := repository.NewTransactionEntryRepository(db)

	transactionService := service.NewTransactionService(
		transactionRepository,
	)

	transactionHandler := handler.NewTransactionHandler(
		transactionService,
	)
	webhookHandler := handler.NewWebhookHandler()

	transferService := service.NewTransferService(
		walletRepository,
		transactionRepository,
		transactionEntryRepository,
		idempotencyRepository,
		redisService,
		db,
	)

	transferHandler := handler.NewTransferHandler(
		transferService,
	)

	walletService := service.NewWalletService(
		walletRepository,
		transactionRepository,
		outboxRepository,
		redisService,
		queueService,
		db,
	)

	authService := service.NewAuthService(
		userRepository,
		walletRepository,
		refreshTokenRepository,
		jwtManager,
		db,
	)

	// Handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	walletHandler := handler.NewWalletHandler(
		walletService,
	)
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

	router.GET(
		"/api/v1/wallet",
		middleware.AuthMiddleware(jwtManager),
		walletHandler.GetWallet,
	)

	router.POST(
		"/api/v1/wallet/topup",
		middleware.AuthMiddleware(jwtManager),
		walletHandler.TopUp,
	)

	router.POST(
		"/api/v1/transfers",
		middleware.AuthMiddleware(jwtManager),
		transferHandler.Transfer,
	)

	router.GET(
		"/api/v1/transactions",
		middleware.AuthMiddleware(jwtManager),
		transactionHandler.GetTransactions,
	)

	router.POST(
		"/api/v1/webhooks/payment",
		webhookHandler.Payment,
	)
	return router
}
