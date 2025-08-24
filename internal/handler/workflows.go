package handler

import (
	"net/http"
	"time"
	"gorm.io/gorm"
	"customer-api/internal/config"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"
	
)

type Stages struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	Name       string         `json:"name" gorm:"not null;unique"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
}


type Workflows struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	Name       string         `json:"name" gorm:"not null;unique"`
	StageID    string         `json:"stage_id" gorm:"not null"`
	FlowOrder int            `json:"flow_order" gorm:"not null"`
	ThresFrom int            `json:"thres_from" gorm:"not null"`
	ThresTo   int            `json:"thres_to" gorm:"not null"`
	Type	  string         `json:"type" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`

	// Relations - hilangkan dari JSON response
	Stage Stages `json:"-" gorm:"foreignKey:StageID"`
}

type WorkflowsDetail struct {
	ID      string `json:"id" gorm:"primaryKey"`
	WorkflowsID string `json:"workflows_id"` // ubah dari uint ke string
	Name    string `json:"name"`
	Sla     int    `json:"sla"`
	Uom     string `json:"uom"`
	IsActive bool   `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Workflow Workflows `json:"-" gorm:"foreignKey:WorkflowsID"`
}

func CreateWorkflows(c *gin.Context) {
	var input Workflows
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid input",
			"data":    err.Error(),
		})
		return
	}

	// Check if workflow name already exists
	var existingWorkflow entity.Workflows
	if result := config.DB.Where("name = ?", input.Name).First(&existingWorkflow); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Workflow name sudah digunakan",
			"data": nil,
		})
		return
	}
	workflow := entity.Workflows{
		Name: input.Name,
		StageID: input.StageID,
		FlowOrder: input.FlowOrder,
		ThresFrom: input.ThresFrom,
		ThresTo: input.ThresTo,
		Type: input.Type,
	}
	if result := config.DB.Create(&workflow); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat workflow",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Workflows berhasil dibuat",
		"data":   workflow,
	})
}


func GetWorkflows(c *gin.Context) {
	var workflows []entity.Workflows
	if result := config.DB.Find(&workflows); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mendapatkan workflows",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mendapatkan workflows",
		"data":    workflows,
	})
}

func GetWorkflow(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	var workflow entity.Workflows
	if result := config.DB.Where("id = ?", id).First(&workflow); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Workflow tidak ditemukan",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mendapatkan workflow",
		"data":    workflow,
	})
}

func UpdateWorkflow(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	var workflow entity.Workflows
	if result := config.DB.Where("id = ?", id).First(&workflow); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Workflow tidak ditemukan",
			"data":    result.Error.Error(),
		})
		return
	}

	var input Workflows
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid input",
			"data":    err.Error(),
		})
		return
	}

	// Check if workflow name already exists (excluding current workflow)
	var existingWorkflow entity.Workflows
	if result := config.DB.Where("name = ? AND id <> ?", input.Name, id).First(&existingWorkflow); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Workflow name sudah digunakan",
			"data": nil,
		})
		return
	}

	workflow.Name = input.Name
	workflow.StageID = input.StageID
	workflow.FlowOrder = input.FlowOrder
	workflow.ThresFrom = input.ThresFrom
	workflow.ThresTo = input.ThresTo
	workflow.Type = input.Type

	if result := config.DB.Save(&workflow); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui workflow",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Workflow berhasil diperbarui",
		"data":    workflow,
	})
}

func DeleteWorkflow(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	var workflow entity.Workflows
	if result := config.DB.Where("id = ?", id).First(&workflow); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Workflow tidak ditemukan",
			"data":    result.Error.Error(),
		})
		return
	}

	if result := config.DB.Delete(&workflow); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menghapus workflow",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Workflow berhasil dihapus",
	})
}

func CreateWorkflowDetail(c *gin.Context) {
	var input WorkflowsDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid input",
			"data":    err.Error(),
		})
		return
	}

	// Check if workflow detail name already exists
	var existingDetail entity.WorkflowsDetail
	if result := config.DB.Where("name = ?", input.Name).First(&existingDetail); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Workflow detail name sudah digunakan",
			"data": nil,
		})
		return
	}
	detail := entity.WorkflowsDetail{
		WorkflowsID: input.WorkflowsID,
		Name: input.Name,
		Sla: input.Sla,
		Uom: input.Uom,
		IsActive: input.IsActive,
	}
	if result := config.DB.Create(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal membuat workflow detail",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Workflow detail berhasil dibuat",
		"data":   detail,
	})
}

func GetWorkflowDetails(c *gin.Context) {
	var details []entity.WorkflowsDetail
	if result := config.DB.Find(&details); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal mendapatkan workflow details",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mendapatkan workflow details",
		"data":    details,
	})
}

func GetWorkflowDetail(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	var detail entity.WorkflowsDetail
	if result := config.DB.Where("id = ?", id).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Workflow detail tidak ditemukan",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil mendapatkan workflow detail",
		"data":    detail,
	})
}

func UpdateWorkflowDetail(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	var detail entity.WorkflowsDetail
	if result := config.DB.Where("id = ?", id).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Workflow detail tidak ditemukan",
			"data":    result.Error.Error(),
		})
		return
	}

	var input WorkflowsDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid input",
			"data":    err.Error(),
		})
		return
	}

	// Check if workflow detail name already exists (excluding current detail)
	var existingDetail entity.WorkflowsDetail
	if result := config.DB.Where("name = ? AND id <> ?", input.Name, id).First(&existingDetail); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"message": "Workflow detail name sudah digunakan",
			"data": nil,
		})
		return
	}

	detail.WorkflowsID = input.WorkflowsID
	detail.Name = input.Name
	detail.Sla = input.Sla
	detail.Uom = input.Uom
	detail.IsActive = input.IsActive

	if result := config.DB.Save(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal memperbarui workflow detail",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Workflow detail berhasil diperbarui",
		"data":    detail,
	})
}

func DeleteWorkflowDetail(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	var detail entity.WorkflowsDetail
	if result := config.DB.Where("id = ?", id).First(&detail); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "Workflow detail tidak ditemukan",
			"data":    result.Error.Error(),
		})
		return
	}

	if result := config.DB.Delete(&detail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Gagal menghapus workflow detail",
			"data":    result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Workflow detail berhasil dihapus",
	})
}
