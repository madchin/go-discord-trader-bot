package item

import (
	"errors"
	"slices"
)

var (
	ErrNameIsEmpty   = errors.New("item name is empty")
	ErrItemsAreEmpty = errors.New("items are empty")
)

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

func (i Item) IsZero() bool {
	return i.name == ""
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

func (items Items) Add(item Item) Items {
	items = append(items, item)
	return items
}

func (items Items) AreEmpty() bool {
	return len(items) == 0
}

func (items Items) Delete(item Item) (Items, error) {
	if items.AreEmpty() {
		return Items{}, ErrItemsAreEmpty
	}
	idx := -1
	for i := 0; i < len(items); i++ {
		if items[i] == item {
			idx = i
			break
		}
	}
	if idx == 0 {
		return items[1:], nil
	}
	if idx == len(items)-1 {
		return items[: idx-1 : idx-1], nil
	}
	return slices.Delete(items, idx, idx+1), nil
}
