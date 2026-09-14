package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Health(c *gin.Context) {
	status := h.Service.Check(c.Request.Context())

	c.JSON(http.StatusOK, status)
}

func (h *Handler) DatabaseHealth(c *gin.Context) {
	status := h.Service.CheckDatabase(c.Request.Context())

	if status.Status == "error" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
