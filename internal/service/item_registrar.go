package service

import (
	"context"
	"fmt"

	"github.com/madchin/trader-bot/internal/domain/item"
)

type itemRegistrar struct {
	itemStorage item.Repository
}

func (itemRegistrar *itemRegistrar) Add(ctx context.Context, item item.Item) error {
	if err := itemRegistrar.itemStorage.Add(ctx, item); err != nil {
		return fmt.Errorf("item registrar add service: %w", err)
	}
	return nil
}

func (itemRegistrar *itemRegistrar) Remove(ctx context.Context, item item.Item) error {
	if err := itemRegistrar.itemStorage.Remove(ctx, item); err != nil {
		return fmt.Errorf("item registrar remove service: %w", err)
	}
	return nil
}

func (itemRegistrar *itemRegistrar) List(ctx context.Context) (item.Items, error) {
	items, err := itemRegistrar.itemStorage.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("item registrar list service: %w", err)
	}
	return items, nil
}
