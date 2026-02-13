package payment

import (
	"time"

	"gorm.io/gorm"

	"upisettle/internal/order"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CreatePaymentRequest struct {
	Channel   string    `json:"channel" binding:"required"` // UPI, CASH, OTHER
	Amount    int64     `json:"amount" binding:"required"`
	Time      time.Time `json:"time" binding:"required"`
	UPIRef    string    `json:"upi_ref"`
	PayerVPA  string    `json:"payer_vpa"`
	PayerName string    `json:"payer_name"`
}

type CreateCashPaymentRequest struct {
	OrderID uint  `json:"order_id" binding:"required"`
	Amount  int64 `json:"amount" binding:"required"`
}

func (s *Service) CreatePayment(merchantID, storeID uint, req CreatePaymentRequest) (Payment, error) {
	p := Payment{
		MerchantID: merchantID,
		StoreID:    storeID,
		Channel:    req.Channel,
		Amount:     req.Amount,
		Time:       req.Time,
		UPIRef:     req.UPIRef,
		PayerVPA:   req.PayerVPA,
		PayerName:  req.PayerName,
	}
	if p.Currency == "" {
		p.Currency = "INR"
	}
	if err := s.db.Create(&p).Error; err != nil {
		return Payment{}, err
	}
	return p, nil
}

// CreateCashPayment records a cash payment and marks the order as paid cash.
func (s *Service) CreateCashPayment(merchantID, storeID uint, req CreateCashPaymentRequest) (Payment, error) {
	var payment Payment

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var o order.Order
		if err := tx.Where("id = ? AND merchant_id = ? AND store_id = ?", req.OrderID, merchantID, storeID).
			First(&o).Error; err != nil {
			return err
		}

		now := time.Now()
		o.Status = order.StatusPaidCash
		o.PaidAt = &now
		if err := tx.Save(&o).Error; err != nil {
			return err
		}

		payment = Payment{
			MerchantID: merchantID,
			StoreID:    storeID,
			Channel:    ChannelCash,
			Amount:     req.Amount,
			Currency:   "INR",
			Time:       now,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

