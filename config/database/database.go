package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// DB is the underlying database connection
var DB *sqlx.DB

// Connect initiate the database connection and migrate all the tables
func Connect(dsn string) {
	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		fmt.Println("[DATABASE]::CONNECTION_ERROR")
		panic(err)
	}

	// Setting the database connection to use in routes
	DB = db

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
