package types

import (
	"time"

	"github.com/jackc/pgtype"
)

type NewJobDTO struct {
	ID        int                `db:"id"`
	File      string             `db:"file"`
	Status    string             `db:"status"`
	Backup    int                `db:"backup"`
	CreatedAt        pgtype.Timestamptz `db:"created_at"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at"`
	DeletedAt        pgtype.Timestamptz `db:"deleted_at"`
}

type SmallJob struct {
	Backup    int       `db:"backup"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"ca"`
}
