package dal

import (
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"

	"gorm.io/gorm"
)

// User struct defines the user
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

// CreateUser create a user entry in the user's table
func CreateUser(user *User) *gorm.DB {
	return database.DB.Create(user)
}

// FindUser searches the user's table with the condition given
func FindUser(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.DB.Model(&User{}).Take(dest, conds...)
}

// FindUserByEmail searches the user's table with the email given
func FindUserByEmail(dest interface{}, email string) *gorm.DB {
	return FindUser(dest, "email = ?", email)
}

func FindUsersById(dest interface{}, userIden interface{}) *gorm.DB {
	return database.DB.Model(&User{}).Find(dest, "id = ?", userIden)
}

func DeleteUser(userIden interface{}) *gorm.DB {
	return database.DB.Unscoped().Delete(&User{}, "id = ?", userIden)
}

func FindAllUsers(dest interface{}) *gorm.DB {
	return database.DB.Model(&User{}).Find(dest)
}

func UpdateUser(userIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&User{}).Where("id = ?", userIden).Updates(data)
}
