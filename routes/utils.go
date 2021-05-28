package routes

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/dongsubkim/myblogv2/data"
	"github.com/foolin/goview"
	"github.com/joho/godotenv"
)

var cld *cloudinary.Cloudinary

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	cld, err = cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_KEY"), os.Getenv("CLOUDINARY_SECRET"))
	if err != nil {
		log.Println("Error loading cloudinary api")
	}
}

// Convenience function to redirect to the error message page
func error_message(w http.ResponseWriter, r *http.Request, msg string) {
	log.Println("ERROR:", msg)
	_ = goview.Render(w, http.StatusBadRequest, "error", goview.M{
		"Authorized":     authorized(r),
		"CategoryNavbar": &CategoryNavbar,
		"err":            msg,
	})
}

// Upload images
func uploadImages(r *http.Request) (images []*data.Image, err error) {
	r.ParseMultipartForm(1024 * 1024 * 32)
	files := r.MultipartForm.File["image"]
	var file multipart.File
	if len(files) > 0 {
		var ctx = context.Background()
		for _, fileHeader := range files {
			file, err = fileHeader.Open()
			if err != nil {
				log.Println("Fail to open fileHeader")
				return
			}
			defer file.Close()

			uploadResult, err := cld.Upload.Upload(
				ctx,
				file,
				uploader.UploadParams{Folder: "MyBlog", AllowedFormats: api.CldAPIArray{"jpg", "jpeg", "png"}})
			if err != nil {
				log.Printf("Failed to upload file, %v", err)
			}
			image := data.Image{
				Uuid:     uploadResult.PublicID,
				Filename: fileHeader.Filename,
				FileUrl:  uploadResult.SecureURL,
			}
			images = append(images, &image)
		}
	}
	return
}

// Delete images from cloudinary server
func deleteImage(uuid string) (err error) {
	destoryResult, err := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: uuid})
	if err != nil || destoryResult.Result != "ok" {
		log.Printf("Faild to destroy image on cloudinary server: %v", err)
	}
	return
}

func replaceImageLink(content string, images []*data.Image) string {
	for _, image := range images {
		content = strings.ReplaceAll(content, fmt.Sprintf("![%s](http://)", image.Filename), fmt.Sprintf("![%s](%s)", image.Filename, image.FileUrl))
	}
	return content
}
