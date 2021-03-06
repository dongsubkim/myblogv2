package routes

import (
	"fmt"
	"net/http"
	"strconv"

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

		r.With(AdminOnly).Get("/edit", renderEdit)
		// Comment router
		r.Route("/comments", CommentRouter)
	})
}

// Get all posts
func getIndex(w http.ResponseWriter, r *http.Request) {
	var posts []data.Post
	var err error
	var page int
	var category string
	var query string
	if p, ok := r.URL.Query()["page"]; ok {
		page, err = strconv.Atoi(p[0])
		if err != nil {
			page = 0
		}
	} else {
		page = 0
	}

	if queries, ok := r.URL.Query()["search"]; ok {
		query = queries[0]
		posts, err = data.PostsBySearch(query, page)
	} else if categories, ok := r.URL.Query()["category"]; ok {
		category = categories[0]
		posts, err = data.PostsByCategory(category, page)
	} else {
		posts, err = data.Posts(page)
	}
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get posts: %v!", err))
		return
	}

	err = goview.Render(w, http.StatusOK, "posts/index", goview.M{
		"Posts":          &posts,
		"CategoryNavbar": &CategoryNavbar,
		"Authorized":     authorized(r),
		"IsFirst":        page == 0,
		"IsLast":         len(posts) < data.PostPerPage,
		"Page":           page,
		"Category":       category,
		"Query":          query,
		"Partials":       []string{"posts/index"},
		"add": func(a int, b int) int {
			return a + b
		},
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
		"Comments":       comments,
		"CategoryNavbar": &CategoryNavbar,
		"Authorized":     authorized(r),
		"Partials":       []string{"posts/show"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render show page error: %v!", err))
		return
	}
}

// create a post
func createPost(w http.ResponseWriter, r *http.Request) {
	images, err := uploadImages(r)
	if err != nil {
		error_message(w, r, "Fail to upload images to cloudinary server")
	}

	title := r.PostFormValue("title")
	category := r.PostFormValue("category")
	content := replaceImageLink(r.PostFormValue("content"), images)

	uuid, err := data.CreatePost(title, category, content, images)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to create a post: %v!", err))
		return
	}

	CategoryNavbar, err = data.UpdateCategory()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update category: %v!", err))
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%v", uuid), http.StatusFound)
}

// update a post
func updatePost(w http.ResponseWriter, r *http.Request) {
	images, err := uploadImages(r)
	if err != nil {
		error_message(w, r, "Fail to upload images to cloudinary server")
	}

	uuid := chi.URLParam(r, "uuid")
	title := r.PostFormValue("title")
	category := r.PostFormValue("category")
	content := replaceImageLink(r.PostFormValue("content"), images)

	err = data.UpdatePost(uuid, title, category, content, images)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update the post: %v!", err))
		return
	}

	deletedImages := r.PostForm["deleteImages"]
	for _, image_uuid := range deletedImages {
		err = data.DeleteImageByUuid(image_uuid)
		if err != nil {
			error_message(w, r, fmt.Sprintf("Fail to remove image of post: %v!", err))
		}
		err = deleteImage(image_uuid)
		if err != nil {
			error_message(w, r, fmt.Sprintf("Error on destory call to cloudinary server: %v!", err))
		}
	}

	CategoryNavbar, err = data.UpdateCategory()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update category: %v!", err))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%v", uuid), http.StatusFound)
}

// delete a post
func deletePost(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	images, err := data.ImagesByPost(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to get images of post: %v!", err))
		return
	}
	for _, image := range images {
		err = deleteImage(image.Uuid)
		if err != nil {
			error_message(w, r, fmt.Sprintf("Error on destory call to cloudinary server: %v!", err))
			return
		}
	}
	err = data.DeletePost(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to delete the post: %v!", err))
		return
	}
	CategoryNavbar, err = data.UpdateCategory()
	if err != nil {
		error_message(w, r, fmt.Sprintf("Fail to update category: %v!", err))
		return
	}
	http.Redirect(w, r, "/post", http.StatusFound)
}

// render edit form
func renderEdit(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	post, err := data.PostByUUID(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get a post: %v!", err))
		return
	}

	images, err := data.ImagesByPost(uuid)
	if err != nil {
		error_message(w, r, fmt.Sprintf("Cannot get images of post: %v", err))
	}

	err = goview.Render(w, http.StatusOK, "posts/edit", goview.M{
		"Post":           &post,
		"Images":         &images,
		"CategoryNavbar": &CategoryNavbar,
		"Authorized":     authorized(r),
		"Partials":       []string{"posts/edit"},
	})
	if err != nil {
		error_message(w, r, fmt.Sprintf("Render index error: %v!", err))
		return
	}
}
