package offer

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrWrongActionType  = errors.New("unable to perform this action for this type of offer")
	ErrNotExistingOffer = errors.New("offer not exists")
)

type Offer struct {
	vendorIdentity VendorIdentity
	product        product
	count          int
}

type Offers []Offer

func NewOffer(vendorIdentity VendorIdentity, product product, count int) Offer {
	return Offer{vendorIdentity, product, count}
}

func (o Offer) Product() product {
	return o.product
}

func (o Offer) Count() int {
	return o.count
}

func (o Offer) VendorIdentity() VendorIdentity {
	return o.vendorIdentity
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

func (o Offers) NotExists() bool {
	return len(o) == 0
}

func (o Offers) ToReadableMessage() string {
	offers := make([]string, len(o))
	for i := 0; i < len(o); i++ {
		price := fmt.Sprintf("%.2f", o[i].product.price)
		offers[i] = fmt.Sprintf("Product: %s, Each Price: %s, Count: %d, Vendor: %s", o[i].product.name, price, o[i].count, "<@!"+o[i].vendorIdentity.id+">")
	}
	return strings.Join(offers, ",\n")
}

func (o Offers) VendorIdentities() VendorIdentities {
	vendorIdentities := make(VendorIdentities, len(o))
	for i := 0; i < len(o); i++ {
		vendorIdentities[i] = o[i].vendorIdentity
	}
	return vendorIdentities
}

func OnOfferAdd(o Offer) error {
	if err := o.validate(); err != nil {
		return err
	}
	if err := o.vendorIdentity.validate(); err != nil {
		return err
	}
	return nil
}

func OnOfferCountUpdate(count int, v VendorIdentity) error {
	tmpOffer := Offer{count: count}
	if err := tmpOffer.validateCount(); err != nil {
		return err
	}
	if err := v.validate(); err != nil {
		return err
	}
	return nil
}

func OnOfferPriceUpdate(price float64, v VendorIdentity) error {
	tmpProd := product{price: price}
	if err := tmpProd.validatePrice(); err != nil {
		return err
	}
	if err := v.validate(); err != nil {
		return err
	}
	return nil
}
