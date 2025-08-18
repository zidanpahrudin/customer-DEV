package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create sosmed for customer
// @Description Create a new social media account for specific customer
// @Tags Social Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param sosmed body entity.Sosmed true "Sosmed data"
// @Success 201 {object} entity.Sosmed
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/sosmeds [post]
// Hapus seluruh fungsi CreateSosmed (baris 11-50)
// Fungsi ini sudah tidak diperlukan karena endpoint POST-nya dihapus
func CreateSosmed(c *gin.Context) {
	customerID := c.Param("id")

	// Check if customer exists
	var customer entity.Customer
	if err := config.DB.First(&customer, customerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var sosmed entity.Sosmed
	if err := c.ShouldBindJSON(&sosmed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set customer ID
	sosmed.CustomerID = customer.ID

	result := config.DB.Create(&sosmed)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sosmed"})
		return
	}

	c.JSON(http.StatusCreated, sosmed)
}

// @Summary Get customer sosmeds
// @Description Get all social media accounts for specific customer
// @Tags Social Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {array} entity.Sosmed
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/sosmeds [get]
func GetCustomerSosmeds(c *gin.Context) {
	customerID := c.Param("id")

	var sosmeds []entity.Sosmed
	result := config.DB.Where("customer_id = ?", customerID).Find(&sosmeds)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sosmeds"})
		return
	}

	c.JSON(http.StatusOK, sosmeds)
}

// @Summary Get sosmed by ID
// @Description Get a specific social media account by ID
// @Tags Social Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sosmed ID"
// @Success 200 {object} entity.Sosmed
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/sosmeds/{id} [get]
func GetSosmed(c *gin.Context) {
	id := c.Param("id")

	var sosmed entity.Sosmed
	result := config.DB.First(&sosmed, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sosmed not found"})
		return
	}

	c.JSON(http.StatusOK, sosmed)
}

// @Summary Update sosmed
// @Description Update an existing social media account
// @Tags Social Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sosmed ID"
// @Param sosmed body entity.Sosmed true "Sosmed data"
// @Success 200 {object} entity.Sosmed
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/sosmeds/{id} [put]
func UpdateSosmed(c *gin.Context) {
	id := c.Param("id")

	var sosmed entity.Sosmed
	result := config.DB.First(&sosmed, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sosmed not found"})
		return
	}

	if err := c.ShouldBindJSON(&sosmed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&sosmed)
	c.JSON(http.StatusOK, sosmed)
}

// @Summary Delete sosmed
// @Description Delete a social media account by ID
// @Tags Social Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sosmed ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/sosmeds/{id} [delete]
func DeleteSosmed(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&entity.Sosmed{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sosmed"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sosmed not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sosmed deleted successfully"})
}

// @Summary Get customer with sosmeds
// @Description Get customer data with all social media accounts
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/with-sosmeds [get]
func GetCustomerWithSosmeds(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.Preload("Sosmeds").First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Get Customer with All Relations (Addresses and Sosmeds)
// @Summary Get customer with all relations
// @Description Get customer data with all relations (addresses and social media accounts)
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/with-all-relations [get]
func GetCustomerWithAddressesAndSosmeds(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.Preload("Addresses").Preload("Sosmeds").First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
