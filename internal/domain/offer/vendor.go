package offer

type vendor struct {
	name string
}

func newVendor(name string) vendor {
	return vendor{name}
}

func (v vendor) Name() string {
	return v.name
}
