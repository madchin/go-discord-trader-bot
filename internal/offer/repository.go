package offer

type BuyRepository interface {
	List(productName string) (Offers, error)
	ListAll() (Offers, error)
	Add(offer Buy) error
	Remove(offer Buy) error
	Update(offer Buy) error
}

type SellRepository interface {
	List(productName string) (Offers, error)
	ListAll() (Offers, error)
	Add(offer Sell) error
	Remove(offer Sell) error
	Update(offer Sell) error
}
