package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterGroupConfig(r *gin.RouterGroup) {
	r.POST("/group-configs", handler.CreateConfigGroup)
	r.GET("/group-configs", handler.GetConfigGroups)
	r.GET("/group-configs/:id", handler.GetConfigGroup)
	r.PUT("/group-configs/:id", handler.UpdateConfigGroup)
	r.DELETE("/group-configs/:id", handler.DeleteConfigGroup)
	// detail group-configs
	r.GET("/group-configs/:id/details", handler.GetConfigGroupDetails)                  // Amb
	r.POST("/group-configs/:id/details", handler.CreateConfigGroupDetail)               // Buat detail group-config baru
	r.GET("/group-configs/:id/details/:detail_id", handler.GetConfigGroupDetail)
	r.PUT("/group-configs/:id/details/:detail_id", handler.UpdateConfigGroupDetail)     // Update detail group-config
	r.DELETE("/group-configs/:id/details/:detail_id", handler.DeleteConfigGroupDetail)

}
