package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
)

// GetInvoices - GET /invoices
func GetInvoices(c *gin.Context) {
	var invoices []entity.Invoice
	var total int64

	// Query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	customerID := c.Query("customer_id")
	status := c.Query("status")

	offset := (page - 1) * limit

	// Build query
	db := config.DB
	query := db.Model(&entity.Invoice{})

	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	if status != "" {
		switch status {
		case "paid":
			query = query.Where("paid_amount >= amount")
		case "unpaid":
			query = query.Where("paid_amount = 0")
		case "partial":
			query = query.Where("paid_amount > 0 AND paid_amount < amount")
		}
	}

	// Count total
	query.Count(&total)

	// Get invoices with pagination
	if err := query.Preload("Customer").Offset(offset).Limit(limit).Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
		return
	}

	// Convert to response format
	var invoiceResponses []dto.InvoiceResponse
	for _, invoice := range invoices {
		invoiceResponses = append(invoiceResponses, convertToInvoiceResponse(invoice))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": invoiceResponses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CreateInvoice - POST /invoices
func CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	// Check if customer exists
	var customer entity.Customer
	if err := db.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found"})
		return
	}

	// Create invoice
	invoice := entity.Invoice{
		CustomerID:    req.CustomerID,
		ProjectID:     req.ProjectID,
		InvoiceNumber: req.InvoiceNumber,
		Amount:        req.Amount,
		IssuedDate:    req.IssuedDate,
		DueDate:       req.DueDate,
		PaidAmount:    req.PaidAmount,
	}

	if err := db.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	// Load customer for response
	db.Preload("Customer").First(&invoice, invoice.ID)

	response := convertToInvoiceResponse(invoice)
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// GetInvoice - GET /invoices/:id
func GetInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice entity.Invoice

	db := config.DB
	if err := db.Preload("Customer").First(&invoice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoice"})
		}
		return
	}

	response := convertToInvoiceResponse(invoice)
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdateInvoice - PUT /invoices/:id
func UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB
	var invoice entity.Invoice
	if err := db.First(&invoice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoice"})
		}
		return
	}

	// Update fields - CustomerID tidak ada di UpdateInvoiceRequest
	if req.ProjectID != nil {
		invoice.ProjectID = *req.ProjectID
	}
	if req.InvoiceNumber != nil {
		invoice.InvoiceNumber = *req.InvoiceNumber
	}
	if req.Amount != nil {
		invoice.Amount = *req.Amount
	}
	if req.IssuedDate != nil {
		invoice.IssuedDate = *req.IssuedDate
	}
	if req.DueDate != nil {
		invoice.DueDate = *req.DueDate
	}
	if req.PaidAmount != nil {
		invoice.PaidAmount = *req.PaidAmount
	}

	if err := db.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}

	// Load customer for response
	db.Preload("Customer").First(&invoice, invoice.ID)

	response := convertToInvoiceResponse(invoice)
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// DeleteInvoice - DELETE /invoices/:id
func DeleteInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice entity.Invoice

	db := config.DB
	if err := db.First(&invoice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoice"})
		}
		return
	}

	if err := db.Delete(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}

// Helper function to convert entity to response
func convertToInvoiceResponse(invoice entity.Invoice) dto.InvoiceResponse {
	balance := invoice.Amount - invoice.PaidAmount
	status := "unpaid"

	if invoice.PaidAmount >= invoice.Amount {
		status = "paid"
	} else if invoice.PaidAmount > 0 {
		status = "partial"
	} else if time.Now().After(invoice.DueDate) {
		status = "overdue"
	}

	response := dto.InvoiceResponse{
		ID:            invoice.ID,
		CustomerID:    invoice.CustomerID,
		ProjectID:     invoice.ProjectID,
		InvoiceNumber: invoice.InvoiceNumber,
		Amount:        invoice.Amount,
		IssuedDate:    invoice.IssuedDate,
		DueDate:       invoice.DueDate,
		PaidAmount:    invoice.PaidAmount,
		Balance:       balance,
		Status:        status,
		CreatedAt:     invoice.CreatedAt,
		UpdatedAt:     invoice.UpdatedAt,
	}

	// Add customer if loaded
	if invoice.Customer.ID != 0 {
		customerResponse := &dto.CustomerResponse{
			ID:          invoice.Customer.ID,
			Name:        invoice.Customer.Name,
			Email:       invoice.Customer.Email,
			Phone:       invoice.Customer.Phone,
			Website:     invoice.Customer.Website,
			Description: invoice.Customer.Description,
			Status:      invoice.Customer.Status,
			Category:    invoice.Customer.Category,
			Rating:      invoice.Customer.Rating,
			AverageCost: invoice.Customer.AverageCost,
		}
		response.Customer = customerResponse
	}

	return response
}
