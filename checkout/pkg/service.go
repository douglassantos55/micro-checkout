package pkg

import "fmt"

type Order struct {
	CustomerID      string       `json:"customer_id" validate:"required"`
	Items           []*OrderItem `json:"items" validate:"required,dive"`
	PaymentMethodID string       `json:"payment_method_id" validate:"required"`
}

type OrderItem struct {
	ProductID string `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required"`
}

type Service interface {
	PlaceOrder(order Order) (*Order, error)
}

type service struct {
	repository Repository
	validator  Validator
}

func NewService(repository Repository, validator Validator) Service {
	return &service{repository, validator}
}

func (s *service) PlaceOrder(order Order) (*Order, error) {
	// validate
	if err := s.validator.Validate(order); err != nil {
		return nil, err
	}

	// save order
	savedOrder, err := s.repository.SaveOrder(order)
	if err != nil {
		return nil, fmt.Errorf("could not place order, please try again")
	}

	// process payment
	// -- s.messageBroker.dispatch("OrderPlaced")

	// reduce stock
	// -- communicate with product service

	return savedOrder, nil
}
