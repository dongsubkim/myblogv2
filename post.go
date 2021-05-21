package main

import "net/http"

func handlePost(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handlePostGet(w, r)
	case "POST":
		err = handlePostPost(w, r)
	case "PUT":
		err = handlePostPut(w, r)
	case "DELETE":
		err = handlePostDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlePostGet(w http.ResponseWriter, r *http.Request) (err error) {
	return
}

func handlePostPost(w http.ResponseWriter, r *http.Request) (err error) {
	return
}

func handlePostPut(w http.ResponseWriter, r *http.Request) (err error) {
	return
}

func handlePostDelete(w http.ResponseWriter, r *http.Request) (err error) {
	return
}
