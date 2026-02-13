package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"upisettle/internal/config"
	"upisettle/internal/merchant"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	db     *gorm.DB
	cfg    config.Config
}

func NewService(db *gorm.DB, cfg config.Config) *Service {
	return &Service{
		db:  db,
		cfg: cfg,
	}
}

type RegisterRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	MerchantName string `json:"merchant_name" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (s *Service) RegisterOwner(req RegisterRequest) (AuthResponse, error) {
	var resp AuthResponse

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resp, err
	}

	var token string
	err = s.db.Transaction(func(tx *gorm.DB) error {
		m := merchant.Merchant{
			Name:           req.MerchantName,
			PrimaryContact: req.Name,
		}
		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		u := User{
			MerchantID: m.ID,
			Name:       req.Name,
			Phone:      req.Phone,
			Email:      req.Email,
			Password:   string(hashed),
			Role:       "owner",
		}
		if err := tx.Create(&u).Error; err != nil {
			return err
		}

		t, err := s.generateToken(u)
		if err != nil {
			return err
		}
		token = t
		return nil
	})
	if err != nil {
		return resp, err
	}

	resp.Token = token
	return resp, nil
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Service) Login(req LoginRequest) (AuthResponse, error) {
	var resp AuthResponse

	var user User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, ErrInvalidCredentials
		}
		return resp, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return resp, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return resp, err
	}
	resp.Token = token
	return resp, nil
}

type Claims struct {
	UserID     uint   `json:"user_id"`
	MerchantID uint   `json:"merchant_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) generateToken(user User) (string, error) {
	claims := Claims{
		UserID:     user.ID,
		MerchantID: user.MerchantID,
		Role:       user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// You can add expiry here if desired.
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}


