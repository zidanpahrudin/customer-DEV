package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Sosmed model - update untuk menambahkan field handle dan active
// Sosmed model
type Sosmed struct {
	ID            string         `json:"id" gorm:"type:char(26);primary_key"`
	CustomerID string           `json:"customer_id" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Platform   string         `json:"platform" gorm:"not null"`
	Handle     string         `json:"handle" gorm:"not null"`
	Username   string         `json:"username"`
	URL        string         `json:"url"`
	Followers  int            `json:"followers" gorm:"default:0"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}



func (s *Sosmed) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}