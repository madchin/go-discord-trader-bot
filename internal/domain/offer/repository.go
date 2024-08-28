package offer

import "context"

type Repository interface {
	ListOffers(ctx context.Context, productName string) (Offers, error)
	ListVendorOffers(ctx context.Context, vendorName string) (Offers, error)
	Add(ctx context.Context, offer Offer) error
	Remove(ctx context.Context, offer Offer) error
	Update(ctx context.Context, oldOffer Offer, updateOffer Offer) error
}
