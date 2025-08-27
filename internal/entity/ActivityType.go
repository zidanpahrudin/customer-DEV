package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type ActivityType struct {
	ID   string `json:"id" gorm:"type:char(36);primary_key"`
	Name string `json:"name" gorm:"not null"`
}

// before saving, ensure ID is set
// before save generate id
func (s *ActivityType) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}