package matching

import (
	"time"

	"gorm.io/gorm"

	"upisettle/internal/order"
	"upisettle/internal/payment"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type ReconcileSummary struct {
	MatchedOrders     int `json:"matched_orders"`
	UnmatchedOrders   int `json:"unmatched_orders"`
	UnmatchedPayments int `json:"unmatched_payments"`
}

// Reconcile performs a simple matching for a given merchant, store and date.
func (s *Service) Reconcile(merchantID, storeID uint, day time.Time) (ReconcileSummary, error) {
	summary := ReconcileSummary{}

	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	var orders []order.Order
	if err := s.db.
		Where("merchant_id = ? AND store_id = ? AND created_at >= ? AND created_at < ? AND status = ?", merchantID, storeID, start, end, order.StatusPending).
		Order("created_at ASC").
		Find(&orders).Error; err != nil {
		return summary, err
	}

	var payments []payment.Payment
	if err := s.db.
		Where("merchant_id = ? AND store_id = ? AND time >= ? AND time < ?", merchantID, storeID, start, end).
		Order("time ASC").
		Find(&payments).Error; err != nil {
		return summary, err
	}

	// Load existing matches to avoid duplicating work.
	var matches []Match
	if len(orders) > 0 {
		orderIDs := make([]uint, 0, len(orders))
		for _, o := range orders {
			orderIDs = append(orderIDs, o.ID)
		}
		if err := s.db.Where("order_id IN ?", orderIDs).Find(&matches).Error; err != nil && err != gorm.ErrRecordNotFound {
			return summary, err
		}
	}

	existingOrderMatched := make(map[uint]bool)
	existingPaymentMatched := make(map[uint]bool)
	for _, m := range matches {
		existingOrderMatched[m.OrderID] = true
		existingPaymentMatched[m.PaymentID] = true
	}

	usedPayment := make(map[uint]bool)

	// For each pending order, find a payment with the same amount within the day that is not yet matched.
	for _, o := range orders {
		if existingOrderMatched[o.ID] {
			continue
		}

		var candidates []payment.Payment
		for _, p := range payments {
			if existingPaymentMatched[p.ID] || usedPayment[p.ID] {
				continue
			}
			if p.Amount == o.Amount {
				candidates = append(candidates, p)
			}
		}

		if len(candidates) == 1 {
			p := candidates[0]
			now := time.Now()
			m := Match{
				OrderID:    o.ID,
				PaymentID:  p.ID,
				Confidence: 1.0,
				MatchedAt:  now,
			}
			if err := s.db.Create(&m).Error; err != nil {
				return summary, err
			}

			usedPayment[p.ID] = true

			// Update order status and paid_at.
			o.Status = order.StatusPaidUPI
			o.PaidAt = &p.Time
			if err := s.db.Save(&o).Error; err != nil {
				return summary, err
			}
			summary.MatchedOrders++
		} else if len(candidates) == 0 {
			// No candidate payment found for this order.
			ex := Exception{
				MerchantID: merchantID,
				StoreID:    storeID,
				OrderID:    &o.ID,
				Type:       ExceptionUnmatchedOrder,
				Reason:     "no payment candidate found for order",
			}
			if err := s.db.Create(&ex).Error; err != nil {
				return summary, err
			}
			summary.UnmatchedOrders++
		} else {
			// Multiple candidates; ambiguous amount.
			ex := Exception{
				MerchantID: merchantID,
				StoreID:    storeID,
				OrderID:    &o.ID,
				Type:       ExceptionAmountMismatch,
				Reason:     "multiple payment candidates with same amount",
			}
			if err := s.db.Create(&ex).Error; err != nil {
				return summary, err
			}
			summary.UnmatchedOrders++
		}
	}

	// Any payments not used or previously matched get an UNMATCHED_PAYMENT exception.
	for _, p := range payments {
		if existingPaymentMatched[p.ID] || usedPayment[p.ID] {
			continue
		}
		ex := Exception{
			MerchantID: merchantID,
			StoreID:    storeID,
			PaymentID:  &p.ID,
			Type:       ExceptionUnmatchedPayment,
			Reason:     "no matching order found for payment",
		}
		if err := s.db.Create(&ex).Error; err != nil {
			return summary, err
		}
		summary.UnmatchedPayments++
	}

	return summary, nil
}

