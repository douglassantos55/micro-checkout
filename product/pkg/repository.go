package pkg

import "fmt"

type Repository interface {
	Save(*Product) error
	Get(id string) (*Product, error)
	GetAll() ([]*Product, error)
}

type memoryRepository struct {
	products map[string]*Product
}

func NewMemoryRepository() Repository {
	products := make(map[string]*Product)
	for i := 0; i < 20; i++ {
		id := fmt.Sprintf("p_%d", i+1)
		products[id] = &Product{
			ID:      id,
			Name:    "Product " + id,
			Price:   float64((i + 1) + 6),
			InStock: 100 - i,
		}
	}
	return &memoryRepository{products}
}

func (r *memoryRepository) GetAll() ([]*Product, error) {
	products := make([]*Product, 0)
	for _, product := range r.products {
		products = append(products, product)
	}
	return products, nil
}

func (r *memoryRepository) Get(id string) (*Product, error) {
	product := r.products[id]
	return product, nil
}

func (r *memoryRepository) Save(product *Product) error {
	r.products[product.ID] = product
	return nil
}
