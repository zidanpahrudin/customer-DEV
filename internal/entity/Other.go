package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Other model - update untuk menggunakan key instead of attribute_name
type Other struct {
	ID         string         `json:"id" gorm:"type:char(36);primary_key"`
	CustomerID string           `json:"customer_id" gorm:"not null"`
	Key        string         `json:"key" gorm:"not null"`
	Value      *string        `json:"value"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}

// BeforeCreate hook - generate ID before create
func (s *Other) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}