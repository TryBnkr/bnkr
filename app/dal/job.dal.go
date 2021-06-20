package dal

import (
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"

	"gorm.io/gorm"
)

// Job struct defines the Job Model
type Job struct {
	gorm.Model
	File   string
	Status string
	Backup *uint `gorm:"not null" gorm:"index"`
	// this is a pointer because int == 0,
}

// CreateJob create a job entry in the job's table
func CreateJob(job *Job) *gorm.DB {
	return database.DB.Create(&job)
}

// FindJob finds a job with given condition
func FindJob(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Take(dest, conds...)
}

// FindJobByUser finds a job with given job and user identifier
func FindJobByUser(dest interface{}, jobIden interface{}, userIden interface{}) *gorm.DB {
	return FindJob(dest, "id = ? AND user = ?", jobIden, userIden)
}

// FindJobsByUser finds the jobs with user's identifier given
func FindJobsByUser(dest interface{}, userIden interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Find(dest, "user = ?", userIden)
}

func FindJobsByBackup(dest interface{}, backupIden interface{}, order interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Order(order).Find(dest, "backup = ?", backupIden)
}

func FindJobsIDByBackup(dest interface{}, backupIden interface{}, order interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Select("id").Order(order).Find(dest, "backup = ?", backupIden)
}

func FindJobsById(dest interface{}, jobIden interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Find(dest, "id = ?", jobIden)
}

func FindAllJobs(dest interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Find(dest)
}

// DeleteJob deletes a job from jobs' table with the given job and user identifier
// func DeleteJob(jobIden interface{}, userIden interface{}) *gorm.DB {
// 	return database.DB.Unscoped().Delete(&Job{}, "id = ? AND user = ?", jobIden, userIden)
// }
func DeleteJob(jobIden interface{}) *gorm.DB {
	return database.DB.Unscoped().Delete(&Job{}, "id = ?", jobIden)
}

// UpdateJob allows to update the job with the given jobID and userID
func UpdateUserJob(jobIden interface{}, userIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Where("id = ? AND user = ?", jobIden, userIden).Updates(data)
}

func UpdateJob(jobIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&Job{}).Where("id = ?", jobIden).Updates(data)
}
