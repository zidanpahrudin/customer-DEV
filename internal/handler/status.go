package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get all statuses
// @Description Get list of all statuses (master data)
// @Tags Statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param active query bool false "Filter by active status"
// @Success 200 {array} dto.StatusResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/statuses [get]
func GetStatuses(c *gin.Context) {
	var statuses []entity.Status
	query := config.DB

	// Filter by active status if provided
	if activeParam := c.Query("active"); activeParam != "" {
		if active, err := strconv.ParseBool(activeParam); err == nil {
			// Assuming we add Active field to Status model later
			_ = active // placeholder for now
		}
	}

	if result := query.Find(&statuses); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statuses"})
		return
	}

	// Convert to response format
	var statusResponses []dto.StatusResponse
	for _, status := range statuses {
		statusResponses = append(statusResponses, dto.StatusResponse{
			ID:         status.ID,
			StatusName: status.StatusName,
			CreatedAt:  status.CreatedAt,
			UpdatedAt:  status.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, statusResponses)
}

// @Summary Get status by ID
// @Description Get a specific status by ID
// @Tags Statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Status ID"
// @Success 200 {object} dto.StatusResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/statuses/{id} [get]
func GetStatus(c *gin.Context) {
	id := c.Param("id")
	var status entity.Status

	if result := config.DB.First(&status, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch status"})
		}
		return
	}

	statusResponse := dto.StatusResponse{
		ID:         status.ID,
		StatusName: status.StatusName,
		CreatedAt:  status.CreatedAt,
		UpdatedAt:  status.UpdatedAt,
	}

	c.JSON(http.StatusOK, statusResponse)
}

// @Summary Create status
// @Description Create a new status (master data)
// @Tags Statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status body dto.CreateStatusRequest true "Status data"
// @Success 201 {object} dto.StatusResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/statuses [post]
func CreateStatus(c *gin.Context) {
	var req dto.CreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if status name already exists
	var existingStatus entity.Status
	if result := config.DB.Where("status_name = ?", req.StatusName).First(&existingStatus); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status name already exists"})
		return
	}

	// Create new status
	status := entity.Status{
		StatusName: req.StatusName,
	}

	if result := config.DB.Create(&status); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create status"})
		return
	}

	statusResponse := dto.StatusResponse{
		ID:         status.ID,
		StatusName: status.StatusName,
		CreatedAt:  status.CreatedAt,
		UpdatedAt:  status.UpdatedAt,
	}

	c.JSON(http.StatusCreated, statusResponse)
}

// @Summary Update status
// @Description Update an existing status
// @Tags Statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Status ID"
// @Param status body dto.UpdateStatusRequest true "Status data"
// @Success 200 {object} dto.StatusResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/statuses/{id} [put]
func UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var status entity.Status

	if result := config.DB.First(&status, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch status"})
		}
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if status name already exists (exclude current status)
	var existingStatus entity.Status
	if result := config.DB.Where("status_name = ? AND id != ?", req.StatusName, id).First(&existingStatus); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status name already exists"})
		return
	}

	// Update status
	status.StatusName = req.StatusName

	if result := config.DB.Save(&status); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	statusResponse := dto.StatusResponse{
		ID:         status.ID,
		StatusName: status.StatusName,
		CreatedAt:  status.CreatedAt,
		UpdatedAt:  status.UpdatedAt,
	}

	c.JSON(http.StatusOK, statusResponse)
}

// @Summary Delete status
// @Description Delete a status (soft delete)
// @Tags Statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Status ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/statuses/{id} [delete]
func DeleteStatus(c *gin.Context) {
	id := c.Param("id")
	var status entity.Status

	// Check if status exists
	if result := config.DB.First(&status, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch status"})
		}
		return
	}

	// Check if status is being used by customers
	var customerCount int64
	config.DB.Model(&entity.Customer{}).Where("status_id = ?", id).Count(&customerCount)
	if customerCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete status that is being used by customers"})
		return
	}

	// Soft delete
	if result := config.DB.Delete(&status); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status deleted successfully"})
}

// @Summary Get customers by status
// @Description Get all customers that have a specific status
// @Tags Statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Status ID"
// @Success 200 {array} dto.CustomerResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/statuses/{id}/customers [get]
func GetCustomersByStatus(c *gin.Context) {
	id := c.Param("id")
	var status entity.Status

	// Check if status exists
	if result := config.DB.First(&status, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch status"})
		}
		return
	}

	// Get customers with this status
	var customers []entity.Customer
	if result := config.DB.Where("status_id = ?", id).Find(&customers); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	// Convert to response format
	var customerResponses []dto.CustomerResponse
	for _, customer := range customers {
		customerResponses = append(customerResponses, dto.CustomerResponse{
			ID:        customer.ID,
			Name:      customer.Name,
			BrandName: customer.BrandName,
			Code:      customer.Code,
			/* 	Email:       customer.Email,
			Phone:       customer.Phone,
			Website:     customer.Website,
			Description: customer.Description, */
			Status:      customer.Status,
			Category:    customer.Category,
			Rating:      customer.Rating,
			AverageCost: customer.AverageCost,
		})
	}

	c.JSON(http.StatusOK, customerResponses)
}
