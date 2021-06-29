package dal

import (
	"database/sql"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
)

// User struct defines the user
type User struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
}

// CreateUser create a user entry in the user's table
func CreateUser(user *User) (sql.Result, error) {
	result, err := database.DB.NamedExec(`INSERT INTO users (created_at, updated_at, name, email, password)
	VALUES (:created_at, :updated_at, :name, :email, :password)`, *user)

	return result, err
}

// FindUserByEmail searches the user's table with the email given
func FindUserByEmail(dest interface{}, email string) error {
	return database.DB.Get(dest, "SELECT * FROM users WHERE email=$1", email)
}

func FindUserById(dest interface{}, userIden interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM users WHERE id=$1", userIden)
}

func DeleteUser(userIden interface{}) (sql.Result, error) {
	return database.DB.Exec("delete from users where id=$1", userIden)
}

func FindAllUsers(dest interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM users")
}

func UpdateUser(data interface{}) (sql.Result, error) {
	return database.DB.NamedExec(`UPDATE users SET (created_at, updated_at, name, email, password)
	= (:created_at, :updated_at, :name, :email, :password) where id=:id`, data)
}
