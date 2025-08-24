package handler

import (
	"net/http"
	"customer-api/internal/config"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"
)




type StageInput struct {
	Name string `json:"name" binding:"required"`
}

type StagesDetail struct {
    ID      string `json:"id" gorm:"primaryKey"`
    StageID string `json:"stage_id"` // ubah dari uint ke string
    Name    string `json:"name"`
    Sla     int    `json:"sla"`
    Uom     string `json:"uom"`
}

func CreateStage(c *gin.Context) {
	var input StageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid input",
			"data":    err.Error(),
		})
		return
	}

	// Check if stage name already exists
	var existingStage entity.Stages
	if result := config.DB.Where("name = ?", input.Name).First(&existingStage); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Stage name sudah digunakan",
			"data": nil,
		})
		return
	}

	stage := entity.Stages{
		Name: input.Name,
	}

	if result := config.DB.Create(&stage); result.Error != nil {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "failed",
		"message": "Gagal membuat stage",
		"data":    result.Error.Error(),
	})
	return
}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Stage berhasil dibuat",
		"stage":   stage,
	})
}

func GetStages(c *gin.Context) {
	var stages []entity.Stages
	if result := config.DB.Find(&stages); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data stages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Stages fetched successfully",
		"data":    stages,
	})
}

func GetStage(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam
	

	var stage entity.Stages
	if result := config.DB.Where("id = ?", id).First(&stage); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Stage fetched successfully",
		"data":    stage,
	})
}

func UpdateStage(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam
	

	var stage entity.Stages
	if result := config.DB.Where("id = ?", id).First(&stage); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	var input StageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Invalid input",
			"data": err.Error(),
		})
		return
	}

	// Check if the new name is already taken by another stage
	var existingStage entity.Stages
	if result := config.DB.Where("name = ? AND id != ?", input.Name, id).First(&existingStage); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Stage name sudah digunakan",
			"data": nil,
		})
		return
	}

	stage.Name = input.Name

	if result := config.DB.Save(&stage); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui stage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Stage berhasil diperbarui",
		"data":    stage,
	})
}


func DeleteStage(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam
	

	var stage entity.Stages
	if result := config.DB.Where("id = ?", id).First(&stage); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stage not found"})
		return
	}

	if result := config.DB.Delete(&stage); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus stage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Stage berhasil dihapus",
	})
}

func CreateStageDetail(c *gin.Context) {
	stageID := c.Param("id")

	var input struct {
		Name string `json:"name" binding:"required"`
		Sla  int    `json:"sla" binding:"required"`
		Uom  string `json:"uom" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid input", "data": err.Error()})
		return
	}

	var stage entity.Stages
	if result := config.DB.Where("id = ?", stageID).First(&stage); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Stage not found"})
		return
	}

	var existingDetail entity.StagesDetail
	if result := config.DB.Where("name = ? AND stage_id = ?", input.Name, stageID).First(&existingDetail); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Stage detail name sudah digunakan"})
		return
	}

	detail := entity.StagesDetail{
		StageID: stage.ID,
		Name:    input.Name,
		Sla:     input.Sla,
		Uom:     input.Uom,
	}

	if result := config.DB.Create(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Gagal membuat stage detail", "data": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Stage detail berhasil dibuat", "data": detail})
}

func GetStageDetails(c *gin.Context) {
	stageID := c.Param("id")
	var details []entity.StagesDetail
	if result := config.DB.Where("stage_id = ?", stageID).Find(&details); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Gagal mengambil data stage details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Stage details fetched successfully", "data": details})
}

func GetStageDetail(c *gin.Context) {
	detailID := c.Param("detail_id")
	var detail entity.StagesDetail
	if result := config.DB.Where("id = ?", detailID).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Stage detail not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Stage detail fetched successfully", "data": detail})
}

func UpdateStageDetail(c *gin.Context) {
	stageID := c.Param("id")
	detailID := c.Param("detail_id")

	var detail entity.StagesDetail
	if result := config.DB.Where("id = ? AND stage_id = ?", detailID, stageID).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Stage detail not found"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
		Sla  int    `json:"sla" binding:"required"`
		Uom  string `json:"uom" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid input", "data": err.Error()})
		return
	}

	var existingDetail entity.StagesDetail
	if result := config.DB.Where("name = ? AND id != ? AND stage_id = ?", input.Name, detailID, stageID).First(&existingDetail); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Stage detail name sudah digunakan"})
		return
	}

	detail.Name = input.Name
	detail.Sla = input.Sla
	detail.Uom = input.Uom

	if result := config.DB.Save(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Gagal memperbarui stage detail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Stage detail berhasil diperbarui", "data": detail})
}

func DeleteStageDetail(c *gin.Context) {
	stageID := c.Param("id")
	detailID := c.Param("detail_id")

	var detail entity.StagesDetail
	if result := config.DB.Where("id = ? AND stage_id = ?", detailID, stageID).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Stage detail not found"})
		return
	}

	if result := config.DB.Delete(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Gagal menghapus stage detail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Stage detail berhasil dihapus"})
}





