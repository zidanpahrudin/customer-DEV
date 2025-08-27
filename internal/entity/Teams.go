package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


// Customer entity


type Teams struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	Name      string         `json:"name" gorm:"not null;unique"`
	TeamLead  string         `json:"team_lead" gorm:"not null"`
	Industry  string         `json:"industry" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`

	// relationships
	TeamLeadUser User       `json:"team_lead_user" gorm:"foreignKey:TeamLead;references:ID"`
	Customers    []Customer `json:"customers" gorm:"foreignKey:TeamID"`
}


// before save generate id
func (s *Teams) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}