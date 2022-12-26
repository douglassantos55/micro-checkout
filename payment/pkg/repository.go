package pkg

type Repository interface {
	GetPaymentMethods() ([]*PaymentMethod, error)
}

type memoryRepository struct {
	paymentMethods map[string]*PaymentMethod
}

func NewMemoryRepository() Repository {
	methods := map[string]*PaymentMethod{
		"cc":      {ID: "cc", Name: "Credit Card"},
		"cash":    {ID: "cash", Name: "Cash"},
		"deposit": {ID: "deposit", Name: "Bank deposit"},
	}
	return &memoryRepository{methods}
}

func (r *memoryRepository) GetPaymentMethods() ([]*PaymentMethod, error) {
	methods := make([]*PaymentMethod, 0)
	for _, method := range r.paymentMethods {
		methods = append(methods, method)
	}
	return methods, nil
}
