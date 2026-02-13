package merchant

import "time"

type Merchant struct {
	ID             uint      `gorm:"primaryKey"`
	Name           string    `gorm:"size:255;not null"`
	PrimaryContact string    `gorm:"size:255"`
	Timezone       string    `gorm:"size:100;default:'Asia/Kolkata'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Merchant) TableName() string {
	return "merchants"
}

type Store struct {
	ID         uint      `gorm:"primaryKey"`
	MerchantID uint      `gorm:"not null;index"`
	Name       string    `gorm:"size:255;not null"`
	Address    string    `gorm:"size:512"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Store) TableName() string {
	return "stores"
}

