package handlers

import (
	"context"
	"net/http"
	"strconv"

	"DeliveryGateway/client"
	proto "DeliveryGateway/proto/delivery"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func CreateDeliveryRequest(grpcClient *client.DeliveryClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Weight        float32 `json:"weight"`
			FromLocation  string  `json:"from_location"`
			ToLocation    string  `json:"to_location"`
			PreferredDate string  `json:"preferred_date"`
			CreatedBy     int32   `json:"created_by"`
			ResponsibleID int32   `json:"responsible_id"`
			CreatedAt     string  `json:"created_at"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		req := &proto.CreateRequestInput{
			Weight:        input.Weight,
			FromLocation:  input.FromLocation,
			ToLocation:    input.ToLocation,
			PreferredDate: input.PreferredDate,
			CreatedBy:     input.CreatedBy,
			ResponsibleId: input.ResponsibleID,
			CreatedAt:     input.CreatedAt,
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.CreateRequest(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"request_id": resp.RequestId,
		})
	}
}

func GetDeliveryRequests(grpcClient *client.DeliveryClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIDParam = c.Query("user_id")
		var statusIDParam = c.Query("status_id")

		req := &proto.GetRequestInput{}

		if userIDParam != "" {
			if userID, err := strconv.Atoi(userIDParam); err == nil {
				req.UserId = &wrapperspb.Int32Value{Value: int32(userID)}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
				return
			}
		}

		if statusIDParam != "" {
			if statusID, err := strconv.Atoi(statusIDParam); err == nil {
				req.StatusId = &wrapperspb.Int32Value{Value: int32(statusID)}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status_id"})
				return
			}
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.GetRequest(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"requests": resp.Requests})
	}
}

func UpdateDeliveryRequest(grpcClient *client.DeliveryClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			RequestID     int32    `json:"request_id" binding:"required"`
			Weight        *float32 `json:"weight"`
			FromLocation  *string  `json:"from_location"`
			ToLocation    *string  `json:"to_location"`
			PreferredDate *string  `json:"preferred_date"`
			ResponsibleID *int32   `json:"responsible_id"`
			StatusID      *int32   `json:"status_id"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req := &proto.UpdateRequestInput{
			RequestId: &wrapperspb.Int32Value{Value: input.RequestID},
		}

		if input.Weight != nil {
			req.Weight = wrapperspb.Float(*input.Weight)
		}
		if input.FromLocation != nil {
			req.FromLocation = wrapperspb.String(*input.FromLocation)
		}
		if input.ToLocation != nil {
			req.ToLocation = wrapperspb.String(*input.ToLocation)
		}
		if input.PreferredDate != nil {
			req.PreferredDate = wrapperspb.String(*input.PreferredDate)
		}
		if input.ResponsibleID != nil {
			req.ResponsibleId = wrapperspb.Int32(*input.ResponsibleID)
		}
		if input.StatusID != nil {
			req.StatusId = wrapperspb.Int32(*input.StatusID)
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.UpdateRequest(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": resp.Success})
	}
}

func DeleteDeliveryRequest(grpcClient *client.DeliveryClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Param("id")
		if requestID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request_id is required"})
			return
		}

		idInt, err := strconv.Atoi(requestID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request_id"})
			return
		}

		req := &proto.DeleteRequestInput{
			RequestId: wrapperspb.Int32(int32(idInt)), // вот так оборачиваем
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.DeleteRequest(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": resp.Success})
	}
}

func ExportDeliveryRequestsCSV(grpcClient *client.DeliveryClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in proto.GetRequestInput
		if v := c.Query("user_id"); v != "" {
			if id, err := strconv.Atoi(v); err == nil {
				in.UserId = &wrapperspb.Int32Value{Value: int32(id)}
			}
		}
		if v := c.Query("status_id"); v != "" {
			if id, err := strconv.Atoi(v); err == nil {
				in.StatusId = &wrapperspb.Int32Value{Value: int32(id)}
			}
		}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.ExportCSV(ctx, &in)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", `attachment; filename="delivery_requests.csv"`)
		c.Data(http.StatusOK, "text/csv; charset=utf-8", resp.Data)
	}
}

func GetRequestStatuses(grpcClient *client.DeliveryClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &emptypb.Empty{}

		token := c.GetHeader("Authorization")
		ctx := grpcClient.WithToken(context.Background(), token)

		resp, err := grpcClient.Client.GetRequestStatuses(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"statuses": resp.Statuses})
	}
}
