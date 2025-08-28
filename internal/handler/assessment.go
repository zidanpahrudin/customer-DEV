package handler

import (
	"customer-api/internal/config"
	"customer-api/internal/dto"
	"customer-api/internal/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Create assessmenet
// @Description Create a new organizational structure for specific customer
// @Tags Structures
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} entity.Structure
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/structures [post]
func CreateAssessment(c *gin.Context) {
	var assessment dto.Assessment
	if err := c.ShouldBindJSON(&assessment); err!= nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
			Errors:  err.Error(),
			}
		})

	db := config.GetDB()
	defer db.Close()


	newAssessment := entity.Assessment{
		Name:       assessment.Name,
		RoleID: assessment.RoleID,
	}

	if err := db.Create(&newAssessment).Error; err!= nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.JSON({
		status: http.StatusCreated,
		message: "Assessment created successfully",
		data: newAssessment,
	})
}


// @Summary Get all assessments
// @Description Get all assessments
// @Tags assessments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []entity.Assessment
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments [get]
func GetAssessments(c *gin.Context) {
	db := config.GetDB()

	var assessments []entity.Assessment
	if err := db.Find(&assessments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.JSON({
		status: http.StatusOK,
		message: "Assessments retrieved successfully",
		data: assessments,
	})
}


// @Summary Get assessment by ID
// @Description Get assessment by ID
// @Tags assessments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assessment ID"
// @Success 200 {object} entity.Assessment
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments/{id} [get]
func GetAssessment(c *gin.Context) {
	db := config.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid assessment ID",
		})
		return
	}

	var assessment entity.Assessment
	if err := db.First(&assessment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Assessment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.JSON({
		status: http.StatusOK,
		message: "Assessment retrieved successfully",
		data: assessment,
	})
	}


// @Summary Update assessment
// @Description Update assessment
//   @Tags assessments
//   @Accept json
//   @Produce json
//   @Security BearerAuth
//   @Param id path int true "Assessment ID"
//   @Param assessment body dto.Assessment true "Assessment object"
//   @Success 200 {object} entity.Assessment
//   @Failure 400 {object} dto.ErrorResponse
//   @Failure 401 {object} dto.ErrorResponse
//   @Failure 404 {object} dto.ErrorResponse
//   @Failure 500 {object} dto.ErrorResponse
//   @Router /api/assessments/{id} [put]
func UpdateAssessment(c *gin.Context) {
	db := config.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid assessment ID",
		})
		return
	}


	var assessment dto.Assessment
	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
			Errors:  err.Error(),
			}
		})

		var dbAssessment entity.Assessment
		if err := db.First(&dbAssessment, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, dto.ErrorResponse{
					Message: "Assessment not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
			return
		}

		dbAssessment.Name = assessment.Name
		dbAssessment.RoleID = assessment.RoleID

		if err := db.Save(&dbAssessment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
			return
		}

		c.JSON({
			status: http.StatusOK,
			message: "Assessment updated successfully",
			data: dbAssessment,
		})
		
}

// @Summary Delete assessment
// @Description Delete assessment
// @Tags assessments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assessment ID"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments/{id} [delete]
func DeleteAssessment(c *gin.Context) {
	db := config.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid assessment ID",
		})
		return
	}

	var assessment entity.Assessment
	if err := db.First(&assessment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Assessment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	if err := db.Delete(&assessment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.Status(http.StatusNoContent)
}


// @Summary Get all assessments by role ID
// @Description Get all assessments by role ID
// @Tags assessments
// @Produce json
// @Security BearerAuth
// @Param role_id path int true "Role ID"
// @Success 200 {object} []entity.Assessment
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments/role/{role_id} [get]
func GetAssessmentsByRoleID(c *gin.Context) {
	db := config.GetDB()
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid role ID",
		})
		return
	}

	var assessments []entity.Assessment
	if err := db.Where("role_id = ?", roleID).Find(&assessments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.JSON({
		status: http.StatusOK,
		message: "Assessments retrieved successfully",
		data: assessments,
	})
	}

// assessment detail
// @Summary Get assessment detail
// @Description Get assessment detail
// @Tags assessments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assessment ID"
// @Success 200 {object} entity.AssessmentDetail
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments/{id}/detail [get]
func GetAssessmentDetail(c *gin.Context) {
	db := config.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid assessment ID",
		})
		return
		}
		var assessmentDetail entity.AssessmentDetail
		if err := db.First(&assessmentDetail, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, dto.ErrorResponse{
					Message: "Assessment detail not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
			return
		}


		c.JSON({
			status: http.StatusOK,
			message: "Assessment detail retrieved successfully",
			data: assessmentDetail,
			})
}

// @Summary Create assessment detail
// @Description Create assessment detail
// @Tags assessments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assessment ID"
// @Param assessment_detail body dto.AssessmentDetail true "Assessment detail object"
// @Success 201 {object} entity.AssessmentDetail
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments/{id}/detail [post]
func CreateAssessmentDetail(c *gin.Context) {
	db := config.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid assessment ID",
			})
			return
			}
			var assessmentDetail dto.AssessmentDetail
			if err := c.ShouldBindJSON(&assessmentDetail); err != nil {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Message: "Invalid request body",
					Errors:  err.Error(),
					}
				})

				var dbAssessment entity.Assessment
				if err := db.First(&dbAssessment, id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, dto.ErrorResponse{
							Message: "Assessment not found",
							})
							return
							}

							c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
								Message: "Internal server error",
							})
							return
						}

						newAssessmentDetail := entity.AssessmentDetail{
							AssessmentID: id,
							Detail:       assessmentDetail.Detail,
						}

						if err := db.Create(&newAssessmentDetail).Error; err != nil {
							c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
								Message: "Internal server error",
							})
							return
						}

						c.JSON({
							status: http.StatusCreated,
							message: "Assessment detail created successfully",
							data: newAssessmentDetail,
						})
						}
						



// assessment detail
// @Summary Update assessment detail
// @Description Update assessment detail
// @Tags assessments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assessment ID"
// @Param assessment_detail body dto.AssessmentDetail true "Assessment detail object"
// @Success 200 {object} entity.AssessmentDetail
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments/{id}/detail [put]
func UpdateAssessmentDetail(c *gin.Context) {
	db := config.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid assessment ID",
		})
		return
	}

	var assessmentDetail dto.AssessmentDetail
	if err := c.ShouldBindJSON(&assessmentDetail); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
			Errors:  err.Error(),
			}
		})
	
	var dbAssessmentDetail entity.AssessmentDetail
	if err := db.First(&dbAssessmentDetail, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Assessment detail not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}
	dbAssessmentDetail.Detail = assessmentDetail.Detail

	if err := db.Save(&dbAssessmentDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.JSON({
		status: http.StatusOK,
		message: "Assessment detail updated successfully",
		data: dbAssessmentDetail,
	})
}

// @Summary Delete assessment detail
// @Description Delete assessment detail
// @Tags assessments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Assessment ID"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/assessments/{id}/detail [delete]
func DeleteAssessmentDetail(c *gin.Context) {
	db := config.GetDB()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid assessment ID",
		})
		return
	}

	var assessmentDetail entity.AssessmentDetail
	if err := db.First(&assessmentDetail, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Message: "Assessment detail not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}
	if err := db.Delete(&assessmentDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.Status(http.StatusNoContent)
}







