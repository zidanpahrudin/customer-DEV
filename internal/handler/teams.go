package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeams(c *gin.Context) {
	var teams []entity.Teams
	if result := config.DB.Find(&teams); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mendapatkan teams",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Teams berhasil ditemukan",
		"data":    teams,
	})
}

// create
func CreateTeam(c *gin.Context) {
	var team entity.Teams
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Gagal membuat team",
			"data":    err.Error(),
		})
		return
	}
	if result := config.DB.Create(&team); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat team",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Team berhasil dibuat",
		"data":    team,
	})

}

// update
func UpdateTeam(c *gin.Context) {
	var team entity.Teams
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui team",
			"data":    err.Error(),
		})
		return
	}
	if result := config.DB.Save(&team); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui team",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Team berhasil diperbarui",
		"data":    team,
	})
}

// delete
func DeleteTeam(c *gin.Context) {
	var team entity.Teams
	if err := config.DB.Where("id = ?", c.Param("id")).First(&team).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Team tidak ditemukan",
			"data":    err.Error(),
		})

		return
	}
	if result := config.DB.Delete(&team); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menghapus team",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Team berhasil dihapus",
		"data":    team,
	})
}

// teams detail
func GetTeamDetails(c *gin.Context) {
	var team entity.Teams
	if err := config.DB.Where("id = ?", c.Param("id")).First(&team).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Team tidak ditemukan",
			"data":    err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Team berhasil ditemukan",
		"data":    team,
	})
}

func GetTeamDetail(c *gin.Context) {
	var teamDetail entity.TeamsDetail
	if err := config.DB.Where("id = ?", c.Param("id")).First(&teamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Team detail tidak ditemukan",
			"data":    err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Team detail berhasil ditemukan",
		"data":    teamDetail,
	})
}

// create teams detail
func CreateTeamDetail(c *gin.Context) {
	var teamDetail entity.TeamsDetail
	if err := c.ShouldBindJSON(&teamDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Gagal membuat team detail",
			"data":    err.Error(),
		})

		return
	}

	if result := config.DB.Create(&teamDetail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat team detail",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Team detail berhasil dibuat",
		"data":    teamDetail,
	})
}

// update teams detail
func UpdateTeamDetail(c *gin.Context) {
	var teamDetail entity.TeamsDetail
	if err := c.ShouldBindJSON(&teamDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui team detail",
			"data":    err.Error(),
		})

		return
	}

	if result := config.DB.Save(&teamDetail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui team detail",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Team detail berhasil diperbarui",
		"data":    teamDetail,
	})
}

// delete teams detail
func DeleteTeamDetail(c *gin.Context) {
	var teamDetail entity.TeamsDetail
	if err := config.DB.Where("id = ?", c.Param("id")).First(&teamDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Team detail tidak ditemukan",
			"data":    err.Error(),
		})

		return
	}

	if result := config.DB.Delete(&teamDetail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menghapus team detail",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Team detail berhasil dihapus",
		"data":    teamDetail,
	})
}
