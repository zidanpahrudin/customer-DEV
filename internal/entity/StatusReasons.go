package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// StatusReasons model - tabel untuk Reason Status customer
type StatusReasons struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	CustomerID string           `json:"customer_id" gorm:"not null"`
	Reason string         	  `json:"reason" gorm:"not null"`
	Status     string         `json:"status" gorm:"type:varchar(20);check:status IN ('active','blocked');not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}


func (sd *StatusReasons) BeforeCreate(tx *gorm.DB) (err error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	sd.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	return
}