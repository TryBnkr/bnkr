package dal

import (
	"database/sql"
	"strconv"

	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/paginator"
)

// Backup struct defines the Backup Model
type Backup struct {
	ID               int            `db:"id"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
	DeletedAt        sql.NullTime   `db:"deleted_at"`
	Name             string         `db:"name"`
	Enable           bool           `db:"enable"`
	Frequency        string         `db:"frequency"`
	Timezone         string         `db:"timezone"`
	CustomFrequency  string         `db:"custom_frequency"`
	Type             string         `db:"type"`
	Bucket           string         `db:"bucket"`
	Region           string         `db:"region"`
	DbName           string         `db:"db_name"`
	DbUser           string         `db:"db_user"`
	DbPassword       string         `db:"db_password"`
	DbHost           string         `db:"db_host"`
	DbPort           string         `db:"db_port"`
	PodLabel         string         `db:"pod_label"`
	PodName          string         `db:"pod_name"`
	DayOfWeek        int            `db:"day_of_week"`
	DayOfMonth       int            `db:"day_of_month"`
	Month            int            `db:"month"`
	Time             string         `db:"time"`
	Container        string         `db:"container"`
	FilesPath        string         `db:"files_path"`
	S3AccessKey      string         `db:"s3_access_key"`
	S3SecretKey      string         `db:"s3_secret_key"`
	StorageDirectory string         `db:"storage_directory"`
	Retention        int            `db:"retention"`
	Emails           string         `db:"emails"`
	URI              sql.NullString `db:"uri"`
	User             int            `db:"user"`
}

// CreateBackup create a backup entry in the backup's table
func CreateBackup(backup *Backup) (int, error) {
	var id int

	rows, err := database.DB.NamedQuery(`INSERT INTO backups (created_at, updated_at, "name", "enable", frequency, timezone, custom_frequency, "type", bucket, region, db_name, db_user, db_password, db_host, db_port, pod_label, pod_name, day_of_week, day_of_month, "month", "time", container, files_path, s3_access_key, s3_secret_key, storage_directory, retention, emails, "user", uri)
	VALUES (:created_at, :updated_at, :name, :enable, :frequency, :timezone, :custom_frequency, :type, :bucket, :region, :db_name, :db_user, :db_password, :db_host, :db_port, :pod_label, :pod_name, :day_of_week, :day_of_month, :month, :time, :container, :files_path, :s3_access_key, :s3_secret_key, :storage_directory, :retention, :emails, :user, :uri) RETURNING id`, *backup)

	if rows.Next() {
		rows.Scan(&id)
	}

	backup.ID = id

	return id, err
}

func FindBackupById(dest interface{}, backupIden interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM backups WHERE id=$1", backupIden)
}

func FindAllBackups(dest interface{}) error {
	return database.DB.Select(dest, "SELECT * FROM backups ORDER BY id ASC")
}

func FindBackups(dest interface{}, order string, p *paginator.Paginator) error {
	return database.DB.Select(dest, "SELECT * FROM backups ORDER BY "+order+" LIMIT "+strconv.Itoa(p.PerPage)+" OFFSET "+strconv.Itoa(p.Offset()))
}

func DeleteBackup(backupIden interface{}) (sql.Result, error) {
	return database.DB.Exec("delete from backups where id=$1", backupIden)
}

func UpdateBackup(data interface{}) (sql.Result, error) {
	result, err := database.DB.NamedExec(`UPDATE backups SET (updated_at, "name", "enable", frequency, timezone, custom_frequency, "type", bucket, region, db_name, db_user, db_password, db_host, db_port, pod_label, pod_name, day_of_week, day_of_month, "month", "time", container, files_path, s3_access_key, s3_secret_key, storage_directory, retention, emails, uri)
	= (:updated_at, :name, :enable, :frequency, :timezone, :custom_frequency, :type, :bucket, :region, :db_name, :db_user, :db_password, :db_host, :db_port, :pod_label, :pod_name, :day_of_week, :day_of_month, :month, :time, :container, :files_path, :s3_access_key, :s3_secret_key, :storage_directory, :retention, :emails, :uri) where id=:id`, data)
	return result, err
}
