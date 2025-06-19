package main

import (
	"DeliveryGateway/client"
	"DeliveryGateway/handlers"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // или "*" если нужно на время открыть для всех
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//transportGrpcClient, err := client.NewTransportServiceClient(os.Getenv("TRANSPORT_MANAGEMENT_SERVICE_URL"))
	transportGrpcClient, err := client.NewTransportServiceClient("trolley.proxy.rlwy.net:27826")
	if err != nil {
		log.Fatalf("failed to connect to transport-manager-service: %v", err)
	}

	//deliveryGrpcClient, err := client.NewDeliveryRequestServiceClient(os.Getenv("DRS_URL"))
	deliveryGrpcClient, err := client.NewDeliveryRequestServiceClient("ballast.proxy.rlwy.net:22723")
	if err != nil {
		log.Fatalf("failed to connect to gRPC service: %v", err)
	}
	authGrpcClient, err := client.NewAuthServiceClient("shinkansen.proxy.rlwy.net:17837")
	if err != nil {
		log.Fatalf("failed to connect to AuthService: %v", err)
	}

	transportApi := r.Group("/transport")
	{
		transportApi.POST("/create", handlers.CreateTransport(transportGrpcClient))
		transportApi.POST("/create-log", handlers.CreateTransportLog(transportGrpcClient))
		transportApi.GET("/", handlers.GetTransportInfo(transportGrpcClient))
		transportApi.GET("/logs", handlers.GetTransportLogInfo(transportGrpcClient))
		transportApi.PUT("/update", handlers.UpdateTransport(transportGrpcClient))
		transportApi.GET("/types", handlers.GetTransportType(transportGrpcClient))
		transportApi.GET("/service-types", handlers.GetServiceType(transportGrpcClient))
	}

	deliveryApi := r.Group("/delivery")
	{
		deliveryApi.POST("/create", handlers.CreateDeliveryRequest(deliveryGrpcClient))
		deliveryApi.GET("/requests", handlers.GetDeliveryRequests(deliveryGrpcClient))
		deliveryApi.GET("/export", handlers.ExportDeliveryRequestsCSV(deliveryGrpcClient))
		deliveryApi.PUT("/update", handlers.UpdateDeliveryRequest(deliveryGrpcClient))
		deliveryApi.DELETE("/delete/:id", handlers.DeleteDeliveryRequest(deliveryGrpcClient))
		deliveryApi.GET("/statuses", handlers.GetRequestStatuses(deliveryGrpcClient))
	}

	authApi := r.Group("/auth")
	{
		authApi.POST("/register", handlers.Register(authGrpcClient))
		authApi.POST("/login", handlers.Login(authGrpcClient))
		authApi.GET("/users", handlers.GetUsers(authGrpcClient))
	}

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
