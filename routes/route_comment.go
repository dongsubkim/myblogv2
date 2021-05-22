package routes

import (
	"net/http"

	"github.com/go-chi/chi"
)

func CommentRouter(r chi.Router) {
	// create a new comment
	r.Post("/", createComment)

	// update a comment
	r.Route("/{commentId}", func(r chi.Router) {
		r.Use(commentAuthorOnly)
		r.Put("/", updateComment)
		r.Delete("/", deleteComment)
	})
}

func commentAuthorOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func createComment(w http.ResponseWriter, r *http.Request) {
	return
}

func updateComment(w http.ResponseWriter, r *http.Request) {
	return
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	return
}
