package offer

import (
	"errors"
)

var (
	ErrWrongActionType  = errors.New("unable to perform this action for this type of offer")
	ErrNotExistingOffer = errors.New("offer not exists")
)

type Offer struct {
	offerType offerType
	action    action
	vendor    vendor
	product   product
	count     int
}

type Offers []Offer

func NewOffer(offerType offerType, action action, vendor string, product product, count int) (Offer, error) {
	vend, err := newVendor(vendor)
	if err != nil {
		return Offer{}, err
	}
	off := Offer{offerType, action, vend, product, count}
	// if err := off.validate(); err != nil {
	// 	return Offer{}, err
	// }
	return off, nil
}

func NewListingOffer(offerType offerType, action action, productName string) Offer {
	return Offer{offerType, action, vendor{}, product{name: productName}, 0}
}

func (o Offer) Product() product {
	return o.product
}

func (o Offer) Vendor() vendor {
	return o.vendor
}

func (o Offer) Action() action {
	return o.action
}

func (o Offer) Type() offerType {
	return o.offerType
}

func (o Offer) Count() int {
	return o.count
}

func (o Offer) isSameOffer(off Offer) bool {
	return o.product.name == off.product.name && o.product.price == off.product.price
}

func (o Offer) merge(off Offer) Offer {
	o.count += off.count
	return o
}

func (o Offers) Contains(candidate Offer) (contains bool) {
	for _, off := range o {
		if off.isSameOffer(candidate) {
			contains = true
		}
	}
	return
}

func (o Offers) MergeSameOffers(candidate Offer) Offer {
	for _, off := range o {
		if off.isSameOffer(candidate) {
			candidate = off.merge(candidate)
		}
	}
	return candidate
}
