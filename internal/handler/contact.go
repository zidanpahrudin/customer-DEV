package handler

import (
	"net/http"
	"time"
	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"

	"github.com/gin-gonic/gin"
	// Hapus import "gorm.io/gorm" karena tidak digunakan
)

// @Summary Create contact for customer
// @Description Create a new contact person for specific customer
// @Tags Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param contact body entity.Contact true "Contact data"
// @Success 201 {object} entity.Contact
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/contacts [post]
// @Summary Create contact
// @Description Create a new contact
// @Tags Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param contact body dto.CreateContactRequest true "Contact data"
// @Success 201 {object} entity.Contact
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/contacts [post]
func CreateContact(c *gin.Context) {
	var req dto.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hapus validasi customer exists karena CustomerID tidak ada di request
	contact := entity.Contact{
		// CustomerID: customer.ID, // Akan diset sesuai kebutuhan
		Name:        req.Name,
		JobPosition: req.JobPosition,
		Email:       req.Email,
		Phone:       req.Phone,
		Mobile:      req.Mobile,
		Main:        req.IsMain,
		Active:      req.Active,
	}

	// Parse birthdate if provided
	if req.Birthdate != "" {
		if birthdate, err := time.Parse("2006-01-02", req.Birthdate); err == nil {
			contact.Birthdate = &birthdate
		}
	}

	result := config.DB.Create(&contact)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

// @Summary Get customer contacts
// @Description Get all contact persons for specific customer
// @Tags Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {array} entity.Contact
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/contacts [get]
func GetCustomerContacts(c *gin.Context) {
	customerID := c.Param("id")

	var contacts []entity.Contact
	result := config.DB.Where("customer_id = ?", customerID).Find(&contacts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// @Summary Get contact by ID
// @Description Get a specific contact person by ID
// @Tags Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Success 200 {object} entity.Contact
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/contacts/{id} [get]
func GetContact(c *gin.Context) {
	id := c.Param("id")

	var contact entity.Contact
	result := config.DB.First(&contact, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// @Summary Update contact
// @Description Update an existing contact person
// @Tags Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Param contact body entity.Contact true "Contact data"
// @Success 200 {object} entity.Contact
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/contacts/{id} [put]
func UpdateContact(c *gin.Context) {
	id := c.Param("id")

	var contact entity.Contact
	result := config.DB.First(&contact, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&contact)
	c.JSON(http.StatusOK, contact)
}

// @Summary Delete contact
// @Description Delete a contact person by ID
// @Tags Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/contacts/{id} [delete]
func DeleteContact(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&entity.Contact{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

// @Summary Get customer with contacts
// @Description Get customer data with all contact persons
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id}/with-contacts [get]
func GetCustomerWithContacts(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.Preload("Contacts").First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}