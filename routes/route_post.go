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

	// create a new post
	r.With(AdminOnly).Post("/", createNewPost)

	r.Route("/{id}", func(r chi.Router) {
		// get a single post
		r.Get("/", getPost)
		// update a post
		r.With(AdminOnly).Put("/", updatePost)
		// delete a post
		r.With(AdminOnly).Delete("/", deletePost)

		// Comment router
		r.Route("/commnets", CommentRouter)
	})
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	err := goview.Render(w, http.StatusOK, "posts/index", goview.M{
		"title":    "Index title!",
		"Partials": []string{"posts/index"},
	})
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := goview.Render(w, http.StatusOK, "posts/show", goview.M{
		"id":       id,
		"Partials": []string{"posts/show"},
	})
	if err != nil {
		fmt.Fprintf(w, "Render show page error: %v!", err)
	}
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
