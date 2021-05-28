package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dongsubkim/myblogv2/data"
	"github.com/foolin/goview"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", loginForm)
	r.Post("/", authenticate)
	r.Route("/register", func(r chi.Router) {
		r.Use(CheckAdminCount)
		r.Get("/", registerForm)
		r.Post("/", createAdmin)
	})
	r.With(AdminOnly).Get("/new", newPostForm)
	r.Get("/logout", logout)
	return r
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("_cookie")
		if err != http.ErrNoCookie {
			session := data.Session{Uuid: cookie.Value}
			valid, err := session.Check()
			if err == nil && valid {
				next.ServeHTTP(w, r)
				return
			}
		}
		log.Println("Admin info not found or expired")
		http.Redirect(w, r, "/admin", http.StatusFound)
	})
}

func CheckAdminCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := data.AdminCount(); err == nil && c == 0 {
			next.ServeHTTP(w, r)
		}
		http.Redirect(w, r, "/post", http.StatusFound)
	})
}

func authorized(r *http.Request) bool {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		valid, err := session.Check()
		if err == nil && valid {
			return true
		}
	}
	return false
}

func newPostForm(w http.ResponseWriter, r *http.Request) {
	err := goview.Render(w, http.StatusOK, "posts/new", goview.M{
		"Authorized":     authorized(r),
		"CategoryNavbar": &CategoryNavbar,
		"Partials":       []string{"posts/new"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
	}
}

func registerForm(w http.ResponseWriter, r *http.Request) {
	err := goview.Render(w, http.StatusOK, "admin/register", goview.M{
		"Authorized":     authorized(r),
		"CategoryNavbar": &CategoryNavbar,
		"Partials":       []string{"admin/register"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
	}
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	if authorized(r) {
		http.Redirect(w, r, "/admin/new", http.StatusFound)
	}
	err := goview.Render(w, http.StatusOK, "admin/login", goview.M{
		"Authorized":     authorized(r),
		"CategoryNavbar": &CategoryNavbar,
		"Partials":       []string{"admin/login"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
	}
}

// Authenticate the admin given the email and password
func authenticate(w http.ResponseWriter, r *http.Request) {
	admin, err := data.AdminByEmail(r.PostFormValue("email"))
	if err != nil {
		error_message(w, r, "Cannot find admin")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(r.PostFormValue("password"))); err == nil {
		session, err := admin.CreateSession()
		if err != nil {
			error_message(w, r, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/admin", http.StatusFound)
	}

}

// Create the admin account
func createAdmin(w http.ResponseWriter, r *http.Request) {
	admin := data.Admin{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	if err := admin.Create(); err != nil {
		error_message(w, r, "Cannot create user")
	}
	http.Redirect(w, r, "/post", http.StatusFound)
}

// GET /logout
// Logs the admin out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/post", http.StatusFound)
}
