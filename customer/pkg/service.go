package pkg

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Filters struct {
	Page int
}

type Service interface {
	GetCustomer(id string) (*Customer, error)
	CreateCustomer(data Customer) (*Customer, error)
	ListCustomers(filters Filters) (QueryResult[*Customer], error)
}

type customerservice struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &customerservice{repo}
}

func (s *customerservice) GetCustomer(id string) (*Customer, error) {
	return s.repository.GetCustomer(id)
}

func (s *customerservice) ListCustomers(filters Filters) (QueryResult[*Customer], error) {
	return s.repository.ListCustomers(filters)
}

func (s *customerservice) CreateCustomer(data Customer) (*Customer, error) {
	return s.repository.CreateCustomer(data)
}
