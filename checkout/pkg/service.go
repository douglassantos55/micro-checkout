package pkg

import (
	"context"
	"fmt"
)

type Order struct {
	ID              string       `json:"id,omitempty"`
	CustomerID      string       `json:"customer_id" validate:"required"`
	Items           []*OrderItem `json:"items" validate:"required,dive"`
	PaymentMethodID string       `json:"payment_method_id" validate:"required"`
}

type OrderItem struct {
	ProductID string `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required"`
}

type Service interface {
	PlaceOrder(context.Context, Order) (*Order, error)
	GetOrders() ([]*Order, error)
}

type service struct {
	repository Repository
	validator  Validator
	broker     MessageBroker
}

func NewService(repository Repository, validator Validator, broker MessageBroker) Service {
	return &service{repository, validator, broker}
}

func (s *service) GetOrders() ([]*Order, error) {
	return s.repository.GetOrders()
}

func (s *service) PlaceOrder(ctx context.Context, order Order) (*Order, error) {
	if err := s.validator.Validate(order); err != nil {
		return nil, err
	}

	savedOrder, err := s.repository.SaveOrder(order)
	if err != nil {
		return nil, fmt.Errorf("could not place order, please try again")
	}

	if err := s.broker.Broadcast(ctx, "order-placed", savedOrder); err != nil {
		return nil, err
	}

	return savedOrder, nil
}
