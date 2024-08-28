package offer

type product struct {
	name  string
	price float64
}

func NewProduct(name string, price float64) (product, error) {
	p := product{name, price}
	// if err := p.validate(); err != nil {
	// 	return product{}, err
	// }
	return p, nil
}

func (p product) Name() string {
	return p.name
}

func (p product) Price() float64 {
	return p.price
}
