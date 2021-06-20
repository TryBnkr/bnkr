package dal

import (
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"

	"gorm.io/gorm"
)

// Queue struct defines the Queue Model
type Queue struct {
	gorm.Model
	Type   string
	Object *uint `gorm:"not null" gorm:"index"`
	// this is a pointer because int == 0,
}

// CreateQueue create a queue entry in the queue's table
func CreateQueue(queue *Queue) *gorm.DB {
	return database.DB.Create(&queue)
}

// FindQueue finds a queue with given condition
func FindQueue(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Take(dest, conds...)
}

// FindQueueByUser finds a queue with given queue and user identifier
func FindQueueByUser(dest interface{}, queueIden interface{}, userIden interface{}) *gorm.DB {
	return FindQueue(dest, "id = ? AND user = ?", queueIden, userIden)
}

func FindQueueTypeAndObject(dest interface{}, t interface{}, o interface{}) *gorm.DB {
	return FindQueue(dest, "type = ? AND object = ?", t, o)
}

// FindQueuesByUser finds the queues with user's identifier given
func FindQueuesByUser(dest interface{}, userIden interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Find(dest, "user = ?", userIden)
}

func FindQueuesByBackup(dest interface{}, backupIden interface{}, order interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Order(order).Find(dest, "backup = ?", backupIden)
}

func FindQueuesByObjectsIdsAndType(dest interface{}, ids interface{}, t string, order interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Order(order).Where("object IN ? AND type = ?", ids, t).Find(dest)
}

func FindQueuesById(dest interface{}, queueIden interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Find(dest, "id = ?", queueIden)
}

func FindAllQueues(dest interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Find(dest)
}

// DeleteQueue deletes a queue from queues' table with the given queue and user identifier
// func DeleteQueue(queueIden interface{}, userIden interface{}) *gorm.DB {
// 	return database.DB.Unscoped().Delete(&Queue{}, "id = ? AND user = ?", queueIden, userIden)
// }
func DeleteQueue(queueIden interface{}) *gorm.DB {
	return database.DB.Unscoped().Delete(&Queue{}, "id = ?", queueIden)
}

// UpdateQueue allows to update the queue with the given queueID and userID
func UpdateUserQueue(queueIden interface{}, userIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Where("id = ? AND user = ?", queueIden, userIden).Updates(data)
}

func UpdateQueue(queueIden interface{}, data interface{}) *gorm.DB {
	return database.DB.Model(&Queue{}).Where("id = ?", queueIden).Updates(data)
}
