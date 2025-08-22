package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterActivityTypeRoutes(r *gin.RouterGroup) {
	r.POST("/activity_types", handler.CreateActivityType)
	r.GET("/activity_types", handler.ReadActivityTypes)
	r.GET("/activity_types/:id", handler.ReadActivityType)
	r.PUT("/activity_types/:id", handler.UpdateActivityType)
	r.DELETE("/activity_types/:id", handler.DeleteActivityType)
	r.GET("/activity_types/:id/activities", handler.ReadActivitiesByActivityType)	
}
