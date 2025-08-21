package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Get all customers
// @Description Get list of all customers with optional status filter
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status" Enums(Active, Inactive, Blocked)
// @Success 200 {object} dto.CustomersResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers [get]
func GetCustomers(c *gin.Context) {
	var customers []entity.Customer
	status := c.Query("status")

	db := config.DB
	if status != "" {
		db = db.Where("status = ?", status)
	}

	result := db.Find(&customers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	// Calculate statistics
	var totalCustomers int64
	config.DB.Model(&entity.Customer{}).Count(&totalCustomers)

	var newCustomers int64
	config.DB.Model(&entity.Customer{}).Where("created_at >= NOW() - INTERVAL '1 year'").Count(&newCustomers)

	var avgCost float64
	config.DB.Model(&entity.Customer{}).Select("COALESCE(AVG(average_cost), 0)").Row().Scan(&avgCost)

	var blockedCustomers int64
	config.DB.Model(&entity.Customer{}).Where("status = ?", "Blocked").Count(&blockedCustomers)

	c.JSON(http.StatusOK, gin.H{
		"customers": customers,
		"stats": gin.H{
			"total_customers":   totalCustomers,
			"new_customers":     newCustomers,
			"avg_cost":          avgCost,
			"blocked_customers": blockedCustomers,
		},
	})
}

// @Summary Create new customer
// @Description Create a new customer record with all related data including addresses, social media, contacts, structures, groups, and other attributes
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
//
//	@Param customer body dto.CreateCustomerRequest true "Customer data" SchemaExample({
//	  "name": "PT Digital Inovasi Indonesia",
//	  "brandName": "DigiInno",
//	  "code": "DIGI",
//	  "accountManagerId": "AM-001",
//	  "logo": null,
//	  "addresses": [
//	    {
//	      "name": "Head Office",
//	      "address": "Jl. Sudirman No. 123, Jakarta Selatan",
//	      "isMain": true,
//	      "active": true
//	    },
//	    {
//	      "name": "Branch Office",
//	      "address": "Jl. Asia Afrika No. 45, Bandung",
//	      "isMain": false,
//	      "active": true
//	    }
//	  ],
//	  "socials": [
//	    {
//	      "platform": "Instagram",
//	      "handle": "@digiinno_id",
//	      "active": true
//	    },
//	    {
//	      "platform": "LinkedIn",
//	      "handle": "digital-inovasi-indonesia",
//	      "active": true
//	    }
//	  ],
//	  "contacts": [
//	    {
//	      "name": "Budi Santoso",
//	      "birthdate": "1985-03-15",
//	      "jobPosition": "CEO",
//	      "email": "budi@digiinno.com",
//	      "phone": "021-5551234",
//	      "mobile": "0812-3456-7890",
//	      "isMain": true,
//	      "active": true
//	    },
//	    {
//	      "name": "Sari Dewi",
//	      "birthdate": "1988-07-22",
//	      "jobPosition": "CTO",
//	      "email": "sari@digiinno.com",
//	      "phone": "021-5551235",
//	      "mobile": "0813-4567-8901",
//	      "isMain": false,
//	      "active": true
//	    }
//	  ],
//	  "structures": [
//	    {
//	      "tempKey": "1",
//	      "parentKey": null,
//	      "name": "Board of Directors",
//	      "level": 1,
//	      "address": "Jakarta",
//	      "active": true
//	    },
//	    {
//	      "tempKey": "2",
//	      "parentKey": "1",
//	      "name": "Technology Division",
//	      "level": 2,
//	      "address": "Jakarta",
//	      "active": true
//	    }
//	  ],
//	  "groups": {
//	    "industryId": "1",
//	    "industryActive": true,
//	    "parentGroupId": "2",
//	    "parentGroupActive": true
//	  },
//	  "others": [
//	    {
//	      "key": "company_size",
//	      "value": "50-100 employees",
//	      "active": true
//	    },
//	    {
//	      "key": "established_year",
//	      "value": "2015",
//	      "active": true
//	    }
//	  ]
//	})
//
// @Success 201 {object} entity.Customer
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers [post]
func CreateCustomer(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Di dalam fungsi CreateCustomer, tambahkan setelah Logo assignment:
	// Create customer entity
	customer := entity.Customer{
		Name:             req.Name,
		BrandName:        req.BrandName,
		Code:             req.Code,
		AccountManagerId: req.AccountManagerId,
		Status:           "Active", // Default status
	}

	// Set logo if provided
	if req.Logo != nil {
		customer.Logo = *req.Logo
	}

	// Set logo_small if provided
	if req.LogoSmall != nil {
		customer.LogoSmall = *req.LogoSmall
	}

	if err := tx.Create(&customer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer: " + err.Error()})
		return
	}

	// Create addresses
	// Create addresses
	for _, addrReq := range req.Addresses {
		address := entity.Address{
			CustomerID: customer.ID,
			// SupplierID field removed as it doesn't exist in entity.Address
			Name:    addrReq.Name,
			Address: addrReq.Address,
			Main:    addrReq.IsMain,
			Active:  addrReq.Active,
		}
		if err := tx.Create(&address).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address: " + err.Error()})
			return
		}
	}

	// Create social media
	for _, socialReq := range req.Socials {
		sosmed := entity.Sosmed{
			CustomerID: customer.ID,
			Name:       socialReq.Platform, // Atau buat field Name di DTO
			Platform:   socialReq.Platform,
			Handle:     socialReq.Handle,
			Active:     socialReq.Active,
		}
		if err := tx.Create(&sosmed).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social media: " + err.Error()})
			return
		}
	}

	// Create contacts
	for _, contactReq := range req.Contacts {
		contact := entity.Contact{
			CustomerID:  customer.ID,
			Name:        contactReq.Name,
			JobPosition: contactReq.JobPosition,
			Email:       contactReq.Email,
			Phone:       contactReq.Phone,
			Mobile:      contactReq.Mobile,
			Main:        contactReq.IsMain,
			Active:      contactReq.Active,
		}

		// Parse birthdate if provided
		if contactReq.Birthdate != "" {
			if birthdate, err := time.Parse("2006-01-02", contactReq.Birthdate); err == nil {
				contact.Birthdate = &birthdate
			}
		}

		if err := tx.Create(&contact).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact: " + err.Error()})
			return
		}
	}

	// Create structures with hierarchy
	tempKeyMap := make(map[string]uint)
	for _, structReq := range req.Structures {
		structure := entity.Structure{
			CustomerID: customer.ID,
			Name:       structReq.Name,
			Level:      structReq.Level,
			Address:    structReq.Address,
			Active:     structReq.Active,
		}

		// Set parent if exists
		if structReq.ParentKey != nil {
			if parentID, exists := tempKeyMap[*structReq.ParentKey]; exists {
				structure.ParentID = &parentID
			}
		}

		if err := tx.Create(&structure).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create structure: " + err.Error()})
			return
		}

		// Store temp key mapping
		tempKeyMap[structReq.TempKey] = structure.ID
	}

	// Create others
	for _, otherReq := range req.Others {
		other := entity.Other{
			CustomerID: customer.ID,
			Key:        otherReq.Key,
			Value:      otherReq.Value,
			Active:     otherReq.Active,
		}
		if err := tx.Create(&other).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create other attribute: " + err.Error()})
			return
		}
	}

	// Handle groups (industry and parent group)
	// Note: This assumes groups already exist in the database
	// Hapus: for _, group := range req.Groups {
	if req.Groups.IndustryID != "" && req.Groups.IndustryActive {
		// Find industry group and associate
		var industryGroup entity.Group
		if err := tx.Where("id = ?", req.Groups.IndustryID).First(&industryGroup).Error; err == nil {
			tx.Model(&customer).Association("Groups").Append(&industryGroup)
		}
	}

	if req.Groups.ParentGroupID != "" && req.Groups.ParentGroupActive {
		// Find parent group and associate
		var parentGroup entity.Group
		if err := tx.Where("id = ?", req.Groups.ParentGroupID).First(&parentGroup).Error; err == nil {
			tx.Model(&customer).Association("Groups").Append(&parentGroup)
		}
	}
	// Hapus: }

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// Load customer with all relations for response
	var createdCustomer entity.Customer
	config.DB.Preload("Addresses").Preload("Sosmeds").Preload("Contacts").Preload("Structures").Preload("Groups").Preload("Others").First(&createdCustomer, customer.ID)

	// Mapping manual untuk response
	response := dto.CustomerResponse{
		ID:               createdCustomer.ID,
		Name:             createdCustomer.Name,
		BrandName:        createdCustomer.BrandName,
		Code:             createdCustomer.Code,
		AccountManagerId: createdCustomer.AccountManagerId,
		/*  Email:            createdCustomer.Email,
		    Phone:            createdCustomer.Phone,
		    Website:          createdCustomer.Website,
		    Description:      createdCustomer.Description, */
		Logo:        createdCustomer.Logo,
		LogoSmall:   createdCustomer.LogoSmall,
		Status:      createdCustomer.Status,
		Category:    createdCustomer.Category,
		Rating:      createdCustomer.Rating,
		AverageCost: createdCustomer.AverageCost,
		CreatedAt:   createdCustomer.CreatedAt,
		UpdatedAt:   createdCustomer.UpdatedAt,
	}

	// Mapping addresses
	for _, addr := range createdCustomer.Addresses {
		response.Addresses = append(response.Addresses, dto.AddressResponse{
			Name:    addr.Name,
			Address: addr.Address,
			IsMain:  addr.Main,
			Active:  addr.Active,
		})
	}

	// Mapping contacts
	for _, contact := range createdCustomer.Contacts {
		response.Contacts = append(response.Contacts, dto.ContactResponse{
			Name:        contact.Name,
			JobPosition: contact.JobPosition,
			Email:       contact.Email,
			Phone:       contact.Phone,
			Active:      contact.Active,
		})
	}

	// Mapping others
	for _, other := range createdCustomer.Others {
		var valueStr string
		if other.Value != nil {
			valueStr = *other.Value
		}
		response.Others = append(response.Others, dto.OtherResponse{
			Key:    other.Key,
			Value:  valueStr,
			Active: other.Active,
		})
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Get customer by ID
// @Description Get a specific customer by ID
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} dto.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/{id} [get]
func GetCustomer(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	var customer entity.Customer
	result := config.DB.First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&customer)
	c.JSON(http.StatusOK, customer)
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&entity.Customer{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func UploadCustomerLogo(c *gin.Context) {
	id := c.Param("id")

	// Check if customer exists
	var customer entity.Customer
	result := config.DB.First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".svg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, PNG, and SVG files are allowed"})
		return
	}

	// Generate unique filename
	filename := "logo_" + id + "_" + time.Now().Format("20060102150405") + ext
	logoPath := "uploads/logos/" + filename

	// Save file
	if err := c.SaveUploadedFile(file, logoPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update customer logo path
	customer.Logo = logoPath
	config.DB.Save(&customer)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Logo uploaded successfully",
		"logo_path": logoPath,
		"customer":  customer,
	})
}

// @Summary Upload customer logo small
// @Description Upload a small logo/icon for customer (PNG, SVG, JPG)
// @Tags Customers
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param logo_small formData file true "Logo small file (PNG, SVG, JPG)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/{id}/logo-small [post]
func UploadCustomerLogoSmall(c *gin.Context) {
	id := c.Param("id")

	// Check if customer exists
	var customer entity.Customer
	result := config.DB.First(&customer, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("logo_small")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".svg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, PNG, and SVG files are allowed"})
		return
	}

	// Validate file size (max 2MB for small logos)
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size must be less than 2MB"})
		return
	}

	// Generate unique filename
	filename := "logo_small_" + id + "_" + time.Now().Format("20060102150405") + ext
	logoSmallPath := "uploads/logos_small/" + filename

	// Save file
	if err := c.SaveUploadedFile(file, logoSmallPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update customer logo_small path
	customer.LogoSmall = logoSmallPath
	config.DB.Save(&customer)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Logo small uploaded successfully",
		"logo_small_path": logoSmallPath,
		"customer":        customer,
	})
}
