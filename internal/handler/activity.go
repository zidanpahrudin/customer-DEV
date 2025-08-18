package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Hapus fungsi helper floatPtrToFloat dan floatToFloatPtr karena tidak diperlukan lagi
func floatPtrToFloat(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}

func floatToFloatPtr(val float64) *float64 {
	return &val
}

// @Summary Get all activities
// @Description Get list of all activities with pagination
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ActivitiesResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities [get]
func GetActivities(c *gin.Context) {
	var activities []entity.Activity
	customerID := c.Query("customer_id")
	status := c.Query("status")
	activityType := c.Query("type")

	db := config.DB.Preload("Customer").Preload("CreatedByUser")

	if customerID != "" {
		db = db.Where("customer_id = ?", customerID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if activityType != "" {
		db = db.Where("type = ?", activityType)
	}

	result := db.Find(&activities)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	// Convert to response format
	var activityResponses []dto.ActivityResponse
	for _, activity := range activities {
		activityResponses = append(activityResponses, dto.ActivityResponse{
			ID:           activity.ID,
			CustomerID:   activity.CustomerID,
			Title:        activity.Title,
			Type:         activity.Type,
			Agenda:       activity.Agenda,
			StartTime:    activity.StartTime.Format(time.RFC3339),
			EndTime:      activity.EndTime.Format(time.RFC3339),
			LocationName: activity.LocationName,
			Status:       activity.Status,
			CreatedBy:    activity.CreatedBy,
			CreatedAt:    activity.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    activity.UpdatedAt.Format(time.RFC3339),
		})
	}

	// Get total count
	var total int64
	config.DB.Model(&entity.Activity{}).Count(&total)

	response := dto.ActivitiesResponse{
		Activities: activityResponses,
		Total:      total,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Create new activity
// @Description Create a new activity for a customer
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param activity body dto.CreateActivityRequest true "Activity data"
// @Success 201 {object} dto.ActivityResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities [post]
func CreateActivity(c *gin.Context) {
	var req dto.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse time strings
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format. Use RFC3339 format (e.g., 2024-01-15T10:00:00Z)"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format. Use RFC3339 format (e.g., 2024-01-15T12:00:00Z)"})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify customer exists
	var customer entity.Customer
	if err := config.DB.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found"})
		return
	}

	activity := entity.Activity{
		CustomerID:   req.CustomerID,
		Title:        req.Title,
		Type:         req.Type,
		Agenda:       req.Agenda,
		StartTime:    startTime,
		EndTime:      endTime,
		LocationName: req.LocationName,
		Status:       "planned",
		CreatedBy:    userID.(uint),
	}

	result := config.DB.Create(&activity)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	// Load relations for response
	config.DB.Preload("Customer").Preload("CreatedByUser").First(&activity, activity.ID)

	activityResponse := dto.ActivityResponse{
		ID:           activity.ID,
		CustomerID:   activity.CustomerID,
		Title:        activity.Title,
		Type:         activity.Type,
		Agenda:       activity.Agenda,
		StartTime:    activity.StartTime.Format(time.RFC3339),
		EndTime:      activity.EndTime.Format(time.RFC3339),
		LocationName: activity.LocationName,
		Status:       activity.Status,
		CreatedBy:    activity.CreatedBy,
		CreatedAt:    activity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    activity.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, activityResponse)
}

// @Summary Get activity by ID
// @Description Get a specific activity by ID
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity ID"
// @Success 200 {object} dto.ActivityResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities/{id} [get]
func GetActivity(c *gin.Context) {
	id := c.Param("id")

	var activity entity.Activity
	result := config.DB.Preload("Customer").Preload("CreatedByUser").Preload("Attendees").First(&activity, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	activityResponse := dto.ActivityResponse{
		ID:           activity.ID,
		CustomerID:   activity.CustomerID,
		Title:        activity.Title,
		Type:         activity.Type,
		Agenda:       activity.Agenda,
		StartTime:    activity.StartTime.Format(time.RFC3339),
		EndTime:      activity.EndTime.Format(time.RFC3339),
		LocationName: activity.LocationName,
		Status:       activity.Status,
		CreatedBy:    activity.CreatedBy,
		CreatedAt:    activity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    activity.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, activityResponse)
}

// @Summary Update activity
// @Description Update an existing activity
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity ID"
// @Param activity body dto.UpdateActivityRequest true "Activity data"
// @Success 200 {object} dto.ActivityResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities/{id} [put]
func UpdateActivity(c *gin.Context) {
	id := c.Param("id")

	var activity entity.Activity
	result := config.DB.First(&activity, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	var req dto.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Title != nil {
		activity.Title = *req.Title
	}
	if req.Type != nil {
		activity.Type = *req.Type
	}
	if req.Agenda != nil {
		activity.Agenda = *req.Agenda
	}
	if req.StartTime != nil {
		startTime, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format. Use RFC3339 format"})
			return
		}
		activity.StartTime = startTime
	}
	if req.EndTime != nil {
		endTime, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format. Use RFC3339 format"})
			return
		}
		activity.EndTime = endTime
	}
	if req.LocationName != nil {
		activity.LocationName = *req.LocationName
	}
	if req.Status != nil {
		activity.Status = *req.Status
	}

	result = config.DB.Save(&activity)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}

	// Load relations for response
	config.DB.Preload("Customer").Preload("CreatedByUser").First(&activity, activity.ID)

	activityResponse := dto.ActivityResponse{
		ID:           activity.ID,
		CustomerID:   activity.CustomerID,
		Title:        activity.Title,
		Type:         activity.Type,
		Agenda:       activity.Agenda,
		StartTime:    activity.StartTime.Format(time.RFC3339),
		EndTime:      activity.EndTime.Format(time.RFC3339),
		LocationName: activity.LocationName,
		Status:       activity.Status,
		CreatedBy:    activity.CreatedBy,
		CreatedAt:    activity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    activity.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, activityResponse)
}

// @Summary Delete activity
// @Description Delete an activity
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities/{id} [delete]
func DeleteActivity(c *gin.Context) {
	id := c.Param("id")

	// Check if activity exists
	var activity entity.Activity
	if err := config.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	result := config.DB.Delete(&activity)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

// @Summary Add attendees to activity
// @Description Add users as attendees to an activity
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity ID"
// @Param attendees body dto.ActivityAttendeeRequest true "Attendee user IDs"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities/{id}/attendees [post]
func AddActivityAttendees(c *gin.Context) {
	id := c.Param("id")
	activityID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	// Check if activity exists
	var activity entity.Activity
	if err := config.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	var req dto.ActivityAttendeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify all users exist
	var users []entity.User
	if err := config.DB.Where("id IN ?", req.UserIDs).Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to verify users"})
		return
	}

	if len(users) != len(req.UserIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more users not found"})
		return
	}

	// Add attendees
	for _, userID := range req.UserIDs {
		attendee := entity.ActivityAttendee{
			ActivityID: uint(activityID),
			UserID:     userID,
		}
		// Use FirstOrCreate to avoid duplicates
		config.DB.FirstOrCreate(&attendee, entity.ActivityAttendee{
			ActivityID: uint(activityID),
			UserID:     userID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendees added successfully"})
}

// @Summary Remove attendees from activity
// @Description Remove users as attendees from an activity
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity ID"
// @Param attendees body dto.ActivityAttendeeRequest true "Attendee user IDs to remove"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities/{id}/attendees [delete]
func RemoveActivityAttendees(c *gin.Context) {
	id := c.Param("id")
	activityID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	var req dto.ActivityAttendeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove attendees
	result := config.DB.Where("activity_id = ? AND user_id IN ?", activityID, req.UserIDs).Delete(&entity.ActivityAttendee{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove attendees"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendees removed successfully"})
}

// @Summary Check-in to activity
// @Description Check-in to an activity with location
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Activity ID"
// @Param checkin body dto.ActivityCheckinRequest true "Check-in data with location"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/activities/{id}/checkin [post]
func CheckinActivity(c *gin.Context) {
	id := c.Param("id")
	activityID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if activity exists
	var activity entity.Activity
	if err := config.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	var req dto.ActivityCheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Karena ActivityCheckinRequest sekarang kosong, kita bisa menghapus binding JSON
	// atau tetap mempertahankannya untuk konsistensi API

	// Check if user is already checked in
	var existingCheckin entity.ActivityCheckin
	if err := config.DB.Where("activity_id = ? AND user_id = ?", activityID, userID).First(&existingCheckin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already checked in to this activity"})
		return
	}

	// Create check-in record
	checkin := entity.ActivityCheckin{
		ActivityID:  uint(activityID),
		UserID:      userID.(uint),
		CheckedInAt: time.Now(),
	}

	result := config.DB.Create(&checkin)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check-in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checked in successfully"})
}
