package entity

import (
	"time"

	"gorm.io/gorm"
)

// Activity model - tabel untuk aktivitas customer
type Activity struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CustomerID   uint           `json:"customer_id" gorm:"not null"`
	Title        string         `json:"title" gorm:"not null"`
	Type         string         `json:"type" gorm:"not null"`
	Agenda       string         `json:"agenda"`
	StartTime    time.Time      `json:"start_time" gorm:"not null"`
	EndTime      time.Time      `json:"end_time" gorm:"not null"`
	LocationName string         `json:"location_name"`
	Status       string         `json:"status" gorm:"default:'Scheduled'"`
	CreatedBy    uint           `json:"created_by" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Customer         Customer          `json:"-" gorm:"foreignKey:CustomerID"`
	Creator          User              `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	ActivityCheckins []ActivityCheckin `json:"activity_checkins,omitempty" gorm:"foreignKey:ActivityID"`
	Attendees        []User            `json:"attendees,omitempty" gorm:"many2many:activity_attendees;"`

}