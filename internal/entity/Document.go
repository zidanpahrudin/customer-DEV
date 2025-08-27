package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type Document struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	Notes       string         `json:"notes" gorm:"not null"`
	Type       string         `json:"type" gorm:"not null"`
	URLFile        string         `json:"url_file" gorm:"not null"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}