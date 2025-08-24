package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStagesRoutes(r *gin.RouterGroup) {
	r.POST("/stages", handler.CreateStage)
	r.GET("/stages", handler.GetStages)
	r.GET("/stages/:id", handler.GetStage)
	r.PUT("/stages/:id", handler.UpdateStage)
	r.DELETE("/stages/:id", handler.DeleteStage)

	// detail stages
	r.GET("/stages/:id/details", handler.GetStageDetails)                  // Ambil semua detail stage
	r.POST("/stages/:id/details", handler.CreateStageDetail)               // Buat detail stage baru
	r.GET("/stages/:id/details/:detail_id", handler.GetStageDetail)        // Ambil detail stage tertentu
	r.PUT("/stages/:id/details/:detail_id", handler.UpdateStageDetail)     // Update detail stage
	r.DELETE("/stages/:id/details/:detail_id", handler.DeleteStageDetail)  // Hapus detail stage

}
