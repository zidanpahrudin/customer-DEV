package entity

import (
	"time"
)

// ActivityAttendee model - tabel pivot untuk attendees aktivitas (many-to-many)
type ActivityAttendee struct {
	ActivityID uint      `json:"activity_id" gorm:"primaryKey;not null"`
	UserID     uint      `json:"user_id" gorm:"primaryKey;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Activity Activity `json:"-" gorm:"foreignKey:ActivityID"`
	User     User     `json:"-" gorm:"foreignKey:UserID"`
}