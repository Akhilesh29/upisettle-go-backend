package order

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CreateOrderRequest struct {
	Amount      int64  `json:"amount" binding:"required"`
	ExternalRef string `json:"external_ref"`
}

func (s *Service) CreateOrder(merchantID, storeID uint, req CreateOrderRequest) (Order, error) {
	order := Order{
		MerchantID:  merchantID,
		StoreID:     storeID,
		ExternalRef: req.ExternalRef,
		Amount:      req.Amount,
		Status:      StatusPending,
	}
	if err := s.db.Create(&order).Error; err != nil {
		return Order{}, err
	}
	return order, nil
}

func (s *Service) ListOrdersByDate(merchantID, storeID uint, day time.Time) ([]Order, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	var orders []Order
	if err := s.db.
		Where("merchant_id = ? AND store_id = ? AND created_at >= ? AND created_at < ?", merchantID, storeID, start, end).
		Order("created_at ASC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

