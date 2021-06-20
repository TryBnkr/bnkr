package render

import (
	"bytes"
	"errors"
	"fmt"

	at "github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/config"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"

	"github.com/justinas/nosurf"

	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"humanDate":          HumanDate,
	"timezonedDate":      TimezonedDate,
	"humanTimezonedDate": HumanTimezonedDate,
	"dayName":            DayName,
	"monthName":          MonthName,
	"formatDate":         FormatDate,
	"iterate":            Iterate,
	"humanFrequency":     HumanFrequency,
	"getBackupJob":       GetBackupJob,
	"getThemes":          GetThemes,
	"add":                Add,
}

var app *config.AppConfig
var pathToTemplates = "./app/templates"

func Add(a, b int) int {
	return a + b
}

// Iterate returns a slice of ints, starting at 1, going to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format(time.ANSIC)
}

func TimezonedDate(t time.Time, tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		fmt.Println(err)
	}
	return t.In(loc)
}

func HumanTimezonedDate(t time.Time, tz string) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		fmt.Println(err)
	}
	return t.In(loc).Format(time.ANSIC)
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

func DayName(d int) string {
	return time.Weekday(d).String()
}

func MonthName(d int) string {
	return time.Month(d).String()
}

func HumanFrequency(b *at.NewBackupDTO) string {
	cron := ""

	switch b.Frequency {
	case "@yearly":
		cron = fmt.Sprintf("%d %s of each year at %s", b.DayOfMonth, MonthName(b.Month), b.Time)
	case "@monthly":
		cron = fmt.Sprintf("%d of each month at %s", b.DayOfMonth, b.Time)
	case "@weekly":
		cron = fmt.Sprintf("Every %s at %s", DayName(b.DayOfWeek), b.Time)
	case "@daily":
		cron = fmt.Sprintf("Daily at %s", b.Time)
	case "@hourly":
		cron = "Hourly"
	}

	return cron
}

func GetBackupJob(b *at.NewBackupDTO, jobs []at.SmallJob) at.SmallJob {
	for _, j := range jobs {
		if j.Backup == int(b.ID) {
			return j
		}
	}

	return at.SmallJob{
		Backup:    0,
		Status:    "",
		CreatedAt: time.Time{},
	}
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *at.TemplateData, r *http.Request) *at.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Theme = utils.GetOptionValue("THEME")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
		userId, _ := app.Session.Get(r.Context(), "user_id").(uint)
		userName, _ := app.Session.Get(r.Context(), "user_name").(string)
		td.UserId = userId
		td.UserName = userName
	}
	return td
}

func GetThemes() []string {
	t := []string{
		"default",
		"slate",
		"solar",
		"sketchy",
		"cerulean",
		"flatly",
		"journal",
		"materia",
		"minty",
		"morph",
		"pulse",
		"quartz",
		"superhero",
		"vapor",
		"united",
		"yeti",
		"sandstone",
		"simplex",
		"lux",
		"darkly",
		"litera",
		"cyborg",
		"cosmo",
		"spacelab",
		"zephyr",
	}

	return t
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *at.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// this is just used for testing, so that we rebuild
		// the cache on every request
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
