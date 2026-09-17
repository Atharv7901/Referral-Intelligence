package main

import (
	"log"

	"github.com/atharvshivarkar/referral-intelligence/internal/company"
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

	//health module
	healthService := health.NewService(db)
	healthHandler := health.NewHandler(healthService)

	//source resolver module
	sourceRepository := company.NewSourceRepository(db)
	atsDetector := company.NewATSDetector()
	sourceResolver := company.NewSourceResolver(&sourceRepository, &atsDetector)

	//company module
	companyRepository := company.NewRepository(db)
	companyService := company.NewService(companyRepository, &sourceResolver)
	companyHandler := company.NewHandler(&companyService)
	r := router.Setup(healthHandler, companyHandler)

	port := cfg.Port

	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("server failed:", err)
	}
}
