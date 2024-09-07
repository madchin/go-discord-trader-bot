package storage_offer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/madchin/trader-bot/internal/domain/offer"
	"github.com/madchin/trader-bot/internal/storage"
)

type offerStorage struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *offerStorage {
	return &offerStorage{db}
}

func (offerStorage *offerStorage) Add(ctx context.Context, offer offer.VendorOffer, onAdd offer.OnVendorOfferAddFunc) error {
	tableName := ctx.Value(storage.CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("storage offer add: %w", err)
	}
	if err := offerStorage.add(ctx, tableName, offer, onAdd); err != nil {
		return fmt.Errorf("storage offer add: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) Remove(ctx context.Context, offer offer.VendorOffer) error {
	tableName := ctx.Value(storage.CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("storage offer remove: %w", err)
	}
	if err := offerStorage.remove(ctx, tableName, offer); err != nil {
		return fmt.Errorf("storage offer remove: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) UpdatePrice(ctx context.Context, offer offer.VendorOffer, updatePrice float64, onUpdatePrice offer.OnVendorOfferUpdatePriceFunc) error {
	tableName := ctx.Value(storage.CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("storage offer update price: %w", err)
	}
	if err := offerStorage.updatePrice(ctx, tableName, offer, updatePrice, onUpdatePrice); err != nil {
		return fmt.Errorf("storage offer update price: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) UpdateCount(ctx context.Context, offer offer.VendorOffer, onUpdateCount offer.OnVendorOfferUpdateCountFunc) error {
	tableName := ctx.Value(storage.CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("storage offer update count: %w", err)
	}
	if err := offerStorage.updateCount(ctx, tableName, offer, onUpdateCount); err != nil {
		return fmt.Errorf("storage offer update count: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) ListOffersByName(ctx context.Context, productName string) (offer.VendorOffers, error) {
	tableName := ctx.Value(storage.CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.VendorOffers{}, fmt.Errorf("storage offer list offers: %w", err)
	}
	offers, err := offerStorage.listOffersByName(ctx, tableName, productName)
	if err != nil {
		return offer.VendorOffers{}, fmt.Errorf("storage offer list offers: %w", err)
	}
	return offers, nil
}

func (offerStorage *offerStorage) ListOffersByIdentity(ctx context.Context, vendorIdentity offer.VendorIdentity) (offer.VendorOffers, error) {
	tableName := ctx.Value(storage.CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.VendorOffers{}, fmt.Errorf("storage offer list vendor offers: %w", err)
	}
	offers, err := offerStorage.listOffersByIdentity(ctx, tableName, vendorIdentity)
	if err != nil {
		return offer.VendorOffers{}, fmt.Errorf("storage offer list vendor offers: %w", err)
	}
	return offers, nil
}

func (offerStorage *offerStorage) add(ctx context.Context, dbTable string, offer offer.VendorOffer, onAdd offer.OnVendorOfferAddFunc) error {
	if err := onAdd(offer); err != nil {
		return fmt.Errorf("validation: %w", err)
	}
	query := fmt.Sprintf("INSERT INTO %s (vendorId,price,productName,count) VALUES ($1,$2,$3,$4)", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		offer.VendorIdentity().RawValue(),
		offer.Product.Price(),
		offer.Product.Name(),
		offer.Count(),
	)
	if err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) remove(ctx context.Context, dbTable string, offer offer.VendorOffer) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE price=$1 AND vendorId=$2 AND productName=$3", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		offer.Product.Price(),
		offer.VendorIdentity().RawValue(),
		offer.Product.Name(),
	)
	if err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) updateCount(ctx context.Context, dbTable string, offer offer.VendorOffer, onUpdateCount offer.OnVendorOfferUpdateCountFunc) error {
	if err := onUpdateCount(offer.Count(), offer.VendorIdentity()); err != nil {
		return fmt.Errorf("validation: %w", err)
	}
	query := fmt.Sprintf("UPDATE %s SET count=$1 WHERE vendorId=$2 AND productName=$3 AND price=$4", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		offer.Count(),
		offer.VendorIdentity().RawValue(),
		offer.Product.Name(),
		offer.Product.Price(),
	)
	if err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) updatePrice(ctx context.Context, dbTable string, offer offer.VendorOffer, updatePrice float64, onUpdatePrice offer.OnVendorOfferUpdatePriceFunc) error {
	if err := onUpdatePrice(updatePrice, offer.VendorIdentity()); err != nil {
		return fmt.Errorf("validation: %w", err)
	}
	query := fmt.Sprintf("UPDATE %s SET price=$1 WHERE vendorId=$2 AND productName=$3 AND price=$4", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		updatePrice,
		offer.VendorIdentity().RawValue(),
		offer.Product.Name(),
		offer.Product.Price(),
	)
	if err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}

func (offerStorage *offerStorage) listOffersByName(ctx context.Context, dbTable string, productName string) (offer.VendorOffers, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE productName=$1 ORDER BY price", dbTable)
	rows, err := offerStorage.db.Query(ctx, query, productName)
	if err != nil {
		return nil, fmt.Errorf("query preparing: %w", err)
	}
	defer rows.Close()
	vendorOffers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (offer.VendorOffer, error) {
		var offModel offerModel
		err := row.Scan(&offModel.id, &offModel.vendorId, &offModel.price, &offModel.productName, &offModel.count)
		if err != nil {
			return offer.VendorOffer{}, fmt.Errorf("during scanning row: %w", err)
		}
		return offModel.mapToDomainVendorOffer(), nil
	})
	if err != nil {
		return offer.VendorOffers{}, fmt.Errorf("collecting rows: %w", err)
	}
	return vendorOffers, nil
}

func (offerStorage *offerStorage) listOffersByIdentity(ctx context.Context, dbTable string, vendorIdentity offer.VendorIdentity) (offer.VendorOffers, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE vendorId=$1 ORDER BY price", dbTable)
	rows, err := offerStorage.db.Query(ctx, query, vendorIdentity.RawValue())
	if err != nil {
		return nil, fmt.Errorf("query preparing: %w", err)
	}
	defer rows.Close()
	vendorOffers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (offer.VendorOffer, error) {
		var offModel offerModel
		err := row.Scan(&offModel.id, &offModel.vendorId, &offModel.price, &offModel.productName, &offModel.count)
		if err != nil {
			return offer.VendorOffer{}, fmt.Errorf("during scanning row: %w", err)
		}
		return offModel.mapToDomainVendorOffer(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("collecting rows: %w", err)
	}

	return vendorOffers, nil
}

func (offerStorage *offerStorage) createTable(ctx context.Context, offerTableName string) error {
	itemTableName := ctx.Value(storage.CtxItemTableDescriptorKey).(string)
	queryItemTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (name TEXT PRIMARY KEY)`, itemTableName)
	queryOfferTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY, 
	vendorId TEXT NOT NULL, 
	price NUMERIC(10,2) NOT NULL, 
	productName TEXT NOT NULL, 
	count INTEGER NOT NULL,
	CONSTRAINT fk_name_item FOREIGN KEY (productName) REFERENCES %s (name))`, offerTableName, itemTableName,
	)
	if _, err := offerStorage.db.Exec(ctx, queryItemTable); err != nil {
		return fmt.Errorf("creating table with name %s: %w", itemTableName, err)
	}
	if _, err := offerStorage.db.Exec(ctx, queryOfferTable); err != nil {
		return fmt.Errorf("creating table with name %s: %w", offerTableName, err)
	}
	return nil
}

// 1. listing all offers with productName = "elo"
// 2. Adding offer
//		a) when offer productName = $productName and price = $price and vendor = $vendor, we only increase count += $count
//		else add new record
// 3. Removing offer
//		a) when productName = $productName and price = $price and vendor = $vendor we remove
//		else do nothing
// 4. Update offer
//		a) when productName = $productName and price = $oldPrice and vendor = $vendor we update $price with $newPrice
//		else do nothing
//											 Tables
//											  Offer
//							count 	    vendor 		name 		price
