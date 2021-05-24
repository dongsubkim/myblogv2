package data

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id        int
	Uuid      string
	Category  []string
	Content   string
	CreatedAt time.Time
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("06/01/02 3:04pm")
}

// Get all Posts in the database and returns it
func Posts() (posts []Post, err error) {
	rows, err := db.Query("SELECT id, uuid, category, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, pq.Array(&post.Category), &post.Content, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Create a new post
func CreatePost(uuid string, category []string, content string) (err error) {
	statement := "insert into posts (uuid, category, content, created_at) values ($1, $2, $3, $4)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Query(uuid, pq.Array(category), content, time.Now())
	return
}

// Update a post
func UpdatePost(uuid string, category []string, content string) (err error) {
	statement := "update posts set category = $2, content = $3, created_at = $4 where uuid = $1"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Query(uuid, pq.Array(category), content, time.Now())
	return
}

// Delete a post
func DeletePost(uuid string) (err error) {
	statement := "delete from posts where uuid = $1"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Query(uuid)
	return
}

// Get Posts by Category
func PostsByCategory(category string) (posts []Post, err error) {
	rows, err := db.Query("SELECT id, category, content, created_at FROM posts WHERE $1 IN category ORDER BY created_at DESC", category)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, pq.Array(&post.Category), &post.Content, &post.CreatedAt); err != nil {
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
	err = db.QueryRow("SELECT id, uuid, category, content, created_at FROM posts WHERE uuid = $1", uuid).
		Scan(&post.Id, &post.Uuid, pq.Array(&post.Category), &post.Content, &post.CreatedAt)
	return
}
