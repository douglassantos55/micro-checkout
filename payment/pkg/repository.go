package pkg

import "fmt"

type Repository interface {
	SaveInvoice(*Invoice) error
	GetPaymentMethods() ([]*PaymentMethod, error)
}

type memoryRepository struct {
	invoices       map[string]*Invoice
	paymentMethods map[string]*PaymentMethod
}

func NewMemoryRepository() Repository {
	methods := map[string]*PaymentMethod{
		"cc":      {ID: "cc", Name: "Credit Card"},
		"cash":    {ID: "cash", Name: "Cash"},
		"deposit": {ID: "deposit", Name: "Bank deposit"},
	}
	return &memoryRepository{
		invoices:       make(map[string]*Invoice),
		paymentMethods: methods,
	}
}

func (r *memoryRepository) GetPaymentMethods() ([]*PaymentMethod, error) {
	methods := make([]*PaymentMethod, 0)
	for _, method := range r.paymentMethods {
		methods = append(methods, method)
	}
	return methods, nil
}

func (r *memoryRepository) SaveInvoice(invoice *Invoice) error {
	id := fmt.Sprintf("invoice_%d", len(r.invoices)+1)
	invoice.ID = id
	r.invoices[id] = invoice
	return nil
}
