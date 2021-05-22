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

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("views", "favicon.ico"))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to my blog. This is index page which shows list of posts"))
	})

	r.Route("/posts", routes.PostRouter)

	r.Mount("/admin", routes.AdminRouter())

	fmt.Println("Serving on Port 3000")
	http.ListenAndServe(":3000", r)
}
