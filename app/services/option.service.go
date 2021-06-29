package services

import (
	"net/http"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/config/database"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/forms"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"

	"gorm.io/gorm/clause"
)

// GetOptions returns the options list
func (m *Repository) GetOptions(w http.ResponseWriter, r *http.Request) {
	options := &[]types.NewOptionDTO{}
	if err := dal.FindAllOptions(&options); err != nil {
		utils.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	values := make(map[string]string)

	for _, option := range *options {
		values[option.Name] = option.Value
	}

	data["values"] = values
	render.Template(w, r, "options.page.html", &types.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostOptions(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	BusinessName := r.Form.Get("BusinessName")
	Theme := r.Form.Get("Theme")
	FromEmail := r.Form.Get("FromEmail")
	SMTPHost := r.Form.Get("SMTPHost")
	SMTPPort := r.Form.Get("SMTPPort")
	SMTPUsername := r.Form.Get("SMTPUsername")
	SMTPPassword := r.Form.Get("SMTPPassword")
	SMTPSecurity := r.Form.Get("SMTPSecurity")

	values := map[string]string{
		"BUSINESS_NAME": BusinessName,
		"THEME": Theme,
		"FROM_EMAIL":    FromEmail,
		"SMTP_HOST":     SMTPHost,
		"SMTP_PORT":     SMTPPort,
		"SMTP_USERNAME": SMTPUsername,
		"SMTP_PASSWORD": SMTPPassword,
		"SMTP_SECURITY": SMTPSecurity,
	}

	data := make(map[string]interface{})
	data["values"] = values

	if !form.Valid() {
		render.Template(w, r, "options.page.html", &types.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	for k, v := range values {
		d := &dal.Option{
			Name:  k,
			Value: v,
		}
		database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&d)
	}

	m.App.Session.Put(r.Context(), "flash", "Options updated!")
	render.Template(w, r, "options.page.html", &types.TemplateData{
		Form: form,
		Data: data,
	})
}
