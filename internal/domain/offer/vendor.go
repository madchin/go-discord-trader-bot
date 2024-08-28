package offer

type vendor struct {
	name string
}

func newVendor(name string) (vendor, error) {
	v := vendor{name}
	if err := v.validate(); err != nil {
		return vendor{}, err
	}
	return v, nil
}

func (v vendor) Name() string {
	return v.name
}
