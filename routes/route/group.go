package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterGroupRoutes(r *gin.RouterGroup) {
	r.GET("/groups", handler.GetGroups)
	r.GET("/groups/:id", handler.GetGroup)
	r.PUT("/groups/:id", handler.UpdateGroup)
	r.DELETE("/groups/:id", handler.DeleteGroup)
	r.GET("/groups/:id/customers", handler.GetGroupCustomers)

	// Nested resource
	r.PUT("/groups/:id/customers/:customer_id", handler.AssignCustomerToGroup)
	r.DELETE("/groups/:id/customers/:customer_id", handler.RemoveCustomerFromGroup)

	r.POST("/groups", handler.CreateGroup)
}
