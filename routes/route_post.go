package routes

import (
	"fmt"
	"net/http"

	"github.com/foolin/goview"
	"github.com/go-chi/chi"
)

func PostRouter(r chi.Router) {
	// get list of posts
	r.Get("/", getIndex)

	// get a single post
	r.Get("/{id}", getPost)

	// create a new post
	r.With(AdminOnly).Post("/", createNewPost)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(AdminOnly)
		// update a post
		r.Put("/", updatePost)
		// delete a post
		r.Delete("/", deletePost)

		// comment router
		r.Route("/comments", CommentRouter)
	})
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	err := goview.Render(w, http.StatusOK, "posts/index", goview.M{
		"title": "Index title!",
		"add": func(a int, b int) int {
			return a + b
		},
		"Partials": []string{"posts/index"},
	})
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	return
}

func createNewPost(w http.ResponseWriter, r *http.Request) {
	return
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	return
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	return
}
