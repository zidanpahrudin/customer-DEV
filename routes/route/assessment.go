package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAssessmentRoutes(r *gin.RouterGroup) {
	r.GET("/assessment", handler.GetAssessments)
	r.GET("/assessment/:id", handler.GetAssessment)
	r.POST("/assessment", handler.CreateAssessment)
	r.PUT("/assessment/:id", handler.UpdateAssessment)
	r.DELETE("/assessment/:id", handler.DeleteAssessment)

	// detail routes
	r.GET("/assessment/:id/details", handler.GetAssessmentDetail)
	r.POST("/assessment/:id/details", handler.CreateAssessmentDetail)
	r.PUT("/assessment/:id/details/:detail_id", handler.UpdateAssessmentDetail)
	r.DELETE("/assessment/:id/details/:detail_id", handler.DeleteAssessmentDetail)

}
