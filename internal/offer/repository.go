package offer

type Repository interface {
	List(productName string) (Offers, error)
	ListAll() (Offers, error)
	Add(offer Offer) error
	Remove(offer Offer) error
	Update(offer Offer) error
}
