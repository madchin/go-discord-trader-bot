package offer

import "context"

type OnOfferAddFunc func(Offer) error
type OnOfferUpdatePriceFunc func(float64, VendorIdentity) error
type OnOfferUpdateCountFunc func(int, VendorIdentity) error

type Repository interface {
	ListOffers(ctx context.Context, productName string) (Offers, error)
	ListVendorOffers(ctx context.Context, vendorIdentity VendorIdentity) (Offers, error)
	Remove(ctx context.Context, offer Offer) error
	Add(ctx context.Context, offer Offer, onAdd OnOfferAddFunc) error
	UpdateCount(ctx context.Context, offer Offer, count int, onUpdate OnOfferUpdateCountFunc) error
	UpdatePrice(ctx context.Context, offer Offer, price float64, onUpdate OnOfferUpdatePriceFunc) error
}
