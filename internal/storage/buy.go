package storage

import (
	"errors"
	"fmt"

	"github.com/madchin/trader-bot/internal/offer"
)

var (
	ErrBuyOffersNotExists = errors.New("sell offers not exists")
)

type Buy struct {
	offers offer.Offers
}

func NewBuy() Buy {
	return Buy{make(offer.Offers, 50)}
}

func (storage *Buy) List(productName string) (offer.Offers, error) {
	specifiedOffers := make(offer.Offers, len(storage.offers))
	for vendor, offers := range storage.offers {
		for i := 0; i < len(offers); i++ {
			if offers[i].Product().Name() == productName {
				specifiedOffers[vendor] = offers
			}
		}
	}
	return specifiedOffers, nil
}

func (storage *Buy) ListAll() (offer.Offers, error) {
	if len(storage.offers) == 0 {
		return nil, fmt.Errorf("storage list all for %w", ErrBuyOffersNotExists)
	}
	return storage.offers, nil
}

func (storage *Buy) Add(buy offer.Buy) error {
	if offers, ok := storage.offers[buy.Vendor()]; ok {
		for i := 0; i < len(offers); i++ {
			if offers[i].IsEqual(buy.Offer()) {
				offers[i] = offers[i].Merge(buy.Offer())
				return nil
			}
		}
		storage.offers[buy.Vendor()] = append(offers, buy.Offer())
		return nil
	}
	storage.offers[buy.Vendor()] = append(make([]offer.Offer, 1), buy.Offer())
	return nil
}

func (storage *Buy) Remove(buy offer.Sell) error {
	if offers, ok := storage.offers[buy.Vendor()]; ok {
		for i := 0; i < len(offers); i++ {
			if offers[i].IsEqual(buy.Offer()) {
				offers[i] = offer.Offer{}
				return nil
			}
		}
	}
	return fmt.Errorf("storage remove offer %s for %w", buy.Offer().Product().Name(), ErrBuyOffersNotExists)
}

func (storage *Buy) Update(buy offer.Sell) error {
	if offers, ok := storage.offers[buy.Vendor()]; ok {
		for i := 0; i < len(offers); i++ {
			if offers[i].IsEqual(buy.Offer()) {
				offers[i] = buy.Offer()
				return nil
			}
		}
	}
	return fmt.Errorf("storage update offer %s for %w", buy.Offer().Product().Name(), ErrBuyOffersNotExists)
}
