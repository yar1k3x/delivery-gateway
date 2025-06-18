package handlers

import (
	"DeliveryGateway/client"
	auth "DeliveryGateway/proto/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func GetUsers(grpcClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIDParam = c.Query("user_id")
		var roleIDParam = c.Query("role_id")

		req := &auth.GetUsersRequest{}

		if userIDParam != "" {
			if userID, err := strconv.Atoi(userIDParam); err == nil {
				req.UserId = &wrapperspb.Int32Value{Value: int32(userID)}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
				return
			}
		}

		if roleIDParam != "" {
			if statusID, err := strconv.Atoi(roleIDParam); err == nil {
				req.RoleId = &wrapperspb.Int32Value{Value: int32(statusID)}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
				return
			}
		}

		resp, err := grpcClient.Client.GetUsers(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": resp.Users})
	}
}
