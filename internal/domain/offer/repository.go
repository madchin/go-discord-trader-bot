package offer

import "context"

type OnVendorOfferAddFunc func(VendorOffer) error
type OnVendorOfferUpdatePriceFunc func(float64, VendorIdentity) error
type OnVendorOfferUpdateCountFunc func(int, VendorIdentity) error

type Repository interface {
	ListOffersByName(ctx context.Context, productName string) (VendorOffers, error)
	ListOffersByIdentity(ctx context.Context, vendorIdentity VendorIdentity) (VendorOffers, error)
	Remove(ctx context.Context, offer VendorOffer) error
	Add(ctx context.Context, offer VendorOffer, onAdd OnVendorOfferAddFunc) error
	UpdateCount(ctx context.Context, offer VendorOffer, onUpdate OnVendorOfferUpdateCountFunc) error
	UpdatePrice(ctx context.Context, offer VendorOffer, onUpdate OnVendorOfferUpdatePriceFunc) error
}
