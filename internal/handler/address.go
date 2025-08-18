package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Create address for customer
// @Description Create a new address for specific customer
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param address body entity.Address true "Address data"
// @Success 201 {object} entity.Address
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/addresses [post]
// @Summary Create address
// @Description Create a new address
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param address body dto.CreateAddressRequest true "Address data"
// @Success 201 {object} entity.Address
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/addresses [post]
func CreateAddress(c *gin.Context) {
	var req dto.CreateAddressRequest
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

	address := entity.Address{
		CustomerID: req.CustomerID,
		Name:       req.Name,
		Address:    req.Address,
		Main:       req.IsMain,
		Active:     req.Active,
	}

	result := config.DB.Create(&address)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, address)
}

// @Summary Get customer addresses
// @Description Get all addresses for specific customer
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {array} entity.Address
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/addresses [get]
func GetCustomerAddresses(c *gin.Context) {
	customerID := c.Param("id")

	var addresses []entity.Address
	result := config.DB.Where("customer_id = ?", customerID).Find(&addresses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// @Summary Get address by ID
// @Description Get a specific address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} entity.Address
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/addresses/{id} [get]
func GetAddress(c *gin.Context) {
	id := c.Param("id")

	var address entity.Address
	result := config.DB.First(&address, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, address)
}

// @Summary Update address
// @Description Update an existing address
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Param address body entity.Address true "Address data"
// @Success 200 {object} entity.Address
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/addresses/{id} [put]
func UpdateAddress(c *gin.Context) {
	id := c.Param("id")

	var address entity.Address
	result := config.DB.First(&address, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	var updateData entity.Address
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If this is set as main address, set all other addresses to false
	if updateData.Main {
		config.DB.Model(&entity.Address{}).Where("customer_id = ? AND id != ?", address.CustomerID, address.ID).Update("main", false)
	}

	// Update the address
	config.DB.Model(&address).Updates(updateData)
	c.JSON(http.StatusOK, address)
}

// @Summary Delete address
// @Description Delete an address by ID
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/addresses/{id} [delete]
func DeleteAddress(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&entity.Address{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

// @Summary Get customer with addresses
// @Description Get customer data with all addresses
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/with-addresses [get]
func GetCustomerWithAddresses(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.Preload("Addresses").First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
