package router

import (
	"github.com/atharvshivarkar/referral-intelligence/internal/health"
	"github.com/gin-gonic/gin"
)

func Setup(healthHandler *health.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthHandler.Health)
	router.GET("/health/db", healthHandler.DatabaseHealth)

	return router
}
