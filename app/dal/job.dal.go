package dal

import (
	"database/sql"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
)

// Job struct defines the Job Model
type Job struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	File      string    `db:"file"`
	Status    string    `db:"status"`
	Backup    int       `db:"backup"`
}

// CreateJob create a job entry in the job's table
func CreateJob(job *Job) (sql.Result, error) {
	result, err := database.DB.NamedExec(`INSERT INTO jobs (created_at, updated_at, file, status, backup)
	VALUES (:created_at, :updated_at, :file, :status, :backup)`, *job)

	return result, err
}

func FindJobsByBackup(dest interface{}, backupIden interface{}, order string) error {
	return database.DB.Select(dest, "SELECT * FROM jobs WHERE backup=$1 ORDER BY "+order, backupIden)
}

func SelectLatestJobForEachBackup(dest interface{}) error {
	return database.DB.Select(dest, `
	SELECT m.backup,ca,m.status FROM (
		SELECT
			backup, MAX(created_at) AS ca
		FROM
			jobs 
		GROUP BY
			backup) t join jobs m on t.backup = m.backup and t.ca = m.created_at;
	`)
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
