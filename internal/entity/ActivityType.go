package entity

import (
	"time"

	"gorm.io/gorm"
)
type ActivityType struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
}

// before saving, ensure ID is set
func (a *ActivityType) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == 0 {
		a.ID = uint(time.Now().UnixNano()) // use nanoseconds for uniqueness
	}
	return
}
