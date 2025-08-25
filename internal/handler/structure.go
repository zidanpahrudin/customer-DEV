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

// @Summary Create structure for customer
// @Description Create a new organizational structure for specific customer
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param structure body entity.Structure true "Structure data"
// @Success 201 {object} entity.Structure
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/structures [post]
// @Summary Create structure
// @Description Create a new structure
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param structure body dto.CreateStructureRequest true "Structure data"
// @Success 201 {object} entity.Structure
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/structures [post]
func CreateStructure(c *gin.Context) {
	var req dto.CreateStructureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hapus validasi customer exists karena CustomerID tidak ada di request
	structure := entity.Structure{
		// CustomerID: customer.ID, // Akan diset sesuai kebutuhan
		Name:    req.Name,
		Level:   req.Level,
		Address: req.Address,
		Active:  req.Active,
	}

	result := config.DB.Create(&structure)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create structure"})
		return
	}

	c.JSON(http.StatusCreated, structure)
}

// @Summary Get customer structures
// @Description Get all organizational structures for specific customer
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {array} entity.Structure
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/structures [get]
func GetCustomerStructures(c *gin.Context) {
	customerID := c.Param("id")

	var structures []entity.Structure
	result := config.DB.Where("customer_id = ?", customerID).Order("level ASC, name ASC").Find(&structures)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch structures"})
		return
	}

	c.JSON(http.StatusOK, structures)
}

// @Summary Get structures by level
// @Description Get organizational structures filtered by level for specific customer
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param level query int false "Structure level filter"
// @Success 200 {array} entity.Structure
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/structures/by-level [get]
func GetStructuresByLevel(c *gin.Context) {
	customerID := c.Param("id")
	levelStr := c.Query("level")

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level parameter"})
		return
	}

	var structures []entity.Structure
	query := config.DB.Where("customer_id = ?", customerID)
	if levelStr != "" {
		query = query.Where("level = ?", level)
	}

	result := query.Order("level ASC, name ASC").Find(&structures)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch structures"})
		return
	}

	c.JSON(http.StatusOK, structures)
}

// @Summary Get structure by ID
// @Description Get a specific organizational structure by ID
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Structure ID"
// @Success 200 {object} entity.Structure
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/structures/{id} [get]
func GetStructure(c *gin.Context) {
	id := c.Param("id")

	var structure entity.Structure
	result := config.DB.Preload("Parent").Preload("Children").First(&structure, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Structure not found"})
		return
	}

	c.JSON(http.StatusOK, structure)
}

// @Summary Update structure
// @Description Update an existing organizational structure
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Structure ID"
// @Param structure body entity.Structure true "Structure data"
// @Success 200 {object} entity.Structure
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/structures/{id} [put]
func UpdateStructure(c *gin.Context) {
	id := c.Param("id")

	var structure entity.Structure
	result := config.DB.First(&structure, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Structure not found"})
		return
	}

	if err := c.ShouldBindJSON(&structure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&structure)
	c.JSON(http.StatusOK, structure)
}

// @Summary Delete structure
// @Description Delete an organizational structure by ID
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Structure ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/structures/{id} [delete]
func DeleteStructure(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&entity.Structure{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete structure"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Structure not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Structure deleted successfully"})
}

// @Summary Get customer with structures
// @Description Get customer data with all organizational structures
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/with-structures [get]
func GetCustomerWithStructures(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.Preload("Structures", func(db *gorm.DB) *gorm.DB {
		return db.Order("level ASC, name ASC")
	}).First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary Get customer with all relations
// @Description Get customer data with all related data (addresses, sosmeds, contacts, structures)
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/with-all [get]
func GetCustomerWithAllRelations(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.Preload("Addresses").
		Preload("Sosmeds").
		Preload("Contacts").
		Preload("Structures", func(db *gorm.DB) *gorm.DB {
			return db.Order("level ASC, name ASC")
		}).
		First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary Get customer full data
// @Description Get complete customer data with all relations
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/full [get]
/* func GetCustomerFull(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.Preload("Addresses").
		Preload("Sosmeds").
		Preload("Contacts").
		Preload("Structures", func(db *gorm.DB) *gorm.DB {
			return db.Order("level ASC, name ASC")
		}).
		First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
} */
