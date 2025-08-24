package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterWorkflowsRoutes(r *gin.RouterGroup) {
	r.POST("/workflows", handler.CreateWorkflows)
	r.GET("/workflows", handler.GetWorkflows)
	r.GET("/workflows/:id", handler.GetWorkflow)
	r.PUT("/workflows/:id", handler.UpdateWorkflow)
	r.DELETE("/workflows/:id", handler.DeleteWorkflow)

	// detail workflows
	r.GET("/workflows/:id/details", handler.GetWorkflowDetails)                  // Ambil semua detail workflow
	r.POST("/workflows/:id/details", handler.CreateWorkflowDetail)               // Buat detail workflow baru
	r.GET("/workflows/:id/details/:detail_id", handler.GetWorkflowDetail)        // Ambil detail workflow tertentu
	r.PUT("/workflows/:id/details/:detail_id", handler.UpdateWorkflowDetail)     // Update detail workflow
	r.DELETE("/workflows/:id/details/:detail_id", handler.DeleteWorkflowDetail)  // Hapus detail workflow
}
