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
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"bytes"
	// "strconv"
	
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
	config.DB.Model(&entity.Customer{}).Where("status = ?", "blocked").Count(&blockedCustomers)

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

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Type assertion ke string
	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID in context has invalid type"})
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
		Name:             *req.Name,
		BrandName:        *req.BrandName,
		Code:             *req.Code,
		AccountManagerId: *req.AccountManagerId,
		Status:           "Active", // Default status
	}
	

	// Set logo if provided
	if req.Logo != "" {
		customer.Logo = req.Logo
	}

	// Set logo_small if provided
	if req.LogoSmall != "" {
		customer.LogoSmall = req.LogoSmall
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
	tempKeyMap := make(map[string]string)
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

	// Insert HistoryCustomer
	history := entity.HistoryCustomer{
		CustomerID: customer.ID,
		UserID:     userID,
		Status:     "Created",
		Notes:     "Created new customer",
	}
	config.DB.Create(&history)



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
	userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Type assertion ke string
	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID in context has invalid type"})
		return
	}
	
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

	// Insert HistoryCustomer
	history := entity.HistoryCustomer{
		CustomerID: customer.ID,
		UserID:     userID,
		Status:     "Updated",
		Notes:     "Updated customer",
	}
	config.DB.Create(&history)

	c.JSON(http.StatusOK, customer)
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Type assertion ke string
	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID in context has invalid type"})
		return
	}

	result := config.DB.Delete(&entity.Customer{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	// Insert HistoryCustomer
	history := entity.HistoryCustomer{
		CustomerID: id,
		UserID:     userID,
		Status:     "Deleted",
		Notes:     "Deleted customer",
	}

	config.DB.Create(&history)

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func UploadCustomerLogo(c *gin.Context) {
	id := c.Param("id")

	userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Type assertion ke string
	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID in context has invalid type"})
		return
	}
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

	// Insert HistoryCustomer
	history := entity.HistoryCustomer{
		CustomerID: customer.ID,
		UserID:     userID,
		Status:     "Logo Uploaded",
		Notes:     "Uploaded logo for customer",
	}
	config.DB.Create(&history)

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

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Type assertion ke string
	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID in context has invalid type"})
		return
	}

	id := c.Param("id")
	status := c.PostForm("status")   // ambil status dari form
	reason := c.PostForm("reason")   // alasan perubahan status
	notes := c.PostForm("notes")     // catatan untuk dokumen

	// Validasi status
	if status != "active" && status != "blocked" && status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Cari customer
	var customer entity.Customer
	if err := config.DB.Where("id = ?", id).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}


	// update status customer
	customer.Status = status
	config.DB.Save(&customer)

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
			UserID:		userID,
			Notes:      notes,
			Type:       "StatusChange",
			URLFile:    filePath,
		}
		config.DB.Create(&document)
	}

	// Insert HistoryCustomer
	history := entity.HistoryCustomer{
		CustomerID: customer.ID,
		UserID:     userID,
		Status:     "Status Changed",
		Notes:     "Changed status to " + status,
	}
	config.DB.Create(&history)

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


// ExportCustomers handles customer export to Excel or PDF
func ExportCustomers(c *gin.Context) {
	exportType := c.Query("type")

	// validasi export type
	if exportType != "excel" && exportType != "pdf" {
		sendError(c, http.StatusBadRequest, "Invalid export type (must be 'excel' or 'pdf')")
		return
	}

	// ambil data customer
	var customers []entity.Customer
	if err := config.DB.Find(&customers).Error; err != nil {
		sendError(c, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}

	switch exportType {
	case "excel":
		file, err := createExcelFile(customers)
		if err != nil {
			sendError(c, http.StatusInternalServerError, "Failed to create excel file: "+err.Error())
			return
		}
		sendFile(c, file, "customers.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	case "pdf":
		file, err := createPDFFile(customers)
		if err != nil {
			sendError(c, http.StatusInternalServerError, "Failed to create pdf file: "+err.Error())
			return
		}
		sendFile(c, file, "customers.pdf", "application/pdf")
	}
}

// helper untuk buat Excel file
func createExcelFile(customers []entity.Customer) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Customers"
	f.NewSheet(sheet)

	// header
	headers := []string{"ID", "Name", "Brand Name", "Code", "Status"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// data
	for row, cust := range customers {
		values := []interface{}{cust.ID, cust.Name, cust.BrandName, cust.Code, cust.Status}
		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2)
			f.SetCellValue(sheet, cell, val)
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// helper untuk buat PDF file
func createPDFFile(customers []entity.Customer) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	// header
	headers := []string{"ID", "Name", "Brand Name", "Code", "Status"}
	for _, h := range headers {
		pdf.CellFormat(40, 10, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// data
	pdf.SetFont("Arial", "", 10)
	for _, cust := range customers {
		pdf.CellFormat(40, 10, cust.ID, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, cust.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, cust.BrandName, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, cust.Code, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, cust.Status, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// helper untuk kirim response error
func sendError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "failed",
		"message": message,
		"data":    nil,
	})
}

// helper untuk kirim file
func sendFile(c *gin.Context, data []byte, filename, contentType string) {
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Header("Expires", "0")
	c.Writer.Write(data)
}

// get historyCustomer by UserID
func getHistoryCustomerByUserID(userID string) ([]entity.HistoryCustomer, error) {
	var history []entity.HistoryCustomer
	if err := config.DB.Where("user_id = ?", userID).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}