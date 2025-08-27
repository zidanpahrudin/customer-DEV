package handler

import (
	"net/http"
	"customer-api/internal/config"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"
	"time"
	"gorm.io/gorm"
)




type GroupConfig struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	Name       string         `json:"name" gorm:"not null;unique"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type GroupConfigDetail struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	GroupConfigID string         `json:"group_config_id" gorm:"not null"` // ULID string
	Name       string         `json:"name" gorm:"not null;unique"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations - hilangkan dari JSON response
	GroupConfig GroupConfig `json:"-" gorm:"foreignKey:GroupConfigID"`
}


func CreateConfigGroup(c *gin.Context) {
	var input entity.GroupConfig
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":   nil,
		})
		return
	}

	if result := config.DB.Create(&input); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat config group",
			"data":   nil,
		})
		return
	}

	// insert group config

	groupConfig := entity.GroupConfig{
		Name: input.Name,
	}
	if result := config.DB.Create(&groupConfig); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat group config",
			"data":   nil,
		})
		return
	}
	

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Config group berhasil dibuat",
		"data": groupConfig,
	})
}


func GetConfigGroups(c *gin.Context) {
	var groups []entity.GroupConfig
	if result := config.DB.Find(&groups); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data config groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data config groups berhasil diambil",
		"data":    groups,
	})
}

func GetConfigGroup(c *gin.Context) {
	id := c.Param("id")
	var group entity.GroupConfig
	if result := config.DB.Where("id = ?", id).First(&group); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config group not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Config group fetched successfully",
		"data":    group,
	})
}

func UpdateConfigGroup(c *gin.Context) {
	id := c.Param("id")
	var input entity.GroupConfig
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":   nil,
		})
		return
	}

	var group entity.GroupConfig
	if result := config.DB.Where("id = ?", id).First(&group); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config group not found"})
		return
	}

	if input.Name != "" {
		group.Name = input.Name
	}

	if result := config.DB.Save(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui config group",
			"data":   nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Config group updated successfully",
		"data":    group,
	})
}

func DeleteConfigGroup(c *gin.Context) {
	id := c.Param("id")
	var group entity.GroupConfig
	if result := config.DB.Where("id = ?", id).First(&group); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config group not found"})
		return
	}

	if result := config.DB.Delete(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menghapus config group",
			"data":   nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Config group deleted successfully",
		"data":    nil,
	})
}

func CreateConfigGroupDetail(c *gin.Context) {
	var input entity.GroupConfigDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":   nil,
		})
		return
	}

	// check if group config exists
	var group entity.GroupConfig
	if result := config.DB.Where("id = ?", input.GroupConfigID).First(&group); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group config not found"})
		return
	}

	if result := config.DB.Create(&input); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat config group detail",
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Config group detail berhasil dibuat",
		"data": input,
	})
}

func GetConfigGroupDetails(c *gin.Context) {
	var details []entity.GroupConfigDetail
	if result := config.DB.Find(&details); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data config group details"})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data config group details berhasil diambil",
		"data":    details,
	})
}

func GetConfigGroupDetail(c *gin.Context) {
	id := c.Param("id")
	var detail entity.GroupConfigDetail
	if result := config.DB.Where("id = ?", id).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config group detail not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Config group detail fetched successfully",
		"data":    detail,
	})
}

func UpdateConfigGroupDetail(c *gin.Context) {
	id := c.Param("id")
	var input entity.GroupConfigDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":   nil,
		})
		return
	}
	var detail entity.GroupConfigDetail
	if result := config.DB.Where("id = ?", id).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config group detail not found"})
		return
	}
	if input.Name != "" {
		detail.Name = input.Name
	}
	if input.IsActive != detail.IsActive {
		detail.IsActive = input.IsActive
	}
	if result := config.DB.Save(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui config group detail",
			"data":   nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Config group detail updated successfully",
		"data":    detail,
	})
}

func DeleteConfigGroupDetail(c *gin.Context) {
	id := c.Param("id")
	var detail entity.GroupConfigDetail
	if result := config.DB.Where("id = ?", id).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config group detail not found"})
		return
	}
	if result := config.DB.Delete(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menghapus config group detail",
			"data":   nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Config group detail deleted successfully",
		"data":    nil,
	})
}




