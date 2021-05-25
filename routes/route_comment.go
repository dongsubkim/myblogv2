package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dongsubkim/myblogv2/data"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
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
		password := r.PostFormValue("comment[password]")
		hashedPassword, err := data.PasswordByComment(chi.URLParam(r, "commentId"))
		if err != nil {
			error_message(w, r, fmt.Sprintf("Fail to parse request body: %v!", err))
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			error_message(w, r, fmt.Sprint("Password not matched!"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func createComment(w http.ResponseWriter, r *http.Request) {
	postUuid := chi.URLParam(r, "uuid")
	log.Println(r.PostFormValue("comment[username]"), r.PostFormValue("comment[password]"), r.PostFormValue("comment[body]"), postUuid)
	err := data.CreateComment(r.PostFormValue("comment[username]"), r.PostFormValue("comment[password]"), r.PostFormValue("comment[body]"), postUuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Error during creating a comment: %v!", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", postUuid), 302)
}

func updateComment(w http.ResponseWriter, r *http.Request) {
	err := data.UpdateComment(chi.URLParam(r, "commentId"), r.PostFormValue("comment[body]"))
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update a comment: %v!", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", chi.URLParam(r, "uuid")), 302)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	err := data.DeleteComment(chi.URLParam(r, "commentId"))
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to delete a comment: %v!", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", chi.URLParam(r, "uuid")), 302)
}
