package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/config"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/forms"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/password"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"

	"gorm.io/gorm"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) GetLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &types.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	pass := r.Form.Get("password")

	values := types.LoginDTO{
		Email:    email,
		Password: pass,
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	data := make(map[string]interface{})
	data["values"] = values

	if !form.Valid() {
		render.Template(w, r, "login.page.html", &types.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Try to log the user in
	u := &types.UserResponse{}

	err2 := dal.FindUserByEmail(u, email).Error

	if errors.Is(err2, gorm.ErrRecordNotFound) {
		m.App.Session.Put(r.Context(), "error", "The email address or password you entered is incorrect, please try again.")
		render.Template(w, r, "login.page.html", &types.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if err2 := password.Verify(u.Password, pass); err2 != nil {
		m.App.Session.Put(r.Context(), "error", "The email address or password you entered is incorrect, please try again.")
		render.Template(w, r, "login.page.html", &types.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "user_id", u.ID)
	m.App.Session.Put(r.Context(), "user_name", u.Name)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/backups", http.StatusSeeOther)
}
