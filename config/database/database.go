package database

import (
	"embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jmoiron/sqlx"
)

// DB is the underlying database connection
var DB *sqlx.DB

//go:embed migrations
var migrations embed.FS

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
	source, _ := httpfs.New(http.FS(migrations), "migrations")
	m, err := migrate.NewWithSourceInstance("httpfs", source, dsn)
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
