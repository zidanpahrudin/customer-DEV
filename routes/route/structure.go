package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStructureRoutes(r *gin.RouterGroup) {
	r.GET("/customers/:id/structures", handler.GetCustomerStructures)
	r.GET("/customers/:id/structures/by-level", handler.GetStructuresByLevel)
	r.GET("/structures/:id", handler.GetStructure)
	r.PUT("/structures/:id", handler.UpdateStructure)
	r.DELETE("/structures/:id", handler.DeleteStructure)
	r.POST("/structures", handler.CreateStructure)
}
