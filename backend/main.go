package main

import (
	"log"

	"fintech-wallet/backend/internal/config"
	"fintech-wallet/backend/internal/database"
	"fintech-wallet/backend/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Database connected successfully")

	err = database.Migrate(db)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	r := router.SetupRouter(db, cfg)

	log.Println("Server running on port :" + cfg.AppPort)

	err = r.Run(":" + cfg.AppPort)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
