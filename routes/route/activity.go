package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterActivityRoutes(r *gin.RouterGroup) {
	r.POST("/activities", handler.CreateActivity)
	r.GET("/activities", handler.GetActivities)
	r.GET("/activities/:id", handler.GetActivity)
	r.PUT("/activities/:id", handler.UpdateActivity)
	r.DELETE("/activities/:id", handler.DeleteActivity)

	// Attendees
	r.POST("/activities/:id/attendees", handler.AddActivityAttendees)
	r.DELETE("/activities/:id/attendees", handler.RemoveActivityAttendees)

	// Check-in
	r.POST("/activities/:id/checkin", handler.CheckinActivity)
}
