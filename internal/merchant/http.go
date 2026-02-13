package merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"upisettle/internal/auth"
)

func RegisterHTTP(rg *gin.RouterGroup, svc *Service) {
	rg.POST("/stores", func(c *gin.Context) {
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

		var req CreateStoreRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		store, err := svc.CreateStore(merchantID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, store)
	})

	rg.GET("/stores", func(c *gin.Context) {
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

		stores, err := svc.ListStores(merchantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stores)
	})
}

