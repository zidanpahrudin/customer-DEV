package entity

import (
	"time"

	"gorm.io/gorm"
)
// Address model - update untuk menambahkan field name dan active
// Address model
type Address struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
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