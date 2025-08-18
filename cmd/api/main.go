package main

import (
	"customer-api/internal/config"
	"customer-api/internal/handler"
	"customer-api/middleware"

	_ "customer-api/cmd/api/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Customer Management API
// @version 1.0
// @description API untuk manajemen customer menggunakan Golang, Gin, dan PostgreSQL
// @host localhost:9000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type \"Bearer\" followed by a space and JWT token.
func main() {
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.String(200, "server sukses berjalan")
	})

	// Initialize Database
	config.ConnectDatabase()

	// Public routes
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Role routes
		protected.POST("/roles", handler.CreateRole)
		protected.GET("/roles", handler.GetRoles)
		protected.GET("/roles/:id", handler.GetRole)
		protected.PUT("/roles/:id", handler.UpdateRole)
		protected.DELETE("/roles/:id", handler.DeleteRole)

		// Customer routes
		protected.POST("/customers", handler.CreateCustomer)
		protected.GET("/customers", handler.GetCustomers)
		protected.GET("/customers/:id", handler.GetCustomer)
		protected.GET("/customers/:id/with-addresses", handler.GetCustomerWithAddresses)
		protected.GET("/customers/:id/with-sosmeds", handler.GetCustomerWithSosmeds)
		protected.GET("/customers/:id/with-contacts", handler.GetCustomerWithContacts)
		protected.GET("/customers/:id/with-structures", handler.GetCustomerWithStructures)
		protected.GET("/customers/:id/with-all", handler.GetCustomerWithAllRelations)
		protected.GET("/customers/:id/full", handler.GetCustomerFull)
		protected.PUT("/customers/:id", handler.UpdateCustomer)
		protected.DELETE("/customers/:id", handler.DeleteCustomer)
		protected.POST("/customers/:id/logo", handler.UploadCustomerLogo)
		protected.POST("/customers/:id/logo-small", handler.UploadCustomerLogoSmall) // Route baru

		// Address routes
		// protected.POST("/customers/:id/addresses", handler.CreateAddress) // DIHAPUS
		protected.GET("/customers/:id/addresses", handler.GetCustomerAddresses)
		protected.GET("/addresses/:id", handler.GetAddress)
		protected.PUT("/addresses/:id", handler.UpdateAddress)
		protected.DELETE("/addresses/:id", handler.DeleteAddress)

		// Sosmed routes
		// protected.POST("/customers/:id/sosmeds", handler.CreateSosmed) // DIHAPUS
		protected.GET("/customers/:id/sosmeds", handler.GetCustomerSosmeds)
		protected.GET("/sosmeds/:id", handler.GetSosmed)
		protected.PUT("/sosmeds/:id", handler.UpdateSosmed)
		protected.DELETE("/sosmeds/:id", handler.DeleteSosmed)

		// Contact routes
		// protected.POST("/customers/:id/contacts", handler.CreateContact) // DIHAPUS
		protected.GET("/customers/:id/contacts", handler.GetCustomerContacts)
		protected.GET("/contacts/:id", handler.GetContact)
		protected.PUT("/contacts/:id", handler.UpdateContact)
		protected.DELETE("/contacts/:id", handler.DeleteContact)

		// Structure routes
		// protected.POST("/customers/:id/structures", handler.CreateStructure) // DIHAPUS
		protected.GET("/customers/:id/structures", handler.GetCustomerStructures)
		protected.GET("/customers/:id/structures/by-level", handler.GetStructuresByLevel)
		protected.GET("/structures/:id", handler.GetStructure)
		protected.PUT("/structures/:id", handler.UpdateStructure)
		protected.DELETE("/structures/:id", handler.DeleteStructure)

		// Group routes
		// protected.POST("/groups", handler.CreateGroup) // DIHAPUS
		protected.GET("/groups", handler.GetGroups)
		protected.GET("/groups/:id", handler.GetGroup)
		protected.PUT("/groups/:id", handler.UpdateGroup)
		protected.DELETE("/groups/:id", handler.DeleteGroup)
		protected.GET("/groups/:id/customers", handler.GetGroupCustomers)
		// Gunakan nested resource approach
		protected.PUT("/groups/:id/customers/:customer_id", handler.AssignCustomerToGroup)
		protected.DELETE("/groups/:id/customers/:customer_id", handler.RemoveCustomerFromGroup)

		// Other routes
		// protected.POST("/customers/:id/others", handler.CreateOther) // DIHAPUS
		protected.GET("/customers/:id/others", handler.GetCustomerOthers)
		protected.GET("/customers/:id/with-others", handler.GetCustomerWithOthers)
		protected.GET("/others/:id", handler.GetOther)
		/* protected.PUT("/others/:id", handler.UpdateOther) */
		protected.DELETE("/others/:id", handler.DeleteOther)
		protected.GET("/others/by-attribute", handler.GetOthersByAttribute)

		// Activity routes
		protected.POST("/activities", handler.CreateActivity)
		protected.GET("/activities", handler.GetActivities)
		protected.GET("/activities/:id", handler.GetActivity)
		protected.PUT("/activities/:id", handler.UpdateActivity)
		protected.DELETE("/activities/:id", handler.DeleteActivity)

		// Activity attendees routes
		protected.POST("/activities/:id/attendees", handler.AddActivityAttendees)
		protected.DELETE("/activities/:id/attendees", handler.RemoveActivityAttendees)

		// Activity check-in routes
		protected.POST("/activities/:id/checkin", handler.CheckinActivity)

		// Invoice routes
		protected.POST("/invoices", handler.CreateInvoice)
		protected.GET("/invoices", handler.GetInvoices)
		protected.GET("/invoices/:id", handler.GetInvoice)
		protected.PUT("/invoices/:id", handler.UpdateInvoice)
		protected.DELETE("/invoices/:id", handler.DeleteInvoice)

		// Payment routes
		protected.POST("/payments", handler.CreatePayment)
		protected.GET("/payments", handler.GetPayments)
		protected.GET("/payments/:id", handler.GetPayment)
		protected.PUT("/payments/:id", handler.UpdatePayment)
		protected.DELETE("/payments/:id", handler.DeletePayment)

		// Invoice-specific payment routes
		protected.GET("/invoices/:id/payments", handler.GetInvoicePayments)
		protected.POST("/invoices/:id/payments", handler.CreateInvoicePayment)

		// Status routes
		protected.POST("/statuses", handler.CreateStatus)
		protected.GET("/statuses", handler.GetStatuses)
		protected.GET("/statuses/:id", handler.GetStatus)
		protected.PUT("/statuses/:id", handler.UpdateStatus)
		protected.DELETE("/statuses/:id", handler.DeleteStatus)

		// Customer-specific status routes
		protected.GET("/customers/:id/statuses", handler.GetCustomersByStatus)

		// Address routes
		protected.POST("/addresses", handler.CreateAddress)

		// Contact routes
		protected.POST("/contacts", handler.CreateContact)

		// Structure routes
		protected.POST("/structures", handler.CreateStructure)

		// Group routes
		protected.POST("/groups", handler.CreateGroup)

		// Other routes
		protected.POST("/others", handler.CreateOther)
	}

	// Serve static files for logos
	r.Static("/uploads", "./uploads")

	// Run the server
	r.Run(":9000")
}
