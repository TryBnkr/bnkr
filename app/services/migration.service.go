package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

func (m *Repository) GetRunningMigrations(w http.ResponseWriter, r *http.Request) {
	migrations := &[]dal.Migration{}
	if err := dal.FindAllMigrations(migrations); err != nil {
		utils.ServerError(w, err)
		return
	}

	var migrationsIds []int
	for _, j := range *migrations {
		if j.Status.String == "running" {
			migrationsIds = append(migrationsIds, j.ID)
		}
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Data: migrationsIds,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
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
	SrcAccess := r.Form.Get("src_access")
	DestAccess := r.Form.Get("dest_access")
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
	SrcKubefile := r.Form.Get("src_kubeconfig")
	DestKubefile := r.Form.Get("dest_kubeconfig")

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
		DestAccess:           DestAccess,
		SrcAccess:            SrcAccess,
		SrcKubeconfig:        SrcKubefile,
		DestKubeconfig:       DestKubefile,
	}

	args := []string{"migrationName", "timezone", "src_type", "dest_type", "src_access", "dest_access"}

	args = append(args, utils.GetRequiredMigTypeFields(SrcType, "src")...)
	args = append(args, utils.GetRequiredMigTypeFields(DestType, "dest")...)

	args = append(args, utils.GetRequiredMigAccessFields(values)...)

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
	}

	http.Redirect(w, r, "/migrations", http.StatusSeeOther)
}

func (m *Repository) CloneMigration(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	migration := &dal.Migration{}

	if err := dal.FindMigrationById(migration, id); err != nil {
		utils.ServerError(w, err)
		return
	}

	migration.Name = "Clone of " + migration.Name
	migration.CreatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	migration.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	migration.User = utils.GetUser(w, r)

	if _, err := dal.CreateMigration(migration); err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Clone created!")

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Clone successfully created!",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) DeleteMigration(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	migration := &dal.Migration{}

	if err := dal.FindMigrationById(migration, id); err != nil {
		utils.ServerError(w, err)
		return
	}

	if migration.Status.String == "running" {
		utils.ErrorJSON(w, errors.New("cant delete a running migration"))
		return
	}

	if _, err := dal.DeleteMigration(id); err != nil {
		utils.ErrorJSON(w, errors.New("unable to delete migration"))
		return
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Migration successfully deleted",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
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
	RunningStatus string
	StartedAt     sql.NullTime
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
		RunningStatus: "running",
		StartedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

func (m *Repository) srcDB(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string
	// Else if direct access is allowed then simply do the dump on Bnkr
	// If the database inside SSH then run dump command on the server using the SSH details then move it to Bnkr
	switch g.SrcAccess {
	// Is the DB in K8S or SSH or we have direct access to it
	case "ssh":
		conn, err := utils.GetSshConn(g.SrcSshUser, g.SrcSshHost, g.SrcSshPort, g.SrcSshKey)
		if err != nil {
			return o, err
		}
		defer conn.Close()

		// Create the DB dump on the server
		err = utils.RunSshCommand(conn, "mysqldump --no-tablespaces -h "+g.SrcDbHost+" -u "+g.SrcDbUser+" --port="+g.SrcDbPort+" -p"+g.SrcDbPassword+" "+g.SrcDbName+" | gzip > /"+c.MigrationName)

		if err != nil {
			return o, err
		}

		// Move the dump DB to Bnkr
		_, err = utils.DownloadFileOverSftp(conn, "/"+c.MigrationName, c.TmpPath+"/"+c.MigrationName)

		if err != nil {
			return o, err
		}

		// Cleanup the server
		err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)

		if err != nil {
			return o, err
		}

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.SrcKubeconfig, "kubeconfig.yml")
		if err != nil {
			return "", err
		}

		// DEBUG
		fmt.Println("start the run", kubeconfigPath)

		// Create MariaDB helper pod
		helperPodName := "bnkr-" + guuid.New().String()
		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--restart=Never", "--image", "mariadb:10.5.9-focal", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		o += `
` + output
		if err != nil {
			// DEBUG
			fmt.Println("check the run", err, output)
			return o, err
		}

		// DEBUG
		fmt.Println("start the wait")

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		o += `
` + output2
		if err != nil {
			// DEBUG
			fmt.Println("check the wait", err, output2)
			return o, err
		}

		// DEBUG
		fmt.Println("Dump the DB in the pod")

		// Dump the DB in the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "mysqldump --no-tablespaces -h " + g.SrcDbHost + " -u " + g.SrcDbUser + " --port=" + g.SrcDbPort + " -p" + g.SrcDbPassword + " " + g.SrcDbName + " | gzip > /" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)

		output3, err := utils.CmdExecutor(cmd)
		o += `
` + output3
		if err != nil {
			// DEBUG
			fmt.Println("check the dump", err, output3)
			return o, err
		}

		// DEBUG
		fmt.Println("start the cp")

		// Move the DB to Bnkr
		args = []string{"cp", "--kubeconfig", kubeconfigPath, helperPodName + ":/" + c.MigrationName, c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output4, err := utils.CmdExecutor(cmd)
		o += `
` + output4
		if err != nil {
			// DEBUG
			fmt.Println("check the cp", err, output4)
			return o, err
		}

		// DEBUG
		fmt.Println("start the delete")

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		o += `
` + output5
		if err != nil {
			// DEBUG
			fmt.Println("check the delete", err, output5)
			return o, err
		}

	case "direct":
		// open the out file for writing
		outfile, err := os.Create(c.MigrationPath)
		if err != nil {
			return "", err
		}
		defer outfile.Close()

		args := []string{"--no-tablespaces", "-h", g.SrcDbHost, "-u", g.SrcDbUser, "--port=" + g.SrcDbPort, "-p" + g.SrcDbPassword, g.SrcDbName}
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

		// case "s3":
		// 	// Download the file from S3
		// 	if err := Repo.DownloadFromS3(g.SrcS3AccessKey, g.SrcS3SecretKey, g.SrcBucket, g.SrcRegion, g.SrcS3File, c.MigrationPath); err != nil {
		// 		return "", err
		// 	}

		// 	o = "Database file download completed successfully!"
	}

	return o, nil
}

func (m *Repository) srcPG(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	uri := g.SrcURI
	if g.SrcType == "bnkr" {
		uri = m.App.DbUri
	}
	// Else if direct access is allowed then simply do the dump on Bnkr
	// If the database inside SSH then run dump command on the server using the SSH details then move it to Bnkr
	switch g.SrcAccess {
	// Is the DB in K8S or SSH or we have direct access to it
	case "ssh":
		conn, err := utils.GetSshConn(g.SrcSshUser, g.SrcSshHost, g.SrcSshPort, g.SrcSshKey)
		if err != nil {
			return o, err
		}
		defer conn.Close()

		// Create the DB dump on the server
		err = utils.RunSshCommand(conn, "pg_dump --dbname="+uri+" | gzip > /"+c.MigrationName)
		if err != nil {
			return o, err
		}

		// Move the dump DB to Bnkr
		_, err = utils.DownloadFileOverSftp(conn, "/"+c.MigrationName, c.TmpPath+"/"+c.MigrationName)
		if err != nil {
			return o, err
		}

		// Cleanup the server
		err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)
		if err != nil {
			return o, err
		}

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.SrcKubeconfig, "kubeconfig.yml")
		if err != nil {
			return "", err
		}

		// Create postgres helper pod
		helperPodName := "bnkr-" + guuid.New().String()
		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--restart=Never", "--image", "postgres:13-alpine", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		o += `
` + output
		if err != nil {
			return o, err
		}

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		o += `
` + output2
		if err != nil {
			return o, err
		}

		// Dump the DB in the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "pg_dump --dbname=" + uri + " | gzip > /" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)

		output3, err := utils.CmdExecutor(cmd)
		o += `
` + output3
		if err != nil {
			return o, err
		}

		// Move the DB to Bnkr
		args = []string{"cp", "--kubeconfig", kubeconfigPath, helperPodName + ":/" + c.MigrationName, c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output4, err := utils.CmdExecutor(cmd)
		o += `
` + output4
		if err != nil {
			return o, err
		}

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		o += `
` + output5
		if err != nil {
			return o, err
		}

	case "direct":
		// open the out file for writing
		outfile, err := os.Create(c.MigrationPath)
		if err != nil {
			return "", err
		}
		defer outfile.Close()

		args := []string{"--dbname=" + uri}
		pg_dump := exec.Command("pg_dump", args...)

		pg_dump.Stderr = os.Stderr

		gzip := exec.Command("gzip")
		gzip.Stdout = outfile

		// Get pg_dump's stdout and attach it to gzip's stdin.
		pipe, err := pg_dump.StdoutPipe()
		if err != nil {
			return "", err
		}
		defer pipe.Close()

		gzip.Stdin = pipe

		// Run pg_dump first.
		err2 := pg_dump.Start()
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

		err = pg_dump.Wait()
		if err != nil {
			return "", err
		}

		o = "Database dump completed successfully!"

		// case "s3":
		// 	// Download the file from S3
		// 	if err := Repo.DownloadFromS3(g.SrcS3AccessKey, g.SrcS3SecretKey, g.SrcBucket, g.SrcRegion, g.SrcS3File, c.MigrationPath); err != nil {
		// 		return "", err
		// 	}

		// 	o = "Database file download completed successfully!"
	}

	return o, nil
}

func (m *Repository) srcMongo(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	switch g.SrcAccess {
	case "ssh":
		conn, err := utils.GetSshConn(g.SrcSshUser, g.SrcSshHost, g.SrcSshPort, g.SrcSshKey)
		if err != nil {
			return o, err
		}
		defer conn.Close()

		// Create the DB dump on the server
		err = utils.RunSshCommand(conn, "mongodump --uri="+g.SrcURI+" --gzip --archive=/"+c.MigrationName)

		if err != nil {
			return o, err
		}

		// Move the dump DB to Bnkr
		_, err = utils.DownloadFileOverSftp(conn, "/"+c.MigrationName, c.TmpPath+"/"+c.MigrationName)
		if err != nil {
			return o, err
		}

		// Cleanup the server
		err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)
		if err != nil {
			return o, err
		}

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.SrcKubeconfig, "kubeconfig.yml")
		if err != nil {
			return "", err
		}

		// Create MongoDB helper pod
		helperPodName := "bnkr-" + guuid.New().String()
		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--restart=Never", "--image", "mongo:5.0.3-focal", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		o += `
` + output
		if err != nil {
			return o, err
		}

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		o += `
` + output2
		if err != nil {
			return o, err
		}

		// Dump the DB in the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "mongodump --uri=" + g.SrcURI, "--gzip", "--archive=/" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)

		output3, err := utils.CmdExecutor(cmd)
		o += `
` + output3
		if err != nil {
			return o, err
		}

		// Move the DB to Bnkr
		args = []string{"cp", "--kubeconfig", kubeconfigPath, helperPodName + ":/" + c.MigrationName, c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output4, err := utils.CmdExecutor(cmd)
		o += `
` + output4
		if err != nil {
			return o, err
		}

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		o += `
` + output5
		if err != nil {
			return o, err
		}
	case "direct":
		// Dump the database
		args := []string{"--uri=" + g.SrcURI, "--gzip", "--archive=" + c.MigrationPath}
		mongodump := exec.Command("mongodump", args...)

		output6, err := utils.CmdExecutor(mongodump)
		o += `
` + output6
		if err != nil {
			return o, err
		}

		// case "s3":
		// 	// Download the file from S3
		// 	if err := Repo.DownloadFromS3(g.SrcS3AccessKey, g.SrcS3SecretKey, g.SrcBucket, g.SrcRegion, g.SrcS3File, c.MigrationPath); err != nil {
		// 		return "", err
		// 	}

		// 	o = "Database file download completed successfully!"
	}

	return o, nil
}

func (m *Repository) srcK8SFiles(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	// DEBUG
	fmt.Println("from srcK8SFiles")

	kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.SrcKubeconfig, "kubeconfig.yml")
	if err != nil {
		// DEBUG
		fmt.Println("Error creating kubeconfig: ", err)
		return o, err
	}

	// DEBUG
	fmt.Println("kubeconfigPath: ", kubeconfigPath)

	podName := g.SrcPodName

	var args []string

	if g.SrcType != "pod" {
		// DEBUG
		fmt.Println("Not pod")

		args = []string{"get", "pod", "--kubeconfig", kubeconfigPath, "-l", g.SrcPodLabel, "-o", "jsonpath={.items[0].metadata.name}"}
		podNameBytes, err := exec.Command("kubectl", args...).Output()
		if err != nil {
			// DEBUG
			fmt.Println("Error getting pod name: ", err)
			return o, err
		}

		podName = string(podNameBytes)
		// DEBUG
		fmt.Println("Found pod name: ", podName)
	}

	// DEBUG
	fmt.Println("Start creating tarball...")

	// Create tarball inside the deployment container
	args = []string{"exec", "-c", g.SrcContainer, podName, "--kubeconfig", kubeconfigPath, "--", "tar -czf /" + c.MigrationName + " -C " + g.SrcFilesPath + " ."}
	cmd := exec.Command("kubectl", args...)

	output, err := utils.CmdExecutor(cmd)
	// DEBUG
	fmt.Println("After exec: ", output, err)
	o += `
` + output
	if err != nil {
		// DEBUG
		fmt.Println("Error creating tarball: ", err)
		return o, err
	}

	// DEBUG
	fmt.Println("Start moving the tarball...")

	// Move the tarball to Bnkr
	args = []string{"cp", "--kubeconfig", kubeconfigPath, podName + ":/" + c.MigrationName, c.MigrationName}
	cmd = exec.Command("kubectl", args...)
	cmd.Dir = c.TmpPath

	output2, err := utils.CmdExecutor(cmd)
	// DEBUG
	fmt.Println("After cp: ", output2, err)
	o += `
` + output2
	if err != nil {
		// DEBUG
		fmt.Println("Oh cp error: ", err)
		return o, err
	}

	// DEBUG
	fmt.Println("Start cleanup...")

	// Cleanup, remove the tarball file from the deployment
	args = []string{"exec", "--kubeconfig", kubeconfigPath, "-c", g.SrcContainer, podName, "--", "rm /" + c.MigrationName}
	cmd = exec.Command("kubectl", args...)
	cmd.Dir = c.TmpPath

	output3, err := utils.CmdExecutor(cmd)
	// DEBUG
	fmt.Println("After exec cleanup: ", output3, err)
	o += `
` + output3
	if err != nil {
		// DEBUG
		fmt.Println("Oh cleanup error: ", err)
		return o, err
	}

	// DEBUG
	fmt.Println("We reached the end!", o)

	return o, nil
}

func (m *Repository) srcSSHFiles(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	conn, err := utils.GetSshConn(g.SrcSshUser, g.SrcSshHost, g.SrcSshPort, g.SrcSshKey)
	if err != nil {
		return o, err
	}
	defer conn.Close()

	// Create the tarball on the server
	err = utils.RunSshCommand(conn, "tar -czf /"+c.MigrationName+" -C "+g.SrcFilesPath+" .")

	if err != nil {
		return o, err
	}

	// Move the tarball DB to Bnkr
	_, err = utils.DownloadFileOverSftp(conn, "/"+c.MigrationName, c.TmpPath+"/"+c.MigrationName)

	if err != nil {
		return o, err
	}

	// Cleanup the server
	err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)

	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *Repository) destK8SFiles(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.DestKubeconfig, "dest-kubeconfig.yml")
	if err != nil {
		m.App.ErrorLog.Println(err)
		return "", err
	}

	podName := g.DestPodName

	var args []string

	if g.DestType != "pod" {
		args = []string{"get", "pod", "--kubeconfig", kubeconfigPath, "-l", g.DestPodLabel, "-o", "jsonpath={.items[0].metadata.name}"}
		podNameBytes, err := exec.Command("kubectl", args...).Output()
		if err != nil {
			return "", err
		}

		podName = string(podNameBytes)
	}

	// Copy the tarball to pod
	args = []string{"cp", "--kubeconfig", kubeconfigPath, c.MigrationName, podName + ":/" + c.MigrationName}
	cmd := exec.Command("kubectl", args...)
	cmd.Dir = c.TmpPath

	output, err := utils.CmdExecutor(cmd)
	o += `
` + output
	if err != nil {
		return o, err
	}

	// Extract the tarball
	args = []string{"exec", podName, "--kubeconfig", kubeconfigPath, "--", "tar -xzf /" + c.MigrationName + " -C " + g.DestFilesPath}
	cmd = exec.Command("kubectl", args...)

	output2, err := utils.CmdExecutor(cmd)
	o += `
` + output2
	if err != nil {
		return o, err
	}

	// Cleanup, remove the tarball file from the pod
	args = []string{"exec", "--kubeconfig", kubeconfigPath, "-c", g.DestContainer, podName, "--", "rm /" + c.MigrationName}
	cmd = exec.Command("kubectl", args...)

	output3, err := utils.CmdExecutor(cmd)
	o += `
` + output3
	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *Repository) destSSHFiles(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	// Create the SSH key file
	conn, err := utils.GetSshConn(g.DestSshUser, g.DestSshHost, g.DestSshPort, g.DestSshKey)
	if err != nil {
		return o, err
	}
	defer conn.Close()

	// Move the tarball to server
	_, err = utils.UploadFileOverSftp(conn, c.TmpPath+"/"+c.MigrationName, "/"+c.MigrationName)
	if err != nil {
		return o, err
	}

	// Restore the files on the server
	err = utils.RunSshCommand(conn, "tar -xzf /"+c.MigrationName+" -C "+g.DestFilesPath)

	if err != nil {
		return o, err
	}

	// Cleanup, remove the tarball file from the server
	err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)

	if err != nil {
		return o, err
	}

	return o, nil
}

func (m *Repository) destDB(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	switch g.DestAccess {
	case "ssh":
		conn, err := utils.GetSshConn(g.DestSshUser, g.DestSshHost, g.DestSshPort, g.DestSshKey)
		if err != nil {
			return o, err
		}
		defer conn.Close()

		// Move the dump DB to server
		_, err = utils.UploadFileOverSftp(conn, c.TmpPath+"/"+c.MigrationName, "/"+c.MigrationName)
		if err != nil {
			return o, err
		}

		// Restore the DB dump on the server
		err = utils.RunSshCommand(conn, "gunzip < /"+c.MigrationName+" | mysql --max_allowed_packet=512M -h "+g.DestDbHost+" -u "+g.DestDbUser+" -p"+g.DestDbPassword+" "+g.DestDbName)

		if err != nil {
			return o, err
		}

		// Cleanup the server
		err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)
		if err != nil {
			return o, err
		}

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.DestKubeconfig, "dest-kubeconfig.yml")
		if err != nil {
			return "", err
		}

		// Create MariaDB helper pod
		helperPodName := "bnkr-" + guuid.New().String()

		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--restart=Never", "--image", "mariadb:10.5.9-focal", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		o += `
` + output
		if err != nil {
			return o, err
		}

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		o += `
` + output2
		if err != nil {
			return o, err
		}

		// Move the DB dump to pod
		args = []string{"cp", "--kubeconfig", kubeconfigPath, c.MigrationName, helperPodName + ":/" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output3, err := utils.CmdExecutor(cmd)
		o += `
` + output3
		if err != nil {
			return o, err
		}

		// Restore the DB on the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "gunzip < /" + c.MigrationName + " | mysql", "--max_allowed_packet=512M", "-h", g.DestDbHost, "-u", g.DestDbUser, "-p" + g.DestDbPassword, g.DestDbName}
		cmd = exec.Command("kubectl", args...)

		output4, err := utils.CmdExecutor(cmd)
		o += `
` + output4
		if err != nil {
			return o, err
		}

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		o += `
` + output5
		if err != nil {
			return o, err
		}

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

		// case "s3":
		// 	if err := Repo.UploadToS3(g.DestS3AccessKey, g.DestS3SecretKey, g.DestBucket, g.DestRegion, c.S3FullPath, c.MigrationPath); err != nil {
		// 		return "", err
		// 	}

		// 	o = "Database upload completed successfully!"
	}

	return o, nil
}

func (m *Repository) destPG(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	uri := g.DestURI
	if g.DestType == "bnkr" {
		uri = m.App.DbUri
	}

	switch g.DestAccess {
	case "ssh":
		conn, err := utils.GetSshConn(g.DestSshUser, g.DestSshHost, g.DestSshPort, g.DestSshKey)
		if err != nil {
			return o, err
		}
		defer conn.Close()

		// Move the dump DB to server
		_, err = utils.UploadFileOverSftp(conn, c.TmpPath+"/"+c.MigrationName, "/"+c.MigrationName)
		if err != nil {
			return o, err
		}

		// Restore the DB dump on the server
		err = utils.RunSshCommand(conn, "gunzip < /"+c.MigrationName+" | psql "+uri)
		if err != nil {
			return o, err
		}

		// Cleanup the server
		err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)
		if err != nil {
			return o, err
		}

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.DestKubeconfig, "dest-kubeconfig.yml")
		if err != nil {
			return "", err
		}

		// Create postgres helper pod
		helperPodName := "bnkr-" + guuid.New().String()
		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--restart=Never", "--image", "postgres:13-alpine", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		o += `
` + output
		if err != nil {
			return o, err
		}

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		o += `
` + output2
		if err != nil {
			return o, err
		}

		// Move the DB dump to pod
		args = []string{"cp", "--kubeconfig", kubeconfigPath, c.MigrationName, helperPodName + ":/" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output3, err := utils.CmdExecutor(cmd)
		o += `
` + output3
		if err != nil {
			return o, err
		}

		// Restore the DB on the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "gunzip < /" + c.MigrationName + " | psql " + uri}
		cmd = exec.Command("kubectl", args...)

		output4, err := utils.CmdExecutor(cmd)
		o += `
` + output4
		if err != nil {
			return o, err
		}

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		o += `
` + output5
		if err != nil {
			return o, err
		}

	case "direct":
		file, err := os.Open(c.MigrationPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		gunzip := exec.Command("gunzip")
		gunzip.Stdin = file

		args := []string{uri}
		psql := exec.Command("psql", args...)

		// Get gunzip's stdout and attach it to psql's stdin.
		pipe, err := gunzip.StdoutPipe()
		if err != nil {
			return "", err
		}
		defer pipe.Close()

		psql.Stdin = pipe

		// Run gunzip first.
		err2 := gunzip.Start()
		if err2 != nil {
			return "", err
		}

		err = psql.Start()
		if err != nil {
			return "", err
		}
		err = psql.Wait()
		if err != nil {
			return "", err
		}

		err = gunzip.Wait()
		if err != nil {
			return "", err
		}

		o = "Database restore completed successfully!"

		// case "s3":
		// 	if err := Repo.UploadToS3(g.DestS3AccessKey, g.DestS3SecretKey, g.DestBucket, g.DestRegion, c.S3FullPath, c.MigrationPath); err != nil {
		// 		return "", err
		// 	}

		// 	o = "Database upload completed successfully!"
	}

	return o, nil
}

func (m *Repository) destMongo(g *dal.Migration, c MigrationCommon) (string, error) {
	var o string

	switch g.DestAccess {
	case "ssh":
		conn, err := utils.GetSshConn(g.DestSshUser, g.DestSshHost, g.DestSshPort, g.DestSshKey)
		if err != nil {
			return o, err
		}
		defer conn.Close()

		// Move the dump DB to server
		_, err = utils.UploadFileOverSftp(conn, c.TmpPath+"/"+c.MigrationName, "/"+c.MigrationName)
		if err != nil {
			return o, err
		}

		// Restore the DB dump on the server
		err = utils.RunSshCommand(conn, "mongorestore --uri="+g.DestURI+" --gzip --drop --archive=/"+c.MigrationName)
		if err != nil {
			return o, err
		}

		// Cleanup the server
		err = utils.RunSshCommand(conn, "rm /"+c.MigrationName)
		if err != nil {
			return o, err
		}

	case "k8s":
		// Create Kubeconfig file
		kubeconfigPath, err := utils.CreateKubeconfigFile(c.TmpPath, g.DestKubeconfig, "dest-kubeconfig.yml")
		if err != nil {
			return "", err
		}

		// Create MongoDB helper pod
		helperPodName := "bnkr-" + guuid.New().String()
		args := []string{"run", helperPodName, "--kubeconfig", kubeconfigPath, "--restart=Never", "--image", "mongo:5.0.3-focal", "--command", "--", "sleep", "infinity"}
		cmd := exec.Command("kubectl", args...)

		output, err := utils.CmdExecutor(cmd)
		o += `
` + output
		if err != nil {
			return o, err
		}

		// Wait for the pod to be ready
		args = []string{"wait", "--kubeconfig", kubeconfigPath, "--for=condition=ready", "pod", helperPodName}
		cmd = exec.Command("kubectl", args...)

		output2, err := utils.CmdExecutor(cmd)
		o += `
` + output2
		if err != nil {
			return o, err
		}

		// Move the DB dump to pod
		args = []string{"cp", "--kubeconfig", kubeconfigPath, c.MigrationName, helperPodName + ":/" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)
		cmd.Dir = c.TmpPath

		output3, err := utils.CmdExecutor(cmd)
		o += `
` + output3
		if err != nil {
			return o, err
		}

		// Restore the DB on the pod
		args = []string{"exec", helperPodName, "--kubeconfig", kubeconfigPath, "--", "sh", "-c", "mongorestore", "--uri=" + g.DestURI, "--gzip", "--drop", "--archive=/" + c.MigrationName}
		cmd = exec.Command("kubectl", args...)

		output4, err := utils.CmdExecutor(cmd)
		o += `
` + output4
		if err != nil {
			return o, err
		}

		// Delete the helper pod
		args = []string{"delete", "--kubeconfig", kubeconfigPath, "pod", helperPodName, "--ignore-not-found"}
		cmd = exec.Command("kubectl", args...)

		output5, err := utils.CmdExecutor(cmd)
		o += `
` + output5
		if err != nil {
			return o, err
		}

	case "direct":
		args := []string{"--uri=" + g.DestURI, "--gzip", "--drop", "--archive=" + c.MigrationPath}
		mongodump := exec.Command("mongorestore", args...)

		output6, err := utils.CmdExecutor(mongodump)
		o = output6
		if err != nil {
			return "", err
		}

		// case "s3":
		// 	if err := Repo.UploadToS3(g.DestS3AccessKey, g.DestS3SecretKey, g.DestBucket, g.DestRegion, c.S3FullPath, c.MigrationPath); err != nil {
		// 		return "", err
		// 	}

		// 	o = "Database upload completed successfully!"
	}

	return o, nil
}

func (m *Repository) TerminateMigration(message string, status string, commons *MigrationCommon, b *dal.Migration, sendMail bool, output string) error {
	if status != "success" {
		commons.Msg.Content = message
		m.App.MailChan <- commons.Msg
	} else {
		if sendMail {
			commons.Msg.Subject = "Migration Succeeded!"
			commons.Msg.Content = fmt.Sprintf("The migration process for the migration %s completed successfully", b.Name)
			m.App.MailChan <- commons.Msg
		}
	}

	dal.UpdateMigrationHelper(b.ID, status, false, true, false, output)
	return os.RemoveAll(commons.TmpPath)
}

func (m *Repository) migrate(id int, migration *dal.Migration) error {
	commons := Repo.PrepareMigration(migration, "", "")

	dal.UpdateMigrationHelper(migration.ID, commons.RunningStatus, true, false, false, "")

	var srcErr, destErr error
	var srcOut, destOut string

	// Create the backup
	switch migration.SrcType {
	// MySQL/MariaDB Database
	case "db":
		srcOut, srcErr = Repo.srcDB(migration, commons)

	case "pg":
	case "bnkr":
		srcOut, srcErr = Repo.srcPG(migration, commons)

	case "mongo":
		srcOut, srcErr = Repo.srcMongo(migration, commons)

	// Files In SSH
	case "ssh":
		srcOut, srcErr = Repo.srcSSHFiles(migration, commons)

	// Files In Deployment or StatefulSet
	case "object":
	case "pod":
		// DEBUG
		fmt.Println("K8S related")
		srcOut, srcErr = Repo.srcK8SFiles(migration, commons)

	case "s3":
		if err := Repo.DownloadFromS3(migration.SrcS3AccessKey, migration.SrcS3SecretKey, migration.SrcBucket, migration.SrcRegion, migration.SrcS3File, commons.MigrationPath); err != nil {
			srcOut, srcErr = "", err
		}
	}

	if srcErr != nil {
		Repo.TerminateMigration("Error while processing the source!", commons.FailedStatus, &commons, migration, true, srcOut+`
`+srcErr.Error())
		return srcErr
	}

	// Restore the backup
	switch migration.DestType {
	// MySQL/MariaDB Database
	case "db":
		destOut, destErr = Repo.destDB(migration, commons)

	case "pg":
	case "bnkr":
		destOut, destErr = Repo.destPG(migration, commons)

	case "mongo":
		destOut, destErr = Repo.destMongo(migration, commons)

	// Files In SSH
	case "ssh":
		destOut, destErr = Repo.destSSHFiles(migration, commons)

	// Files In Deployment or StatefulSet
	case "object":
	case "pod":
		destOut, destErr = Repo.destK8SFiles(migration, commons)

	// S3
	case "s3":
		if err := Repo.UploadToS3(migration.DestS3AccessKey, migration.DestS3SecretKey, migration.DestBucket, migration.DestRegion, commons.S3FullPath, commons.MigrationPath); err != nil {
			destOut, destErr = "", err
		}
	}

	if destErr != nil {
		Repo.TerminateMigration("Error while processing the destination!", commons.FailedStatus, &commons, migration, true, destOut+`
`+destErr.Error())
		return destErr
	}

	Repo.TerminateMigration("", commons.SuccessStatus, &commons, migration, true, srcOut+`
`+destOut)

	return nil
}

func (m *Repository) MigrateNow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	migration := &dal.Migration{}

	if err := dal.FindMigrationById(migration, id); err != nil {
		m.App.ErrorLog.Println(err)
		utils.ErrorJSON(w, err)
		return
	}

	if migration.Status.String == "running" {
		utils.ErrorJSON(w, errors.New("already running"))
		return
	}

	go m.migrate(id, migration)

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Migration run!",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
