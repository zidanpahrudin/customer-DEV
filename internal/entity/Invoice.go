package entity

import (
	"time"

	"gorm.io/gorm"
)

// Invoice model - tabel untuk invoice
type Invoice struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	CustomerID    uint           `json:"customer_id" gorm:"not null"`
	ProjectID     string         `json:"project_id"`
	InvoiceNumber string         `json:"invoice_number" gorm:"unique;not null"`
	Amount        float64        `json:"amount" gorm:"not null"`
	IssuedDate    time.Time      `json:"issued_date" gorm:"not null"`
	DueDate       time.Time      `json:"due_date" gorm:"not null"`
	PaidAmount    float64        `json:"paid_amount" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Customer Customer  `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:InvoiceID"` // Tambahkan ini
}