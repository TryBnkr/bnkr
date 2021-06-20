package dal

import (
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"

	"gorm.io/gorm"
)

// Option struct defines the Option Model
type Option struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;not null"`
	Value string
}

// CreateOption create a option entry in the option's table
func CreateOption(option *Option) *gorm.DB {
	return database.DB.Create(option)
}

// FindOption finds a option with given condition
func FindOption(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.DB.Model(&Option{}).Take(dest, conds...)
}

// FindOptionByUser finds a option with given option and user identifier
func FindOptionByUser(dest interface{}, optionIden interface{}, userIden interface{}) *gorm.DB {
	return FindOption(dest, "id = ? AND user = ?", optionIden, userIden)
}

// FindOptionsByUser finds the options with user's identifier given
func FindOptionsByUser(dest interface{}, userIden interface{}) *gorm.DB {
	return database.DB.Model(&Option{}).Find(dest, "user = ?", userIden)
}

func FindOptionsByName(dest interface{}, optionName interface{}) *gorm.DB {
	return database.DB.Model(&Option{}).Find(dest, "name = ?", optionName)
}

// DeleteOption deletes a option from options' table with the given option and user identifier
func DeleteOption(optionIden interface{}, userIden interface{}) *gorm.DB {
	return database.DB.Unscoped().Delete(&Option{}, "id = ? AND user = ?", optionIden, userIden)
}

// UpdateOption allows to update the option with the given optionID and userID
func UpdateOption(optionIden interface{}, userIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&Option{}).Where("id = ? AND user = ?", optionIden, userIden).Updates(data)
}

func FindAllOptions(dest interface{}) *gorm.DB {
	return database.DB.Model(&Option{}).Find(dest)
}
