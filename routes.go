package main

import (
	"embed"
	"net/http"

	"github.com/MohammedAl-Mahdawi/bnkr/app/services"
	"github.com/MohammedAl-Mahdawi/bnkr/config"
	pm "github.com/MohammedAl-Mahdawi/bnkr/utils/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(pm.NoSurf)
	mux.Use(pm.SessionLoad)

	var staticFS = http.FS(staticFiles)
	fileServer := http.FileServer(staticFS)
	mux.Handle("/static/*", fileServer)

	mux.Route("/auth", func(mux chi.Router) {
		mux.Get("/login", services.Repo.GetLogin)
		mux.Post("/login", services.Repo.PostLogin)
		mux.Get("/logout", services.Repo.Logout)
	})

	mux.Route("/users", func(mux chi.Router) {
		mux.Use(pm.Auth)
		mux.Get("/", services.Repo.GetUsers)
		mux.Get("/new", services.Repo.GetNewUser)
		mux.Post("/new", services.Repo.PostNewUser)
		mux.Post("/{id}", services.Repo.PostNewUser)
		mux.Get("/{id}", services.Repo.GetNewUser)
	})

	mux.Route("/options", func(mux chi.Router) {
		mux.Use(pm.Auth)
		mux.Get("/", services.Repo.GetOptions)
		mux.Post("/", services.Repo.PostOptions)
	})

	mux.Route("/jobs", func(mux chi.Router) {
		mux.Use(pm.Auth)
		mux.Get("/backups/{id}", services.Repo.GetJobsByBackup)
	})

	mux.Route("/backups", func(mux chi.Router) {
		mux.Use(pm.Auth)
		mux.Get("/", services.Repo.GetBackups)
		mux.Get("/new", services.Repo.GetNewBackup)
		mux.Post("/new", services.Repo.PostNewBackup)
		mux.Post("/{id}", services.Repo.PostNewBackup)
		mux.Get("/{id}", services.Repo.GetNewBackup)
	})

	mux.Route("/migrations", func(mux chi.Router) {
		mux.Use(pm.Auth)
		mux.Get("/", services.Repo.GetMigrations)
		mux.Get("/new", services.Repo.GetNewMigration)
	})

	mux.Route("/json", func(mux chi.Router) {
		mux.Use(pm.CsrfVerifier)
		mux.Use(pm.Auth)
		mux.Post("/backups/clone/{id}", services.Repo.CloneBackup)
		mux.Delete("/backups/{id}", services.Repo.DeleteBackup)
		mux.Delete("/users/{id}", services.Repo.DeleteUser)
		mux.Post("/jobs/backup/{id}", services.Repo.PostJob)
		mux.Post("/jobs/restore/{id}", services.Repo.RestoreJob)
		mux.Delete("/jobs/{id}", services.Repo.DeleteJob)
		mux.Post("/jobs/running/backups", services.Repo.GetRunningBackups)
		mux.Post("/jobs/download/{jid}/{bid}", services.Repo.DownloadFile)
		mux.Post("/jobs/running/{id}", services.Repo.GetRunningJobs)
	})

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/backups", http.StatusSeeOther)
	})

	return mux
}
