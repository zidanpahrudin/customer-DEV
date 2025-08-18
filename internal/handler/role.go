package handler

import (
	"net/http"
	"strconv"

	"customer-api/internal/config"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"
)

type RoleInput struct {
	RoleName string `json:"role_name" binding:"required"`
}

// CreateRole - Hapus swagger annotations
func CreateRole(c *gin.Context) {
	var input RoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role name already exists
	var existingRole entity.Role
	if result := config.DB.Where("role_name = ?", input.RoleName).First(&existingRole); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role name sudah digunakan"})
		return
	}

	role := entity.Role{
		RoleName: input.RoleName,
	}

	if result := config.DB.Create(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role berhasil dibuat",
		"role":    role,
	})
}

// @Summary Get all roles
// @Description Get list of all roles
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Role
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/roles [get]
func GetRoles(c *gin.Context) {
	var roles []entity.Role
	if result := config.DB.Find(&roles); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// @Summary Get role by ID
// @Description Get a specific role by ID
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} entity.Role
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/roles/{id} [get]
func GetRole(c *gin.Context) {
	id := c.Param("id")
	var role entity.Role

	if result := config.DB.First(&role, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// UpdateRole - Hapus swagger annotations
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role entity.Role

	if result := config.DB.First(&role, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	var input RoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role name already exists (exclude current role)
	var existingRole entity.Role
	if result := config.DB.Where("role_name = ? AND id != ?", input.RoleName, id).First(&existingRole); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role name sudah digunakan"})
		return
	}

	// Update role
	role.RoleName = input.RoleName

	if result := config.DB.Save(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role berhasil diupdate",
		"role":    role,
	})
}

// DeleteRole - Hapus swagger annotations
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	var role entity.Role

	if result := config.DB.First(&role, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	// Check if role is being used by users
	var userCount int64
	config.DB.Model(&entity.User{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak dapat dihapus karena masih digunakan oleh user"})
		return
	}

	// Prevent deletion of default roles (ID 1 and 2)
	roleID, _ := strconv.Atoi(id)
	if roleID == 1 || roleID == 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role default tidak dapat dihapus"})
		return
	}

	if result := config.DB.Delete(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role berhasil dihapus"})
}

// SetupDefaultRoles - Hapus swagger annotations
func SetupDefaultRoles(c *gin.Context) {
	// Create default roles if they don't exist
	defaultRoles := []entity.Role{
		{ID: 1, RoleName: "admin"},
		{ID: 2, RoleName: "user"},
	}

	for _, role := range defaultRoles {
		var existingRole entity.Role
		if result := config.DB.First(&existingRole, role.ID); result.Error != nil {
			config.DB.Create(&role)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default roles berhasil dibuat"})
}