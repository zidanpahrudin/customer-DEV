package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Customer model - update untuk menambahkan field baru
type HistoryCustomer struct {
	ID		string         		`json:"id" gorm:"primaryKey;size:26"`
	CustomerID      string         `json:"customer_id" gorm:"not null"`
	UserID          string         `json:"user_id" gorm:"not null"`
	Status          string         `json:"status" gorm:"default:'Active'"` // Status internal
	Notes      string         `json:"notes"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`

	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID;references:ID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`

}

// BeforeCreate hook - generate ID before create
// before save generate id
func (s *HistoryCustomer) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}