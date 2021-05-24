package routes

import (
	"fmt"
	"net/http"

	"github.com/foolin/goview"
	"github.com/go-chi/chi"
)

// A completely separate router for administrator routes
func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", adminIndex)
	r.Get("/new", renderNew)
	// r.Get("/accounts", adminListAccounts)
	return r
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//   ctx := r.Context()
		// //   perm, ok := ctx.Value("acl.permission").(YourPermissionType)
		//   if !ok || !perm.IsAdmin() {
		// 	http.Error(w, http.StatusText(403), 403)
		// 	return
		//   }
		next.ServeHTTP(w, r)
	})
}

func renderNew(w http.ResponseWriter, r *http.Request) {
	err := goview.Render(w, http.StatusOK, "posts/new", goview.M{
		"Partials": []string{"posts/new"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
	}
}

func adminIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Amdin index page"))
}
