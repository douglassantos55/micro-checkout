package pkg

type Repository interface {
	GetCustomer(id string) (*Customer, error)
	CreateCustomer(data Customer) (*Customer, error)
	ListCustomers(filters Filters) ([]*Customer, error)
}

type inMemoryRepository struct {
	customers map[string]*Customer
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		customers: make(map[string]*Customer),
	}
}

func (r *inMemoryRepository) GetCustomer(id string) (*Customer, error) {
	customer, found := r.customers[id]
	if !found {
		return nil, makeError(404, "customer not found")
	}
	return customer, nil
}

func (r *inMemoryRepository) ListCustomers(filters Filters) ([]*Customer, error) {
	customers := make([]*Customer, 0)
	for _, customer := range r.customers {
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *inMemoryRepository) CreateCustomer(data Customer) (*Customer, error) {
	r.customers[data.ID] = &data
	return &data, nil
}
