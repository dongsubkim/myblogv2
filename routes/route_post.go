package routes

import (
	"fmt"
	"net/http"

	"github.com/dongsubkim/myblogv2/data"
	"github.com/foolin/goview"
	"github.com/go-chi/chi"
)

func PostRouter(r chi.Router) {
	// get list of posts
	r.Get("/", getIndex)

	// create a new post
	r.With(AdminOnly).Post("/", createPost)

	r.Route("/{uuid}", func(r chi.Router) {
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

// Get all posts
func getIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := data.Posts()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get posts: %v!", err))
	} else {
		err = goview.Render(w, http.StatusOK, "posts/index", goview.M{
			"title":    "Index title!",
			"posts":    posts,
			"Partials": []string{"posts/index"},
		})
		if err != nil {
			error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
			// fmt.Fprintf(w, "Render index error: %v!", err)
		}
	}
}

// get a single post
func getPost(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	post, err := data.PostByUUID(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get a post: %v!", err))
	}
	err = goview.Render(w, http.StatusOK, "posts/show", goview.M{
		"post":         post,
		"lastModified": post.CreatedAtDate(),
		"Partials":     []string{"posts/show"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render show page error: %v!", err))
		// fmt.Fprintf(w, "Render show page error: %v!", err)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uuid := r.FormValue("uuid")
	category := []string{r.FormValue("category")}
	content := r.FormValue("content")
	err := data.CreatePost(uuid, category, content)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to create a post: %v!", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", uuid), 302)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uuid := r.FormValue("uuid")
	category := []string{r.FormValue("category")}
	content := r.FormValue("content")
	err := data.UpdatePost(uuid, category, content)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update the post: %v!", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", uuid), 302)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	err := data.DeletePost(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to delete the post: %v!", err))
	}
	http.Redirect(w, r, "/post", 302)
}
