package storage_offer

import (
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type offerModel struct {
	id          int
	vendorId    string
	price       float64
	productName string
	count       int
}

func (o offerModel) mapToDomainVendorOffer() offer.VendorOffer {
	product := offer.NewProduct(o.productName, o.price)
	off := offer.NewOffer(product, o.count)
	identity := offer.NewVendorIdentity(o.vendorId)
	return offer.NewVendorOffer(identity, off)
}
