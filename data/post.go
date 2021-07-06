package data

import (
	"bytes"
	"html/template"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
	stripmd "github.com/writeas/go-strip-markdown"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const PostPerPage = 5

type Post struct {
	Id        int
	Uuid      string
	Title     string
	Category  []string
	Content   string
	CreatedAt time.Time
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("06/01/02 3:04pm")
}

func (post *Post) PopulateCategory() string {
	return strings.Join(post.Category, ", ")
}

func (post *Post) ThumbnailImage() string {
	image, err := ImageByPost(post.Uuid)
	if err == nil {
		return image.SquareURL()
	}
	return ""
}

func (post *Post) ParseContent() template.HTML {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("vs"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(post.Content), &buf); err != nil {
		panic(err)
	}
	return template.HTML(buf.Bytes())
}

func (post *Post) SanitizedContent() string {
	stripped := stripmd.Strip(bluemonday.StrictPolicy().Sanitize(post.Content))
	if len(stripped) < 200 {
		return stripped
	}
	return stripped[:200] + "..."
}

// Get all Posts in the database and returns it
func Posts(page int) (posts []Post, err error) {
	rows, err := db.Query("SELECT id, uuid, title, category, content, created_at FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2", PostPerPage, page*PostPerPage)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Title, pq.Array(&post.Category), &post.Content, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

// Create a new post
func CreatePost(title, categoryRaw, content string, images []*Image) (uuid string, err error) {
	statement := "insert into posts (uuid, title, category, content, created_at) values ($1, $2, $3, $4, $5)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	uuid = createUUID()
	category := strings.Split(categoryRaw, ", ")
	_, err = stmt.Exec(uuid, title, pq.Array(category), content, time.Now())
	if err != nil {
		return
	}
	err = insertImages(images, uuid)
	return
}

// Update a post
func UpdatePost(uuid, title, categoryRaw, content string, images []*Image) (err error) {
	category := strings.Split(categoryRaw, ", ")
	_, err = db.Exec("UPDATE posts SET title = $2, category = $3, content = $4, created_at = $5 where uuid = $1", uuid, title, pq.Array(category), content, time.Now())
	if err != nil {
		return
	}
	err = insertImages(images, uuid)
	return
}

// Delete a post
func DeletePost(uuid string) (err error) {
	err = DeleteImagesByPost(uuid)
	if err != nil {
		return
	}
	err = DeleteCommentsByPost(uuid)
	if err != nil {
		return
	}
	_, err = db.Exec("delete from posts where uuid = $1", uuid)
	return
}

// Get Posts by Search
func PostsBySearch(query string, page int) (posts []Post, err error) {
	var stmt = `SELECT id, uuid, title, category, content, created_at FROM posts WHERE title LIKE '%` + query + `%' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := db.Query(stmt, PostPerPage, PostPerPage*page)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Title, pq.Array(&post.Category), &post.Content, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

// Get Posts by Category
func PostsByCategory(category string, page int) (posts []Post, err error) {
	rows, err := db.Query("SELECT id, uuid, title, category, content, created_at FROM posts WHERE $1=any(category) ORDER BY created_at DESC LIMIT $2 OFFSET $3", category, PostPerPage, PostPerPage*page)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Title, pq.Array(&post.Category), &post.Content, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

// Get a thread by the UUID
func PostByUUID(uuid string) (post Post, err error) {
	post = Post{}
	err = db.QueryRow("SELECT id, uuid, title, category, content, created_at FROM posts WHERE uuid = $1", uuid).
		Scan(&post.Id, &post.Uuid, &post.Title, pq.Array(&post.Category), &post.Content, &post.CreatedAt)
	return
}

// Update category of navbar
func UpdateCategory() (category map[string]int, err error) {
	category = make(map[string]int)
	rows, err := db.Query("SELECT category, count(*) FROM (SELECT unnest(category) AS category FROM posts) AS foo GROUP BY category ORDER BY category")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var cat string
		var count int
		if err = rows.Scan(&cat, &count); err != nil {
			return
		}
		category[cat] = count
	}
	return
}
