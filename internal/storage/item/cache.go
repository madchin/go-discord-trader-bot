package storage_item

import (
	"context"
	"errors"
	"fmt"

	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/storage"
	"github.com/patrickmn/go-cache"
)

var (
	errCacheMiss   = errors.New("cache missed")
	errTypeParse   = errors.New("unable to parse to item type")
	errItemMissing = errors.New("item is missing")
)

type itemCache struct {
	cache *cache.Cache
}

func (itemCache *itemCache) Add(ctx context.Context, incomingItem item.Item) error {
	return itemCache.add(ctx, incomingItem)
}

func (itemCache *itemCache) Remove(ctx context.Context, incomingItem item.Item) error {
	return itemCache.remove(ctx, incomingItem)
}

func (itemCache *itemCache) ListByName(ctx context.Context, incomingItem item.Item) (item.Item, error) {
	return itemCache.listByName(ctx, incomingItem)
}

func (itemCache *itemCache) List(ctx context.Context) (item.Items, error) {
	return itemCache.list(ctx)
}

func (itemCache *itemCache) listByName(ctx context.Context, incomingItem item.Item) (item.Item, error) {
	items, err := itemCache.list(ctx)
	if err != nil {
		return item.Item{}, errCacheMiss
	}
	if !items.Contains(incomingItem) {
		return item.Item{}, errItemMissing
	}
	return item.New(incomingItem.Name()), nil
}

func (itemCache *itemCache) remove(ctx context.Context, incomingItem item.Item) error {
	items, err := itemCache.list(ctx)
	if err != nil && err != errCacheMiss {
		return fmt.Errorf("retrieving cache: %w", err)
	}
	if items.AreEmpty() {
		itemCache.cache.Delete(cacheKey(ctx))
		return nil
	}
	items, err = items.Delete(incomingItem)
	if err != nil {
		return fmt.Errorf("retrieving cache: %w", err)
	}
	itemCache.cache.Set(cacheKey(ctx), items, cache.DefaultExpiration)
	return nil
}

func (itemCache *itemCache) add(ctx context.Context, incomingItem item.Item) error {
	items, err := itemCache.list(ctx)
	if err != nil && err != errCacheMiss {
		return fmt.Errorf("retrieving cache: %w", err)
	}
	items = items.Add(incomingItem)
	itemCache.cache.Set(cacheKey(ctx), items, cache.DefaultExpiration)
	return nil
}

func (itemCache *itemCache) override(ctx context.Context, items item.Items) {
	itemCache.cache.Set(cacheKey(ctx), items, cache.DefaultExpiration)
}

func (itemCache *itemCache) list(ctx context.Context) (item.Items, error) {
	cacheItems, ok := itemCache.cache.Get(cacheKey(ctx))
	if !ok {
		return nil, errCacheMiss
	}
	items, ok := cacheItems.(item.Items)
	if !ok {
		return nil, errTypeParse
	}
	return items, nil
}

func cacheKey(ctx context.Context) string {
	val, ok := ctx.Value(storage.CtxItemTableDescriptorKey).(string)
	if !ok {
		return ""
	}
	return val
}
