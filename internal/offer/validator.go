package offer

import (
	"errors"
	"fmt"
)

var (
	ErrProductNameEmpty          = errors.New("product name is empty")
	ErrProductPriceLessThanZero  = errors.New("product price is less than 0")
	ErrVendorIsEmpty             = errors.New("offer vendor is empty")
	ErrOfferCountLessOrEqualZero = errors.New("offer count is less than 0")
)

func validationError(wrapped error) error {
	return fmt.Errorf("validation error: %w", wrapped)
}

func (o Offer) validate() error {
	if o.count <= 0 {
		return validationError(ErrOfferCountLessOrEqualZero)
	}
	return nil
}

func (p product) validate() error {
	if p.name == "" {
		return validationError(ErrProductNameEmpty)
	}
	if p.price < 0 {
		return validationError(ErrProductPriceLessThanZero)
	}
	return nil
}

func (v Vendor) validate() error {
	if v.name == "" {
		return validationError(ErrVendorIsEmpty)
	}
	return nil
}
