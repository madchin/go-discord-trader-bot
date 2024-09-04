package storage

import (
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type itemModel struct {
	name string
}

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

func (i itemModel) mapToDomainItem() item.Item {
	return item.New(i.name)
}
