package entity

import (
	"time"
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type Document struct {
	ID		string         		`json:"id" gorm:"primaryKey;size:26"`
	CustomerID string           `json:"customer_id" gorm:"not null"`
	Notes       string         `json:"notes" gorm:"not null"`
	Type       string         `json:"type" gorm:"not null"`
	URLFile        string         `json:"url_file" gorm:"not null"`
	UserID     string           `json:"user_id" gorm:"not null"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	// Relations
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// before create hook
func (c *Document) BeforeCreate(tx *gorm.DB) error {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return err
	}
	c.ID = id.String()
	return nil
}