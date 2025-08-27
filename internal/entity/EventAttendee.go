package entity

import (
	"time"
)
// EventAttendee model - tabel pivot untuk attendees event (many-to-many)
type EventAttendee struct {
	EventID uint      `json:"event_id" gorm:"primaryKey;not null"`
	UserID  uint      `json:"user_id" gorm:"primaryKey;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// Relations
	Event Event `json:"-" gorm:"foreignKey:EventID"`
	User  User  `json:"-" gorm:"foreignKey:UserID"`
}
