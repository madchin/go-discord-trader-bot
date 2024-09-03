package item

import "errors"

var ErrNameIsEmpty = errors.New("item name is empty")

type Item struct {
	name string
}

type Items []Item

func New(name string) Item {
	return Item{name}
}

func (i Item) Validate() error {
	if i.name == "" {
		return ErrNameIsEmpty
	}
	return nil
}

func (i Item) Name() string {
	return i.name
}

func (items Items) String() string {
	readable := items[0].name
	for i := 1; i < len(items); i++ {
		readable += ", " + items[i].name
	}
	return readable
}

func (items Items) Contains(item Item) bool {
	for i := 0; i < len(items); i++ {
		if items[i] == item {
			return true
		}
	}
	return false
}
