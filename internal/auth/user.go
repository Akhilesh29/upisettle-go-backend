package auth

import "time"

// User represents a user of the system.
// Typically the owner or staff of a merchant.
type User struct {
	ID         uint      `gorm:"primaryKey"`
	MerchantID uint      `gorm:"not null;index"`
	Name       string    `gorm:"size:255;not null"`
	Phone      string    `gorm:"size:20;index"`
	Email      string    `gorm:"size:255;uniqueIndex"`
	Password   string    `gorm:"size:255;not null"` // bcrypt hash
	Role       string    `gorm:"size:50;not null"`  // e.g. owner, staff
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (User) TableName() string {
	return "users"
}


