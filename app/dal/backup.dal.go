package dal

import (
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"

	"gorm.io/gorm"
)

// Backup struct defines the Backup Model
type Backup struct {
	gorm.Model
	Name             string `gorm:"not null"`
	Frequency        string `gorm:"not null"`
	Timezone         string
	CustomFrequency  string
	Type             string `gorm:"not null"`
	Bucket           string
	Region           string
	DbName           string
	DbUser           string
	DbPassword       string
	DbHost           string
	DbPort           string
	PodLabel         string
	PodName          string
	DayOfWeek        *int
	DayOfMonth       int
	Month            int
	Time             string
	Container        string
	FilesPath        string
	S3AccessKey      string
	S3SecretKey      string
	StorageDirectory string
	Retention        uint
	Emails           string
	User             *uint `gorm:"not null" gorm:"index"`
	// this is a pointer because int == 0,
}

// CreateBackup create a backup entry in the backup's table
func CreateBackup(backup *Backup) *gorm.DB {
	return database.DB.Create(backup)
}

// FindBackup finds a backup with given condition
func FindBackup(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.DB.Model(&Backup{}).Take(dest, conds...)
}

// FindBackupByUser finds a backup with given backup and user identifier
func FindBackupByUser(dest interface{}, backupIden interface{}, userIden interface{}) *gorm.DB {
	return FindBackup(dest, "id = ? AND user = ?", backupIden, userIden)
}

// FindBackupsByUser finds the backups with user's identifier given
func FindBackupsByUser(dest interface{}, userIden interface{}) *gorm.DB {
	return database.DB.Model(&Backup{}).Find(dest, "user = ?", userIden)
}

func FindBackupsById(dest interface{}, backupIden interface{}) *gorm.DB {
	return database.DB.Model(&Backup{}).Find(dest, "id = ?", backupIden)
}

func FindAllBackups(dest interface{}) *gorm.DB {
	return database.DB.Model(&Backup{}).Find(dest)
}

// DeleteBackup deletes a backup from backups' table with the given backup and user identifier
// func DeleteBackup(backupIden interface{}, userIden interface{}) *gorm.DB {
// 	return database.DB.Unscoped().Delete(&Backup{}, "id = ? AND user = ?", backupIden, userIden)
// }
func DeleteBackup(backupIden interface{}) *gorm.DB {
	return database.DB.Unscoped().Delete(&Backup{}, "id = ?", backupIden)
}

// UpdateBackup allows to update the backup with the given backupID and userID
func UpdateUserBackup(backupIden interface{}, userIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&Backup{}).Where("id = ? AND user = ?", backupIden, userIden).Updates(data)
}

func UpdateBackup(backupIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&Backup{}).Where("id = ?", backupIden).Updates(data)
}
