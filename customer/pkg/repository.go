package pkg

import (
	"fmt"
	"math"
)

type QueryResult[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page  int `json:"curr"`
}

const PAGINATION_SIZE = 20

type Repository interface {
	GetCustomer(id string) (*Customer, error)
	CreateCustomer(data Customer) (*Customer, error)
	UpdateCustomer(id string, data Customer) (*Customer, error)
	DeleteCustomer(id string) error
	ListCustomers(filters Filters) (QueryResult[*Customer], error)
}

type inMemoryRepository struct {
	customers map[string]*Customer
}

func NewInMemoryRepository() Repository {
	customers := make(map[string]*Customer, 0)
	for i := 0; i < 100; i++ {
		id := "c_" + fmt.Sprintf("%d", i)
		customers[id] = &Customer{
			ID:    id,
			Name:  "customer " + id,
			Email: "customer" + id + "@email.com",
		}
	}
	return &inMemoryRepository{
		customers: customers,
	}
}

func (r *inMemoryRepository) GetCustomer(id string) (*Customer, error) {
	return r.customers[id], nil
}

func (r *inMemoryRepository) ListCustomers(filters Filters) (QueryResult[*Customer], error) {
	customers := make([]*Customer, 0)
	for _, customer := range r.customers {
		customers = append(customers, customer)
	}

	total := len(customers)
	if total > PAGINATION_SIZE {
		start := (filters.Page - 1) * PAGINATION_SIZE
		end := filters.Page * PAGINATION_SIZE

		customers = customers[start:end]
	}

	pages := math.Max(1, float64(total/PAGINATION_SIZE))

	return QueryResult[*Customer]{
		Items: customers,
		Total: total,
		Pages: int(pages),
		Page:  filters.Page,
	}, nil
}

func (r *inMemoryRepository) CreateCustomer(data Customer) (*Customer, error) {
	r.customers[data.ID] = &data
	return &data, nil
}

func (r *inMemoryRepository) UpdateCustomer(id string, data Customer) (*Customer, error) {
	r.customers[id] = &Customer{
		ID:    id,
		Name:  data.Name,
		Email: data.Email,
	}
	return r.customers[id], nil
}

func (r *inMemoryRepository) DeleteCustomer(id string) error {
	delete(r.customers, id)
	return nil
}
