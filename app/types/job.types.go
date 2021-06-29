package types

import (
	"time"

	"github.com/jackc/pgtype"
)

type NewJobDTO struct {
	ID        int    `db:"id"`
	File      string `db:"file"`
	Status    string `db:"status"`
	Backup    int    `db:"backup"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
}

type SmallJob struct {
	Backup    int
	Status    string
	CreatedAt time.Time
}
