package storage

import (
	"errors"
	"fmt"
	"log"

	"github.com/madchin/trader-bot/internal/offer"
)

var (
	ErrOffersNotExists = errors.New("off offers not exists")
)

type OfferStorage struct {
	offers offer.Offers
}

func New() *OfferStorage {
	return &OfferStorage{make(offer.Offers, 50)}
}

// FIXME mitigate from inmemory storage
func (storage *OfferStorage) List(productName string) (offer.Offers, error) {
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
func (storage *OfferStorage) ListAll() (offer.Offers, error) {
	if len(storage.offers) == 0 {
		return nil, fmt.Errorf("storage list all for %w", ErrOffersNotExists)
	}
	return storage.offers, nil
}
func (storage *OfferStorage) Add(off offer.Offer) error {
	if offers, ok := storage.offers[off.Vendor()]; ok {
		offer, mergedIndexes, err := off.Add(offers)
		if err != nil {
			return fmt.Errorf("storage offer add %s for vendor %v %w", off.Product().Name(), off.Vendor(), err)
		}
		mergeCount := storage.removeMergedOffers(offer, mergedIndexes)
		log.Printf("merged count %d", mergeCount)
		storage.offers[off.Vendor()] = append(offers, offer)
	}
	storage.offers[off.Vendor()] = append(make([]offer.Offer, 1), off)
	log.Printf("adding %v", off)
	//data, _ := storage.ListAll()
	//log.Printf("actual storage list: %v", data)
	return nil
}
func (storage *OfferStorage) Remove(off offer.Offer) error {
	if offers, ok := storage.offers[off.Vendor()]; ok {
		removeIndice, err := off.Remove(offers)
		if err != nil {
			return fmt.Errorf("storage offer remove %s for vendor %v %w", off.Product().Name(), off.Vendor(), err)
		}
		offers[removeIndice] = offer.Offer{}
		return nil
	}
	return fmt.Errorf("storage remove offer %s for vendor %v %w", off.Product().Name(), off.Vendor(), ErrOffersNotExists)
}
func (storage *OfferStorage) Update(off offer.Offer) error {
	if offers, ok := storage.offers[off.Vendor()]; ok {
		offer, updateIndice, err := off.Update(offers)
		if err != nil {
			return fmt.Errorf("storage update offer %s for vendor %v %w", off.Product().Name(), off.Vendor(), err)
		}
		offers[updateIndice] = offer
	}
	return fmt.Errorf("storage update offer %s for vendor %v %w", off.Product().Name(), off.Vendor(), ErrOffersNotExists)
}

func (storage *OfferStorage) removeMergedOffers(off offer.Offer, mergedIndexes []int) (merged int) {
	if offers, ok := storage.offers[off.Vendor()]; ok {
		for _, idx := range mergedIndexes {
			offers[idx] = offer.Offer{}
			merged++
		}
	}
	return
}
