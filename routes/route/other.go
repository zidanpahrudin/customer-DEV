package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterOtherRoutes(r *gin.RouterGroup) {
	r.GET("/others/:id", handler.GetOther)
	r.DELETE("/others/:id", handler.DeleteOther)
	r.GET("/others/by-attribute", handler.GetOthersByAttribute)
	r.POST("/others", handler.CreateOther)
}
