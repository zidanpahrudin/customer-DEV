package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(r *gin.RouterGroup) {
	r.POST("/roles", handler.CreateRole)
	r.GET("/roles", handler.GetRoles)
	r.GET("/roles/:id", handler.GetRole)
	r.PUT("/roles/:id", handler.UpdateRole)
	r.DELETE("/roles/:id", handler.DeleteRole)
}
