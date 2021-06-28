package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB is the underlying database connection
var DB *pgxpool.Pool

// Connect initiate the database connection and migrate all the tables
func Connect(dsn string) {
	dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Println("[DATABASE]::CONNECTION_ERROR")
		panic(err)
	}

	// Setting the database connection to use in routes
	DB = dbpool

	fmt.Println("[DATABASE]::CONNECTED")
}

// Migrate migrates all the database tables
func Migrate(dsn string) error {
	m, err := migrate.New(
		"file://migrations",
		dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if !(errors.Is(err, migrate.ErrNoChange)) {
			return err
		}
	}

	return nil
}
