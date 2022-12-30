package pkg

import "fmt"

type Repository interface {
	GetOrders() ([]*Order, error)
	SaveOrder(order *Order) (*Order, error)
}

type memoryRepository struct {
	orders map[string]*Order
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		orders: make(map[string]*Order),
	}
}

func (r *memoryRepository) GetOrders() ([]*Order, error) {
	orders := make([]*Order, 0)
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *memoryRepository) SaveOrder(order *Order) (*Order, error) {
	id := fmt.Sprintf("order_%d", len(r.orders)+1)
	order.ID = id
	r.orders[id] = order
	return order, nil
}
