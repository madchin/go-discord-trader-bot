package storage_item

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/madchin/trader-bot/internal/domain/item"
	"github.com/madchin/trader-bot/internal/storage"
)

type itemPersistence struct {
	db *pgx.Conn
}

func NewDatabase(db *pgx.Conn) *itemPersistence {
	return &itemPersistence{db}
}

func (itemPersistence *itemPersistence) Add(ctx context.Context, item item.Item) error {
	tableName := ctx.Value(storage.CtxItemTableDescriptorKey).(string)
	if err := itemPersistence.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("database item add: %w", err)
	}
	if err := itemPersistence.add(ctx, tableName, item); err != nil {
		return fmt.Errorf("database item add: %w", err)
	}
	return nil
}

func (itemPersistence *itemPersistence) Remove(ctx context.Context, item item.Item) error {
	tableName := ctx.Value(storage.CtxItemTableDescriptorKey).(string)
	if err := itemPersistence.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("database item remove: %w", err)
	}
	if err := itemPersistence.remove(ctx, tableName, item); err != nil {
		return fmt.Errorf("database item remove: %w", err)
	}
	return nil
}

func (itemPersistence *itemPersistence) List(ctx context.Context) (item.Items, error) {
	tableName := ctx.Value(storage.CtxItemTableDescriptorKey).(string)
	if err := itemPersistence.createTable(ctx, tableName); err != nil {
		return nil, fmt.Errorf("database item list: %w", err)
	}
	items, err := itemPersistence.list(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("database item list: %w", err)
	}
	return items, nil
}

func (itemPersistence *itemPersistence) ListByName(ctx context.Context, incomingItem item.Item) (item.Item, error) {
	tableName := ctx.Value(storage.CtxItemTableDescriptorKey).(string)
	if err := itemPersistence.createTable(ctx, tableName); err != nil {
		return item.Item{}, fmt.Errorf("database item list by name: %w", err)
	}
	domainItem, err := itemPersistence.listByName(ctx, tableName, incomingItem)
	if err != nil {
		return item.Item{}, fmt.Errorf("database item list by name: %w", err)
	}
	return domainItem, nil
}

func (itemPersistence *itemPersistence) add(ctx context.Context, tableName string, item item.Item) error {
	query := fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1)`, tableName)
	if _, err := itemPersistence.db.Exec(ctx, query, item.Name()); err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (itemPersistence *itemPersistence) list(ctx context.Context, tableName string) (item.Items, error) {
	query := fmt.Sprintf("SELECT name FROM %s", tableName)
	rows, err := itemPersistence.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query preparing: %w", err)
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (item.Item, error) {
		var itemModel itemModel
		if err := row.Scan(&itemModel.name); err != nil {
			return item.Item{}, fmt.Errorf("during row scanning: %w", err)
		}
		return itemModel.mapToDomainItem(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("during rows collecting: %w", err)
	}

	return items, nil
}

func (itemPersistence *itemPersistence) remove(ctx context.Context, tableName string, item item.Item) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE name = $1", tableName)
	if _, err := itemPersistence.db.Exec(ctx, query, item.Name()); err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (itemPersistence *itemPersistence) listByName(ctx context.Context, tableName string, incomingItem item.Item) (item.Item, error) {
	query := fmt.Sprintf("SELECT name FROM %s WHERE name = $1", tableName)
	row := itemPersistence.db.QueryRow(ctx, query, incomingItem.Name())
	var itemModel itemModel
	if err := row.Scan(&itemModel.name); err != nil && err != pgx.ErrNoRows {
		return item.Item{}, fmt.Errorf("during scanning: %w", err)
	}
	return itemModel.mapToDomainItem(), nil
}

func (itemPersistence *itemPersistence) createTable(ctx context.Context, name string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (name TEXT PRIMARY KEY)`, name)
	if _, err := itemPersistence.db.Exec(ctx, query); err != nil {
		return fmt.Errorf("creating table with name %s: %w", name, err)
	}
	return nil
}
