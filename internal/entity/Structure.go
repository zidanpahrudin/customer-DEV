package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Structure model - update untuk menambahkan field address dan active
type Structure struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Level      int            `json:"level" gorm:"not null"`
	ParentID   *uint          `json:"parent_id"`
	Address    string         `json:"address"`
	Position   int            `json:"position" gorm:"default:0"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer    `json:"-" gorm:"foreignKey:CustomerID"`
	Parent   *Structure  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Structure `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}