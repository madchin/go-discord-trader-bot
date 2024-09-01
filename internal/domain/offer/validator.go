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
	ErrProductPriceTooLarge      = errors.New("product price is too large")
	ErrVendorIdentityIsEmpty     = errors.New("vendor identity is empty")
)

var (
	MinPrice float64 = 0.01
	MaxPrice float64 = 100_000_000
	MinCount int     = 0
)

func validationError(wrapped error) error {
	return fmt.Errorf("validation error: %w", wrapped)
}

func (o Offer) validate() error {
	if err := o.validateCount(); err != nil {
		return err
	}
	if err := o.product.validate(); err != nil {
		return err
	}
	return nil
}

func (p product) validate() error {
	if p.name == "" {
		return validationError(ErrProductNameEmpty)
	}
	if err := p.validatePrice(); err != nil {
		return err
	}
	return nil
}

func (v VendorIdentity) validate() error {
	if v.id == "" {
		return validationError(ErrVendorIdentityIsEmpty)
	}
	return nil
}

func (p product) validatePrice() error {
	if p.price < MinPrice {
		return validationError(ErrProductPriceLessThanZero)
	}
	if p.price > MaxPrice {
		return validationError(ErrProductPriceTooLarge)
	}
	return nil
}

func (o Offer) validateCount() error {
	if o.count <= MinCount {
		return validationError(ErrOfferCountLessOrEqualZero)
	}
	return nil
}
