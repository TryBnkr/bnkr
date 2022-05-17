package types

import (
	"database/sql"
)

type NewBackupDTO struct {
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
