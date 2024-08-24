package offer

type product struct {
	name  string
	price float32
}

func NewProduct(name string, price float32) (product, error) {
	p := product{name, price}
	if err := p.validate(); err != nil {
		return product{}, err
	}
	return p, nil
}

func (p product) Name() string {
	return p.name
}

func (p product) Price() float32 {
	return p.price
}
