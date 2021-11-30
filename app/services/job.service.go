package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/forms"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/paginator"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"

	guuid "github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/go-chi/chi/v5"
)

func (m *Repository) PostJob(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	m.App.Queue <- types.NewQueueDTO{
		ID:      id,
		Process: "backup",
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Backup process queued!",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) RestoreJob(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	m.App.Queue <- types.NewQueueDTO{
		ID:      id,
		Process: "restore",
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Backup restoreation queued!",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) GetS3PreSignedURL(b *types.NewBackupDTO, j *types.NewJobDTO) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(b.Region),
		Credentials: credentials.NewStaticCredentials(b.S3AccessKey, b.S3SecretKey, ""),
	})

	if err != nil {
		return "", err
	}

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(j.File),
	})
	urlStr, err := req.Presign(3 * time.Hour)

	if err != nil {
		return "", err
	}

	return urlStr, nil
}

func (m *Repository) DownloadFile(w http.ResponseWriter, r *http.Request) {
	jid, _ := strconv.Atoi(chi.URLParam(r, "jid"))
	bid, _ := strconv.Atoi(chi.URLParam(r, "bid"))

	job := &types.NewJobDTO{}
	if err := dal.FindJobById(job, jid); err != nil {
		utils.ServerError(w, err)
		return
	}

	backup := &types.NewBackupDTO{}
	if err := dal.FindBackupById(backup, bid); err != nil {
		utils.ServerError(w, err)
		return
	}

	urlStr, err := Repo.GetS3PreSignedURL(backup, job)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Data: urlStr,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) GetRunningJobs(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	jobs := &[]types.NewJobDTO{}
	if err := dal.FindJobsIDByBackup(jobs, id, "created_at desc"); err != nil {
		utils.ServerError(w, err)
		return
	}

	var jobsIds []int
	for _, j := range *jobs {
		jobsIds = append(jobsIds, j.ID)
	}

	queues := &[]dal.Queue{}
	if len(jobsIds) > 0 {
		if err := dal.FindQueuesByObjectsIdsAndType(queues, jobsIds, "job", "created_at desc"); err != nil {
			utils.ServerError(w, err)
			return
		}
	}

	jobsIds = []int{}
	for _, j := range *queues {
		jobsIds = append(jobsIds, j.Object)
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Data: jobsIds,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) GetJobsByBackup(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	page, _ := strconv.Atoi(r.URL.Query().Get("p"))

	backup := &types.NewBackupDTO{}
	if err := dal.FindBackupById(backup, id); err != nil {
		utils.ServerError(w, err)
		return
	}

	var jobsCount int
	if err := dal.Count(&jobsCount, "jobs", " WHERE backup="+strconv.Itoa(id)); err != nil {
		utils.ServerError(w, err)
		return
	}

	cp := 1
	if page > 1 {
		cp = page
	}

	p := &paginator.Paginator{
		CurrentPage: cp,
		PerPage:     20,
		TotalCount:  jobsCount,
	}

	jobs := &[]types.NewJobDTO{}
	if err := dal.FindJobsByBackup(jobs, id, "created_at desc", p); err != nil {
		utils.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["jobs"] = jobs
	data["backup"] = backup
	data["pagination"] = p
	render.Template(w, r, "jobs.page.tmpl", &types.TemplateData{
		Data: data,
	})
}

func (m *Repository) GetNewJob(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := make(map[string]interface{})
	if id != 0 {
		job := &types.NewJobDTO{}

		if err := dal.FindJobById(job, id); err != nil {
			utils.ServerError(w, err)
			return
		}
		data["id"] = id
		data["values"] = job
	}

	render.Template(w, r, "jobs.new.page.tmpl", &types.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	job := &types.NewJobDTO{}
	if err := dal.FindJobById(job, id); err != nil {
		utils.ServerError(w, err)
		return
	}

	backup := &types.NewBackupDTO{}
	if err := dal.FindBackupById(backup, job.Backup); err != nil {
		utils.ServerError(w, err)
		return
	}

	if job.File != "" {
		if err := Repo.DeleteS3Object(backup, job); err != nil {
			utils.ServerError(w, err)
			return
		}
	}

	if _, err := dal.DeleteJob(id); err != nil {
		utils.ServerError(w, errors.New("unable to delete job"))
		return
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Job successfully deleted",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) PostNewJob(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/jobs", http.StatusSeeOther)
}

func (m *Repository) DeleteS3Object(b *types.NewBackupDTO, j *types.NewJobDTO) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(b.Region),
		Credentials: credentials.NewStaticCredentials(b.S3AccessKey, b.S3SecretKey, ""),
	})

	if err != nil {
		return err
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Delete the item
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(b.Bucket), Key: aws.String(j.File)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(j.File),
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Repository) DeleteExtraBackups(b *types.NewBackupDTO) error {
	jobs := &[]types.NewJobDTO{}

	if err := dal.FindAllJobsByBackup(jobs, b.ID, "created_at asc"); err != nil {
		return err
	}

	currentBackups := len(*jobs)
	retention := b.Retention

	if currentBackups >= retention {
		backupsMustDeleted := (*jobs)[0:(currentBackups - (retention - 1))]
		for _, j := range backupsMustDeleted {
			// Delete the file & the db record
			if j.File != "" {
				err := Repo.DeleteS3Object(b, &j)
				if err != nil {
					return err
				}
			}
			if _, err := dal.DeleteJob(j.ID); err != nil {
				return errors.New("unable to delete job")
			}
		}
	}

	return nil
}

func (m *Repository) UploadToS3(S3AccessKey string, S3SecretKey string, Bucket string, Region string, S3FullPath string, SrcFilePath string) error {
	file, err := os.Open(SrcFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewStaticCredentials(S3AccessKey, S3SecretKey, ""),
	})

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(Bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(S3FullPath),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		return err
	}

	return nil
}

func (m *Repository) DownloadFromS3(S3AccessKey string, S3SecretKey string, Bucket string, Region string, S3FullPath string, DestFilePath string) error {
	file, err := os.Create(DestFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewStaticCredentials(S3AccessKey, S3SecretKey, ""),
	})

	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(Bucket),
			Key:    aws.String(S3FullPath),
		})
	if err != nil {
		return err
	}

	return nil
}

type BackupCommon struct {
	BackupPath    string
	BackupName    string
	S3FullPath    string
	S3Path        string
	Dir           string
	Msg           types.MailData
	FailedStatus  string
	SuccessStatus string
	StartedAt     time.Time
}

func (m *Repository) PrepareBackup(b *types.NewBackupDTO, backupName string, s3FullPath string) BackupCommon {
	loc, _ := time.LoadLocation(b.Timezone)
	currentTime := time.Now().In(loc).Format("2006.01.02-150405")
	// Example ./bnkr/ad21d8b9-3663-4bfb-8978-30d0ec51a1b8
	dir := "./bnkr/" + guuid.New().String()

	os.MkdirAll(dir, os.ModePerm)

	if backupName == "" {
		if b.Type == "db" {
			backupName = b.DbName + "." + currentTime + ".sql.gz"
		} else if b.Type == "mongo" {
			backupName = currentTime + ".gz"
		} else if b.Type == "pg" || b.Type == "bnkr" {
			backupName = currentTime + ".psql.gz"
		} else {
			// Files
			backupName = currentTime + ".tar.gz"
		}
	}

	msg := types.MailData{
		To:   strings.Split(b.Emails, ","),
		From: utils.GetOptionValue("FROM_EMAIL"),
	}

	if b.Type == "db" || b.Type == "mongo" || b.Type == "pg" || b.Type == "bnkr" {
		msg.Subject = fmt.Sprintf("Database backup %s failed!", b.Name)
	} else {
		msg.Subject = fmt.Sprintf("Files backup %s failed!", b.Name)
	}

	backupPath := dir + "/" + backupName

	s3Path := "/"
	if b.StorageDirectory != "" {
		s3Path = "/" + b.StorageDirectory + "/"
	}

	if s3FullPath == "" {
		s3FullPath = s3Path + backupName
	}

	return BackupCommon{
		BackupPath:    backupPath,
		BackupName:    backupName,
		S3FullPath:    s3FullPath,
		S3Path:        s3Path,
		Dir:           dir,
		Msg:           msg,
		FailedStatus:  "fail",
		SuccessStatus: "success",
		StartedAt:     time.Now(),
	}
}

func (m *Repository) TerminateBackup(message string, status string, commons *BackupCommon, b *types.NewBackupDTO, sendMail bool) (*dal.Job, error) {
	if status != "success" {
		commons.Msg.Content = message
		m.App.MailChan <- commons.Msg
	} else {
		if sendMail {
			commons.Msg.Subject = "Backup Succeeded!"
			commons.Msg.Content = fmt.Sprintf("The backup process for the backup %s completed successfully", b.Name)
			m.App.MailChan <- commons.Msg
		}
	}

	os.RemoveAll(commons.Dir)
	return Repo.SaveJob(commons, status, b)
}

func (m *Repository) TerminateRestore(message string, status string, commons *BackupCommon, b *types.NewBackupDTO, err error) error {
	if status != "success" {
		commons.Msg.Subject = "Restoration failed!"
		commons.Msg.Content = message
		m.App.MailChan <- commons.Msg
	} else {
		commons.Msg.Subject = "Restoration Succeeded!"
		commons.Msg.Content = fmt.Sprintf("The restoration process for the backup %s completed successfully", b.Name)
		m.App.MailChan <- commons.Msg
	}

	os.RemoveAll(commons.Dir)
	return err
}

func (m *Repository) FilesRestore(b *types.NewBackupDTO, j *types.NewJobDTO) error {
	commons := Repo.PrepareBackup(b, "", j.File)

	if err := Repo.DownloadFromS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateRestore("Can't download from S3", commons.FailedStatus, &commons, b, err)
	}

	podName := b.PodName

	var args []string

	if b.Type != "pod" {
		args = []string{"get", "pod", "-l", b.PodLabel, "-o", "jsonpath={.items[0].metadata.name}"}
		podNameBytes, err := exec.Command("kubectl", args...).Output()
		if err != nil {
			return Repo.TerminateRestore("Can't get pod name!", commons.FailedStatus, &commons, b, err)
		}
		podName = string(podNameBytes)
	}

	// Copy the tarball to current container
	args = []string{"cp", commons.BackupPath, podName + ":/" + commons.BackupName, "-c", b.Container}
	kubectlCp := exec.Command("kubectl", args...)

	if _, err := kubectlCp.Output(); err != nil {
		return Repo.TerminateRestore("Can't copy the tarball to current container", commons.FailedStatus, &commons, b, err)
	}

	// Unzip
	// TODO some commads miss the container like this one
	args = []string{"exec", podName, "--", "sh", "-c", "cd / ; tar -xzf " + commons.BackupName + " -C " + b.FilesPath}
	unzip := exec.Command("kubectl", args...)

	if _, err := unzip.Output(); err != nil {
		return Repo.TerminateRestore("Unzip error!", commons.FailedStatus, &commons, b, err)
	}

	// Cleanup, remove the tarball file from the deployment
	args = []string{"exec", podName, "--", "sh", "-c", "cd / ; rm " + commons.BackupName}
	cleanup := exec.Command("kubectl", args...)

	if _, err := cleanup.Output(); err != nil {
		return Repo.TerminateRestore("Can't clean up!", commons.FailedStatus, &commons, b, err)
	}

	return Repo.TerminateRestore("", commons.SuccessStatus, &commons, b, nil)
}

func (m *Repository) DbRestore(b *types.NewBackupDTO, j *types.NewJobDTO) error {
	commons := Repo.PrepareBackup(b, "", j.File)

	if err := Repo.DownloadFromS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateRestore("Can't download from S3", commons.FailedStatus, &commons, b, err)
	}

	file, err := os.Open(commons.BackupPath)
	if err != nil {
		return Repo.TerminateRestore("Can't open file", commons.FailedStatus, &commons, b, err)
	}
	defer file.Close()

	gunzip := exec.Command("gunzip")
	gunzip.Stdin = file
	args := []string{"-h", b.DbHost, "-u", b.DbUser, "-p" + b.DbPassword, b.DbName}
	mysql := exec.Command("mysql", args...)

	// Get gunzip's stdout and attach it to mysql's stdin.
	pipe, err := gunzip.StdoutPipe()
	if err != nil {
		return Repo.TerminateRestore("Unzip pip error", commons.FailedStatus, &commons, b, err)
	}
	defer pipe.Close()

	mysql.Stdin = pipe

	// Run gunzip first.
	err2 := gunzip.Start()
	if err2 != nil {
		return Repo.TerminateRestore("Can't start unziping!", commons.FailedStatus, &commons, b, err2)
	}

	err = mysql.Start()
	if err != nil {
		return Repo.TerminateRestore("Can't start mysql command", commons.FailedStatus, &commons, b, err)
	}
	err = mysql.Wait()
	if err != nil {
		return Repo.TerminateRestore("Error while executing mysql command", commons.FailedStatus, &commons, b, err)
	}

	err = gunzip.Wait()
	if err != nil {
		return Repo.TerminateRestore("Error while unzipping", commons.FailedStatus, &commons, b, err)
	}

	return Repo.TerminateRestore("", commons.SuccessStatus, &commons, b, nil)
}

func (m *Repository) PgRestore(b *types.NewBackupDTO, j *types.NewJobDTO) error {
	commons := Repo.PrepareBackup(b, "", j.File)

	if err := Repo.DownloadFromS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateRestore("Can't download from S3", commons.FailedStatus, &commons, b, err)
	}

	file, err := os.Open(commons.BackupPath)
	if err != nil {
		return Repo.TerminateRestore("Can't open file", commons.FailedStatus, &commons, b, err)
	}
	defer file.Close()

	gunzip := exec.Command("gunzip")
	gunzip.Stdin = file

	uri := b.URI.String
	if b.Type == "bnkr" {
		uri = m.App.DbUri
	}

	args := []string{uri}
	psql := exec.Command("psql", args...)

	// Get gunzip's stdout and attach it to psql's stdin.
	pipe, err := gunzip.StdoutPipe()
	if err != nil {
		return Repo.TerminateRestore("Unzip pip error", commons.FailedStatus, &commons, b, err)
	}
	defer pipe.Close()

	psql.Stdin = pipe

	// Run gunzip first.
	err2 := gunzip.Start()
	if err2 != nil {
		return Repo.TerminateRestore("Can't start unziping!", commons.FailedStatus, &commons, b, err2)
	}

	err = psql.Start()
	if err != nil {
		return Repo.TerminateRestore("Can't start psql command", commons.FailedStatus, &commons, b, err)
	}
	err = psql.Wait()
	if err != nil {
		return Repo.TerminateRestore("Error while executing psql command", commons.FailedStatus, &commons, b, err)
	}

	err = gunzip.Wait()
	if err != nil {
		return Repo.TerminateRestore("Error while unzipping", commons.FailedStatus, &commons, b, err)
	}

	return Repo.TerminateRestore("", commons.SuccessStatus, &commons, b, nil)
}

func (m *Repository) MongoDbRestore(b *types.NewBackupDTO, j *types.NewJobDTO) error {
	commons := Repo.PrepareBackup(b, "", j.File)

	if err := Repo.DownloadFromS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateRestore("Can't download from S3", commons.FailedStatus, &commons, b, err)
	}

	// Restore the database
	args := []string{"--uri=" + b.URI.String, "--gzip", "--drop", "--archive=" + commons.BackupPath}
	mongodump := exec.Command("mongorestore", args...)

	if _, err := mongodump.Output(); err != nil {
		return Repo.TerminateRestore("Can't excute mongorestore command", commons.FailedStatus, &commons, b, err)
	}

	return Repo.TerminateRestore("", commons.SuccessStatus, &commons, b, nil)
}

func (m *Repository) DbBackup(b *types.NewBackupDTO, sendMail bool) (*dal.Job, error) {
	commons := Repo.PrepareBackup(b, "", "")
	// check the retention number and remove if necessary
	err := Repo.DeleteExtraBackups(b)
	if err != nil {
		return Repo.TerminateBackup("Cant delete extra backups!", commons.FailedStatus, &commons, b, sendMail)
	}

	// open the out file for writing
	outfile, err := os.Create(commons.BackupPath)
	if err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to create %s file!", commons.BackupPath), commons.FailedStatus, &commons, b, sendMail)
	}
	defer outfile.Close()

	args := []string{"--max_allowed_packet=512M", "--single-transaction", "--quick", "-h", b.DbHost, "-u", b.DbUser, "--port=" + b.DbPort, "-p" + b.DbPassword, b.DbName}
	mysqldump := exec.Command("mysqldump", args...)

	mysqldump.Stderr = os.Stderr

	gzip := exec.Command("gzip")
	gzip.Stdout = outfile

	// Get mysqldump's stdout and attach it to gzip's stdin.
	pipe, err := mysqldump.StdoutPipe()
	if err != nil {
		return Repo.TerminateBackup("Failed to get mysqldump pipe!", commons.FailedStatus, &commons, b, sendMail)
	}
	defer pipe.Close()

	gzip.Stdin = pipe

	// Run mysqldump first.
	err2 := mysqldump.Start()
	if err2 != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to start command: %s", strings.Replace("mysqldump "+strings.Join(args, " "), b.DbPassword, "******", 1)), commons.FailedStatus, &commons, b, sendMail)
	}

	err = gzip.Start()
	if err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to start command: %s", "gzip"), commons.FailedStatus, &commons, b, sendMail)
	}
	err = gzip.Wait()
	if err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to excute command: %s", "gzip"), commons.FailedStatus, &commons, b, sendMail)
	}

	err = mysqldump.Wait()
	if err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to excute command: %s", strings.Replace("mysqldump "+strings.Join(args, " "), b.DbPassword, "******", 1)), commons.FailedStatus, &commons, b, sendMail)
	}

	if err := Repo.UploadToS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateBackup("Cant upload to S3", commons.FailedStatus, &commons, b, sendMail)
	}

	return Repo.TerminateBackup("", commons.SuccessStatus, &commons, b, sendMail)
}

func (m *Repository) PgBackup(b *types.NewBackupDTO, sendMail bool) (*dal.Job, error) {
	commons := Repo.PrepareBackup(b, "", "")
	// check the retention number and remove if necessary
	err := Repo.DeleteExtraBackups(b)
	if err != nil {
		return Repo.TerminateBackup("Cant delete extra backups!", commons.FailedStatus, &commons, b, sendMail)
	}

	// open the out file for writing
	outfile, err := os.Create(commons.BackupPath)
	if err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to create %s file!", commons.BackupPath), commons.FailedStatus, &commons, b, sendMail)
	}
	defer outfile.Close()

	uri := b.URI.String
	if b.Type == "bnkr" {
		uri = m.App.DbUri
	}

	args := []string{"-c", "--dbname=" + uri}
	pg_dump := exec.Command("pg_dump", args...)

	pg_dump.Stderr = os.Stderr

	gzip := exec.Command("gzip")
	gzip.Stdout = outfile

	// Get pg_dump's stdout and attach it to gzip's stdin.
	pipe, err := pg_dump.StdoutPipe()
	if err != nil {
		return Repo.TerminateBackup("Failed to get pg_dump pipe!", commons.FailedStatus, &commons, b, sendMail)
	}
	defer pipe.Close()

	gzip.Stdin = pipe

	// Run pg_dump first.
	err2 := pg_dump.Start()
	if err2 != nil {
		return Repo.TerminateBackup("Failed to start pg_dump command!", commons.FailedStatus, &commons, b, sendMail)
	}

	err = gzip.Start()
	if err != nil {
		return Repo.TerminateBackup("Failed to start gzip command!", commons.FailedStatus, &commons, b, sendMail)
	}
	err = gzip.Wait()
	if err != nil {
		return Repo.TerminateBackup("Failed to excute gzip command!", commons.FailedStatus, &commons, b, sendMail)
	}

	err = pg_dump.Wait()
	if err != nil {
		return Repo.TerminateBackup("Failed to excute pg_dump command!", commons.FailedStatus, &commons, b, sendMail)
	}

	if err := Repo.UploadToS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateBackup("Cant upload to S3", commons.FailedStatus, &commons, b, sendMail)
	}

	return Repo.TerminateBackup("", commons.SuccessStatus, &commons, b, sendMail)
}

func (m *Repository) MongoDbBackup(b *types.NewBackupDTO, sendMail bool) (*dal.Job, error) {
	commons := Repo.PrepareBackup(b, "", "")
	// check the retention number and remove if necessary
	err := Repo.DeleteExtraBackups(b)
	if err != nil {
		return Repo.TerminateBackup("Cant delete extra backups!", commons.FailedStatus, &commons, b, sendMail)
	}

	// Dump the database
	args := []string{"--uri=" + b.URI.String, "--gzip", "--archive=" + commons.BackupPath}
	mongodump := exec.Command("mongodump", args...)

	if _, err := mongodump.Output(); err != nil {
		return Repo.TerminateBackup("Failed to execute mongodump command!", commons.FailedStatus, &commons, b, sendMail)
	}

	// Upload to S3
	if err := Repo.UploadToS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateBackup("Cant upload to S3", commons.FailedStatus, &commons, b, sendMail)
	}

	return Repo.TerminateBackup("", commons.SuccessStatus, &commons, b, sendMail)
}

func (m *Repository) SaveJob(commons *BackupCommon, status string, b *types.NewBackupDTO) (*dal.Job, error) {
	o := &dal.Job{
		CreatedAt: commons.StartedAt,
		UpdatedAt: commons.StartedAt,
		File:      commons.S3FullPath,
		CompletedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Status: status,
		Backup: b.ID,
	}

	if _, err := dal.CreateJob(o); err != nil {
		return nil, err
	}

	return o, nil
}

func (m *Repository) FilesBackup(b *types.NewBackupDTO, sendMail bool) (*dal.Job, error) {
	commons := Repo.PrepareBackup(b, "", "")
	// check the retention number and remove if necessary
	err := Repo.DeleteExtraBackups(b)
	if err != nil {
		return Repo.TerminateBackup("Cant delete extra backups!", commons.FailedStatus, &commons, b, sendMail)
	}

	podName := b.PodName

	var args []string

	if b.Type != "pod" {
		args = []string{"get", "pod", "-l", b.PodLabel, "-o", "jsonpath={.items[0].metadata.name}"}
		podNameBytes, err := exec.Command("kubectl", args...).Output()
		if err != nil {
			return Repo.TerminateBackup(fmt.Sprintf("Failed to execute command: %s", "kubectl "+strings.Join(args, " ")), commons.FailedStatus, &commons, b, sendMail)
		}

		podName = string(podNameBytes)
	}

	// Create tarball inside the deployment container
	args = []string{"exec", "-c", b.Container, podName, "--", "sh", "-c", "cd / ; tar -czf " + commons.BackupName + " -C " + b.FilesPath + " ."}
	tarball := exec.Command("kubectl", args...)

	if _, err := tarball.Output(); err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to execute command: %s", "kubectl "+strings.Join(args, " ")), commons.FailedStatus, &commons, b, sendMail)
	}

	// Copy the tarball to current container
	args = []string{"cp", podName + ":/" + commons.BackupName, commons.BackupPath, "-c", b.Container}
	copyToCont := exec.Command("kubectl", args...)

	if _, err := copyToCont.Output(); err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to execute command: %s", "kubectl "+strings.Join(args, " ")), commons.FailedStatus, &commons, b, sendMail)
	}

	// Cleanup, remove the tarball file from the deployment
	args = []string{"exec", "-c", b.Container, podName, "--", "sh", "-c", "cd / ; rm " + commons.BackupName}
	cleanup := exec.Command("kubectl", args...)

	cleanup.Stderr = os.Stdout
	if _, err := cleanup.Output(); err != nil {
		return Repo.TerminateBackup(fmt.Sprintf("Failed to execute command: %s", "kubectl "+strings.Join(args, " ")), commons.FailedStatus, &commons, b, sendMail)
	}

	// Upload to S3
	if err := Repo.UploadToS3(b.S3AccessKey, b.S3SecretKey, b.Bucket, b.Region, commons.S3FullPath, commons.BackupPath); err != nil {
		return Repo.TerminateBackup("Cant upload to S3", commons.FailedStatus, &commons, b, sendMail)
	}

	return Repo.TerminateBackup("", commons.SuccessStatus, &commons, b, sendMail)
}

func (m *Repository) CreateNewJob(b *types.NewBackupDTO, sendMail bool) (*dal.Job, error) {
	if b.Type == "db" {
		return Repo.DbBackup(b, sendMail)
	} else if b.Type == "mongo" {
		return Repo.MongoDbBackup(b, sendMail)
	} else if b.Type == "pg" || b.Type == "bnkr" {
		return Repo.PgBackup(b, sendMail)
	} else {
		return Repo.FilesBackup(b, sendMail)
	}
}

func (m *Repository) RestoreBackup(b *types.NewBackupDTO, j *types.NewJobDTO) error {
	if b.Type == "db" {
		return Repo.DbRestore(b, j)
	} else if b.Type == "mongo" {
		return Repo.MongoDbRestore(b, j)
	} else if b.Type == "pg" || b.Type == "bnkr" {
		return Repo.PgRestore(b, j)
	} else {
		return Repo.FilesRestore(b, j)
	}
}
