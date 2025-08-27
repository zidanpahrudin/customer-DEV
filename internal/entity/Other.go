package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Other model - update untuk menggunakan key instead of attribute_name
type Other struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	Key        string         `json:"key" gorm:"not null"`
	Value      *string        `json:"value"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}