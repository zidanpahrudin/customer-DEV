package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"fmt"
	"gorm.io/gorm"
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

	// Create customer
	customer := entity.Customer{
		Name:             req.Name,
		BrandName:        req.BrandName,
		Code:             req.Code,
		AccountManagerId: req.AccountManagerId,
		Status:           req.StatusName, // Gunakan req.StatusName dari request body
	}

	if req.Logo != nil {
		customer.Logo = *req.Logo
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
	for _, group := range req.Groups {
		if group.IndustryID != "" && group.IndustryActive {
			// Find industry group and associate
			var industryGroup entity.Group
			if err := tx.Where("id = ?", group.IndustryID).First(&industryGroup).Error; err == nil {
				tx.Model(&customer).Association("Groups").Append(&industryGroup)
			}
		}

		if group.ParentGroupID != "" && group.ParentGroupActive {
			// Find parent group and associate
			var parentGroup entity.Group
			if err := tx.Where("id = ?", group.ParentGroupID).First(&parentGroup).Error; err == nil {
				tx.Model(&customer).Association("Groups").Append(&parentGroup)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// Load customer with all relations for response
	var createdCustomer entity.Customer
	config.DB.Preload("Addresses").Preload("Sosmeds").Preload("Contacts").Preload("Structures").Preload("Groups").Preload("Others").First(&createdCustomer, customer.ID)

	c.JSON(http.StatusCreated, createdCustomer)
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


// @Summary Update block Customer status
// @Description Update the status of a customer to Blocked
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} dto.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/status/{id} [post]
func UpdateCustomerStatus(c *gin.Context) {

	// Ambil user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	id := c.Param("id")
	status := c.PostForm("status")   // ambil status dari form
	reason := c.PostForm("reason")   // alasan perubahan status
	notes := c.PostForm("notes")     // catatan untuk dokumen

	// Validasi status
	if status != "active" && status != "blocked" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Cari customer
	var customer entity.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Simpan reason ke StatusReasons
	statusReason := entity.StatusReasons{
		CustomerID: customer.ID,
		Reason:     reason,
		Status:     status,
	}
	config.DB.Create(&statusReason)

	// Handle file upload
	file, err := c.FormFile("file")
	var document entity.Document
	if err == nil { // kalau ada file
		// Simpan file ke folder uploads/
		filePath := fmt.Sprintf("uploads/documents/%d_%s", customer.ID, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		// Simpan record document
		document = entity.Document{
			CustomerID: customer.ID,
			UserID:		uid,
			Notes:      notes,
			Type:       "StatusChange",
			URLFile:    filePath,
		}
		config.DB.Create(&document)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message":       "Customer status updated to " + status,
		"customer":      customer,
		"status_reason": statusReason,
		"document":      document,
	})
}


// @Summary Get Customer Status Reason And Document
// @Description Get the status reason and document for a customer
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} dto.CustomerStatusResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/customers/status/{id} [get]
func GetCustomerStatus(c *gin.Context) {
	id := c.Param("id")

	// Cari customer
	var customer entity.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Ambil status reasons
	var statusReasons []entity.StatusReasons
	if err := config.DB.Where("customer_id = ? AND is_active = ?", customer.ID, true).Find(&statusReasons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch status reasons"})
		return
	}

	// Ambil documents terkait status perubahan
	var documents []entity.Document
	if err := config.DB.Where("customer_id = ? AND is_active = ?", customer.ID, true).Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Success Get Customer status",
		"customer":      customer,
		"StatusReasons": statusReasons,
		"document":      documents,
	})
}


// @Summary Get customer statistics
// @Description Get statistics about customers including total count, new customers in the last year, average cost, and blocked customers
// @Tags Customers
// @Param status query string false "Filter by status" Enums(Active, Inactive, Blocked)
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CustomerStatsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/customers/stats [get]
func GetCustomerStats(c *gin.Context) {
	status := c.Query("status")

	// Helper: apply filter & range tahun
	queryWithFilter := func(base *gorm.DB, status string, start, end *time.Time) *gorm.DB {
		q := base
		if status != "" {
			q = q.Where("status = ?", status)
		}
		if start != nil && end != nil {
			q = q.Where("created_at BETWEEN ? AND ?", *start, *end)
		}
		return q
	}

	// Periode
	now := time.Now()
	thisYearStart := now.AddDate(-1, 0, 0) // 1 tahun ke belakang
	lastYearStart := now.AddDate(-2, 0, 0) // 2 tahun ke belakang
	lastYearEnd := now.AddDate(-1, 0, 0)

	var totalThisYear, totalLastYear int64
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, &thisYearStart, &now).
		Count(&totalThisYear)
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, &lastYearStart, &lastYearEnd).
		Count(&totalLastYear)

	var newThisYear, newLastYear int64
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, &thisYearStart, &now).
		Count(&newThisYear)
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, &lastYearStart, &lastYearEnd).
		Count(&newLastYear)

	var avgThisYear, avgLastYear float64
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, &thisYearStart, &now).
		Select("COALESCE(AVG(average_cost), 0)").Row().Scan(&avgThisYear)
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, &lastYearStart, &lastYearEnd).
		Select("COALESCE(AVG(average_cost), 0)").Row().Scan(&avgLastYear)

	var churnThisYear, churnLastYear int64
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, nil, nil).
		Where("last_transaction_at < ?", thisYearStart).Count(&churnThisYear)
	queryWithFilter(config.DB.Model(&entity.Customer{}), status, nil, nil).
		Where("last_transaction_at BETWEEN ? AND ?", lastYearStart, lastYearEnd).Count(&churnLastYear)

	calcGrowth := func(thisYear, lastYear float64) float64 {
		if lastYear == 0 {
			if thisYear > 0 {
				return 100.0
			}
			return 0
		}
		return ((thisYear - lastYear) / lastYear) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Customer statistics fetched successfully",
		"data": gin.H{
			"total_customers": gin.H{
				"value":  totalThisYear,
				"growth": calcGrowth(float64(totalThisYear), float64(totalLastYear)),
			},
			"new_customers": gin.H{
				"value":  newThisYear,
				"growth": calcGrowth(float64(newThisYear), float64(newLastYear)),
			},
			"avg_revenue": gin.H{
				"value":  avgThisYear,
				"growth": calcGrowth(avgThisYear, avgLastYear),
			},
			"churn_customers": gin.H{
				"value":  churnThisYear,
				"growth": calcGrowth(float64(churnThisYear), float64(churnLastYear)),
			},
		},
	})
}

