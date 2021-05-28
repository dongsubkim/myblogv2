package data

import (
	"log"
	"strings"
	"time"
)

type Image struct {
	Id        int
	Uuid      string
	Filename  string
	FileUrl   string
	PostUuid  string
	CreatedAt time.Time
}

func (image *Image) SquareURL() (url string) {
	return strings.Replace(image.FileUrl, "/upload", "/upload/w_180,ar_1:1,c_fill,g_auto", 1)
}

func insertImages(images []*Image, postUuid string) (err error) {
	for _, image := range images {
		image.PostUuid = postUuid
		err = image.CreateImage()
		if err != nil {
			log.Println("Error during CreateImage", err)
			return
		}
	}
	return
}

// CreateImage inserts the image to postgres
func (image *Image) CreateImage() (err error) {
	statement := "insert into images (uuid, filename, file_url, post_uuid, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(image.Uuid, image.Filename, image.FileUrl, image.PostUuid, time.Now()).Scan(&image.Id, &image.Uuid, &image.CreatedAt)
	return
}

// DeleteImage removes image from posgres db
func (image *Image) DeleteImage() (err error) {
	_, err = db.Exec("DELETE FROM images where uuid = $1", image.Uuid)
	return
}

// DeleteImageByUuid removes image fromo db by uuid
func DeleteImageByUuid(uuid string) (err error) {
	_, err = db.Exec("DELETE FROM images where uuid = $1", uuid)
	return
}

// DeleteByPost removes images of post
func DeleteImagesByPost(postUuid string) (err error) {
	_, err = db.Exec("Delete from images where post_uuid = $1", postUuid)
	return
}

// ImagesByPost returns images of the post by post uuid
func ImagesByPost(postUuid string) (images []Image, err error) {
	rows, err := db.Query("SELECT id, uuid, filename, file_url, post_uuid, created_at from images where post_uuid = $1", postUuid)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		image := Image{}
		if err = rows.Scan(&image.Id, &image.Uuid, &image.Filename, &image.FileUrl, &image.PostUuid, &image.CreatedAt); err != nil {
			return
		}
		images = append(images, image)
	}
	return
}

func ImageByPost(postUuid string) (image Image, err error) {
	stmt := "SELECT file_url from images where post_uuid = $1 LIMIT 1"
	err = db.QueryRow(stmt, postUuid).Scan(&image.FileUrl)
	return
}
