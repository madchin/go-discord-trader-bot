package offer

import (
	"errors"
	"fmt"
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

type Offers map[vendor][]Offer

func NewOffer(offerType offerType, action action, vendor string, product product, count int) (Offer, error) {
	vend, err := newVendor(vendor)
	if err != nil {
		return Offer{}, err
	}
	off := Offer{offerType, action, vend, product, count}
	if err := off.validate(); err != nil {
		return Offer{}, err
	}
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

func (offer Offer) Add(actualVendorOffers []Offer) (_ Offer, mergedIndices []int, _ error) {
	if offer.action != Add {
		return Offer{}, []int{}, fmt.Errorf("offer add %w", ErrWrongActionType)
	}
	for idx, candidate := range actualVendorOffers {
		if offer.isSameOffer(candidate) {
			offer = offer.merge(candidate)
			mergedIndices = append(mergedIndices, idx)
		}
	}
	return offer, mergedIndices, nil
}

func (offer Offer) Remove(actualVendorOffers []Offer) (removedIndice int, err error) {
	if offer.action != Remove {
		return -1, fmt.Errorf("offer remove %w", ErrWrongActionType)
	}
	for idx, existingOffer := range actualVendorOffers {
		if offer.isSameOffer(existingOffer) {
			return idx, nil
		}
	}
	return -1, ErrNotExistingOffer
}

func (offer Offer) Update(actualVendorOffers []Offer) (_ Offer, updateIndice int, _ error) {
	if offer.action != Update {
		return Offer{}, -1, fmt.Errorf("offer update %w", ErrWrongActionType)
	}
	for indice, existingOffer := range actualVendorOffers {
		if offer.isSameOffer(existingOffer) {
			o := offer.merge(existingOffer)
			return o, indice, nil
		}
	}
	return Offer{}, -1, fmt.Errorf("offer update %w", ErrNotExistingOffer)
}

func (o Offer) isSameOffer(off Offer) bool {
	return o.product.name == off.product.name && o.product.price == off.product.price
}

func (o Offer) merge(off Offer) Offer {
	o.count += off.count
	return o
}
