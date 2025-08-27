package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type OthersConfigDetail struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	ConfigID  string         `json:"config_id" gorm:"index"`
	Icon 	  string         `json:"icon"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// relationships
	Config OthersConfig `json:"config" gorm:"foreignKey:ConfigID;references:ID"`
}


// before save generate id
func (s *OthersConfigDetail) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}