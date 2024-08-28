package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Connect(user, password, dbName string) (*pgx.Conn, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@db:5432/%s", user, password, dbName)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database err: %v", err)
	}
	return conn, nil
}
