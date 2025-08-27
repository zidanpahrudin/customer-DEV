package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// Role model for user roles
type Role struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	RoleName  string         `json:"role_name" gorm:"unique;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Users []User `json:"users,omitempty" gorm:"foreignKey:RoleID"`
}
