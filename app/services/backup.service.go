package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/forms"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"

	"github.com/go-chi/chi/v5"
)

// GetBackups returns the backups list
func (m *Repository) GetBackups(w http.ResponseWriter, r *http.Request) {
	// Get all backups
	backups := &[]types.NewBackupDTO{}
	if err := dal.FindAllBackups(backups); err != nil {
		utils.ServerError(w, err)
		return
	}

	// Get latest job foreach backup
	var jobs []types.SmallJob
	dal.SelectLatestJobForEachBackup(&jobs)

	data := make(map[string]interface{})
	data["backups"] = backups
	data["jobs"] = jobs
	render.Template(w, r, "backups.page.html", &types.TemplateData{
		Data: data,
	})
}

func (m *Repository) CloneBackup(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	backup := &types.NewBackupDTO{}

	if err := dal.FindBackupById(backup, id); err != nil {
		utils.ServerError(w, err)
		return
	}

	nb := dal.Backup(*backup)

	nb.Name = "Clone of " + nb.Name
	nb.CreatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	nb.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if _, err := dal.CreateBackup(&nb); err != nil {
		utils.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Clone created!")
	// Insert cron
	// TODO check error of UpdateOrInsertCron
	m.UpdateOrInsertCron(nb.ID, "create")

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Clone successfully created!",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) GetNewBackup(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := make(map[string]interface{})
	if id != 0 {
		backup := &types.NewBackupDTO{}

		if err := dal.FindBackupById(backup, id); err != nil {
			utils.ServerError(w, err)
			return
		}
		data["id"] = id
		data["values"] = backup
	}

	data["timezones"] = utils.GetTimeZones()
	data["times"] = utils.GetTimes()

	render.Template(w, r, "backups.new.page.html", &types.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) DeleteBackup(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if _, err := dal.DeleteBackup(id); err != nil {
		utils.ServerError(w, errors.New("unable to delete backup"))
		return
	}

	// Delete cron
	m.App.Cron.Remove(m.App.CronIds[id])
	delete(m.App.CronIds, id)

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "Backup successfully deleted",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) UpdateOrInsertCron(id int, typ string) error {
	if typ == "update" {
		m.App.Cron.Remove(m.App.CronIds[id])
		delete(m.App.CronIds, id)
	}

	backup := &types.NewBackupDTO{}
	if err := dal.FindBackupById(backup, id); err != nil {
		return err
	}

	cron := render.ConstructCron(backup)
	cronId, _ := m.App.Cron.AddFunc("CRON_TZ="+backup.Timezone+" "+cron, func() { m.CreateNewJob(backup, false) })

	m.App.CronIds[backup.ID] = cronId

	return nil
}

func (m *Repository) PostNewBackup(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	backupName := r.Form.Get("backupName")
	frequency := r.Form.Get("frequency")
	timezone := r.Form.Get("timezone")
	backupTime := r.Form.Get("time")
	region := r.Form.Get("region")
	customFrequency := r.Form.Get("customFrequency")
	backupType := r.Form.Get("type")
	bucket := r.Form.Get("bucket")
	container := r.Form.Get("container")
	dbName := r.Form.Get("dbName")
	URI := sql.NullString{
		String: r.Form.Get("uri"),
		Valid:  true,
	}
	dbUser := r.Form.Get("dbUser")
	dbPassword := r.Form.Get("dbPassword")
	dbHost := r.Form.Get("dbHost")
	dbPort := r.Form.Get("dbPort")
	podLabel := r.Form.Get("podLabel")
	podName := r.Form.Get("podName")
	filesPath := r.Form.Get("filesPath")
	accessKey := r.Form.Get("accessKey")
	secretKey := r.Form.Get("secretKey")
	dayOfWeek, _ := strconv.Atoi(r.Form.Get("dayOfWeek"))
	dayOfMonth, _ := strconv.Atoi(r.Form.Get("dayOfMonth"))
	month, _ := strconv.Atoi(r.Form.Get("month"))
	notificationEmail := r.Form.Get("notificationEmail")
	storageDirectory := r.Form.Get("storageDirectory")
	backupRetention, err := strconv.Atoi(r.Form.Get("backupRetention"))
	if err != nil {
		form.Errors.Add("backupRetention", "Invalid data type")
	}

	values := types.NewBackupDTO{
		Name:             backupName,
		Frequency:        frequency,
		CustomFrequency:  customFrequency,
		Timezone:         timezone,
		Time:             backupTime,
		Region:           region,
		Type:             backupType,
		Bucket:           bucket,
		Container:        container,
		DbName:           dbName,
		DbUser:           dbUser,
		DbPassword:       dbPassword,
		DbHost:           dbHost,
		DbPort:           dbPort,
		PodLabel:         podLabel,
		PodName:          podName,
		FilesPath:        filesPath,
		DayOfWeek:        dayOfWeek,
		DayOfMonth:       dayOfMonth,
		Month:            month,
		URI:              URI,
		S3AccessKey:      accessKey,
		S3SecretKey:      secretKey,
		StorageDirectory: storageDirectory,
		Retention:        backupRetention,
		Emails:           notificationEmail,
	}

	args := []string{"backupName", "frequency", "timezone", "type", "bucket", "accessKey", "secretKey", "region"}

	if backupType == "db" {
		args = append(args, "dbName", "dbUser", "dbPassword", "dbHost", "dbPort")
	} else if backupType == "object" {
		args = append(args, "podLabel", "filesPath", "container")
	} else if backupType == "mongo" {
		args = append(args, "uri")
	} else {
		args = append(args, "filesPath", "container", "podName")
	}

	if frequency == "custom" {
		args = append(args, "customFrequency")
	} else if frequency == "@daily" {
		args = append(args, "time")
	} else if frequency == "@weekly" {
		args = append(args, "time", "dayOfWeek")
	} else if frequency == "@monthly" {
		args = append(args, "time", "dayOfMonth")
	} else if frequency == "@yearly" {
		args = append(args, "month", "dayOfMonth", "time")
	}

	form.Required(args...)

	if frequency == "custom" {
		form.IsCron("customFrequency")
	}

	data := make(map[string]interface{})
	data["values"] = values

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	if id != 0 {
		data["id"] = id
	}

	data["timezones"] = utils.GetTimeZones()
	data["times"] = utils.GetTimes()

	if !form.Valid() {
		render.Template(w, r, "backups.new.page.html", &types.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	d := &dal.Backup{
		Name:             values.Name,
		Frequency:        values.Frequency,
		CustomFrequency:  values.CustomFrequency,
		Timezone:         values.Timezone,
		Time:             values.Time,
		Type:             values.Type,
		Bucket:           values.Bucket,
		DbName:           values.DbName,
		DbUser:           values.DbUser,
		DbPassword:       values.DbPassword,
		DbHost:           values.DbHost,
		DbPort:           values.DbPort,
		PodLabel:         values.PodLabel,
		PodName:          values.PodName,
		Container:        values.Container,
		DayOfWeek:        values.DayOfWeek,
		DayOfMonth:       values.DayOfMonth,
		Month:            values.Month,
		FilesPath:        values.FilesPath,
		S3AccessKey:      values.S3AccessKey,
		S3SecretKey:      values.S3SecretKey,
		Region:           values.Region,
		URI:              values.URI,
		StorageDirectory: values.StorageDirectory,
		Retention:        values.Retention,
		Emails:           values.Emails,
		User:             utils.GetUser(w, r),
	}

	if id != 0 {
		d.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		d.ID = id
		if _, err := dal.UpdateBackup(d); err != nil {
			utils.ServerError(w, err)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "Backup updated")
		// Update cron
		// TODO check error of UpdateOrInsertCron
		m.UpdateOrInsertCron(id, "update")
	} else {
		d.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		d.CreatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if _, err := dal.CreateBackup(d); err != nil {
			utils.ServerError(w, err)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "Backup created")
		// Insert cron
		// TODO check error of UpdateOrInsertCron
		m.UpdateOrInsertCron(d.ID, "create")
	}

	http.Redirect(w, r, "/backups", http.StatusSeeOther)
}
