package routes

import (
	"log"
	"net/http"

	"github.com/foolin/goview"
)

// Convenience function to redirect to the error message page
func error_message(w http.ResponseWriter, r *http.Request, msg string) {
	log.Println("ERROR:", msg)
	_ = goview.Render(w, http.StatusBadRequest, "error", goview.M{
		"err": msg,
	})
	return
}
