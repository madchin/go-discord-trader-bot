package offer

import "errors"

var ErrOfferTypeNotDefined = errors.New("offer type is neither Buy or Sell")

type offerType struct {
	t string
}

var (
	Buy  = offerType{"buy"}
	Sell = offerType{"sell"}
)

func OfferTypeFromCandidate(candidate string) (offerType, error) {
	if Buy.t == candidate {
		return Buy, nil
	}
	if Sell.t == candidate {
		return Sell, nil
	}
	return offerType{}, ErrOfferTypeNotDefined
}
