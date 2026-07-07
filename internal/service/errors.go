package service

import (
	"errors"
)

var (
	ErrCustomerRequired = errors.New("customer is required")
	ErrProductRequired = errors.New("product is required")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
)