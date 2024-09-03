package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/madchin/trader-bot/internal/domain/item"
)

type ItemStorage struct {
	db *pgx.Conn
}

func NewItemStorage(db *pgx.Conn) *ItemStorage {
	return &ItemStorage{db}
}

func (itemStorage *ItemStorage) Add(ctx context.Context, item item.Item) error {
	tableName := ctx.Value(CtxItemTableDescriptorKey).(string)
	if err := itemStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("storage item add: %w", err)
	}
	if err := itemStorage.add(ctx, tableName, item); err != nil {
		return fmt.Errorf("storage item add: %w", err)
	}
	return nil
}

func (itemStorage *ItemStorage) Remove(ctx context.Context, item item.Item) error {
	tableName := ctx.Value(CtxItemTableDescriptorKey).(string)
	if err := itemStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("storage item remove: %w", err)
	}
	if err := itemStorage.remove(ctx, tableName, item); err != nil {
		return fmt.Errorf("storage item remove: %w", err)
	}
	return nil
}

func (itemStorage *ItemStorage) List(ctx context.Context) (item.Items, error) {
	tableName := ctx.Value(CtxItemTableDescriptorKey).(string)
	if err := itemStorage.createTable(ctx, tableName); err != nil {
		return nil, fmt.Errorf("storage item list: %w", err)
	}
	items, err := itemStorage.list(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("storage item list: %w", err)
	}
	return items, nil
}

func (itemStorage *ItemStorage) add(ctx context.Context, tableName string, item item.Item) error {
	query := fmt.Sprintf(`INSERT INTO %s name VALUES ($1)`, tableName)
	if _, err := itemStorage.db.Exec(ctx, query, item.Name()); err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (itemStorage *ItemStorage) list(ctx context.Context, tableName string) (item.Items, error) {
	query := fmt.Sprintf("SELECT DISTINCT id,name FROM %s", tableName)
	rows, err := itemStorage.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query preparing: %w", err)
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (item.Item, error) {
		var itemModel itemModel
		if err := row.Scan(&itemModel.id, &itemModel.name); err != nil {
			return item.Item{}, fmt.Errorf("during row scanning: %w", err)
		}
		return itemModel.mapToDomainItem(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("during rows collecting: %w", err)
	}

	return items, nil
}

func (itemStorage *ItemStorage) remove(ctx context.Context, tableName string, item item.Item) error {
	query := fmt.Sprintf("REMOVE FROM %s WHERE $1", tableName)
	if _, err := itemStorage.db.Exec(ctx, query, item.Name()); err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (itemStorage *ItemStorage) createTable(ctx context.Context, name string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (id SERIAL PRIMARY KEY, name TEXT NOT NULL)`, name)
	if _, err := itemStorage.db.Exec(ctx, query); err != nil {
		return fmt.Errorf("creating table with name %s: %w", name, err)
	}
	return nil
}
