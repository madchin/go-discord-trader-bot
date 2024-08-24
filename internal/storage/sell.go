package storage

import (
	"errors"
	"fmt"

	"github.com/madchin/trader-bot/internal/offer"
)

var (
	ErrSellOffersNotExists = errors.New("sell offers not exists")
)

type Sell struct {
	offers offer.Offers
}

func NewSell() Sell {
	return Sell{make(offer.Offers, 50)}
}

func (storage *Sell) List(productName string) (offer.Offers, error) {
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
func (storage *Sell) ListAll() (offer.Offers, error) {
	if len(storage.offers) == 0 {
		return nil, fmt.Errorf("storage list all for %w", ErrSellOffersNotExists)
	}
	return storage.offers, nil
}
func (storage *Sell) Add(sell offer.Sell) error {
	if offers, ok := storage.offers[sell.Vendor()]; ok {
		for i := 0; i < len(offers); i++ {
			if offers[i].IsEqual(sell.Offer()) {
				offers[i] = offers[i].Merge(sell.Offer())
			}
		}
		storage.offers[sell.Vendor()] = append(offers, sell.Offer())
	}
	storage.offers[sell.Vendor()] = append(make([]offer.Offer, 1), sell.Offer())
	return nil
}
func (storage *Sell) Remove(sell offer.Sell) error {
	if offers, ok := storage.offers[sell.Vendor()]; ok {
		for i := 0; i < len(offers); i++ {
			if offers[i].IsEqual(sell.Offer()) {
				offers[i] = offer.Offer{}
				return nil
			}
		}
	}
	return fmt.Errorf("storage remove offer %s for %w", sell.Offer().Product().Name(), ErrSellOffersNotExists)
}
func (storage *Sell) Update(sell offer.Sell) error {
	if offers, ok := storage.offers[sell.Vendor()]; ok {
		for i := 0; i < len(offers); i++ {
			if offers[i].IsEqual(sell.Offer()) {
				offers[i] = sell.Offer()
				return nil
			}
		}
	}
	return fmt.Errorf("storage update offer %s for %w", sell.Offer().Product().Name(), ErrSellOffersNotExists)
}
