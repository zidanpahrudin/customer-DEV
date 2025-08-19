package entity

import (
	"time"

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

// User model for authentication
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	RoleID    uint           `json:"role_id" gorm:"default:2"` // Default to regular user role
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Update struct User, tambahkan di bagian Relations:
	// Relations
	Role                Role              `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Activities          []Activity        `json:"activities,omitempty" gorm:"foreignKey:CreatedBy"`
	ActivityCheckins    []ActivityCheckin `json:"activity_checkins,omitempty" gorm:"foreignKey:UserID"`
	AttendingActivities []Activity        `json:"attending_activities,omitempty" gorm:"many2many:activity_attendees;"`
}

// Customer model - update untuk menambahkan field baru
// Customer model - update untuk menambahkan field baru
type Customer struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"not null"`
	BrandName        string         `json:"brand_name"`
	Code             string         `json:"code" gorm:"unique"`
	AccountManagerId string         `json:"account_manager_id"`
	Email            string         `json:"email"`
	Phone            string         `json:"phone"`
	Website          string         `json:"website"`
	Description      string         `json:"description"`
	Logo             string         `json:"logo"`
	LogoSmall        string         `json:"logo_small"` // Field baru untuk logo kecil
	Status           string         `json:"status" gorm:"default:'Active'"` // Status internal
	Category         string         `json:"category"`
	Rating           float64        `json:"rating" gorm:"default:0"`
	AverageCost      float64        `json:"average_cost" gorm:"default:0"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Addresses  []Address   `json:"addresses,omitempty" gorm:"foreignKey:CustomerID"`
	Sosmeds    []Sosmed    `json:"sosmeds,omitempty" gorm:"foreignKey:CustomerID"`
	Contacts   []Contact   `json:"contacts,omitempty" gorm:"foreignKey:CustomerID"`
	Structures []Structure `json:"structures,omitempty" gorm:"foreignKey:CustomerID"`
	Groups     []Group     `json:"groups,omitempty" gorm:"many2many:customer_groups;"`
	Others     []Other     `json:"others,omitempty" gorm:"foreignKey:CustomerID"`
	Activities []Activity  `json:"activities,omitempty" gorm:"foreignKey:CustomerID"`
}

// Address model - update untuk menambahkan field name dan active
// Address model
type Address struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	// SupplierID *uint          `json:"supplier_id"` // HAPUS field ini
	Name       string         `json:"name" gorm:"not null"`
	Street     string         `json:"street"`
	Address    string         `json:"address" gorm:"not null"`
	City       string         `json:"city"`
	State      string         `json:"state"`
	Country    string         `json:"country"`
	PostalCode string         `json:"postal_code"`
	Main       bool           `json:"main" gorm:"default:false"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}

// Sosmed model - update untuk menambahkan field handle dan active
// Sosmed model
type Sosmed struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Platform   string         `json:"platform" gorm:"not null"`
	Handle     string         `json:"handle" gorm:"not null"`
	Username   string         `json:"username"`
	URL        string         `json:"url"`
	Followers  int            `json:"followers" gorm:"default:0"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}

// Contact model - update untuk menambahkan field baru
type Contact struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CustomerID  uint           `json:"customer_id" gorm:"not null"`
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

// Structure model - update untuk menambahkan field address dan active
type Structure struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Level      int            `json:"level" gorm:"not null"`
	ParentID   *uint          `json:"parent_id"`
	Address    string         `json:"address"`
	Position   int            `json:"position" gorm:"default:0"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer    `json:"-" gorm:"foreignKey:CustomerID"`
	Parent   *Structure  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Structure `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// Other model - update untuk menggunakan key instead of attribute_name
type Other struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id" gorm:"not null"`
	Key        string         `json:"key" gorm:"not null"`
	Value      *string        `json:"value"`
	Active     bool           `json:"active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations - hilangkan dari JSON response
	Customer Customer `json:"-" gorm:"foreignKey:CustomerID"`
}

// Group model
type Group struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	NameGroup string         `json:"name_group" gorm:"not null"`
	Value     string         `json:"value"`
	Active    bool           `json:"active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Customers []Customer `json:"customers,omitempty" gorm:"many2many:customer_groups;"`
}

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

// ActivityCheckin model - tabel untuk check-in aktivitas
type ActivityCheckin struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ActivityID  uint           `json:"activity_id" gorm:"not null"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	CheckedInAt time.Time      `json:"checked_in_at" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Activity Activity `json:"-" gorm:"foreignKey:ActivityID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ActivityAttendee model - tabel pivot untuk attendees aktivitas (many-to-many)
type ActivityAttendee struct {
	ActivityID uint      `json:"activity_id" gorm:"primaryKey;not null"`
	UserID     uint      `json:"user_id" gorm:"primaryKey;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Activity Activity `json:"-" gorm:"foreignKey:ActivityID"`
	User     User     `json:"-" gorm:"foreignKey:UserID"`
}

// Invoice model - tabel untuk invoice
// Hapus seluruh struct Invoice (baris 237-253)
// type Invoice struct {
//     ID            uint           `json:"id" gorm:"primaryKey"`
//     CustomerID    uint           `json:"customer_id" gorm:"not null"`
//     ProjectID     string         `json:"project_id"`
//     InvoiceNumber string         `json:"invoice_number" gorm:"unique;not null"`
//     Amount        float64        `json:"amount" gorm:"not null"`
//     IssuedDate    time.Time      `json:"issued_date" gorm:"not null"`
//     DueDate       time.Time      `json:"due_date" gorm:"not null"`
//     PaidAmount    float64        `json:"paid_amount" gorm:"default:0"`
//     CreatedAt     time.Time      `json:"created_at"`
//     UpdatedAt     time.Time      `json:"updated_at"`
//     DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
//
//     // Relations
//     Customer Customer  `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
//     Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:InvoiceID"`
// }

// Hapus seluruh struct Payment (baris 256-268)
// type Payment struct {
//     ID        uint           `json:"id" gorm:"primaryKey"`
//     InvoiceID uint           `json:"invoice_id" gorm:"not null"`
//     Amount    float64        `json:"amount" gorm:"not null"`
//     PaidAt    time.Time      `json:"paid_at" gorm:"not null"`
//     CreatedAt time.Time      `json:"created_at"`
//     UpdatedAt time.Time      `json:"updated_at"`
//     DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
//
//     // Relations
//     Invoice Invoice `json:"invoice,omitempty" gorm:"foreignKey:InvoiceID"`
// }

// Status model - tabel untuk status customer
type Status struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	StatusName string         `json:"status_name" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
