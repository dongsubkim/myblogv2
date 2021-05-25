package routes

import (
	"fmt"
	"net/http"

	"github.com/dongsubkim/myblogv2/data"
	"github.com/foolin/goview"
	"github.com/go-chi/chi"
)

var CategoryNavbar map[string]int

type PostRequest struct {
	Title    string
	Category string
	Content  string
}

func init() {
	CategoryNavbar, _ = data.UpdateCategory()
}

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

		r.Get("/edit", renderEdit)
		// Comment router
		r.Route("/comments", CommentRouter)
	})
}

// Get all posts
func getIndex(w http.ResponseWriter, r *http.Request) {
	var posts []data.Post
	var err error
	if category, ok := r.URL.Query()["category"]; ok {
		posts, err = data.PostsByCategory(category[0])
	} else {
		posts, err = data.Posts()
	}
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get posts: %v!", err))
		return
	}
	err = goview.Render(w, http.StatusOK, "posts/index", goview.M{
		"Posts":          &posts,
		"CategoryNavbar": &CategoryNavbar,
		"Partials":       []string{"posts/index"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
		return
	}
}

// get a single post
func getPost(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	post, err := data.PostByUUID(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get a post: %v!", err))
		return
	}
	comments, err := post.Comments()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get post's comments: %v!", err))
		return
	}
	err = goview.Render(w, http.StatusOK, "posts/show", goview.M{
		"Post":           &post,
		"PostUuid":       func() string { return post.Uuid },
		"Comments":       comments,
		"CategoryNavbar": &CategoryNavbar,
		"Partials":       []string{"posts/show"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render show page error: %v!", err))
		return
	}
}

// create a post
func createPost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10485760)
	title := r.FormValue("title")
	category := r.FormValue("category")
	content := r.FormValue("content")
	uuid, err := data.CreatePost(title, category, content)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to create a post: %v!", err))
		return
	}
	CategoryNavbar, err = data.UpdateCategory()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update category: %v!", err))
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", uuid), 302)
}

// update a post
func updatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10485760)
	uuid := chi.URLParam(r, "uuid")
	title := r.PostFormValue("title")
	category := r.PostFormValue("category")
	content := r.PostFormValue("content")
	err := data.UpdatePost(uuid, title, category, content)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update the post: %v!", err))
		return
	}
	CategoryNavbar, err = data.UpdateCategory()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update category: %v!", err))
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", uuid), 302)
}

// delete a post
func deletePost(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	err := data.DeletePost(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to delete the post: %v!", err))
		return
	}
	CategoryNavbar, err = data.UpdateCategory()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update category: %v!", err))
		return
	}
	http.Redirect(w, r, "/post", 302)
}

// render edit form
func renderEdit(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	post, err := data.PostByUUID(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get a post: %v!", err))
		return
	}
	err = goview.Render(w, http.StatusOK, "posts/edit", goview.M{
		"post":     &post,
		"Partials": []string{"posts/edit"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
		return
	}
}
