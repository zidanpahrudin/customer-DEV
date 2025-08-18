package handler

import (
	"net/http"
	"strconv"

	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OtherInput struct {
	AttributeName string `json:"attribute_name" binding:"required"`
	Value         string `json:"value"`
	Active        bool   `json:"active"`
}

// @Summary Create other attribute for customer
// @Description Create a new other attribute for specific customer
// @Tags Others
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param other body OtherInput true "Other attribute data"
// @Success 201 {object} entity.Other
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/others [post]
// Fungsi CreateOther sudah di-comment out (baris 32-76)
// Hapus seluruh blok komentar ini karena sudah tidak diperlukan

// @Summary Get customer others
// @Description Get all other attributes for specific customer
// @Tags Others
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param active query bool false "Filter by active status"
// @Param attribute_name query string false "Filter by attribute name"
// @Success 200 {array} entity.Other
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/others [get]
func GetCustomerOthers(c *gin.Context) {
	customerID := c.Param("id")
	query := config.DB.Where("customer_id = ?", customerID)

	// Filter by active status if provided
	if activeParam := c.Query("active"); activeParam != "" {
		if active, err := strconv.ParseBool(activeParam); err == nil {
			query = query.Where("active = ?", active)
		}
	}

	// Filter by attribute name if provided
	if attributeName := c.Query("attribute_name"); attributeName != "" {
		query = query.Where("attribute_name ILIKE ?", "%"+attributeName+"%")
	}

	var others []entity.Other
	result := query.Find(&others)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch other attributes"})
		return
	}

	c.JSON(http.StatusOK, others)
}

// @Summary Get other by ID
// @Description Get a specific other attribute by ID
// @Tags Others
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Other ID"
// @Success 200 {object} entity.Other
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/others/{id} [get]
func GetOther(c *gin.Context) {
	id := c.Param("id")
	var other entity.Other

	if result := config.DB.First(&other, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Other attribute not found"})
		return
	}

	c.JSON(http.StatusOK, other)
}

// @Summary Update other
// @Description Update an existing other attribute
// @Tags Others
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Other ID"
// @Param other body OtherInput true "Other attribute data"
// @Success 200 {object} entity.Other
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/others/{id} [put]
/* func UpdateOther(c *gin.Context) {
	id := c.Param("id")
	var other entity.Other

	if result := config.DB.First(&other, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Other attribute not found"})
		return
	}

	var input OtherInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update other
	other.AttributeName = input.AttributeName
	other.Value = input.Value
	other.Active = input.Active

	if result := config.DB.Save(&other); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update other attribute"})
		return
	}

	c.JSON(http.StatusOK, other)
} */

// @Summary Delete other
// @Description Delete an other attribute by ID
// @Tags Others
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Other ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/others/{id} [delete]
func DeleteOther(c *gin.Context) {
	id := c.Param("id")
	var other entity.Other

	if result := config.DB.First(&other, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Other attribute not found"})
		return
	}

	if result := config.DB.Delete(&other); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete other attribute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Other attribute deleted successfully"})
}

// @Summary Get customer with others
// @Description Get customer data with all other attributes
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/with-others [get]
func GetCustomerWithOthers(c *gin.Context) {
	id := c.Param("id")
	var customer entity.Customer

	if result := config.DB.Preload("Others").First(&customer, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary Get others by attribute name
// @Description Get all other attributes filtered by attribute name across all customers
// @Tags Others
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param attribute_name query string true "Attribute name to filter"
// @Param active query bool false "Filter by active status"
// @Success 200 {array} entity.Other
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/others/by-attribute [get]
func GetOthersByAttribute(c *gin.Context) {
	attributeName := c.Query("attribute_name")
	if attributeName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "attribute_name parameter is required"})
		return
	}

	query := config.DB.Where("attribute_name ILIKE ?", "%"+attributeName+"%")

	// Filter by active status if provided
	if activeParam := c.Query("active"); activeParam != "" {
		if active, err := strconv.ParseBool(activeParam); err == nil {
			query = query.Where("active = ?", active)
		}
	}

	var others []entity.Other
	result := query.Preload("Customer").Find(&others)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch other attributes"})
		return
	}

	c.JSON(http.StatusOK, others)
}

// @Summary Create other
// @Description Create a new other field
// @Tags Others
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param other body dto.CreateOtherRequest true "Other data"
// @Success 201 {object} entity.Other
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/others [post]
func CreateOther(c *gin.Context) {
	var req dto.CreateOtherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate customer exists
	var customer entity.Customer
	if err := config.DB.First(&customer, req.CustomerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate customer"})
		}
		return
	}

	other := entity.Other{
		CustomerID: req.CustomerID,
		Key:        req.Key,
		Value:      req.Value,
		Active:     req.Active,
	}

	result := config.DB.Create(&other)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create other field"})
		return
	}

	c.JSON(http.StatusCreated, other)
}
