package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterContactRoutes(r *gin.RouterGroup) {
	r.GET("/customers/:id/contacts", handler.GetCustomerContacts)
	r.GET("/contacts/:id", handler.GetContact)
	r.PUT("/contacts/:id", handler.UpdateContact)
	r.DELETE("/contacts/:id", handler.DeleteContact)
	r.POST("/contacts", handler.CreateContact)
}
