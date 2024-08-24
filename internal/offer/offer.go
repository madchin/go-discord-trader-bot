package offer

type Offer struct {
	product product
	count   int
}

type Vendor struct {
	name string
}

type Sell struct {
	vendor Vendor
	offer  Offer
}

type Buy struct {
	vendor Vendor
	offer  Offer
}

type Offers map[Vendor][]Offer

func NewOffer(product product, count int) (Offer, error) {
	off := Offer{product, count}
	if err := off.validate(); err != nil {
		return Offer{}, err
	}
	return off, nil
}

func NewVendor(name string) (Vendor, error) {
	v := Vendor{name}
	if err := v.validate(); err != nil {
		return Vendor{}, err
	}
	return v, nil
}

func NewBuy(vendor string, product product, count int) (Buy, error) {
	offer, err := NewOffer(product, count)
	if err != nil {
		return Buy{}, err
	}
	vend, err := NewVendor(vendor)
	if err != nil {
		return Buy{}, err
	}
	return Buy{vend, offer}, nil
}

func NewSell(vendor string, product product, count int) (Sell, error) {
	offer, err := NewOffer(product, count)
	if err != nil {
		return Sell{}, err
	}
	vend, err := NewVendor(vendor)
	if err != nil {
		return Sell{}, err
	}
	return Sell{vend, offer}, nil
}

func (s Sell) Vendor() Vendor {
	return s.vendor
}

func (s Sell) Offer() Offer {
	return s.offer
}

func (b Buy) Vendor() Vendor {
	return b.vendor
}

func (b Buy) Offer() Offer {
	return b.offer
}

func (o Offer) Product() product {
	return o.product
}

func (o Offer) IsEqual(off Offer) bool {
	return o.product.name == off.product.name && o.product.price == off.product.price
}

func (o Offer) Merge(off Offer) Offer {
	o.count += off.count
	return o
}
