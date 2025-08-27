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
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
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