package dal

import (
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
)

func FindAllMigrations(dest interface{}) error {
	return database.DB.Select(dest, "SELECT * FROM migrations ORDER BY id DESC")
}

func FindMigrationById(dest interface{}, migrationIden interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM migrations WHERE id=$1", migrationIden)
}