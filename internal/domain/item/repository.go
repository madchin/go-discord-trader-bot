package item

import "context"

type Repository interface {
	Add(ctx context.Context, item Item) error
	List(ctx context.Context) (Items, error)
	Remove(ctx context.Context, item Item) error
}
