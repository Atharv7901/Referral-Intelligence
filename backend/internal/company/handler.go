package company

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type createCompanyRequest struct {
	Name          string  `json:"name" binding:"required"`
	Slug          string  `json:"slug" binding:"required"`
	Website       *string `json:"website"`
	LinkedInUrl   *string `json:"linkedin_url"`
	CareerPageUrl *string `json:"career_page_url"`
	CompanyType   string  `json:"company_type" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createCompanyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	company := &Company{
		Name:          req.Name,
		Slug:          req.Slug,
		Website:       req.Website,
		LinkedInUrl:   req.LinkedInUrl,
		CareerPageUrl: req.CareerPageUrl,
		CompanyType:   req.CompanyType,
		IsActive:      true,
	}

	if err := h.service.Create(c.Request.Context(), company); err != nil {
		switch {
		case errors.Is(err, ErrInvalidCompanyName),
			errors.Is(err, ErrInvalidCompanySlug),
			errors.Is(err, ErrInvalidCompanyType):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create company",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, company)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid company id",
		})
		return
	}

	company, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "company not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get company",
		})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *Handler) List(c *gin.Context) {
	search := c.Query("search")

	companies, err := h.service.List(
		c.Request.Context(),
		search,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list companies",
		})
		return
	}

	c.JSON(http.StatusOK, companies)
}
