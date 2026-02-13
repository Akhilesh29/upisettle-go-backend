package order

import "time"

// Status constants for orders.
const (
	StatusPending   = "PENDING"
	StatusPaidUPI   = "PAID_UPI"
	StatusPaidCash  = "PAID_CASH"
	StatusPartial   = "PARTIAL"
	StatusCancelled = "CANCELLED"
)

type Order struct {
	ID          uint      `gorm:"primaryKey"`
	MerchantID  uint      `gorm:"not null;index"`
	StoreID     uint      `gorm:"not null;index"`
	ExternalRef string    `gorm:"size:255"` // optional link to POS ref
	Amount      int64     `gorm:"not null"` // store in smallest currency unit (paise)
	Currency    string    `gorm:"size:10;default:'INR'"`
	Status      string    `gorm:"size:32;not null;default:'PENDING'"`
	CreatedAt   time.Time
	PaidAt      *time.Time
	UpdatedAt   time.Time
}

func (Order) TableName() string {
	return "orders"
}

