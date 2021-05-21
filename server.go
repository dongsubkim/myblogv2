package main

import "net/http"

func main() {
	server := http.Server{
		Addr: "127.0.0.1:3000",
	}
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/posts", handlePost)

	server.ListenAndServe()
}
