package entity

import (
	"time"
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


// Contact model - update untuk menambahkan field baru
type Contact struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	CustomerID  string           `json:"customer_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Birthdate   *time.Time     `json:"birthdate"`
	JobPosition string         `json:"job_position"`
	Position    string         `json:"position"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	Mobile      string         `json:"mobile"`
	Department  string         `json:"department"`
	Main        bool           `json:"main" gorm:"default:false"`
	Active      bool           `json:"active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}

// BeforeCreate hook - generate ID before create
func (c *Contact) BeforeCreate(tx *gorm.DB) error {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return err
	}
	c.ID = id.String()
	return nil
}