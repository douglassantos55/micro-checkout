package pkg

import "fmt"

type Repository interface {
	SaveOrder(order Order) (*Order, error)
}

type memoryRepository struct {
	orders map[string]*Order
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		orders: make(map[string]*Order),
	}
}

func (r *memoryRepository) SaveOrder(order Order) (*Order, error) {
	id := fmt.Sprintf("order_%d", len(r.orders)+1)
	r.orders[id] = &order
	return &order, nil
}
