package dal

import (
	"database/sql"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
	"github.com/jmoiron/sqlx"
)

// Queue struct defines the Queue Model
type Queue struct {
	ID        int          `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Type      string       `db:"type"`
	Object    int          `db:"object"`
}

// CreateQueue create a queue entry in the queue's table
func CreateQueue(queue *Queue) (int, error) {
	var id int

	rows, err := database.DB.NamedQuery(`INSERT INTO queues (created_at, updated_at, type, object)
	VALUES (:created_at, :updated_at, :type, :object) RETURNING id`, *queue)

	if rows.Next() {
		rows.Scan(&id)
	}

	queue.ID = id

	return id, err
}

func FindQueueByTypeAndObject(dest interface{}, t interface{}, o interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM queues WHERE type=$1 AND object=$2", t, o)
}

func FindQueuesByObjectsIdsAndType(dest interface{}, ids interface{}, t string, order string) error {
	query, args, err := sqlx.In("SELECT * FROM queues WHERE object IN (?) AND type = (?) ORDER BY "+order+";", ids, t)

	if err != nil {
		return err
	}

	return database.DB.Select(dest, database.DB.Rebind(query), args...)
}

func DeleteQueue(queueIden interface{}) (sql.Result, error) {
	return database.DB.Exec("delete from queues where id=$1", queueIden)
}
