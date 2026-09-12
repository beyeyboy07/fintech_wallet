package main

import (
	"log"
	"time"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/redis"
	"fintech-wallet/backend/internal/router"
)

func main() {
	cfg := config.Load()

	// Connect to Redis
	redisClient := redis.NewClient(cfg)

	if err := redis.Ping(redisClient); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected successfully")

	// Test Redis SET
	err := redis.Set(
		redisClient,
		"test:key",
		"hello redis",
		5*time.Minute,
	)

	if err != nil {
		log.Fatal("Redis SET failed:", err)
	}

	// Test Redis GET
	value, err := redis.Get(
		redisClient,
		"test:key",
	)

	if err != nil {
		log.Fatal("Redis GET failed:", err)
	}

	log.Println("Redis value:", value)

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
	r := router.SetupRouter(db, cfg)

	// Start server
	log.Println("Server running on port :" + cfg.AppPort)

	err = r.Run(":" + cfg.AppPort)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
