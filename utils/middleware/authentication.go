package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/MohammedAl-Mahdawi/bnkr/config"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"

	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewMiddleware(a *config.AppConfig) {
	app = a
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	p, _ := strconv.ParseBool(os.Getenv("PRODUCTION"))

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   p,
		SameSite: http.SameSiteLaxMode,
	})

	csrfHandler.ExemptRegexp("/json/*")

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsAuthenticated(r) {
			app.Session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CsrfVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrf := r.Header.Get("csrf-token")
		if !nosurf.VerifyToken(nosurf.Token(r), csrf) {
			http.Error(w, "Incorrect CSRF token!", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
