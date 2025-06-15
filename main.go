package main

import (
	"DeliveryGateway/client"
	"DeliveryGateway/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//grpcClient, err := client.NewTransportServiceClient(os.Getenv("TRANSPORT_MANAGEMENT_SERVICE_URL"))
	transportGrpcClient, err := client.NewTransportServiceClient("trolley.proxy.rlwy.net:27826")
	if err != nil {
		log.Fatalf("failed to connect to transport-manager-service: %v", err)
	}

	deliveryGrpcClient, err := client.NewDeliveryRequestServiceClient("ballast.proxy.rlwy.net:22723")
	if err != nil {
		log.Fatalf("failed to connect to gRPC service: %v", err)
	}

	transportApi := r.Group("/transport")
	{
		transportApi.POST("/create", handlers.CreateTransport(transportGrpcClient))
		transportApi.POST("/create-log", handlers.CreateTransportLog(transportGrpcClient))
		transportApi.GET("/", handlers.GetTransportInfo(transportGrpcClient))
		transportApi.GET("/logs", handlers.GetTransportLogInfo(transportGrpcClient))
		transportApi.PUT("/update", handlers.UpdateTransport(transportGrpcClient))
	}

	deliveryApi := r.Group("/delivery")
	{
		deliveryApi.POST("/create", handlers.CreateDeliveryRequest(deliveryGrpcClient))
		deliveryApi.GET("/requests", handlers.GetDeliveryRequests(deliveryGrpcClient))
		deliveryApi.PUT("/update", handlers.UpdateDeliveryRequest(deliveryGrpcClient))
		deliveryApi.DELETE("/delete/:id", handlers.DeleteDeliveryRequest(deliveryGrpcClient))
	}

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
