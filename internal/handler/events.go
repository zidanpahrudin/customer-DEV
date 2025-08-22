package handler

import (
	"net/http"
	"strconv"

	"customer-api/internal/config"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" 

	"customer-api/internal/dto"
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
func ReadEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if limit == 0 {
		limit = 10 // Default limit
	}
	if page == 0 {
		page = 1 // Default page
	}

	var events []entity.Event
	offset := (page - 1) * limit

	if result := config.DB.Limit(limit).Offset(offset).Find(&events); result.Error != nil {

		c.JSON(http.StatusNotFound, dto.Response{
			Status:  http.StatusNotFound,
			Message: "No events found",
			Data:    []entity.Event{},
		})
		return
	}
	if len(events) == 0 {
		c.JSON(http.StatusOK, dto.Response{
			Status:  http.StatusOK,
			Message: "No events found",
			Data:    events,
		})
		return
	}


	c.JSON(http.StatusOK, dto.Response{
			Status:  http.StatusOK,
			Message: "Events retrieved successfully",
			Data:    events,
		})
}


// @Summary Create a new Event
// @Description Create a new event
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body entity.Event true "Event"
// @Success 201 {object} entity.Event
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/events [post]
func CreateEvents(c *gin.Context) {
	var event entity.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}


// @Summary Get a Event by ID
// @Description Get a event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 200 {object} entity.Event
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/events/{id} [get]
func ReadOneEvents(c *gin.Context) {
	var event entity.Event
	id := c.Param("id")
	if result := config.DB.First(&event, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data event"})
		}
		return
	}
	c.JSON(http.StatusOK, event)
}


// @Summary Update a Event by ID
// @Description Update a event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID"
// @Param event body entity.Event true "Event"
// @Success 200 {object} entity.Event
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/events/{id} [put]
func UpdateEvents(c *gin.Context) {
	var event entity.Event
	id := c.Param("id")
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := config.DB.First(&event, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data event"})
		}
		return
	}

	if err := config.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data event"})
		return
	}

	c.JSON(http.StatusOK, event)
}


// @Summary Delete a Event by ID
// @Description Delete a event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/events/{id} [delete]
func DeleteEvents(c *gin.Context) {
	var event entity.Event
	id := c.Param("id")
	if result := config.DB.First(&event, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data event"})
		}
		return
	}

	if err := config.DB.Delete(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data event"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// @Summary Get Events by Customer ID
// @Description Get all events for a specific customer
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {array} entity.Event
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/events [get]
func GetCustomerEvents(c *gin.Context) {
	customerID := c.Param("id")
	var events []entity.Event

	if result := config.DB.Where("customer_id = ?", customerID).Find(&events); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No events found for this customer"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data events"})
		}
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get Events by Type
// @Description Get all events of a specific type
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "Event Type"
// @Success 200 {array} entity.Event
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/events [get]
func GetEventType(c *gin.Context) {
	eventType := c.Query("type")
	var events []entity.Event

	if result := config.DB.Where("type = ?", eventType).Find(&events); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No events found for this type"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data events"})
		}
		return
	}

	c.JSON(http.StatusOK, events)
}


