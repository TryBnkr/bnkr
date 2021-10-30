package dal

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/paginator"
	"github.com/jmoiron/sqlx"
)

// Job struct defines the Job Model
type Job struct {
	ID          int          `db:"id"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	CompletedAt sql.NullTime `db:"completed_at"`
	File        string       `db:"file"`
	Status      string       `db:"status"`
	Backup      int          `db:"backup"`
}

// CreateJob create a job entry in the job's table
func CreateJob(job *Job) (sql.Result, error) {
	result, err := database.DB.NamedExec(`INSERT INTO jobs (created_at, updated_at, completed_at, file, status, backup)
	VALUES (:created_at, :updated_at, :completed_at, :file, :status, :backup)`, *job)

	return result, err
}

func FindAllJobsByBackup(dest interface{}, backupIden interface{}, order string) error {
	return database.DB.Select(dest, "SELECT * FROM jobs WHERE backup=$1 ORDER BY "+order, backupIden)
}

func FindJobsByBackup(dest interface{}, backupIden interface{}, order string, p *paginator.Paginator) error {
	return database.DB.Select(dest, "SELECT * FROM jobs WHERE backup=$1 ORDER BY "+order+" LIMIT "+strconv.Itoa(p.PerPage)+" OFFSET "+strconv.Itoa(p.Offset()), backupIden)
}

func SelectPaginatedJobsByBackup(dest interface{}, backupIden interface{}, order string) error {
	return database.DB.Select(dest, "SELECT * FROM jobs WHERE backup=$1 ORDER BY "+order, backupIden)
}

func SelectLatestJobForEachBackup(backupsIds []int) ([]types.SmallJob, error) {
	dest := []types.SmallJob{}
	
	query, args, err := sqlx.In(`
	SELECT m.backup,ca,m.status FROM (
		SELECT
			backup, MAX(created_at) AS ca
		FROM
			jobs 
Where backup in (?)
		GROUP BY
			backup) t join jobs m on t.backup = m.backup and t.ca = m.created_at;
		`, backupsIds)

	if err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(database.DB.Rebind(query), args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var CreatedAt time.Time
		var Bakup int
		var Status string

		err = rows.Scan(&Bakup, &CreatedAt, &Status)
		if err != nil {
			return nil, err
		}
		dest = append(dest, types.SmallJob{
			CreatedAt: CreatedAt,
			Backup:    Bakup,
			Status:    Status,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func FindJobsIDByBackup(dest interface{}, backupIden interface{}, order string) error {
	return database.DB.Select(dest, "SELECT id FROM jobs WHERE backup=$1 ORDER BY "+order, backupIden)
}

func FindJobById(dest interface{}, jobIden interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM jobs WHERE id=$1", jobIden)
}

func DeleteJob(jobIden interface{}) (sql.Result, error) {
	result, err := database.DB.Exec("delete from jobs where id=$1", jobIden)
	return result, err
}
