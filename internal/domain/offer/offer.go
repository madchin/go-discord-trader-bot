package offer

import (
	"errors"
)

var (
	ErrWrongActionType  = errors.New("unable to perform this action for this type of offer")
	ErrNotExistingOffer = errors.New("offer not exists")
)

type Offer struct {
	product product
	count   int
}

type Offers []Offer

func NewOffer(product product, count int) Offer {
	return Offer{product, count}
}

func (o Offer) Product() product {
	return o.product
}

func (o Offer) Count() int {
	return o.count
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
