package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStatusRoutes(r *gin.RouterGroup) {
	r.POST("/statuses", handler.CreateStatus)
	r.GET("/statuses", handler.GetStatuses)
	r.GET("/statuses/:id", handler.GetStatus)
	r.PUT("/statuses/:id", handler.UpdateStatus)
	r.DELETE("/statuses/:id", handler.DeleteStatus)
}
