package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// ActivityAttendee model - tabel pivot untuk attendees aktivitas (many-to-many)
type ActivityAttendee struct {
	ID         string    `json:"id" gorm:"type:char(36);primary_key"`
	ActivityID uint      `json:"activity_id" gorm:"primaryKey;not null"`
	UserID     uint      `json:"user_id" gorm:"primaryKey;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Activity Activity `json:"-" gorm:"foreignKey:ActivityID"`
	User     User     `json:"-" gorm:"foreignKey:UserID"`
}

// before save generate id
func (s *ActivityAttendee) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}