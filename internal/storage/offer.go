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

func New(db *pgx.Conn) *offerStorage {
	return &offerStorage{db}
}

func (offerStorage *offerStorage) Add(ctx context.Context, offer offer.Offer) error {
	tableName := ctx.Value(DbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("database error add for product %s vendor %s price %f count %d: err %w",
			offer.Product().Name(),
			offer.Vendor().Name(),
			offer.Product().Price(),
			offer.Count(),
			err,
		)
	}
	return offerStorage.add(ctx, tableName, offer)
}

func (offerStorage *offerStorage) Remove(ctx context.Context, offer offer.Offer) error {
	tableName := ctx.Value(DbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("database error remove for product %s vendor %s price %f count %d: err %w",
			offer.Product().Name(),
			offer.Vendor().Name(),
			offer.Product().Price(),
			offer.Count(),
			err,
		)
	}
	return offerStorage.remove(ctx, tableName, offer)
}

func (offerStorage *offerStorage) UpdatePrice(ctx context.Context, offer offer.Offer, price float64) error {
	tableName := ctx.Value(DbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("database error update for product %s vendor %s price %f count %d: err %w",
			offer.Product().Name(),
			offer.Vendor().Name(),
			offer.Product().Price(),
			offer.Count(),
			err,
		)
	}
	return offerStorage.updatePrice(ctx, tableName, offer, price)
}

func (offerStorage *offerStorage) UpdateCount(ctx context.Context, offer offer.Offer, count int) error {
	tableName := ctx.Value(DbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("database error update for product %s vendor %s price %f count %d: err %w",
			offer.Product().Name(),
			offer.Vendor().Name(),
			offer.Product().Price(),
			offer.Count(),
			err,
		)
	}
	return offerStorage.updateCount(ctx, tableName, offer, count)
}

func (offerStorage *offerStorage) ListOffers(ctx context.Context, productName string) (offer.Offers, error) {
	tableName := ctx.Value(DbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.Offers{}, fmt.Errorf("database error list offers for product %s: err %w",
			productName,
			err,
		)
	}
	return offerStorage.listOffers(ctx, tableName, productName)
}

func (offerStorage *offerStorage) ListVendorOffers(ctx context.Context, vendorName string) (offer.Offers, error) {
	tableName := ctx.Value(DbTableDescriptorKey).(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.Offers{}, fmt.Errorf("database error list offers for vendor %s: err %w",
			vendorName,
			err,
		)
	}
	return offerStorage.listVendorOffers(ctx, tableName, vendorName)
}

func (offerStorage *offerStorage) add(ctx context.Context, dbTable string, offer offer.Offer) error {
	query := fmt.Sprintf("INSERT INTO %s (vendor,price,productName,count) VALUES ($1,$2,$3,$4)", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		offer.Vendor().Name(),
		offer.Product().Price(),
		offer.Product().Name(),
		offer.Count(),
	)
	if err != nil {
		return fmt.Errorf("database error add for product %s vendor %s price %f count %d: err %w",
			offer.Product().Name(),
			offer.Vendor().Name(),
			offer.Product().Price(),
			offer.Count(),
			err,
		)
	}
	return nil
}

func (offerStorage *offerStorage) remove(ctx context.Context, dbTable string, offer offer.Offer) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE price=$1 AND vendor=$2 AND productName=$3", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		offer.Product().Price(),
		offer.Vendor().Name(),
		offer.Product().Name(),
	)
	if err != nil {
		return fmt.Errorf("database error remove for product %s vendor %s price %f count %d: %w",
			offer.Product().Name(),
			offer.Vendor().Name(),
			offer.Product().Price(),
			offer.Count(),
			err,
		)
	}
	return nil
}

func (offerStorage *offerStorage) updateCount(ctx context.Context, dbTable string, offer offer.Offer, count int) error {
	query := fmt.Sprintf("UPDATE %s SET count=$1 WHERE vendor=$2 AND productName=$3 AND price=$4", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		count,
		offer.Vendor().Name(),
		offer.Product().Name(),
		offer.Product().Price(),
	)
	if err != nil {
		return fmt.Errorf("update count for offer with product name %s and vendor %s failed: err: %w", offer.Product().Name(), offer.Vendor().Name(), err)
	}
	return nil
}

func (offerStorage *offerStorage) updatePrice(ctx context.Context, dbTable string, offer offer.Offer, price float64) error {
	query := fmt.Sprintf("UPDATE %s SET price=$1 WHERE vendor=$2 AND productName=$3 AND price=$4", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		price,
		offer.Vendor().Name(),
		offer.Product().Name(),
		offer.Product().Price(),
	)
	if err != nil {
		return fmt.Errorf("update count for offer with product name %s and vendor %s failed: err: %w", offer.Product().Name(), offer.Vendor().Name(), err)
	}
	return nil
}

// FIXME scanning offers to offer.Offer as db type is different than domain one
func (offerStorage *offerStorage) listOffers(ctx context.Context, dbTable string, productName string) (offer.Offers, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE productName=$1", dbTable)
	rows, err := offerStorage.db.Query(ctx, query, productName)
	if err != nil {
		return nil, fmt.Errorf("database error sell offers listing for product %s %w", productName, err)
	}
	offers := make(offer.Offers, 0, 5)
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("database error sell offers listing during retrieving row: %w", err)
		}
		var offerModel offerModel
		err := rows.Scan(&offerModel.id, &offerModel.vendor, &offerModel.price, &offerModel.productName, &offerModel.count)
		if err != nil {
			return nil, fmt.Errorf("database error sell offers listing during retrieving row: %w", err)
		}
		offers = append(offers, offerModel.mapToDomainOffer())
	}
	return offers, nil
}

func (offerStorage *offerStorage) listVendorOffers(ctx context.Context, dbTable, vendorName string) (offer.Offers, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE vendor=$1", dbTable)
	rows, err := offerStorage.db.Query(ctx, query, vendorName)
	if err != nil {
		return nil, fmt.Errorf("database error sell offers listing for vendor %s %w", vendorName, err)
	}
	offers := make(offer.Offers, 0, 5)
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("database error sell offers listing during retrieving row: %w", err)
		}
		val, _ := rows.Values()
		fmt.Printf("rows model is %v\n", val)
		var offerModel offerModel
		err := rows.Scan(&offerModel.id, &offerModel.vendor, &offerModel.price, &offerModel.productName, &offerModel.count)
		if err != nil {
			return nil, fmt.Errorf("database error sell offers listing during retrieving row: %w", err)
		}
		fmt.Printf("offer model is %v\n", offerModel)
		offers = append(offers, offerModel.mapToDomainOffer())
	}
	return offers, nil
}

// FIXME validation for numeric price
func (offerStorage *offerStorage) createTable(ctx context.Context, name string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY, 
	vendor TEXT NOT NULL, 
	price NUMERIC(10,2) NOT NULL, 
	productName TEXT NOT NULL, 
	count INTEGER NOT NULL)`, name,
	)
	_, err := offerStorage.db.Exec(ctx, query)
	if err != nil {
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
