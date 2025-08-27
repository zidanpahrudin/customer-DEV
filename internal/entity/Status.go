package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Status model - tabel untuk status customer
type Status struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	StatusName string         `json:"status_name" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (sd *Status) BeforeCreate(tx *gorm.DB) (err error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	sd.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	return
}