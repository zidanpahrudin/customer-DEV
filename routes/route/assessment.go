package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAssessmentRoutes(r *gin.RouterGroup) {
	r.GET("/assessment", handler.GetAssessment)
	r.POST("/assessment", handler.CreateAssessment)
	r.PUT("/assessment/:id", handler.UpdateAssessment)
	r.DELETE("/assessment/:id", handler.DeleteAssessment)

	// detail routes
	r.GET("/assessment/:id/detail", handler.GetAssessmentDetail)
	r.POST("/assessment/:id/detail", handler.CreateAssessmentDetail)
	r.PUT("/assessment/:id/detail/:detail_id", handler.UpdateAssessmentDetail)
	r.DELETE("/assessment/:id/detail/:detail_id", handler.DeleteAssessmentDetail)
	
}
