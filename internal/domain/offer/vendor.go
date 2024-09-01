package offer

import (
	"fmt"
	"log"
)

type VendorIdentity struct {
	id string
}

type VendorIdentities []VendorIdentity

func NewVendorIdentity(identity string) VendorIdentity {
	return VendorIdentity{identity}
}

func (v VendorIdentity) RawValue() string {
	return v.id
}

func (v VendorIdentities) String() string {
	rawValue := fmt.Sprintf("'%s'", v[0].RawValue())
	for i := 1; i < len(v); i++ {
		rawValue += fmt.Sprintf(", '%s'", v[i].RawValue())
	}
	log.Printf("vendor values %s", rawValue)
	return rawValue
}
