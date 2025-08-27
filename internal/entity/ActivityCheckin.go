package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)
// ActivityCheckin model - tabel untuk check-in aktivitas
type ActivityCheckin struct {
	ID		string         		`json:"id" gorm:"primaryKey;size:26"`
	ActivityID  uint           `json:"activity_id" gorm:"not null"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	CheckedInAt time.Time      `json:"checked_in_at" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Activity Activity `json:"-" gorm:"foreignKey:ActivityID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}