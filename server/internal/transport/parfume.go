package transport

import (
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

	parfumes := unauthorized.Group("/parfumes")
	{
		parfumes.POST("", h.Create)
		parfumes.GET("", h.GetAll)
		parfumes.GET("/:id", h.GetByID)
		parfumes.PUT("/:id", h.Update)
		parfumes.DELETE("/:id", h.Delete)
	}

}

// POST /parfumes

func (h *ParfumeHandler) Create(c *gin.Context) {

	var input models.ParfumeCreate

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	parfume, err := h.service.Create(&input)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"parfume": parfume,
	})
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

// GET /parfumes?limit=10&offset=0

func (h *ParfumeHandler) GetAll(c *gin.Context) {

	limit, err := strconv.Atoi(
		c.DefaultQuery("limit", "10"),
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid limit",
		})

		return
	}

	offset, err := strconv.Atoi(
		c.DefaultQuery("offset", "0"),
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid offset",
		})

		return
	}

	parfumes, err := h.service.GetAll(limit, offset)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parfumes": parfumes,
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
