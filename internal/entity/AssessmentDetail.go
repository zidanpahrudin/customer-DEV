package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type AssessmentDetail struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	AssessmentID string         `json:"assessment_id" gorm:"not null"` // ULID string
	Name       string         `json:"name" gorm:"not null;unique"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations - hilangkan dari JSON response
	Assessment Assessment `json:"-" gorm:"foreignKey:AssessmentID"`
}



// before save generate id
func (s *AssessmentDetail) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}