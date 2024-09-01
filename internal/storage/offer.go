package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/madchin/trader-bot/internal/domain/offer"
)

type offerStorage struct {
	db *pgx.Conn
}

func NewOffer(db *pgx.Conn) *offerStorage {
	return &offerStorage{db}
}

func (offerStorage *offerStorage) Add(ctx context.Context, offer offer.Offer, onAdd offer.OnOfferAddFunc) error {
	tableName := ctx.Value(CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return err
	}
	return offerStorage.add(ctx, tableName, offer, onAdd)
}

func (offerStorage *offerStorage) Remove(ctx context.Context, offer offer.Offer) error {
	tableName := ctx.Value(CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return err
	}
	return offerStorage.remove(ctx, tableName, offer)
}

func (offerStorage *offerStorage) UpdatePrice(ctx context.Context, offer offer.Offer, price float64, onUpdatePrice offer.OnOfferUpdatePriceFunc) error {
	tableName := ctx.Value(CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return err
	}
	return offerStorage.updatePrice(ctx, tableName, offer, price, onUpdatePrice)
}

func (offerStorage *offerStorage) UpdateCount(ctx context.Context, offer offer.Offer, count int, onUpdateCount offer.OnOfferUpdateCountFunc) error {
	tableName := ctx.Value(CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return err
	}
	return offerStorage.updateCount(ctx, tableName, offer, count, onUpdateCount)
}

func (offerStorage *offerStorage) ListOffers(ctx context.Context, productName string) (offer.Offers, error) {
	tableName := ctx.Value(CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.Offers{}, err
	}
	return offerStorage.listOffers(ctx, tableName, productName)
}

func (offerStorage *offerStorage) ListVendorOffers(ctx context.Context, vendorIdentity offer.VendorIdentity) (offer.Offers, error) {
	tableName := ctx.Value(CtxBuySellDbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.Offers{}, err
	}
	return offerStorage.listVendorOffers(ctx, tableName, vendorIdentity)
}

func (offerStorage *offerStorage) add(ctx context.Context, dbTable string, offer offer.Offer, onAdd offer.OnOfferAddFunc) error {
	if err := onAdd(offer); err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (vendorId,price,productName,count) VALUES ($1,$2,$3,$4)", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		offer.VendorIdentity().RawValue(),
		offer.Product().Price(),
		offer.Product().Name(),
		offer.Count(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (offerStorage *offerStorage) remove(ctx context.Context, dbTable string, offer offer.Offer) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE price=$1 AND vendorId=$2 AND productName=$3", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		offer.Product().Price(),
		offer.VendorIdentity().RawValue(),
		offer.Product().Name(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (offerStorage *offerStorage) updateCount(ctx context.Context, dbTable string, offer offer.Offer, count int, onUpdateCount offer.OnOfferUpdateCountFunc) error {
	if err := onUpdateCount(count, offer.VendorIdentity()); err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET count=$1 WHERE vendorId=$2 AND productName=$3 AND price=$4", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		count,
		offer.VendorIdentity().RawValue(),
		offer.Product().Name(),
		offer.Product().Price(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (offerStorage *offerStorage) updatePrice(ctx context.Context, dbTable string, offer offer.Offer, price float64, onUpdatePrice offer.OnOfferUpdatePriceFunc) error {
	if err := onUpdatePrice(price, offer.VendorIdentity()); err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET price=$1 WHERE vendorId=$2 AND productName=$3 AND price=$4", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		price,
		offer.VendorIdentity().RawValue(),
		offer.Product().Name(),
		offer.Product().Price(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (offerStorage *offerStorage) listOffers(ctx context.Context, dbTable string, productName string) (offer.Offers, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE productName=$1 ORDER BY price", dbTable)
	rows, err := offerStorage.db.Query(ctx, query, productName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offerModels, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (offerModel, error) {
		var offModel offerModel
		err := row.Scan(&offModel.id, &offModel.vendorId, &offModel.price, &offModel.productName, &offModel.count)
		if err != nil {
			return offerModel{}, fmt.Errorf("during scanning row: %w", err)
		}
		return offModel, nil
	})
	if err != nil {
		return offer.Offers{}, fmt.Errorf("collecting rows %w", err)
	}
	return mapStorageOffersToDomainOffers(offerModels), nil
}

func (offerStorage *offerStorage) listVendorOffers(ctx context.Context, dbTable string, vendorIdentity offer.VendorIdentity) (offer.Offers, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE vendorId=$1 ORDER BY price", dbTable)
	rows, err := offerStorage.db.Query(ctx, query, vendorIdentity.RawValue())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offerModels, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (offerModel, error) {
		var offModel offerModel
		err := row.Scan(&offModel.id, &offModel.vendorId, &offModel.price, &offModel.productName, &offModel.count)
		if err != nil {
			return offerModel{}, fmt.Errorf("during scanning row: %w", err)
		}
		return offModel, nil
	})
	if err != nil {
		return nil, fmt.Errorf("list vendor offers collecting rows: %w", err)
	}

	return mapStorageOffersToDomainOffers(offerModels), nil
}

func (offerStorage *offerStorage) createTable(ctx context.Context, name string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY, 
	vendorId TEXT NOT NULL, 
	price NUMERIC(10,2) NOT NULL, 
	productName TEXT NOT NULL, 
	count INTEGER NOT NULL)`, name,
	)
	if _, err := offerStorage.db.Exec(ctx, query); err != nil {
		return fmt.Errorf("error occured during creating table with name %s err: %v", name, err)
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
