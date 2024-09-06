package storage_item

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/patrickmn/go-cache"
)

type ItemStorage struct {
	persistence *itemPersistence
	cache       *itemCache
}

func New(dbConn *pgx.Conn, cache *cache.Cache) *ItemStorage {
	return &ItemStorage{
		persistence: &itemPersistence{dbConn},
		cache:       &itemCache{cache},
	}
}

func (itemStorage *ItemStorage) Add(ctx context.Context, incomingItem item.Item) error {
	if err := itemStorage.persistence.Add(ctx, incomingItem); err != nil {
		return fmt.Errorf("item storage: %w", err)
	}
	log.Printf("adding to persistence: aa %v", incomingItem)
	if _, err := itemStorage.cache.List(ctx); err == errCacheMiss {
		log.Printf("cache miss: %v", incomingItem)
		if err := itemStorage.updateCacheWithItemsFromPersistenceLayer(ctx); err != nil {
			log.Printf("update cache: %v", incomingItem)
			return fmt.Errorf("item storage: %w", err)
		}
		return nil
	}
	if err := itemStorage.cache.Add(ctx, incomingItem); err != nil {
		_ = itemStorage.persistence.Remove(ctx, incomingItem)
		return fmt.Errorf("item storage: %w", err)
	}
	return nil
}

func (itemStorage *ItemStorage) List(ctx context.Context) (item.Items, error) {
	items, err := itemStorage.cache.List(ctx)
	if err == nil {
		return items, nil
	}
	items, err = itemStorage.persistence.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("item storage: %w", err)
	}
	itemStorage.cache.override(ctx, items)

	return items, nil
}

func (itemStorage *ItemStorage) ListByName(ctx context.Context, incomingItem item.Item) (item.Item, error) {
	domainItem, err := itemStorage.cache.ListByName(ctx, incomingItem)
	if err == nil {
		return domainItem, nil
	}
	domainItem, err = itemStorage.persistence.ListByName(ctx, incomingItem)
	if err != nil {
		return item.Item{}, fmt.Errorf("item storage: %w", err)
	}
	if err := itemStorage.updateCacheWithItemsFromPersistenceLayer(ctx); err != nil {
		return item.Item{}, fmt.Errorf("item storage: %w", err)
	}
	return domainItem, nil
}

func (itemStorage *ItemStorage) Remove(ctx context.Context, incomingItem item.Item) error {
	if err := itemStorage.persistence.Remove(ctx, incomingItem); err != nil {
		return fmt.Errorf("item storage: %w", err)
	}
	if err := itemStorage.cache.Remove(ctx, incomingItem); err != nil {
		_ = itemStorage.persistence.Add(ctx, incomingItem)
		return fmt.Errorf("item storage: %w", err)
	}
	return nil
}

func (itemStorage *ItemStorage) updateCacheWithItemsFromPersistenceLayer(ctx context.Context) error {
	items, err := itemStorage.persistence.List(ctx)
	if err != nil {
		return fmt.Errorf("update cache: %w", err)
	}
	itemStorage.cache.override(ctx, items)
	return nil
}
