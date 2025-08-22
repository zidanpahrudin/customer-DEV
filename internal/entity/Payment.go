package entity

import (
	"time"

	"gorm.io/gorm"
)

// Payment model - tabel untuk pembayaran invoice
type Payment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	InvoiceID uint           `json:"invoice_id" gorm:"not null"`
	Amount    float64        `json:"amount" gorm:"not null"`
	PaidAt    time.Time      `json:"paid_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Invoice Invoice `json:"invoice,omitempty" gorm:"foreignKey:InvoiceID"`
}