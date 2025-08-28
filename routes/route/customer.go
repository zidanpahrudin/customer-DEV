package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCustomerRoutes(r *gin.RouterGroup) {
	
	r.POST("/customers", handler.CreateCustomer)
	r.GET("/customers", handler.GetCustomers)

	
	r.GET("/customers/statistics", handler.GetCustomerStats)
	// export data
	r.GET("/customers/export", handler.ExportCustomers)

	// history customer
	r.GET("/customers/:id/history", handler.getHistoryCustomerByUserID)


	r.GET("/customers/:id", handler.GetCustomer)
	r.GET("/customers/:id/with-addresses", handler.GetCustomerWithAddresses)
	r.GET("/customers/:id/with-sosmeds", handler.GetCustomerWithSosmeds)
	r.GET("/customers/:id/with-contacts", handler.GetCustomerWithContacts)
	r.GET("/customers/:id/with-structures", handler.GetCustomerWithStructures)
	r.GET("/customers/:id/with-all", handler.GetCustomerWithAllRelations)
	// r.GET("/customers/:id/full", handler.GetCustomerFull)
	r.PUT("/customers/:id", handler.UpdateCustomer)
	r.DELETE("/customers/:id", handler.DeleteCustomer)
	r.POST("/customers/:id/logo", handler.UploadCustomerLogo)

	// Customer status
	r.POST("/customers/:id/status", handler.UpdateCustomerStatus)
	r.GET("/customers/:id/status", handler.GetCustomerStatus)

	// Customer relations
	r.GET("/customers/:id/others", handler.GetCustomerOthers)
	r.GET("/customers/:id/with-others", handler.GetCustomerWithOthers)
	r.GET("/customers/:id/statuses", handler.GetCustomersByStatus)

	

	
}
