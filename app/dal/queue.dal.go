package dal

import (
	"database/sql"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
)

// Queue struct defines the Queue Model
type Queue struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Type      string    `db:"type"`
	Object    int       `db:"object"`
}

// CreateQueue create a queue entry in the queue's table
func CreateQueue(queue *Queue) (sql.Result, error) {
	result, err := database.DB.NamedExec(`INSERT INTO queues (created_at, updated_at, type, object)
	VALUES (:created_at, :updated_at, :type, :object)`, *queue)

	return result, err
}

func FindQueueByTypeAndObject(dest interface{}, t interface{}, o interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM queues WHERE type=$1 AND object=$2", t, o)
}

func FindQueuesByObjectsIdsAndType(dest interface{}, ids interface{}, t string, order string) error {
	return database.DB.Get(dest, "SELECT * FROM queues WHERE object IN $1 AND type = $2 ORDER BY "+order, ids, t)
}

func DeleteQueue(queueIden interface{}) (sql.Result, error) {
	result, err := database.DB.Exec("delete from queues where id=$1", queueIden)
	return result, err
}
