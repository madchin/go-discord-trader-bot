package storage

import "github.com/madchin/trader-bot/internal/domain/offer"

type offerModel struct {
	id          int
	vendor      string
	price       float64
	productName string
	count       int
}

func (o offerModel) mapToDomainOffer() offer.Offer {
	product := offer.NewProduct(o.productName, o.price)
	return offer.NewOffer(o.vendor, product, o.count)
}
