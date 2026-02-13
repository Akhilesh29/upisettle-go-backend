package payment

import "time"

const (
	ChannelUPI  = "UPI"
	ChannelCash = "CASH"
	ChannelOther = "OTHER"
)

type Payment struct {
	ID           uint      `gorm:"primaryKey"`
	MerchantID   uint      `gorm:"not null;index"`
	StoreID      uint      `gorm:"not null;index"`
	Channel      string    `gorm:"size:16;not null"`
	Amount       int64     `gorm:"not null"` // paise
	Currency     string    `gorm:"size:10;default:'INR'"`
	Time         time.Time `gorm:"not null;index"`
	UPIRef       string    `gorm:"size:128;index"`
	PayerVPA     string    `gorm:"size:255"`
	PayerName    string    `gorm:"size:255"`
	RawMessageID string    `gorm:"size:255"` // SMS/email source id if applicable
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Payment) TableName() string {
	return "payments"
}

