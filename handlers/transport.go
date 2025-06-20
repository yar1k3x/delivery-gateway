package handlers

import (
	"DeliveryGateway/client"
	proto "DeliveryGateway/proto/transport"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func CreateTransport(grpcClient *client.TransportClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Number          string `json:"number"`
			TransportName   string `json:"transport_name"`
			TypeID          int32  `json:"type_id"`
			IsActive        int32  `json:"is_active"`
			CurrentDriverID int32  `json:"current_driver_id"`
			ImageUrl        string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)
		log.Println("Creating transport", req)

		resp, err := grpcClient.Client.CreateTransport(ctx, &proto.CreateTransportRequest{
			Number:          req.Number,
			TransportName:   req.TransportName,
			TypeId:          req.TypeID,
			IsActive:        req.IsActive,
			CurrentDriverId: req.CurrentDriverID,
			ImageUrl:        req.ImageUrl,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":      resp.Success,
			"transport_id": resp.TransportId,
		})
	}
}

func CreateTransportLog(grpcClient *client.TransportClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TransportId   int32  `json:"transport_id"`
			ServiceTypeId int32  `json:"service_type_id"`
			ServiceDate   string `json:"service_date"`
			Description   string `json:"description"`
			Mileage       int32  `json:"mileage"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.CreateTransportLog(ctx, &proto.CreateTransportLogRequest{
			TransportId:   req.TransportId,
			ServiceTypeId: req.ServiceTypeId,
			ServiceDate:   req.ServiceDate,
			Description:   req.Description,
			Mileage:       req.Mileage,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
		})
	}
}

func GetTransportInfo(grpcClient *client.TransportClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req proto.GetTransportInfoRequest

		if val, ok := c.GetQuery("transport_id"); ok {
			if id, err := strconv.Atoi(val); err == nil {
				req.TransportId = wrapperspb.Int32(int32(id))
			}
		}
		if val, ok := c.GetQuery("is_active"); ok {
			if active, err := strconv.Atoi(val); err == nil {
				req.IsActive = wrapperspb.Int32(int32(active))
			}
		}
		if val, ok := c.GetQuery("current_driver_id"); ok {
			if id, err := strconv.Atoi(val); err == nil {
				req.CurrentDriverId = wrapperspb.Int32(int32(id))
			}
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.GetTransportInfo(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"transports": resp.Transports})
	}
}

func GetTransportLogInfo(grpcClient *client.TransportClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req proto.GetTransportLogsInfoRequest

		if val, ok := c.GetQuery("transport_id"); ok {
			if id, err := strconv.Atoi(val); err == nil {
				req.TransportId = wrapperspb.Int32(int32(id))
			}
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.GetTransportLogsInfo(ctx, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"transportLogs": resp.TransportLogs})
	}
}

func UpdateTransport(grpcClient *client.TransportClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			TransportID     *int32  `json:"transport_id"`
			IsActive        *int32  `json:"is_active"`
			CurrentDriverID *int32  `json:"current_driver_id"`
			ImageUrl        *string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&json); err != nil || json.TransportID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing transport_id"})
			return
		}

		req := &proto.UpdateTransportRequest{
			TransportId: wrapperspb.Int32(*json.TransportID),
		}
		if json.IsActive != nil {
			req.IsActive = wrapperspb.Int32(*json.IsActive)
		}
		if json.CurrentDriverID != nil {
			req.CurrentDriverId = wrapperspb.Int32(*json.CurrentDriverID)
		}
		if json.ImageUrl != nil {
			req.ImageUrl = wrapperspb.String(*json.ImageUrl)
		}
		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.UpdateTransport(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": resp.Success,
		})
	}
}

func GetTransportType(grpcClient *client.TransportClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &emptypb.Empty{}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.GetTransportType(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"types": resp.Types})
	}
}

func GetServiceType(grpcClient *client.TransportClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &emptypb.Empty{}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.GetServiceType(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"types": resp.Types})
	}
}
