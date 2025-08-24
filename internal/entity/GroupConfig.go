package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type GroupConfig struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	Name       string         `json:"name" gorm:"not null;unique"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	IsDeleted bool           `json:"is_deleted" gorm:"default:false"`
}


// before save generate id
func (s *GroupConfig) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}