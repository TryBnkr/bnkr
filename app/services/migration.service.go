package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/forms"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"
	"github.com/go-chi/chi/v5"
	guuid "github.com/google/uuid"
)

func (m *Repository) GetMigrations(w http.ResponseWriter, r *http.Request) {
	// Get all migrations
	migrations := &[]dal.Migration{}
	if err := dal.FindAllMigrations(migrations); err != nil {
		utils.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["migrations"] = migrations
	render.Template(w, r, "migrations.page.tmpl", &types.TemplateData{
		Data: data,
	})
}

func (m *Repository) GetNewMigration(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := make(map[string]interface{})
	if id != 0 {
		migration := &dal.Migration{}

		if err := dal.FindMigrationById(migration, id); err != nil {
			utils.ServerError(w, err)
			return
		}
		data["id"] = id
		data["values"] = migration
	}

	data["timezones"] = utils.GetTimeZones()

	render.Template(w, r, "migrations.new.page.tmpl", &types.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostNewMigration(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	migrationName := r.Form.Get("migrationName")
	timezone := r.Form.Get("timezone")
	SrcRegion := r.Form.Get("src_region")
	DestRegion := r.Form.Get("dest_region")
	SrcType := r.Form.Get("src_type")
	DestType := r.Form.Get("dest_type")
	SrcBucket := r.Form.Get("src_bucket")
	DestBucket := r.Form.Get("dest_bucket")
	SrcContainer := r.Form.Get("src_container")
	DestContainer := r.Form.Get("dest_container")
	SrcDbName := r.Form.Get("src_db_name")
	DestDbName := r.Form.Get("dest_db_name")
	SrcSshHost := r.Form.Get("src_ssh_host")
	SrcSshPort := r.Form.Get("src_ssh_port")
	SrcSshUser := r.Form.Get("src_ssh_user")
	SrcSshKey := r.Form.Get("src_ssh_key")
	DestSshHost := r.Form.Get("dest_ssh_host")
	DestSshPort := r.Form.Get("dest_ssh_port")
	DestSshUser := r.Form.Get("dest_ssh_user")
	DestSshKey := r.Form.Get("dest_ssh_key")
	SrcURI := r.Form.Get("src_uri")
	DestURI := r.Form.Get("dest_uri")
	SrcDbUser := r.Form.Get("src_db_user")
	DestDbUser := r.Form.Get("dest_db_user")
	DestDbPassword := r.Form.Get("dest_db_password")
	SrcDbPassword := r.Form.Get("src_db_password")
	SrcDbHost := r.Form.Get("src_db_host")
	DestDbHost := r.Form.Get("dest_db_host")
	DestDbPort := r.Form.Get("dest_db_port")
	SrcDbPort := r.Form.Get("src_db_port")
	SrcPodLabel := r.Form.Get("src_pod_label")
	DestPodLabel := r.Form.Get("dest_pod_label")
	DestPodName := r.Form.Get("dest_pod_name")
	SrcPodName := r.Form.Get("src_pod_name")
	SrcFilesPath := r.Form.Get("src_files_path")
	DestFilesPath := r.Form.Get("dest_files_path")
	DestS3AccessKey := r.Form.Get("dest_s3_access_key")
	SrcS3AccessKey := r.Form.Get("src_s3_access_key")
	SrcS3SecretKey := r.Form.Get("src_s3_secret_key")
	DestS3SecretKey := r.Form.Get("dest_s3_secret_key")
	Emails := r.Form.Get("emails")
	DestStorageDirectory := r.Form.Get("dest_storage_directory")
	SrcS3File := r.Form.Get("src_s3_file")

	values := &dal.Migration{
		Name:                 migrationName,
		Timezone:             timezone,
		Emails:               Emails,
		SrcType:              SrcType,
		SrcBucket:            SrcBucket,
		SrcRegion:            SrcRegion,
		SrcDbName:            SrcDbName,
		SrcDbUser:            SrcDbUser,
		SrcDbPassword:        SrcDbPassword,
		SrcDbHost:            SrcDbHost,
		SrcDbPort:            SrcDbPort,
		SrcPodLabel:          SrcPodLabel,
		SrcPodName:           SrcPodName,
		SrcContainer:         SrcContainer,
		SrcFilesPath:         SrcFilesPath,
		SrcS3AccessKey:       SrcS3AccessKey,
		SrcS3SecretKey:       SrcS3SecretKey,
		SrcS3File:            SrcS3File,
		SrcURI:               SrcURI,
		SrcSshHost:           SrcSshHost,
		SrcSshPort:           SrcSshPort,
		SrcSshUser:           SrcSshUser,
		SrcSshKey:            SrcSshKey,
		DestType:             DestType,
		DestBucket:           DestBucket,
		DestRegion:           DestRegion,
		DestDbName:           DestDbName,
		DestDbUser:           DestDbUser,
		DestDbPassword:       DestDbPassword,
		DestDbHost:           DestDbHost,
		DestDbPort:           DestDbPort,
		DestPodLabel:         DestPodLabel,
		DestPodName:          DestPodName,
		DestContainer:        DestContainer,
		DestFilesPath:        DestFilesPath,
		DestS3AccessKey:      DestS3AccessKey,
		DestS3SecretKey:      DestS3SecretKey,
		DestStorageDirectory: DestStorageDirectory,
		DestURI:              DestURI,
		DestSshHost:          DestSshHost,
		DestSshPort:          DestSshPort,
		DestSshUser:          DestSshUser,
		DestSshKey:           DestSshKey,
	}

	args := []string{"migrationName", "timezone", "src_type", "dest_type", "dest_bucket", "src_bucket", "src_s3_access_key", "dest_s3_access_key", "dest_s3_secret_key", "src_s3_secret_key", "dest_region", "src_region"}

	args = append(args, utils.GetRequiredMigTypeFields(SrcType, "src")...)
	args = append(args, utils.GetRequiredMigTypeFields(DestType, "dest")...)

	form.Required(args...)

	data := make(map[string]interface{})
	data["values"] = *values

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	if id != 0 {
		data["id"] = id
	}

	data["timezones"] = utils.GetTimeZones()

	if !form.Valid() {
		render.Template(w, r, "migrations.new.page.tmpl", &types.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if id != 0 {
		values.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		values.ID = id
		if _, err := dal.UpdateMigration(values); err != nil {
			utils.ServerError(w, err)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "Migration updated")
		// Update cron
		// TODO check error of UpdateOrInsertCron
		m.UpdateOrInsertCron(id, "update")
	} else {
		values.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		values.CreatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		values.User = utils.GetUser(w, r)
		if _, err := dal.CreateMigration(values); err != nil {
			utils.ServerError(w, err)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "Migration created")
		// Insert cron
		// TODO check error of UpdateOrInsertCron
		m.UpdateOrInsertCron(values.ID, "create")
	}

	http.Redirect(w, r, "/migrations", http.StatusSeeOther)
}

type MigrationCommon struct {
	MigrationPath string
	MigrationName string
	S3FullPath    string
	S3Path        string
	TmpPath       string
	Msg           types.MailData
	FailedStatus  string
	SuccessStatus string
	StartedAt     time.Time
}

func (m *Repository) PrepareMigration(b *dal.Migration, migrationName string, s3FullPath string) MigrationCommon {
	loc, _ := time.LoadLocation(b.Timezone)
	currentTime := time.Now().In(loc).Format("2006.01.02-150405")
	// Example /bnkr/ad21d8b9-3663-4bfb-8978-30d0ec51a1b8
	dir, _ := filepath.Abs("./bnkr/" + guuid.New().String())

	os.MkdirAll(dir, os.ModePerm)

	if migrationName == "" {
		if b.SrcType == "db" {
			migrationName = b.SrcDbName + "." + currentTime + ".sql.gz"
		} else if b.SrcType == "mongo" {
			migrationName = currentTime + ".gz"
		} else if b.SrcType == "pg" || b.SrcType == "bnkr" {
			migrationName = currentTime + ".psql.gz"
		} else {
			// Files
			migrationName = currentTime + ".tar.gz"
		}
	}

	msg := types.MailData{
		To:   strings.Split(b.Emails, ","),
		From: utils.GetOptionValue("FROM_EMAIL"),
	}

	if b.SrcType == "db" || b.SrcType == "mongo" || b.SrcType == "pg" || b.SrcType == "bnkr" {
		msg.Subject = fmt.Sprintf("Database migration %s failed!", b.Name)
	} else {
		msg.Subject = fmt.Sprintf("Files migration %s failed!", b.Name)
	}

	migrationPath := dir + "/" + migrationName

	s3Path := "/"
	if b.DestStorageDirectory != "" {
		s3Path = "/" + b.DestStorageDirectory + "/"
	}

	if s3FullPath == "" {
		s3FullPath = s3Path + migrationName
	}

	return MigrationCommon{
		MigrationPath: migrationPath,
		MigrationName: migrationName,
		S3FullPath:    s3FullPath,
		S3Path:        s3Path,
		TmpPath:       dir,
		Msg:           msg,
		FailedStatus:  "fail",
		SuccessStatus: "success",
		StartedAt:     time.Now(),
	}
}

func (m *Repository) srcDB(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string
	// Else if direct access is allowed then simply do the dump on Bnkr
	// If the database inside SSH then run dump command on the server using the SSH details then move it to Bnkr
	switch g.SrcAccess {
	// Is the DB in K8S or SSH or we have direct access to it
	case "ssh":
		// Create the SSH key file
		sshKeyPath := c.TmpPath + "/" + "id_rsa"
		err := ioutil.WriteFile(sshKeyPath, []byte(g.SrcSshKey), 0600)
		if err != nil {
			m.App.ErrorLog.Println(err)
			return "", err
		}

		// Create the DB dump on the server
		args := []string{"-i", sshKeyPath, "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", g.SrcSshUser + "@" + g.SrcSshHost, "cd /; mysqldump -h " + g.SrcDbHost + " -u " + g.SrcDbUser + " --port=" + g.SrcDbPort + " -p" + g.SrcDbPassword + " " + g.SrcDbName, "|", "gzip", ">", c.MigrationName}
		cmd := exec.Command("ssh", args...)

		output, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Move the dump DB to Bnkr
		args2 := []string{"-i", sshKeyPath, "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", g.SrcSshUser + "@" + g.SrcSshHost + ":/" + c.MigrationName, c.MigrationName}
		cmd2 := exec.Command("scp", args2...)
		cmd2.Dir = c.TmpPath

		output2, err := utils.CmdExecutor(cmd2)
		if err != nil {
			return "", err
		}

		o = output + `
` + output2

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath := c.TmpPath + "/" + "kubeconfig.yml"
		err := ioutil.WriteFile(kubeconfigPath, []byte(g.SrcKubeconfig), 0600)
		if err != nil {
			m.App.ErrorLog.Println(err)
			return "", err
		}

		// Create MariaDB helper pod
		helperPodName := "bnkr-" + guuid.New().String()
		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--rm", "--restart=Never", "--image", "mariadb:10.5.9-focal", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Dump the DB in the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "cd /; mysqldump -h " + g.SrcDbHost + " -u " + g.SrcDbUser + " --port=" + g.SrcDbPort + " -p" + g.SrcDbPassword + " " + g.SrcDbName, "|", "gzip", ">", c.MigrationName}
		cmd = exec.Command("kubectl", args...)

		output3, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}
		// Move the DB to Bnkr
		args = []string{"cp", "--kubeconfig", kubeconfigPath, helperPodName + ":/" + c.MigrationName, c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output4, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		o = output + `

` + output2 + `

` + output3 + `

` + output4 + `

` + output5
	case "direct":
		// open the out file for writing
		outfile, err := os.Create(c.MigrationPath)
		if err != nil {
			return "", err
		}
		defer outfile.Close()

		args := []string{"-h", g.SrcDbHost, "-u", g.SrcDbUser, "--port=" + g.SrcDbPort, "-p" + g.SrcDbPassword, g.SrcDbName}
		mysqldump := exec.Command("mysqldump", args...)

		mysqldump.Stderr = os.Stderr

		gzip := exec.Command("gzip")
		gzip.Stdout = outfile

		// Get mysqldump's stdout and attach it to gzip's stdin.
		pipe, err := mysqldump.StdoutPipe()
		if err != nil {
			return "", err
		}
		defer pipe.Close()

		gzip.Stdin = pipe

		// Run mysqldump first.
		err2 := mysqldump.Start()
		if err2 != nil {
			return "", err
		}

		err = gzip.Start()
		if err != nil {
			return "", err
		}
		err = gzip.Wait()
		if err != nil {
			return "", err
		}

		err = mysqldump.Wait()
		if err != nil {
			return "", err
		}

		o = "Database dump completed successfully!"

	case "s3":
		// Download the file from S3
		if err := Repo.DownloadFromS3(g.SrcS3AccessKey, g.SrcS3SecretKey, g.SrcBucket, g.SrcRegion, g.SrcS3File, c.MigrationPath); err != nil {
			return "", err
		}

		o = "Database file download completed successfully!"
	}

	return o, nil
}

func (m *Repository) destDB(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	switch g.DestAccess {
	case "ssh":
		// Create the SSH key file
		sshKeyPath := c.TmpPath + "/" + "id_rsa"
		err := ioutil.WriteFile(sshKeyPath, []byte(g.DestSshKey), 0600)
		if err != nil {
			m.App.ErrorLog.Println(err)
			return "", err
		}

		// Move the dump DB to server
		args2 := []string{"-i", sshKeyPath, "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", c.MigrationName, g.DestSshUser + "@" + g.DestSshHost + ":/" + c.MigrationName}
		cmd2 := exec.Command("scp", args2...)
		cmd2.Dir = c.TmpPath

		output, err := utils.CmdExecutor(cmd2)
		if err != nil {
			return "", err
		}

		// Restore the DB dump on the server
		args := []string{"-i", sshKeyPath, "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", g.DestSshUser + "@" + g.DestSshHost, "cd /; gunzip", "<", c.MigrationName, "|", "mysql", "--max_allowed_packet=512M", "-h", g.DestDbHost, "-u", g.DestDbUser, "-p" + g.DestDbPassword, g.DestDbName}
		cmd := exec.Command("ssh", args...)

		output2, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		o = output + `
` + output2

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath := c.TmpPath + "/" + "kubeconfig.yml"
		err := ioutil.WriteFile(kubeconfigPath, []byte(g.DestKubeconfig), 0600)
		if err != nil {
			m.App.ErrorLog.Println(err)
			return "", err
		}

		// Create MariaDB helper pod
		helperPodName := "bnkr-" + guuid.New().String()

		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--rm", "--restart=Never", "--image", "mariadb:10.5.9-focal", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Move the DB dump to pod
		args = []string{"cp", "--kubeconfig", kubeconfigPath, c.MigrationName, helperPodName + ":/" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output3, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Restore the DB on the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "cd /; gunzip", "<", c.MigrationName, "|", "mysql", "--max_allowed_packet=512M", "-h", g.DestDbHost, "-u", g.DestDbUser, "-p" + g.DestDbPassword, g.DestDbName}
		cmd = exec.Command("kubectl", args...)

		output4, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		if err != nil {
			return "", err
		}

		o = output + `

` + output2 + `

` + output3 + `

` + output4 + `

` + output5
	case "direct":
		file, err := os.Open(c.MigrationPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		gunzip := exec.Command("gunzip")
		gunzip.Stdin = file
		args := []string{"-h", g.DestDbHost, "-u", g.DestDbUser, "-p" + g.DestDbPassword, g.DestDbName}
		mysql := exec.Command("mysql", args...)

		// Get gunzip's stdout and attach it to mysql's stdin.
		pipe, err := gunzip.StdoutPipe()
		if err != nil {
			return "", err
		}
		defer pipe.Close()

		mysql.Stdin = pipe

		// Run gunzip first.
		err2 := gunzip.Start()
		if err2 != nil {
			return "", err
		}

		err = mysql.Start()
		if err != nil {
			return "", err
		}
		err = mysql.Wait()
		if err != nil {
			return "", err
		}

		err = gunzip.Wait()
		if err != nil {
			return "", err
		}

		o = "Database restore completed successfully!"

	case "s3":
		if err := Repo.UploadToS3(g.DestS3AccessKey, g.DestS3SecretKey, g.DestBucket, g.DestRegion, c.S3FullPath, c.MigrationPath); err != nil {
			return "", err
		}

		o = "Database upload completed successfully!"
	}

	return o, nil
}

func (m *Repository) migrate(id int) {
	migration := &dal.Migration{}

	if err := dal.FindMigrationById(migration, id); err != nil {
		m.App.ErrorLog.Println(err)
		return
	}

	commons := Repo.PrepareMigration(migration, "", "")

	// Create the backup
	switch migration.SrcType {
	case "db":
		Repo.srcDB(migration, commons)
	case "ssh":

	}

	// Restore the backup
	switch migration.DestType {
	case "db":
		Repo.destDB(migration, commons)
	case "ssh":

	}
}

func (m *Repository) MigrateNow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	go m.migrate(id)

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Migration process queued!",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
