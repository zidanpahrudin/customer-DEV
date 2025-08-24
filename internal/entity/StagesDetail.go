package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type StagesDetail struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	StageID string         `json:"stage_id" gorm:"not null"` // ULID string
	Name       string         `json:"name" gorm:"not null;unique"`
	Sla		int            `json:"sla" gorm:"not null"` // in hours
	Uom 	 string         `json:"uom" gorm:"not null"` // unit of measure
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	// Relations - hilangkan dari JSON response
	Stage Stages `json:"-" gorm:"foreignKey:StageID"`
}


// before save generate id
func (sd *StagesDetail) BeforeCreate(tx *gorm.DB) (err error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	sd.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	return
}