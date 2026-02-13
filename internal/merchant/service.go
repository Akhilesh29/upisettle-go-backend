package merchant

import (
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CreateStoreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func (s *Service) CreateStore(merchantID uint, req CreateStoreRequest) (Store, error) {
	store := Store{
		MerchantID: merchantID,
		Name:       req.Name,
		Address:    req.Address,
	}
	if err := s.db.Create(&store).Error; err != nil {
		return Store{}, err
	}
	return store, nil
}

func (s *Service) ListStores(merchantID uint) ([]Store, error) {
	var stores []Store
	if err := s.db.Where("merchant_id = ?", merchantID).Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

