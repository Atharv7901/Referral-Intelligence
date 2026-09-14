package main

import (
	"log"

	"github.com/atharvshivarkar/referral-intelligence/internal/config"
	"github.com/atharvshivarkar/referral-intelligence/internal/database"
	"github.com/atharvshivarkar/referral-intelligence/internal/health"
	"github.com/atharvshivarkar/referral-intelligence/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}

	defer db.Close()

	healthService := health.NewService(db)
	healthHandler := health.NewHandler(healthService)

	r := router.Setup(healthHandler)

	port := cfg.Port

	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("server failed:", err)
	}
}
