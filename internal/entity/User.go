package entity

import (
	"time"

	"gorm.io/gorm"
)
// User model for authentication
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	RoleID    uint           `json:"role_id" gorm:"default:2"` // Default to regular user role
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Update struct User, tambahkan di bagian Relations:
	// Relations
	Role                Role              `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Activities          []Activity        `json:"activities,omitempty" gorm:"foreignKey:CreatedBy"`
	ActivityCheckins    []ActivityCheckin `json:"activity_checkins,omitempty" gorm:"foreignKey:UserID"`
	AttendingActivities []Activity        `json:"attending_activities,omitempty" gorm:"many2many:activity_attendees;"`
}