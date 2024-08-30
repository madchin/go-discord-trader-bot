package offer

import (
	"errors"
)

var (
	ErrWrongActionType  = errors.New("unable to perform this action for this type of offer")
	ErrNotExistingOffer = errors.New("offer not exists")
)

type Offer struct {
	vendor  vendor
	product product
	count   int
}

type Offers []Offer

func NewOffer(vendor string, product product, count int) Offer {
	return Offer{newVendor(vendor), product, count}
}

func (o Offer) Product() product {
	return o.product
}

func (o Offer) Vendor() vendor {
	return o.vendor
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

func OnOfferAdd(o Offer) error {
	return o.validate()
}

func OnOfferCountUpdate(o Offer) error {
	return o.validateCount()
}

func OnOfferPriceUpdate(o Offer) error {
	return o.product.validatePrice()
}
