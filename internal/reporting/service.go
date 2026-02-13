package reporting

import (
	"time"

	"gorm.io/gorm"

	"upisettle/internal/matching"
	"upisettle/internal/order"
	"upisettle/internal/payment"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type DailySummary struct {
	Date              string `json:"date"`
	TotalOrders       int    `json:"total_orders"`
	TotalSalesAmount  int64  `json:"total_sales_amount"`
	UPITotalAmount    int64  `json:"upi_total_amount"`
	CashTotalAmount   int64  `json:"cash_total_amount"`
	MatchedOrders     int    `json:"matched_orders"`
	UnmatchedOrders   int    `json:"unmatched_orders"`
	ExceptionsCount   int    `json:"exceptions_count"`
	ExceptionsAmount  int64  `json:"exceptions_amount"`
}

func (s *Service) GetDailySummary(merchantID, storeID uint, day time.Time) (DailySummary, error) {
	summary := DailySummary{
		Date: day.Format("2006-01-02"),
	}

	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	var orders []order.Order
	if err := s.db.
		Where("merchant_id = ? AND store_id = ? AND created_at >= ? AND created_at < ?", merchantID, storeID, start, end).
		Find(&orders).Error; err != nil {
		return summary, err
	}

	for _, o := range orders {
		summary.TotalOrders++
		summary.TotalSalesAmount += o.Amount
		if o.Status == order.StatusPending {
			summary.UnmatchedOrders++
		} else {
			summary.MatchedOrders++
		}
	}

	var payments []payment.Payment
	if err := s.db.
		Where("merchant_id = ? AND store_id = ? AND time >= ? AND time < ?", merchantID, storeID, start, end).
		Find(&payments).Error; err != nil {
		return summary, err
	}

	for _, p := range payments {
		switch p.Channel {
		case payment.ChannelUPI:
			summary.UPITotalAmount += p.Amount
		case payment.ChannelCash:
			summary.CashTotalAmount += p.Amount
		}
	}

	var exceptions []matching.Exception
	if err := s.db.
		Where("merchant_id = ? AND store_id = ? AND created_at >= ? AND created_at < ?", merchantID, storeID, start, end).
		Find(&exceptions).Error; err != nil {
		return summary, err
	}

	summary.ExceptionsCount = len(exceptions)

	// Approximate exceptions amount: sum associated order or payment amounts.
	for _, ex := range exceptions {
		if ex.OrderID != nil {
			var o order.Order
			if err := s.db.First(&o, *ex.OrderID).Error; err == nil {
				summary.ExceptionsAmount += o.Amount
				continue
			}
		}
		if ex.PaymentID != nil {
			var p payment.Payment
			if err := s.db.First(&p, *ex.PaymentID).Error; err == nil {
				summary.ExceptionsAmount += p.Amount
			}
		}
	}

	return summary, nil
}

type ExceptionDTO struct {
	ID         uint      `json:"id"`
	Type       string    `json:"type"`
	Reason     string    `json:"reason"`
	OrderID    *uint     `json:"order_id,omitempty"`
	PaymentID  *uint     `json:"payment_id,omitempty"`
	Resolved   bool      `json:"resolved"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *Service) ListExceptions(merchantID, storeID uint, day time.Time) ([]ExceptionDTO, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	var exceptions []matching.Exception
	if err := s.db.
		Where("merchant_id = ? AND store_id = ? AND created_at >= ? AND created_at < ?", merchantID, storeID, start, end).
		Order("created_at ASC").
		Find(&exceptions).Error; err != nil {
		return nil, err
	}

	result := make([]ExceptionDTO, 0, len(exceptions))
	for _, ex := range exceptions {
		result = append(result, ExceptionDTO{
			ID:        ex.ID,
			Type:      ex.Type,
			Reason:    ex.Reason,
			OrderID:   ex.OrderID,
			PaymentID: ex.PaymentID,
			Resolved:  ex.Resolved,
			CreatedAt: ex.CreatedAt,
		})
	}
	return result, nil
}

