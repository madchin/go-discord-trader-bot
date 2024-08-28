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
	tableName := ctx.Value("dbTableDescriptor").(string)
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
	tableName := ctx.Value("dbTableDescriptor").(string)
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

func (offerStorage *offerStorage) Update(ctx context.Context, oldOffer offer.Offer, updateOffer offer.Offer) error {
	tableName := ctx.Value("dbTableDescriptor").(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return fmt.Errorf("database error update for product %s vendor %s price %f count %d: err %w",
			oldOffer.Product().Name(),
			oldOffer.Vendor().Name(),
			oldOffer.Product().Price(),
			oldOffer.Count(),
			err,
		)
	}
	return offerStorage.update(ctx, tableName, oldOffer, updateOffer)
}

func (offerStorage *offerStorage) ListOffers(ctx context.Context, productName string) (offer.Offers, error) {
	tableName := ctx.Value("dbTableDescriptor").(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.Offers{}, fmt.Errorf("database error list offers for product %s: err %w",
			productName,
			err,
		)
	}
	return offerStorage.listOffers(ctx, tableName, productName)
}

func (offerStorage *offerStorage) ListVendorOffers(ctx context.Context, vendorName string) (offer.Offers, error) {
	tableName := ctx.Value("dbTableDescriptor").(string)
	if err := offerStorage.createTable(ctx, tableName); err != nil {
		return offer.Offers{}, fmt.Errorf("database error list offers for vendor %s: err %w",
			vendorName,
			err,
		)
	}
	return offerStorage.listVendorOffers(ctx, tableName, vendorName)
}

func (offerStorage *offerStorage) add(ctx context.Context, dbTable string, offer offer.Offer) error {
	offers, _ := offerStorage.listVendorOffers(ctx, dbTable, offer.Vendor().Name())
	if offers.Contains(offer) {
		offer = offers.MergeSameOffers(offer)
		err := offerStorage.updateCount(ctx, dbTable, offer, offer.Count())
		if err != nil {
			return fmt.Errorf("database error add, wanted to update offer because you have already same %w", err)
		}
	}
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
	offers, _ := offerStorage.listOffers(ctx, dbTable, offer.Product().Name())
	if !offers.Contains(offer) {
		return fmt.Errorf("unable to remove offer because you do not have this one")
	}
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

func (offerStorage *offerStorage) update(ctx context.Context, dbTable string, oldOffer offer.Offer, updateOffer offer.Offer) error {
	offers, _ := offerStorage.listVendorOffers(ctx, dbTable, oldOffer.Vendor().Name())
	if !offers.Contains(oldOffer) {
		return fmt.Errorf("unable to update offer because you do not have it")
	}
	query := fmt.Sprintf("UPDATE %s SET price=$1, count=$2 WHERE vendor=$3 AND productName=$4 AND price=$5", dbTable)
	_, err := offerStorage.db.Exec(ctx, query,
		updateOffer.Product().Price(),
		updateOffer.Count(),
		oldOffer.Vendor().Name(),
		oldOffer.Product().Name(),
		oldOffer.Product().Price(),
	)
	if err != nil {
		return fmt.Errorf("update for offer with product name %s and vendor %s failed: err: %w", oldOffer.Product().Name(), oldOffer.Vendor().Name(), err)
	}
	return nil
}

func (offerStorage *offerStorage) updateCount(ctx context.Context, dbTable string, offer offer.Offer, count int) error {
	query := fmt.Sprintf("UPDATE %s SET count=$2 WHERE vendor=$3 AND productName=$4 AND price=$5", dbTable)
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
		var offer offer.Offer
		err := rows.Scan(&offer)
		if err != nil {
			return nil, fmt.Errorf("database error sell offers listing during retrieving row: %w", err)
		}
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
		var offer offer.Offer
		err := rows.Scan(&offer)
		if err != nil {
			return nil, fmt.Errorf("database error sell offers listing during retrieving row: %w", err)
		}
	}
	return offers, nil
}

func (offerStorage *offerStorage) createTable(ctx context.Context, name string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (
	id SERIAL PRIMARY KEY, 
	vendor TEXT NOT NULL, 
	price NUMERIC(7,5) NOT NULL, 
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
