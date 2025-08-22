package entity

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ActivityTypeId  uint       `json:"activity_type_id" gorm:"not null"`
	ScheduledAt    time.Time  `json:"scheduled_at" gorm:"not null"`
	ScheduledTime time.Time  `json:"scheduled_time" gorm:"not null"`
	CustomerID  uint           `json:"customer_id" gorm:"not null"`
	ProjectID   uint           `json:"project_id" gorm:"not null"`
	Attendees  []User         `json:"attendees,omitempty" gorm:"many2many:event_attendees;"`
	Location 	string         `json:"location"`
	Agenda 	string         `json:"agenda"`
	Status 	string         `json:"status" gorm:"default:'upcoming'"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`

	// Relations
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	ActivityType ActivityType `json:"activity_type,omitempty" gorm:"foreignKey:ActivityTypeId"`
	EventAttendees []EventAttendee `json:"event_attendees,omitempty" gorm:"foreignKey:EventID"`
	Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}