package types

import (
	"database/sql"
	"time"
)

type NewJobDTO struct {
	ID          int          `db:"id"`
	File        string       `db:"file"`
	Status      string       `db:"status"`
	Backup      int          `db:"backup"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
	CompletedAt sql.NullTime `db:"completed_at"`
}

type SmallJob struct {
	Backup    int       `db:"backup"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"ca"`
}
