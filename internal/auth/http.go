package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHTTP wires auth HTTP handlers.
func RegisterHTTP(rg *gin.RouterGroup, svc *Service) {
	rg.POST("/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := svc.RegisterOwner(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	})

	rg.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := svc.Login(req)
		if err != nil {
			if err == ErrInvalidCredentials {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
}

