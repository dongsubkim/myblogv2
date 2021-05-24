package data

import (
	"time"

	"github.com/lib/pq"
)

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

// Get all Posts in the database and returns it
func Posts() (posts []Post, err error) {
	rows, err := db.Query("SELECT id, uuid, title, category, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Title, pq.Array(&post.Category), &post.Content, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Create a new post
func CreatePost(uuid, title, content string, category []string) (err error) {
	statement := "insert into posts (uuid, title, category, content, created_at) values ($1, $2, $3, $4, $5)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Query(uuid, title, pq.Array(category), content, time.Now())
	return
}

// Update a post
func UpdatePost(uuid, title, content string, category []string) (err error) {
	_, err = db.Exec("update posts set title = $2, category = $3, content = $4, created_at = $5 where uuid = $1", uuid, title, pq.Array(category), content, time.Now())
	return
}

// Delete a post
func DeletePost(uuid string) (err error) {
	_, err = db.Exec("delete from posts where uuid = $1", uuid)
	return
}

// Get Posts by Category
func PostsByCategory(category string) (posts []Post, err error) {
	rows, err := db.Query("SELECT id, title, category, content, created_at FROM posts WHERE $1=any(category) ORDER BY created_at DESC", category)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Title, pq.Array(&post.Category), &post.Content, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Get a thread by the UUID
func PostByUUID(uuid string) (post Post, err error) {
	post = Post{}
	err = db.QueryRow("SELECT id, uuid, title, category, content, created_at FROM posts WHERE uuid = $1", uuid).
		Scan(&post.Id, &post.Uuid, &post.Title, pq.Array(&post.Category), &post.Content, &post.CreatedAt)
	return
}
