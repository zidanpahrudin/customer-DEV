package entity

import (
	"time"
	"crypto/rand"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


// User model for authentication
type User struct {
	ID       string `json:"id" gorm:"primaryKey;size:26"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	RoleID    string           `json:"role_id" gorm:"default:2"` // Default to regular user role
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

// BeforeCreate hook for generating ID before inserting to database
func (u *User) BeforeCreate(tx *gorm.DB) error {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return err
	}
	u.ID = id.String()
	return nil
}