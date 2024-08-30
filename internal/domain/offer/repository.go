package offer

import "context"

type Repository interface {
	ListOffers(ctx context.Context, productName string) (Offers, error)
	ListVendorOffers(ctx context.Context, vendorName string) (Offers, error)
	Remove(ctx context.Context, offer Offer) error
	Add(ctx context.Context, offer Offer, onAdd func(Offer) error) error
	UpdateCount(ctx context.Context, offer Offer, count int, onUpdate func(Offer) error) error
	UpdatePrice(ctx context.Context, offer Offer, price float64, onUpdate func(Offer) error) error
}
