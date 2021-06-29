package data

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Comment struct {
	Id        int
	Uuid      string
	Username  string
	Password  string
	Body      string
	PostUuid  string
	CreatedAt time.Time
}

func (comment *Comment) CreatedAtDate() string {
	return comment.CreatedAt.Format("06/01/02 3:04pm")
}

// get comments to a post
func (post *Post) Comments() (comments []Comment, err error) {
	rows, err := db.Query("SELECT id, uuid, username, body, post_uuid, created_at FROM comments where post_uuid = $1 ORDER BY created_at", post.Uuid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		comment := Comment{}
		if err = rows.Scan(&comment.Id, &comment.Uuid, &comment.Username, &comment.Body, &comment.PostUuid, &comment.CreatedAt); err != nil {
			return
		}
		comments = append(comments, comment)
	}
	return
}

// Create a new comment
func CreateComment(username, password, body, postUuid string) (err error) {
	statement := "insert into comments (uuid, username, password, body, post_uuid, created_at) values ($1, $2, $3, $4, $5, $6)"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	uuid := createUUID()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return
	}
	rows, err := stmt.Query(uuid, username, string(hashedPassword), body, postUuid, time.Now())
	if err != nil {
		return
	}
	defer rows.Close()
	return
}

// Get a comment by commentId
func CommentByUuid(uuid string) (comment Comment, err error) {
	err = db.QueryRow("SELECT id, uuid, username, password, body, post_uuid, created_at from comments where uuid = $1", uuid).
		Scan(&comment.Id, &comment.Uuid, &comment.Username, &comment.Password, &comment.Body, &comment.PostUuid, &comment.CreatedAt)
	return
}

// Update comment
func UpdateComment(uuid, body string) (err error) {
	_, err = db.Exec("update comments set body = $2 where uuid = $1", uuid, body)
	return
}

// Delete comment
func DeleteComment(uuid string) (err error) {
	_, err = db.Exec("delete from comments where uuid = $1", uuid)
	return
}

// Get a password of comment
func PasswordByComment(uuid string) (password string, err error) {
	err = db.QueryRow("select password from comments where uuid = $1", uuid).Scan(&password)
	return
}

// DeleteCommentsByPost removes comments of a post
func DeleteCommentsByPost(postUuid string) (err error) {
	_, err = db.Exec("delete from comments where post_uuid = $1", postUuid)
	return
}
