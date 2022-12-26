package pkg

type PaymentMethod struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	GetPaymentMethods() ([]*PaymentMethod, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) GetPaymentMethods() ([]*PaymentMethod, error) {
	return s.repository.GetPaymentMethods()
}
