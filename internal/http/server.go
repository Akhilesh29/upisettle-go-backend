package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"upisettle/internal/auth"
	"upisettle/internal/config"
	"upisettle/internal/logger"
	"upisettle/internal/matching"
	"upisettle/internal/merchant"
	"upisettle/internal/order"
	"upisettle/internal/payment"
	"upisettle/internal/reporting"
)

type Server struct {
	cfg     config.Config
	log     logger.Logger
	db      *gorm.DB
	engine  *gin.Engine
	authSvc *auth.Service
}

func NewServer(cfg config.Config, log logger.Logger, db *gorm.DB) *Server {
	gin.SetMode(gin.ReleaseMode)
	if cfg.Env == "development" {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	s := &Server{
		cfg:     cfg,
		log:     log,
		db:      db,
		engine:  engine,
		authSvc: auth.NewService(db, cfg),
	}

	s.registerRoutes()
	return s
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	s.log.Printf("starting HTTP server on %s", addr)
	return s.engine.Run(addr)
}

func (s *Server) registerRoutes() {
	// Health check
	s.engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := s.engine.Group("/api/v1")

	// Auth routes (no auth required)
	authGroup := api.Group("/auth")
	auth.RegisterHTTP(authGroup, s.authSvc)

	// Protected routes
	protected := api.Group("")
	protected.Use(auth.AuthMiddleware(s.cfg.JWTSecret))

	merchantSvc := merchant.NewService(s.db)
	orderSvc := order.NewService(s.db)
	paymentSvc := payment.NewService(s.db)
	matchingSvc := matching.NewService(s.db)
	reportingSvc := reporting.NewService(s.db)

	merchant.RegisterHTTP(protected, merchantSvc)
	order.RegisterHTTP(protected, orderSvc)
	payment.RegisterHTTP(protected, paymentSvc)
	matching.RegisterHTTP(protected, matchingSvc)
	reporting.RegisterHTTP(protected, reportingSvc)
}

