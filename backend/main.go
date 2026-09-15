package main

import (
	_ "fintech-wallet/backend/docs"
	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/redis"
	"fintech-wallet/backend/internal/router"
	"fintech-wallet/backend/internal/service"
	"fintech-wallet/backend/internal/worker"
	"log"
)

// @title Fintech Wallet API
// @version 1.0
// @description REST API for Fintech Wallet learning project.
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	// Connect to Redis
	redisClient := redis.NewClient(cfg)

	if err := redis.Ping(redisClient); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	queueClient := worker.NewClient(
		cfg.RedisHost + ":" + cfg.RedisPort,
	)

	defer queueClient.Close()

	log.Println("Redis connected successfully")

	// Test Redis SET
	// err := redis.Set(
	// 	redisClient,
	// 	"test:key",
	// 	"hello redis",
	// 	5*time.Minute,
	// )

	// if err != nil {
	// 	log.Fatal("Redis SET failed:", err)
	// }

	// // Test Redis GET
	// value, err := redis.Get(
	// 	redisClient,
	// 	"test:key",
	// )

	// if err != nil {
	// 	log.Fatal("Redis GET failed:", err)
	// }

	// log.Println("Redis value:", value)

	// Connect to database
	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Database connected successfully")

	// Run migration
	err = database.Migrate(db)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Setup router
	redisService := service.NewRedisService(redisClient)

	r := router.SetupRouter(
		db,
		cfg,
		redisService,
		queueClient,
	)

	// Start server
	log.Println("Server running on port :" + cfg.AppPort)

	err = r.Run(":" + cfg.AppPort)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
