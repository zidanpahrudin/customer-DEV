package entity

import (
	"time"
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Role model for user roles
type Role struct {
	ID        string         `json:"id" gorm:"type:char(36);primary_key"`
	RoleName  string         `json:"role_name" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Users []User `json:"users,omitempty" gorm:"foreignKey:RoleID"`
}


func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
    r.ID = ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
    return
}