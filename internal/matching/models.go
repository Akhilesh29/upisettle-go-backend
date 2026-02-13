package matching

import "time"

type Match struct {
	ID         uint      `gorm:"primaryKey"`
	OrderID    uint      `gorm:"not null;index"`
	PaymentID  uint      `gorm:"not null;index"`
	Confidence float64   `gorm:"not null"`
	MatchedAt  time.Time `gorm:"not null"`
}

func (Match) TableName() string {
	return "matches"
}

const (
	ExceptionUnmatchedOrder   = "UNMATCHED_ORDER"
	ExceptionUnmatchedPayment = "UNMATCHED_PAYMENT"
	ExceptionAmountMismatch   = "AMOUNT_MISMATCH"
)

type Exception struct {
	ID         uint       `gorm:"primaryKey"`
	MerchantID uint       `gorm:"not null;index"`
	StoreID    uint       `gorm:"not null;index"`
	OrderID    *uint      `gorm:"index"`
	PaymentID  *uint      `gorm:"index"`
	Type       string     `gorm:"size:64;not null"`
	Reason     string     `gorm:"size:512"`
	Resolved   bool       `gorm:"not null;default:false"`
	ResolvedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Exception) TableName() string {
	return "exceptions"
}

