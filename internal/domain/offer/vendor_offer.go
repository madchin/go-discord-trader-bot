package offer

import (
	"fmt"
	"strings"
)

type VendorIdentity struct {
	id string
}

type VendorOffer struct {
	vendorIdentity VendorIdentity
	Offer
}

type VendorOffers []VendorOffer

func NewVendorIdentity(identity string) VendorIdentity {
	return VendorIdentity{identity}
}

func NewVendorOffer(vendorIdentity VendorIdentity, offer Offer) VendorOffer {
	return VendorOffer{vendorIdentity, offer}
}

func (v VendorIdentity) RawValue() string {
	return v.id
}

func (o VendorOffer) VendorIdentity() VendorIdentity {
	return o.vendorIdentity
}

func (o VendorOffers) ToReadableMessage() string {
	offers := make([]string, len(o))
	for i := 0; i < len(o); i++ {
		price := fmt.Sprintf("%.2f", o[i].Product.price)
		offers[i] = fmt.Sprintf("Product: %s, Each Price: %s, Count: %d, Vendor: %s", o[i].Product.name, price, o[i].count, "<@!"+o[i].vendorIdentity.id+">")
	}
	return strings.Join(offers, ",\n")
}

func (o VendorOffers) Contains(candidate VendorOffer) (contains bool) {
	for _, off := range o {
		if off.isSameOffer(candidate) {
			contains = true
		}
	}
	return
}

func (o VendorOffers) NotExists() bool {
	return len(o) == 0
}

func (o VendorOffers) MergeSameOffers(candidate VendorOffer) VendorOffer {
	for _, off := range o {
		if off.isSameOffer(candidate) {
			candidate = off.merge(candidate)
		}
	}
	return candidate
}

func OnVendorOfferAdd(o VendorOffer) error {
	if err := o.validate(); err != nil {
		return err
	}
	if err := o.vendorIdentity.validate(); err != nil {
		return err
	}
	return nil
}

func (o VendorOffer) merge(off VendorOffer) VendorOffer {
	o.count += off.count
	return o
}

func (o VendorOffer) isSameOffer(off VendorOffer) bool {
	return o.Product.name == off.Product.name && o.Product.price == off.Product.price
}
