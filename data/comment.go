package data

import (
	"log"
	"time"
)

type Comment struct {
	Id        int
	Name      string
	Password  string
	Body      string
	CommentId int
	CreatedAt time.Time
}

func (comment *Comment) CreatedAtDate() string {
	return comment.CreatedAt.Format("06/01/02 3:04pm")
}

// get comments to a post
func (post *Post) Comments() (comments []Comment, err error) {
	rows, err := db.Query("SELECT id, username, body, created_at FROM comments where post_id = $1", post.Id)
	if err != nil {
		log.Panicln("Fail to get comments", err)
		return
	}
	for rows.Next() {
		comment := Comment{}
		if err = rows.Scan(&comment.Id, &comment.Name, &comment.Body, &comment.CreatedAt); err != nil {
			return
		}
		comments = append(comments, comment)
	}
	rows.Close()
	return
}

// Create a new post
func CreateComment(cred Credentials, post Post, body string) (comment Comment, err error) {
	statement := "insert into posts (username, password, body, comment_id, created_at) values ($1, $2, $3, $4, $5) returning id, username, password, body, comment_id, created_at"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(cred.Username, cred.Password, body, post.Id, time.Now()).Scan(&comment.Id, &comment.Name, &comment.Password, &comment.Body, &comment.CommentId, &comment.CreatedAt)
	return
}
