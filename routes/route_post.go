package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dongsubkim/myblogv2/data"
	"github.com/foolin/goview"
	"github.com/go-chi/chi"
)

var CategoryNavbar map[string]int

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
		r.Route("/commnets", CommentRouter)
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
	err = goview.Render(w, http.StatusOK, "posts/show", goview.M{
		"Post":           &post,
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
	uuid := createUUID()
	title := r.FormValue("title")
	category := strings.Split(r.FormValue("category"), ", ")
	content := r.FormValue("content")
	err := data.CreatePost(uuid, title, content, category)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to create a post: %v!", err))
		return
	}
	CategoryNavbar, err = data.UpdateCategory()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update category: %v!", err))
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%v", uuid), 302)
}

// update a post
func updatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uuid := r.FormValue("uuid")
	title := r.FormValue("title")
	category := []string{r.FormValue("category")}
	content := r.FormValue("content")
	err := data.UpdatePost(uuid, title, content, category)
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
