package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"nota-parfume/internal/models"
	"nota-parfume/internal/service"

	"github.com/gin-gonic/gin"
)

type ParfumeHandler struct {
	service service.ParfumeService
}

func NewParfumeHandler(service service.ParfumeService) *ParfumeHandler {
	return &ParfumeHandler{
		service: service,
	}
}

// parfumeHandler := transport.NewParfumeHandler(parfumeService)
func (h *ParfumeHandler) ParfumeRegisterRoutes(authorized *gin.RouterGroup, unauthorized *gin.RouterGroup) {
	publicParfumes := unauthorized.Group("/parfumes")
	{
		publicParfumes.GET("/:id", h.GetByID)
		publicParfumes.GET("", h.GetAll)
	}

	protectedParfumes := unauthorized.Group("/parfumes")
	{
		protectedParfumes.POST("", h.Create)
		protectedParfumes.PUT("/:id", h.Update)
		protectedParfumes.DELETE("/:id", h.Delete)
	}

}

func (h *ParfumeHandler) GetAll(c *gin.Context) {

	filter := models.ParfumeFilter{
		Brand:    c.Query("brand"),
		Category: c.Query("category"),
		Search:   c.Query("search"),
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	parfumes, total, err := h.service.GetAll(
		filter,
		page,
		limit,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"parfumes": parfumes,
		"page":     page,
		"limit":    limit,
		"total":    total,
	})
}

// POST /parfumes
func (h *ParfumeHandler) Create(c *gin.Context) {

	var req models.ParfumeCreate

	req.Name = c.PostForm("name")
	req.Description = c.PostForm("description")
	req.Brand = c.PostForm("brand")
	req.Category = c.PostForm("category")
	req.Badge = c.PostForm("badge")

	// PricePerMl
	price, err := strconv.ParseInt(
		c.PostForm("price_per_ml"),
		10,
		64,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid price_per_ml",
		})
		return
	}

	req.PricePerMl = price

	// Notes
	if err := json.Unmarshal(
		[]byte(c.PostForm("notes")),
		&req.Notes,
	); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid notes format",
		})
		return
	}

	// AvailableVolumes
	if err := json.Unmarshal(
		[]byte(c.PostForm("available_volumes")),
		&req.AvailableVolumes,
	); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid available_volumes format",
		})
		return
	}

	// IsActive
	isActive := c.PostForm("is_active")

	if isActive != "" {

		value, err := strconv.ParseBool(isActive)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid is_active",
			})
			return
		}

		req.IsActive = &value
	}

	// Image
	var imageUrl string

	file, err := c.FormFile("image")

	if err == nil {

		imageUrl, err = h.service.UploadImage(file)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	parfume, err := h.service.Create(&req, imageUrl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, parfume)
}

// GET /parfumes/:id

func (h *ParfumeHandler) GetByID(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})

		return
	}

	parfume, err := h.service.GetByID(uint(id))

	if err != nil {

		if errors.Is(err, service.ErrParfumeNotFound) {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "parfume not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parfume": parfume,
	})
}

// PUT /parfumes/:id

func (h *ParfumeHandler) Update(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})

		return
	}

	var input models.ParfumeUpdate

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	parfume, err := h.service.Update(
		uint(id),
		&input,
	)

	if err != nil {

		if errors.Is(err, service.ErrParfumeNotFound) {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "parfume not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parfume": parfume,
	})
}

// DELETE /parfumes/:id

func (h *ParfumeHandler) Delete(c *gin.Context) {

	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})

		return
	}

	err = h.service.Delete(uint(id))

	if err != nil {

		if errors.Is(err, service.ErrParfumeNotFound) {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "parfume not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "parfume deleted",
	})
}
