package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/dal"
	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/forms"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/paginator"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/password"
	"github.com/MohammedAl-Mahdawi/bnkr/utils/render"

	"github.com/go-chi/chi/v5"
)

// GetUsers returns the users list
func (m *Repository) GetUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("p"))

	var usersCount int
	if err := dal.Count(&usersCount, "users", ""); err != nil {
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
		TotalCount:  usersCount,
	}

	users := &[]types.NewUserDTO{}
	if err := dal.FindUsers(users, "created_at desc", p); err != nil {
		utils.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["users"] = users
	data["pagination"] = p
	render.Template(w, r, "users.page.tmpl", &types.TemplateData{
		Data: data,
	})
}

func (m *Repository) GetNewUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := make(map[string]interface{})
	if id != 0 {
		user := &types.NewUserDTO{}

		if err := dal.FindUserById(user, id); err != nil {
			utils.ServerError(w, err)
			return
		}
		data["id"] = id
		values := types.NewUserForm{
			Name:  user.Name,
			Email: user.Email,
		}
		data["values"] = values
	}

	render.Template(w, r, "users.new.page.tmpl", &types.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if _, err := dal.DeleteUser(id); err != nil {
		utils.ServerError(w, errors.New("unable to delete user"))
		return
	}

	out, _ := json.Marshal(&types.MsgResponse{
		Message: "User successfully deleted",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) PostNewUser(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	userName := r.Form.Get("userName")
	email := r.Form.Get("email")
	pass := r.Form.Get("password")
	passwordConfirmation := r.Form.Get("passwordConfirmation")

	values := types.NewUserForm{
		Name:                 userName,
		Password:             pass,
		PasswordConfirmation: passwordConfirmation,
		Email:                email,
	}

	form.Required("email")
	form.IsEmail("email")
	form.PasswordConfirmation(pass, passwordConfirmation)

	data := make(map[string]interface{})
	data["values"] = values

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	if id != 0 {
		data["id"] = id
	}

	if !form.Valid() {
		render.Template(w, r, "users.new.page.tmpl", &types.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	d := &dal.User{
		Name:  values.Name,
		Email: values.Email,
	}

	if values.Password != "" {
		d.Password = password.Generate(values.Password)
	}

	if id != 0 {
		d.ID = id
		d.UpdatedAt = time.Now()
		if _, err := dal.UpdateUser(d); err != nil {
			utils.ServerError(w, err)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "User updated")
	} else {
		d.UpdatedAt = time.Now()
		d.CreatedAt = time.Now()
		if _, err := dal.CreateUser(d); err != nil {
			if utils.IsDuplicateKeyError(err) {
				form.Errors.Add("email", "User already exist!")
				render.Template(w, r, "users.new.page.tmpl", &types.TemplateData{
					Form: form,
					Data: data,
				})
			} else {
				utils.ServerError(w, err)
			}
			return
		}
		m.App.Session.Put(r.Context(), "flash", "User created")
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
