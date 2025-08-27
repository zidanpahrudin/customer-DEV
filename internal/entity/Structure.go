package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Structure model - update untuk menambahkan field address dan active
type Structure struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	CustomerID string           `json:"customer_id" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Level      int            `json:"level" gorm:"not null"`
	ParentID   *string               `json:"parent_id"`
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

func (sd *Structure) BeforeCreate(tx *gorm.DB) (err error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	sd.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	return
}