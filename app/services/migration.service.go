package services

import (
	"net/http"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/utils"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"
)

func (m *Repository) GetMigrations(w http.ResponseWriter, r *http.Request) {
	// Get all migrations
	migrations := &[]types.NewMigrationDTO{}
	if err := dal.FindAllMigrations(migrations); err != nil {
		utils.ServerError(w, err)
		return
	}

	// Get latest job foreach migration
	var jobs []types.SmallJob
	dal.SelectLatestJobForEachMigration(&jobs)

	ce := make(map[int]time.Time)
	for _, e := range m.App.Cron.Entries() {
		for b, ci := range m.App.CronIds {
			if ci == e.ID {
				ce[b] = e.Next
			}
		}
	}

	data := make(map[string]interface{})
	data["migrations"] = migrations
	data["jobs"] = jobs
	data["nextOcc"] = ce
	render.Template(w, r, "migrations.page.tmpl", &types.TemplateData{
		Data: data,
	})
}
