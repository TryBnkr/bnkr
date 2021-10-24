package dal

import (
	"database/sql"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
)

type Migration struct {
	ID                   int            `db:"id"`
	CreatedAt            sql.NullTime   `db:"created_at"`
	UpdatedAt            sql.NullTime   `db:"updated_at"`
	DeletedAt            sql.NullTime   `db:"deleted_at"`
	StartedAt            sql.NullTime   `db:"started_at"`
	CompletedAt          sql.NullTime   `db:"completed_at"`
	Name                 string         `db:"name"`
	User                 int            `db:"user"`
	Timezone             string         `db:"timezone"`
	Emails               string         `db:"emails"`
	SrcType              string         `db:"src_type"`
	SrcAccess            string         `db:"src_access"`
	SrcBucket            string         `db:"src_bucket"`
	SrcRegion            string         `db:"src_region"`
	SrcDbName            string         `db:"src_db_name"`
	SrcDbUser            string         `db:"src_db_user"`
	SrcDbPassword        string         `db:"src_db_password"`
	SrcDbHost            string         `db:"src_db_host"`
	SrcDbPort            string         `db:"src_db_port"`
	SrcPodLabel          string         `db:"src_pod_label"`
	SrcPodName           string         `db:"src_pod_name"`
	SrcContainer         string         `db:"src_container"`
	SrcFilesPath         string         `db:"src_files_path"`
	SrcS3AccessKey       string         `db:"src_s3_access_key"`
	SrcS3SecretKey       string         `db:"src_s3_secret_key"`
	SrcS3File            string         `db:"src_s3_file"`
	SrcURI               string         `db:"src_uri"`
	SrcSshHost           string         `db:"src_ssh_host"`
	SrcSshPort           string         `db:"src_ssh_port"`
	SrcSshUser           string         `db:"src_ssh_user"`
	SrcSshKey            string         `db:"src_ssh_key"`
	SrcKubeconfig        string         `db:"src_kubeconfig"`
	DestType             string         `db:"dest_type"`
	DestAccess           string         `db:"dest_access"`
	DestBucket           string         `db:"dest_bucket"`
	DestRegion           string         `db:"dest_region"`
	DestDbName           string         `db:"dest_db_name"`
	DestDbUser           string         `db:"dest_db_user"`
	DestDbPassword       string         `db:"dest_db_password"`
	DestDbHost           string         `db:"dest_db_host"`
	DestDbPort           string         `db:"dest_db_port"`
	DestPodLabel         string         `db:"dest_pod_label"`
	DestPodName          string         `db:"dest_pod_name"`
	DestContainer        string         `db:"dest_container"`
	DestFilesPath        string         `db:"dest_files_path"`
	DestS3AccessKey      string         `db:"dest_s3_access_key"`
	DestS3SecretKey      string         `db:"dest_s3_secret_key"`
	DestStorageDirectory string         `db:"dest_storage_directory"`
	DestURI              string         `db:"dest_uri"`
	DestSshHost          string         `db:"dest_ssh_host"`
	DestSshPort          string         `db:"dest_ssh_port"`
	DestSshUser          string         `db:"dest_ssh_user"`
	DestSshKey           string         `db:"dest_ssh_key"`
	DestKubeconfig       string         `db:"dest_kubeconfig"`
	Status               sql.NullString `db:"status"`
	Output               sql.NullString `db:"output"`
}

func CreateMigration(migration *Migration) (int, error) {
	var id int

	rows, err := database.DB.NamedQuery(`INSERT INTO migrations (created_at, updated_at, "name", timezone, src_type, dest_type, src_access, dest_access, src_bucket, dest_bucket, src_region, dest_region, src_db_name, dest_db_name, src_db_user, dest_db_user, src_db_password, dest_db_password, src_db_host, dest_db_host, src_db_port, dest_db_port, src_pod_label, dest_pod_label, src_pod_name, dest_pod_name, src_container, dest_container, src_files_path, dest_files_path, src_s3_access_key, dest_s3_access_key, src_s3_secret_key, dest_s3_secret_key, src_s3_file, dest_storage_directory, src_ssh_host, dest_ssh_host, src_ssh_port, dest_ssh_port, src_ssh_user, dest_ssh_user, src_ssh_key, dest_ssh_key, src_uri, dest_uri, dest_kubeconfig, src_kubeconfig, emails, "user")
	VALUES (:created_at, :updated_at, :name, :timezone, :src_type, :dest_type, :src_access, :dest_access, :src_bucket, :dest_bucket, :src_region, :dest_region, :src_db_name, :dest_db_name, :src_db_user, :dest_db_user, :src_db_password, :dest_db_password, :src_db_host, :dest_db_host, :src_db_port, :dest_db_port, :src_pod_label, :dest_pod_label, :src_pod_name, :dest_pod_name, :src_container, :dest_container, :src_files_path, :dest_files_path, :src_s3_access_key, :dest_s3_access_key, :src_s3_secret_key, :dest_s3_secret_key, :src_s3_file, :dest_storage_directory, :src_ssh_host, :dest_ssh_host, :src_ssh_port, :dest_ssh_port, :src_ssh_user, :dest_ssh_user, :src_ssh_key, :dest_ssh_key, :src_uri, :dest_uri, :dest_kubeconfig, :src_kubeconfig, :emails, :user) RETURNING id`, *migration)

	if rows.Next() {
		rows.Scan(&id)
	}

	migration.ID = id

	return id, err
}

func UpdateMigration(data interface{}) (sql.Result, error) {
	result, err := database.DB.NamedExec(`UPDATE migrations SET (updated_at, "name", timezone, src_type, dest_type, src_access, dest_access, src_bucket, dest_bucket, src_region, dest_region, src_db_name, dest_db_name, src_db_user, dest_db_user, src_db_password, dest_db_password, src_db_host, dest_db_host, src_db_port, dest_db_port, src_pod_label, dest_pod_label, src_pod_name, dest_pod_name, src_container, dest_container, src_files_path, dest_files_path, src_s3_access_key, dest_s3_access_key, src_s3_secret_key, dest_s3_secret_key, src_s3_file, dest_storage_directory, src_ssh_host, dest_ssh_host, src_ssh_port, dest_ssh_port, src_ssh_user, dest_ssh_user, src_ssh_key, dest_ssh_key, src_uri, dest_uri, dest_kubeconfig, src_kubeconfig, emails)
	= (:updated_at, :name, :timezone, :src_type, :dest_type, :src_access, :dest_access, :src_bucket, :dest_bucket, :src_region, :dest_region, :src_db_name, :dest_db_name, :src_db_user, :dest_db_user, :src_db_password, :dest_db_password, :src_db_host, :dest_db_host, :src_db_port, :dest_db_port, :src_pod_label, :dest_pod_label, :src_pod_name, :dest_pod_name, :src_container, :dest_container, :src_files_path, :dest_files_path, :src_s3_access_key, :dest_s3_access_key, :src_s3_secret_key, :dest_s3_secret_key, :src_s3_file, :dest_storage_directory, :src_ssh_host, :dest_ssh_host, :src_ssh_port, :dest_ssh_port, :src_ssh_user, :dest_ssh_user, :src_ssh_key, :dest_ssh_key, :src_uri, :dest_uri, :dest_kubeconfig, :src_kubeconfig, :emails) where id=:id`, data)
	return result, err
}

func FindAllMigrations(dest interface{}) error {
	return database.DB.Select(dest, "SELECT * FROM migrations ORDER BY id DESC")
}

func FindMigrationById(dest interface{}, migrationIden interface{}) error {
	return database.DB.Get(dest, "SELECT * FROM migrations WHERE id=$1", migrationIden)
}

func UpdateMigrationDatesAndStatuses(migration *Migration, uf string, pf string) (sql.Result, error) {
	result, err := database.DB.NamedExec(`UPDATE migrations SET (updated_at, status`+uf+`)
	= (:updated_at, :status`+pf+`) where id=:id`, migration)
	return result, err
}

func UpdateMigrationHelper(id int, status string, SetStartedAt bool, SetCompletedAt bool, nullValues bool, output string) {
	o := Migration{
		ID: id,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Status: sql.NullString{
			String: status,
			Valid:  true,
		},
	}

	uf := ""
	pf := ""

	if SetStartedAt {
		uf += ", started_at"
		pf += ", :started_at"
		o.StartedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	if SetCompletedAt {
		uf += ", completed_at"
		pf += ", :completed_at"
		o.CompletedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	if nullValues {
		uf += ", started_at, completed_at"
		pf += ", :started_at, :completed_at"
		o.StartedAt = sql.NullTime{}
		o.CompletedAt = sql.NullTime{}
	}

	if output != "" {
		uf += ", output"
		pf += ", :output"
		o.Output = sql.NullString{
			String: output,
			Valid:  true,
		}
	}

	if _, err := UpdateMigrationDatesAndStatuses(&o, uf, pf); err != nil {
		// TODO handle error
		return
	}
}

func DeleteMigration(migrationIden interface{}) (sql.Result, error) {
	return database.DB.Exec("delete from migrations where id=$1", migrationIden)
}
