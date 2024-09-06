package storage_item

import "github.com/madchin/trader-bot/internal/domain/item"

type itemModel struct {
	name string
}

func (i itemModel) mapToDomainItem() item.Item {
	return item.New(i.name)
}
