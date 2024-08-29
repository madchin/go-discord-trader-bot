package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type dbtableDescriptorKey struct{}

var DbTableDescriptorKey = dbtableDescriptorKey{}

func Connect(user, password, dbName string) (*pgx.Conn, error) {
	dbUrl := fmt.Sprintf("postgres://postgres:%s@db:5432/postgres", password)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database err: %v", err)
	}
	return conn, nil
}

func DbTableDescriptorValue(name, guildId string) string {
	return name + "_" + guildId
}
