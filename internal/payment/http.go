package payment

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"upisettle/internal/auth"
)

func RegisterHTTP(rg *gin.RouterGroup, svc *Service) {
	rg.POST("/stores/:storeId/payments", func(c *gin.Context) {
		rawMerchantID, ok := c.Get(auth.ContextMerchantIDKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing merchant context"})
			return
		}
		merchantID, ok := rawMerchantID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid merchant context"})
			return
		}

		storeIDParam := c.Param("storeId")
		storeIDUint64, err := strconv.ParseUint(storeIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid storeId"})
			return
		}
		storeID := uint(storeIDUint64)

		var req CreatePaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// If client omits time, default to now (since binding requires it, this is just a safeguard)
		if req.Time.IsZero() {
			req.Time = time.Now()
		}

		p, err := svc.CreatePayment(merchantID, storeID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, p)
	})

	rg.POST("/stores/:storeId/cash-payments", func(c *gin.Context) {
		rawMerchantID, ok := c.Get(auth.ContextMerchantIDKey)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing merchant context"})
			return
		}
		merchantID, ok := rawMerchantID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid merchant context"})
			return
		}

		storeIDParam := c.Param("storeId")
		storeIDUint64, err := strconv.ParseUint(storeIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid storeId"})
			return
		}
		storeID := uint(storeIDUint64)

		var req CreateCashPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		p, err := svc.CreateCashPayment(merchantID, storeID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, p)
	})
}

