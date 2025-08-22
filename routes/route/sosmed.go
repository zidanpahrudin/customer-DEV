package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterSosmedRoutes(r *gin.RouterGroup) {
	r.GET("/customers/:id/sosmeds", handler.GetCustomerSosmeds)
	r.GET("/sosmeds/:id", handler.GetSosmed)
	r.PUT("/sosmeds/:id", handler.UpdateSosmed)
	r.DELETE("/sosmeds/:id", handler.DeleteSosmed)
}
