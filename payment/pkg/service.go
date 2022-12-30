package pkg

type PaymentMethod struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID    string  `json:"id"`
	Total float64 `json:"total"`
}

type Invoice struct {
	ID      string  `json:"id"`
	Total   float64 `json:"total"`
	Status  string  `json:"status"`
	OrderID string  `json:"order_id"`
}

type Service interface {
	GetInvoices() ([]*Invoice, error)
	ProcessPayment(order Order) (*Invoice, error)
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

func (s *service) GetInvoices() ([]*Invoice, error) {
	return s.repository.GetInvoices()
}

func (s *service) ProcessPayment(order Order) (*Invoice, error) {
	invoice := &Invoice{
		Total:   order.Total,
		Status:  "pending",
		OrderID: order.ID,
	}
	if err := s.repository.SaveInvoice(invoice); err != nil {
		return nil, err
	}
	return invoice, nil
}
