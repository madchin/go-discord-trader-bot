package storage

import "github.com/madchin/trader-bot/internal/domain/offer"

type offerModel struct {
	id          int
	vendorId    string
	price       float64
	productName string
	count       int
}

func (o offerModel) mapToDomainOffer() offer.Offer {
	product := offer.NewProduct(o.productName, o.price)
	identity := offer.NewVendorIdentity(o.vendorId)
	return offer.NewOffer(identity, product, o.count)
}

func mapStorageOffersToDomainOffers(storageOffers []offerModel) offer.Offers {
	offers := make(offer.Offers, len(storageOffers))
	for i := 0; i < len(storageOffers); i++ {
		offers[i] = storageOffers[i].mapToDomainOffer()
	}
	return offers
}
