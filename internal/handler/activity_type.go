package handler

import (
	"net/http"
	"strconv"

	"customer-api/internal/config"
	"customer-api/internal/entity"
	"github.com/gin-gonic/gin"

	"customer-api/internal/dto"
)



// @Summary Get all Activitiy Types
// @Description Get list of all activity types
// @Tags Activity Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {array} entity.ActivityType
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activity-types [get]
func ReadActivityTypes(c *gin.Context) {

		db := config.DB
		limit, _ := strconv.Atoi(c.Query("limit"))
		page, _ := strconv.Atoi(c.Query("page"))
		if limit == 0 {
			limit = 10 // Default limit
		}

		if page == 0 {
			page = 1 // Default page
		}
		var activityTypes []entity.ActivityType
		offset := (page - 1) * limit
		if result := db.Limit(limit).Offset(offset).Find(&activityTypes); result.Error != nil {
			c.JSON(http.StatusNotFound, dto.Response{
				Status:  http.StatusNotFound,
				Message: "No activity types found",
				Data:    []entity.ActivityType{},
			})
			return
		}

		if len(activityTypes) == 0 {
			c.JSON(http.StatusOK, dto.Response{
				Status:  http.StatusOK,
				Message: "No activity types found",
				Data:    activityTypes,
			})
			return
		}

		
		if result := db.Limit(limit).Offset(offset).Find(&activityTypes); result.Error != nil {
			c.JSON(http.StatusInternalServerError, dto.Response{
				Status:  http.StatusInternalServerError,
				Message: "Failed to retrieve activity types",
				Data:    []entity.ActivityType{},
			})
			return
		}
		

		c.JSON(http.StatusOK, dto.Response{
			Status:  http.StatusOK,
			Message: "Activity types found",
			Data:    activityTypes,
		})

	

}


// @Summary Create a new Activity Type
// @Description Create a new activity type
// @Tags Activity Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param activityType body entity.ActivityType true "Activity Type"
// @Success 201 {object} entity.ActivityType
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activity-types [post]
func CreateActivityType(c *gin.Context) {
	
		var activityType entity.ActivityType
		if err := c.ShouldBindJSON(&activityType); err != nil {
			c.JSON(http.StatusBadRequest, dto.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid input",
				Data:    nil,
			})
			return
		}
		db := config.DB
		if result := db.Create(&activityType); result.Error != nil {
			c.JSON(http.StatusInternalServerError, dto.Response{
				Status:  http.StatusInternalServerError,
				Message: "Failed to create activity type",
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusCreated, dto.Response{
			Status:  http.StatusCreated,
			Message: "Activity type created",
			Data:    activityType,
		})
	
}


// @Summary Update an Activity Type
// @Description Update an activity type
// @Tags Activity Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity Type ID"
// @Param activityType body entity.ActivityType true "Activity Type"
// @Success 200 {object} entity.ActivityType
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activity-types/{id} [put]
func UpdateActivityType(c *gin.Context) {
	// Parse ID dari path
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil || idInt <= 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
		})
		return
	}
	id := uint(idInt)

	// Bind JSON ke struct
	var input entity.ActivityType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    nil,
		})
		return
	}

	db := config.DB

	// Pastikan record ada
	var activityType entity.ActivityType
	if result := db.First(&activityType, id); result.Error != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Status:  http.StatusNotFound,
			Message: "Activity type not found",
			Data:    nil,
		})
		return
	}

	// Update fields dari input (gunakan map atau struct)
	if result := db.Model(&activityType).Updates(input); result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update activity type",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: "Activity type updated",
		Data:    activityType,
	})
}



// @Summary Delete an Activity Type
// @Description Delete an activity type
// @Tags Activity Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity Type ID"
// @Success 204 {object} dto.Response
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activity-types/{id} [delete]
func DeleteActivityType(c *gin.Context) {
	id := c.Param("id")
	if  id == "" {
		c.JSON(http.StatusBadRequest, dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
		})
		return
	}


	db := config.DB
	if result := db.Delete(&entity.ActivityType{ID: id}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete activity type",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusNoContent, dto.Response{
		Status:  http.StatusNoContent,
		Message: "Activity type deleted",
		Data:    nil,
	})
}


// @Summary Get an Activity Type by ID
// @Description Get an activity type by ID
// @Tags Activity Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity Type ID"
// @Success 200 {object} entity.ActivityType
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activity-types/{id} [get]
func ReadActivityType(c *gin.Context) {
	
		id, _ := strconv.Atoi(c.Param("id"))
		db := config.DB
		var activityType entity.ActivityType
		if result := db.First(&activityType, id); result.Error != nil {
			c.JSON(http.StatusNotFound, dto.Response{
				Status:  http.StatusNotFound,
				Message: "Activity type not found",
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusOK, dto.Response{
			Status:  http.StatusOK,
			Message: "Activity type found",
			Data:    activityType,
		})
	
}

// @Summary Get all Activities by Activity Type ID
// @Description Get list of all activities by activity type ID
// @Tags Activity Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity Type ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {array} entity.Activity
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activity-types/{id}/activities [get]
func ReadActivitiesByActivityType(c *gin.Context) {

		id, _ := strconv.Atoi(c.Param("id"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		page, _ := strconv.Atoi(c.Query("page"))
		if limit == 0 {
			limit = 10 // Default limit
		}
		if page == 0 {
			page = 1 // Default page
		}
		db := config.DB
		var activities []entity.Activity
		offset := (page - 1) * limit
		if result := db.Limit(limit).Offset(offset).Find(&activities, "activity_type_id = ?", id); result.Error != nil {
			c.JSON(http.StatusNotFound, dto.Response{
				Status:  http.StatusNotFound,
				Message: "No activities found",
				Data:    []entity.Activity{},
			})
			return
		}
		c.JSON(http.StatusOK, dto.Response{
			Status:  http.StatusOK,
			Message: "Activities found",
			Data:    activities,
		})
		if len(activities) == 0 {
			c.JSON(http.StatusOK, dto.Response{
				Status:  http.StatusOK,
				Message: "No activities found",
				Data:    activities,
			})
			return
		}
		if result := db.Limit(limit).Offset(offset).Find(&activities, "activity_type_id = ?", id); result.Error != nil {
			c.JSON(http.StatusInternalServerError, dto.Response{
				Status:  http.StatusInternalServerError,
				Message: "Failed to retrieve activities",
				Data:    []entity.Activity{},
			})
			return
		}
		c.JSON(http.StatusOK, dto.Response{
			Status:  http.StatusOK,
			Message: "Activities found",
			Data:    activities,
		})
}


