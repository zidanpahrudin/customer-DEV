package routes

import (
	"customer-api/internal/handler"
	"customer-api/middleware"
	"customer-api/routes/route"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// Public routes
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// Register all modules
	route.RegisterRoleRoutes(protected)
	route.RegisterCustomerRoutes(protected)
	route.RegisterAddressRoutes(protected)
	route.RegisterSosmedRoutes(protected)
	route.RegisterContactRoutes(protected)
	route.RegisterStructureRoutes(protected)
	route.RegisterGroupRoutes(protected)
	route.RegisterOtherRoutes(protected)
	route.RegisterActivityRoutes(protected)
	route.RegisterInvoiceRoutes(protected)
	route.RegisterPaymentRoutes(protected)
	route.RegisterStatusRoutes(protected)
	route.RegisterEventsRoutes(protected)
	route.RegisterActivityTypeRoutes(protected)

	
}
