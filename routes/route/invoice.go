package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterInvoiceRoutes(r *gin.RouterGroup) {
	r.POST("/invoices", handler.CreateInvoice)
	r.GET("/invoices", handler.GetInvoices)
	r.GET("/invoices/:id", handler.GetInvoice)
	r.PUT("/invoices/:id", handler.UpdateInvoice)
	r.DELETE("/invoices/:id", handler.DeleteInvoice)

	// Invoice-specific payments
	r.GET("/invoices/:id/payments", handler.GetInvoicePayments)
	r.POST("/invoices/:id/payments", handler.CreateInvoicePayment)
}
