package pkg

import "fmt"

type Product struct {
	ID      string
	Name    string
	Price   float64
	InStock int
}

type Service interface {
	GetProduct(id string) (*Product, error)
	ListProducts() ([]*Product, error)
	ReduceStock(string, int) error
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetProduct(id string) (*Product, error) {
	return s.repository.Get(id)
}

func (s *service) ListProducts() ([]*Product, error) {
	return s.repository.GetAll()
}

func (s *service) ReduceStock(id string, qty int) error {
	product, err := s.repository.Get(id)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("product %v not found", id)
	}
	product.InStock -= qty
	if err := s.repository.Save(product); err != nil {
		return fmt.Errorf("could not reduce product stock")
	}
	return nil
}
