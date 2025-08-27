package entity

import (
	"time"
	"math/rand"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)


// Group model
type Group struct {
	ID		string         		`json:"id" gorm:"primaryKey;size:26"`
	NameGroup string         `json:"name_group" gorm:"not null"`
	Value     string         `json:"value"`
	Active    bool           `json:"active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Customers []Customer `json:"customers,omitempty" gorm:"many2many:customer_groups;"`
}


// auto generate id
// before save generate id
func (s *Group) BeforeCreate(tx *gorm.DB) (err error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    s.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
    return
}
