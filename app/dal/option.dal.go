package dal

import (
	"database/sql"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
)

// Option struct defines the Option Model
type Option struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
	Value     string    `db:"value"`
}

// CreateOption create a option entry in the option's table
func CreateOption(option *Option) (sql.Result, error) {
	result, err := database.DB.NamedExec(`INSERT INTO options (created_at, updated_at, name, value)
	VALUES (:created_at, :updated_at, :name, :value)`, *option)

	return result, err
}

func CreateOrUpdateOption(option *Option) (sql.Result, error) {
	return database.DB.NamedExec(`
	INSERT INTO options (name, value) 
	VALUES (:name, :value) 
	ON CONFLICT (name) DO UPDATE 
	SET value = :value
		`, *option)
}

func FindOptionByName(dest interface{}, optionName interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM options WHERE name=$1", optionName)
}

func FindAllOptions(dest interface{}) error {
	return database.DB.Select(dest, "SELECT * FROM options ORDER BY id ASC")
}
