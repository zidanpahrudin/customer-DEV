package entity

import (
	"time"

	"gorm.io/gorm"
)

// StatusReasons model - tabel untuk Reason Status customer
type StatusReasons struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	Reason string         	  `json:"reason" gorm:"not null"`
	Status     string         `json:"status" gorm:"type:varchar(20);check:status IN ('active','blocked');not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}