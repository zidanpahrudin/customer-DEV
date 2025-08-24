package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


type TeamsDetail struct {
	ID        string         `json:"id" gorm:"primaryKey;size:26"`
	TeamsID string         `json:"teams_id" gorm:"not null"` // ULID string
	JobPosition	   string         `json:"job_position" gorm:"not null;unique"`
	EmployeeName string         `json:"employee_name" gorm:"not null"`
	PhoneNumber string         `json:"phone_number" gorm:"not null;unique"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Teams Teams `json:"-" gorm:"foreignKey:TeamsID"`

}


// before save generate id
func (s *TeamsDetail) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}