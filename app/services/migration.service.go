package services

import (
	"net/http"
	"strconv"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/forms"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"
	"github.com/go-chi/chi/v5"
)

func (m *Repository) GetMigrations(w http.ResponseWriter, r *http.Request) {
	// Get all migrations
	migrations := &[]types.NewMigrationDTO{}
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
		migration := &types.NewMigrationDTO{}

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
