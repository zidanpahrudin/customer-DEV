package handler

import (
	"net/http"
	"strconv"

	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"
)

type GroupInput struct {
	NameGroup string `json:"name_group" binding:"required"`
	Value     string `json:"value"`
	Active    bool   `json:"active"`
}

// @Summary Create group
// @Description Create a new group
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group body dto.CreateGroupsRequest true "Group data"
// @Success 201 {object} entity.Group
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/groups [post]
func CreateGroup(c *gin.Context) {
	var req dto.CreateGroupsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Industry Group if provided
	if req.IndustryID != "" {
		industryGroup := entity.Group{
			NameGroup: "Industry",
			Value:     req.IndustryID,
			Active:    req.IndustryActive,
		}

		result := config.DB.Create(&industryGroup)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create industry group"})
			return
		}

		c.JSON(http.StatusCreated, industryGroup)
		return
	}

	// Create Parent Group if provided
	if req.ParentGroupID != "" {
		parentGroup := entity.Group{
			NameGroup: "Parent Group",
			Value:     req.ParentGroupID,
			Active:    req.ParentGroupActive,
		}

		result := config.DB.Create(&parentGroup)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create parent group"})
			return
		}

		c.JSON(http.StatusCreated, parentGroup)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Either IndustryId or ParentGroupId must be provided"})
}

// @Summary Get all groups
// @Description Get list of all groups
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param active query bool false "Filter by active status"
// @Success 200 {array} entity.Group
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/groups [get]
func GetGroups(c *gin.Context) {
	var groups []entity.Group
	query := config.DB

	// Filter by active status if provided
	if activeParam := c.Query("active"); activeParam != "" {
		if active, err := strconv.ParseBool(activeParam); err == nil {
			query = query.Where("active = ?", active)
		}
	}

	if result := query.Find(&groups); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// @Summary Get group by ID
// @Description Get a specific group by ID
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} entity.Group
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/groups/{id} [get]
func GetGroup(c *gin.Context) {
	id := c.Param("id")
	var group entity.Group

	if result := config.DB.First(&group, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// @Summary Update group
// @Description Update an existing group
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Param group body GroupInput true "Group data"
// @Success 200 {object} entity.Group
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/groups/{id} [put]
func UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var group entity.Group

	if result := config.DB.First(&group, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group tidak ditemukan"})
		return
	}

	var input GroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if group name already exists (exclude current group)
	var existingGroup entity.Group
	if result := config.DB.Where("name_group = ? AND id != ?", input.NameGroup, id).First(&existingGroup); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group name sudah digunakan"})
		return
	}

	// Update group
	group.NameGroup = input.NameGroup
	group.Value = input.Value
	group.Active = input.Active

	if result := config.DB.Save(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Group berhasil diupdate",
		"group":   group,
	})
}

// @Summary Delete group
// @Description Delete a group by ID
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/groups/{id} [delete]
func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	var group entity.Group

	if result := config.DB.First(&group, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group tidak ditemukan"})
		return
	}

	// Remove all customer-group associations first
	config.DB.Model(&group).Association("Customers").Clear()

	if result := config.DB.Delete(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group berhasil dihapus"})
}

// @Summary Assign customer to group
// @Description Assign a customer to a specific group
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group_id path int true "Group ID"
// @Param customer_id path int true "Customer ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/groups/{group_id}/customers/{customer_id} [post]
func AssignCustomerToGroup(c *gin.Context) {
	groupID := c.Param("group_id")
	customerID := c.Param("customer_id")

	var group entity.Group
	if result := config.DB.First(&group, groupID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group tidak ditemukan"})
		return
	}

	var customer entity.Customer
	if result := config.DB.First(&customer, customerID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	// Add customer to group
	if err := config.DB.Model(&group).Association("Customers").Append(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan customer ke group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil ditambahkan ke group"})
}

// @Summary Remove customer from group
// @Description Remove a customer from a specific group
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param group_id path int true "Group ID"
// @Param customer_id path int true "Customer ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/groups/{group_id}/customers/{customer_id} [delete]
func RemoveCustomerFromGroup(c *gin.Context) {
	groupID := c.Param("group_id")
	customerID := c.Param("customer_id")

	var group entity.Group
	if result := config.DB.First(&group, groupID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group tidak ditemukan"})
		return
	}

	var customer entity.Customer
	if result := config.DB.First(&customer, customerID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	// Remove customer from group
	if err := config.DB.Model(&group).Association("Customers").Delete(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus customer dari group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil dihapus dari group"})
}

// @Summary Get group customers
// @Description Get all customers in a specific group
// @Tags Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Group ID"
// @Success 200 {array} entity.Customer
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/groups/{id}/customers [get]
func GetGroupCustomers(c *gin.Context) {
	id := c.Param("id")
	var group entity.Group

	if result := config.DB.Preload("Customers").First(&group, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, group.Customers)
}