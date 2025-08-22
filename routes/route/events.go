package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventsRoutes(r *gin.RouterGroup) {
	r.POST("/events", handler.CreateEvents)
	r.GET("/events", handler.ReadEvents)
	r.GET("/events/:id", handler.ReadOneEvents)
	r.PUT("/events/:id", handler.UpdateEvents)
	r.DELETE("/events/:id", handler.DeleteEvents)
	r.GET("/customers/:id/events", handler.GetCustomerEvents)
	r.GET("/event/type/:type", handler.GetEventType)
	
}
