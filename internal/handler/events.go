package handler

import (
	"net/http"
	"strconv"

	"customer-api/internal/config"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"
)


// @Summary Get all Events
// @Description Get list of all events
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {array} entity.Event
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/events [get]
func Read(c *gin.Context) {
	var events []entity.Event
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	if limit <= 0 {
		limit = 10 // Default limit
	}
	if page <= 0 {
		page = 1 // Default page
	}

	offset := (page - 1) * limit

	if result := config.DB.Limit(limit).Offset(offset).Find(&events); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
