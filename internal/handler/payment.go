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

// GetPayments - GET /payments
func GetPayments(c *gin.Context) {
	var payments []entity.Payment
	var total int64

	// Query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	invoiceID := c.Query("invoice_id")

	offset := (page - 1) * limit

	// Build query
	db := config.DB
	query := db.Model(&entity.Payment{})

	if invoiceID != "" {
		query = query.Where("invoice_id = ?", invoiceID)
	}

	// Count total
	query.Count(&total)

	// Get payments with pagination
	if err := query.Preload("Invoice").Offset(offset).Limit(limit).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	// Convert to response format
	var paymentResponses []dto.PaymentResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, convertToPaymentResponse(payment))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": paymentResponses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CreatePayment - POST /payments
func CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	// Check if invoice exists
	var invoice entity.Invoice
	if err := db.First(&invoice, req.InvoiceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice not found"})
		return
	}

	// Create payment
	payment := entity.Payment{
		InvoiceID: req.InvoiceID,
		Amount:    req.Amount,
		PaidAt:    req.PaidAt,
	}

	if err := db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	// Update invoice paid amount
	invoice.PaidAmount += req.Amount
	db.Save(&invoice)

	// Load invoice for response
	db.Preload("Invoice").First(&payment, payment.ID)

	response := convertToPaymentResponse(payment)
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// GetPayment - GET /payments/:id
func GetPayment(c *gin.Context) {
	id := c.Param("id")
	var payment entity.Payment

	db := config.DB
	if err := db.Preload("Invoice").First(&payment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment"})
		}
		return
	}

	response := convertToPaymentResponse(payment)
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// UpdatePayment - PUT /payments/:id
func UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB
	var payment entity.Payment
	if err := db.First(&payment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment"})
		}
		return
	}

	// Store old amount for invoice update
	oldAmount := payment.Amount

	// Update fields
	if req.Amount != nil {
		payment.Amount = *req.Amount
	}
	if req.PaidAt != nil {
		payment.PaidAt = *req.PaidAt
	}

	if err := db.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	// Update invoice paid amount if amount changed
	if req.Amount != nil {
		var invoice entity.Invoice
		if err := db.First(&invoice, payment.InvoiceID).Error; err == nil {
			invoice.PaidAmount = invoice.PaidAmount - oldAmount + payment.Amount
			db.Save(&invoice)
		}
	}

	// Load invoice for response
	db.Preload("Invoice").First(&payment, payment.ID)

	response := convertToPaymentResponse(payment)
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// DeletePayment - DELETE /payments/:id
func DeletePayment(c *gin.Context) {
	id := c.Param("id")
	var payment entity.Payment

	db := config.DB
	if err := db.First(&payment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment"})
		}
		return
	}

	// Update invoice paid amount
	var invoice entity.Invoice
	if err := db.First(&invoice, payment.InvoiceID).Error; err == nil {
		invoice.PaidAmount -= payment.Amount
		db.Save(&invoice)
	}

	if err := db.Delete(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}

// GetInvoicePayments - GET /invoices/:id/payments
func GetInvoicePayments(c *gin.Context) {
	invoiceID := c.Param("id")
	var payments []entity.Payment

	db := config.DB
	if err := db.Where("invoice_id = ?", invoiceID).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	// Convert to response format
	var paymentResponses []dto.PaymentResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, convertToPaymentResponse(payment))
	}

	c.JSON(http.StatusOK, gin.H{"data": paymentResponses})
}

// CreateInvoicePayment - POST /invoices/:id/payments
func CreateInvoicePayment(c *gin.Context) {
	invoiceID := c.Param("id")
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Override invoice ID from URL
	invoiceIDUint, err := strconv.ParseUint(invoiceID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}
	req.InvoiceID = uint(invoiceIDUint)

	db := config.DB

	// Check if invoice exists
	var invoice entity.Invoice
	if err := db.First(&invoice, req.InvoiceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice not found"})
		return
	}

	// Create payment
	payment := entity.Payment{
		InvoiceID: req.InvoiceID,
		Amount:    req.Amount,
		PaidAt:    req.PaidAt,
	}

	if err := db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	// Update invoice paid amount
	invoice.PaidAmount += req.Amount
	db.Save(&invoice)

	// Load invoice for response
	db.Preload("Invoice").First(&payment, payment.ID)

	response := convertToPaymentResponse(payment)
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// Helper function to convert entity to response
func convertToPaymentResponse(payment entity.Payment) dto.PaymentResponse {
	response := dto.PaymentResponse{
		ID:        payment.ID,
		InvoiceID: payment.InvoiceID,
		Amount:    payment.Amount,
		PaidAt:    payment.PaidAt,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}

	// Add invoice if loaded
	if payment.Invoice.ID != 0 {
		invoiceResponse := &dto.InvoiceResponse{
			ID:            payment.Invoice.ID,
			CustomerID:    payment.Invoice.CustomerID,
			ProjectID:     payment.Invoice.ProjectID,
			InvoiceNumber: payment.Invoice.InvoiceNumber,
			Amount:        payment.Invoice.Amount,
			IssuedDate:    payment.Invoice.IssuedDate,
			DueDate:       payment.Invoice.DueDate,
			PaidAmount:    payment.Invoice.PaidAmount,
			Balance:       payment.Invoice.Amount - payment.Invoice.PaidAmount,
			CreatedAt:     payment.Invoice.CreatedAt,
			UpdatedAt:     payment.Invoice.UpdatedAt,
		}

		// Set status
		if payment.Invoice.PaidAmount >= payment.Invoice.Amount {
			invoiceResponse.Status = "paid"
		} else if payment.Invoice.PaidAmount > 0 {
			invoiceResponse.Status = "partial"
		} else if time.Now().After(payment.Invoice.DueDate) {
			invoiceResponse.Status = "overdue"
		} else {
			invoiceResponse.Status = "unpaid"
		}

		response.Invoice = invoiceResponse
	}

	return response
}