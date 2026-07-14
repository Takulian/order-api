package grpcserver

import (
	"errors"
	"order-api/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, service.ErrCustomerRequired),
		errors.Is(err, service.ErrProductRequired),
		errors.Is(err, service.ErrInvalidQuantity):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, service.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}
