package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAddressRoutes(r *gin.RouterGroup) {
	r.GET("/customers/:id/addresses", handler.GetCustomerAddresses)
	r.GET("/addresses/:id", handler.GetAddress)
	r.PUT("/addresses/:id", handler.UpdateAddress)
	r.DELETE("/addresses/:id", handler.DeleteAddress)
	r.POST("/addresses", handler.CreateAddress)
}
