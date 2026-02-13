package order

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"upisettle/internal/auth"
)

func RegisterHTTP(rg *gin.RouterGroup, svc *Service) {
	rg.POST("/stores/:storeId/orders", func(c *gin.Context) {
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

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order, err := svc.CreateOrder(merchantID, storeID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, order)
	})

	rg.GET("/stores/:storeId/orders", func(c *gin.Context) {
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

		dateStr := c.Query("date")
		if dateStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "date query param is required (YYYY-MM-DD)"})
			return
		}
		day, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected YYYY-MM-DD"})
			return
		}

		orders, err := svc.ListOrdersByDate(merchantID, storeID, day)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	})
}

