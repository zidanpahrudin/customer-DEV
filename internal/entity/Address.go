package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Address model - update untuk menambahkan field name dan active
// Address model
type Address struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	CustomerID string           `json:"customer_id" gorm:"not null"`
	// SupplierID *uint          `json:"supplier_id"` // HAPUS field ini
	Name       string         `json:"name" gorm:"not null"`
	Street     string         `json:"street"`
	Address    string         `json:"address" gorm:"not null"`
	City       string         `json:"city"`
	State      string         `json:"state"`
	Country    string         `json:"country"`
	PostalCode string         `json:"postal_code"`
	Main       bool           `json:"main" gorm:"default:false"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}

// BeforeCreate hook - generate ulid for ID
// before save generate id
func (s *Address) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}