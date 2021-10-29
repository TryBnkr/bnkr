package dal

import "github.com/MohammedAl-Mahdawi/bnkr/config/database"

func Count(dest interface{}, table string, where string) error {
	return database.DB.Get(dest, "SELECT COUNT(*) FROM "+table+where)
}
