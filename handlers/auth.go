package handlers

import (
	"DeliveryGateway/client"
	auth "DeliveryGateway/proto/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		resp, err := authClient.Client.Register(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": resp.Token, "error": resp.Error})
	}
}

func Login(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		resp, err := authClient.Client.Login(c, &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": resp.Token, "error": resp.Error})
	}
}
