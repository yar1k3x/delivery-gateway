package client

import (
	auth_proto "DeliveryGateway/proto/auth"
	delivery_proto "DeliveryGateway/proto/delivery"
	transport_proto "DeliveryGateway/proto/transport"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TransportClient struct {
	Client transport_proto.TransportServiceClient
}
type DeliveryClient struct {
	Client delivery_proto.DeliveryRequestServiceClient
}
type AuthClient struct {
	Client auth_proto.AuthServiceClient
}

func NewTransportServiceClient(addr string) (*TransportClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &TransportClient{
		Client: transport_proto.NewTransportServiceClient(conn),
	}, nil
}

func NewDeliveryRequestServiceClient(addr string) (*DeliveryClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &DeliveryClient{
		Client: delivery_proto.NewDeliveryRequestServiceClient(conn),
	}, nil
}

func NewAuthServiceClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		Client: auth_proto.NewAuthServiceClient(conn),
	}, nil
}

func (tc *TransportClient) WithToken(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}
func (tc *DeliveryClient) WithToken(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}
