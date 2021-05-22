package main

import "time"

type Post struct {
	Id        int
	Author    string
	Category  string
	Content   string
	CreatedAt time.Time
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("06/01/02 3:04pm")
}
