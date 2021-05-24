package main

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/dongsubkim/myblogv2/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func setHttpMethod(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if m, ok := r.URL.Query()["_method"]; ok && len(m) == 1 {
				r.Method = m[0]
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(setHttpMethod)
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("views", "favicon.ico"))
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/post", 302)
		})
	})

	r.Route("/post", routes.PostRouter)

	r.Mount("/admin", routes.AdminRouter())

	fmt.Println("Serving on Port 3000")
	http.ListenAndServe(":3000", r)
}
