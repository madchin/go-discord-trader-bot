package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type ctxItembTableDescriptorKey struct{}

var CtxItemTableDescriptorKey = &ctxItembTableDescriptorKey{}

type ctxBuySellDbTableDescriptorKey struct{}

var CtxBuySellDbTableDescriptorKey = &ctxBuySellDbTableDescriptorKey{}

type dbCreds struct {
	databaseName string
	password     string
	user         string
}

func Connect(dbCreds dbCreds) (*pgx.Conn, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@db:5432/%s", dbCreds.user, dbCreds.password, dbCreds.databaseName)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("connection to database err: %v", err)
	}
	return conn, nil
}

func TableWithGuildIdSuffix(name, guildId string) string {
	return name + "_" + guildId
}

func LoadCredentials() (dbCreds, error) {
	dbNameFilePath := os.Getenv("DB_NAME_FILE")
	if dbNameFilePath == "" {
		return dbCreds{}, errors.New("DB_NAME_FILE environment variable not provided. It needs to be set with path to .db.name.env file")
	}
	dbPasswordFilePath := os.Getenv("DB_PASSWORD_FILE")
	if dbPasswordFilePath == "" {
		return dbCreds{}, errors.New("DB_PASSWORD_FILE environment variable not provided. It needs to be set with path to .db.password.env file")
	}
	dbUserFilePath := os.Getenv("DB_USER_FILE")
	if dbUserFilePath == "" {
		return dbCreds{}, errors.New("DB_USER_FILE environment variable not provided. It needs to be set with path to .db.user.env file")
	}
	databaseName, err := os.ReadFile(dbNameFilePath)
	if err != nil {
		panic(err)
	}
	password, err := os.ReadFile(dbPasswordFilePath)
	if err != nil {
		panic(err)
	}
	user, err := os.ReadFile(dbUserFilePath)
	if err != nil {
		panic(err)
	}
	return dbCreds{databaseName: string(databaseName), password: string(password), user: string(user)}, nil
}
