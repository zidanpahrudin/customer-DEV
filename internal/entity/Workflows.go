package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type Workflows struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	Name       string         `json:"name" gorm:"not null;unique"`
	StageID    string         `json:"stage_id" gorm:"not null"`
	FlowOrder int            `json:"flow_order" gorm:"not null"`
	ThresFrom int            `json:"thres_from" gorm:"not null"`
	ThresTo   int            `json:"thres_to" gorm:"not null"`
	Type	  string         `json:"type" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	// Relations - hilangkan dari JSON response
	Stage Stages `json:"-" gorm:"foreignKey:StageID"`
}


// before save generate id
func (s *Workflows) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}