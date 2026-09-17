package router

import (
	"github.com/atharvshivarkar/referral-intelligence/internal/company"
	"github.com/atharvshivarkar/referral-intelligence/internal/health"
	"github.com/gin-gonic/gin"
)

func Setup(healthHandler *health.Handler, companyHandler *company.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthHandler.Health)
	router.GET("/health/db", healthHandler.DatabaseHealth)

	api := router.Group("/api/v1")
	{
		api.GET("/companies", companyHandler.List)
		api.GET("/companies/:id", companyHandler.GetByID)
		api.POST("/companies", companyHandler.Create)
		api.POST("/companies/:id/sources/resolve", companyHandler.ResolveSource)
	}

	return router
}
