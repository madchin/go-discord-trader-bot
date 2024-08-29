package offer

type product struct {
	name  string
	price float64
}

func NewProduct(name string, price float64) product {
	return product{name, price}
}

func (p product) Name() string {
	return p.name
}

func (p product) Price() float64 {
	return p.price
}
