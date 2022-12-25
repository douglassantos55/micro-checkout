package pkg

type Customer struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type Filters struct {
	Page int
}

type Service interface {
	GetCustomer(id string) (*Customer, error)
	ListCustomers(filters Filters) (QueryResult[*Customer], error)
	CreateCustomer(data Customer) (*Customer, error)
	UpdateCustomer(id string, data Customer) (*Customer, error)
	DeleteCustomer(id string) error
}

type customerservice struct {
	validator  Validator
	repository Repository
}

func NewService(repo Repository, validator Validator) Service {
	return &customerservice{
		repository: repo,
		validator:  validator,
	}
}

func (s *customerservice) GetCustomer(id string) (*Customer, error) {
	customer, err := s.repository.GetCustomer(id)
	if customer == nil {
		return nil, makeError(404, "customer not found")
	}
	return customer, err
}

func (s *customerservice) ListCustomers(filters Filters) (QueryResult[*Customer], error) {
	return s.repository.ListCustomers(filters)
}

func (s *customerservice) CreateCustomer(data Customer) (*Customer, error) {
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}
	return s.repository.CreateCustomer(data)
}

func (s *customerservice) UpdateCustomer(id string, data Customer) (*Customer, error) {
	customer, err := s.repository.GetCustomer(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, makeError(404, "customer not found")
	}
	if err := s.validator.Validate(data); err != nil {
		return nil, err
	}
	return s.repository.UpdateCustomer(id, data)
}

func (s *customerservice) DeleteCustomer(id string) error {
	customer, err := s.repository.GetCustomer(id)
	if err != nil {
		return err
	}
	if customer == nil {
		return makeError(404, "customer not found")
	}
	return s.repository.DeleteCustomer(id)
}
