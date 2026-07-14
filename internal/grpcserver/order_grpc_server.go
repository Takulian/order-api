package grpcserver

import (
	"context"
	"order-api/internal/dto"
	"order-api/internal/model"
	"order-api/internal/service"
	orderv1 "order-api/proto/order/v1"
)

type OrderGRPCServer struct {
	orderv1.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewOrderGRPCServer(service *service.OrderService) *OrderGRPCServer {
	return &OrderGRPCServer{
		service: service,
	}
}

func toProtoOrder(o model.Order) *orderv1.Order {
	return &orderv1.Order{
		Id:       int32(o.ID),
		Customer: o.Customer,
		Product:  o.Product,
		Quantity: int32(o.Quantity),
		Status:   o.Status,
	}
}

func (s *OrderGRPCServer) GetOrders(ctx context.Context, req *orderv1.GetOrdersRequest) (*orderv1.GetOrdersResponse, error) {
	orders, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	protoOrders := make([]*orderv1.Order, 0, len(orders))
	for _, o := range orders {
		protoOrders = append(protoOrders, toProtoOrder(o))
	}

	return &orderv1.GetOrdersResponse{
		Orders: protoOrders,
	}, nil
}

func (s *OrderGRPCServer) GetOrderByID(ctx context.Context, req *orderv1.GetOrderByIDRequest) (*orderv1.GetOrderByIDResponse, error) {
	order, err := s.service.GetByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &orderv1.GetOrderByIDResponse{
		Order: toProtoOrder(order),
	}, nil
}

func (s *OrderGRPCServer) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	dtoReq := dto.CreateOrderRequest{
		Customer: req.Customer,
		Product:  req.Product,
		Quantity: int(req.Quantity),
	}

	order, err := s.service.Create(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &orderv1.CreateOrderResponse{
		Order: toProtoOrder(order),
	}, nil
}

func (s *OrderGRPCServer) UpdateOrder(ctx context.Context, req *orderv1.UpdateOrderRequest) (*orderv1.UpdateOrderResponse, error) {
	dtoReq := dto.UpdateOrderRequest{
		Customer: req.Customer,
		Product:  req.Product,
		Quantity: int(req.Quantity),
	}

	order, err := s.service.Update(ctx, int(req.Id), dtoReq)
	if err != nil {
		return nil, err
	}

	return &orderv1.UpdateOrderResponse{
		Order: toProtoOrder(order),
	}, nil
}

func (s *OrderGRPCServer) DeleteOrder(ctx context.Context, req *orderv1.DeleteOrderRequest) (*orderv1.DeleteOrderResponse, error) {
	err := s.service.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &orderv1.DeleteOrderResponse{}, nil
}
