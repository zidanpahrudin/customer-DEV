package route

import (
	"customer-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterTeamsRoutes(r *gin.RouterGroup) {
	r.POST("/teams", handler.CreateTeam)
	r.GET("/teams", handler.GetTeams)
	r.GET("/teams/:id", handler.GetTeams)
	r.PUT("/teams/:id", handler.UpdateTeam)
	r.DELETE("/teams/:id", handler.DeleteTeam)

	// detail teams
	r.POST("/teams/detail", handler.CreateTeamDetail)
	r.GET("/teams/detail", handler.GetTeamDetails)
	r.GET("/teams/detail/:id", handler.GetTeamDetail)
	r.PUT("/teams/detail/:id", handler.UpdateTeamDetail)
	r.DELETE("/teams/detail/:id", handler.DeleteTeamDetail)
}
