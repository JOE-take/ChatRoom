package main

import (
	"net/http"
)

var content []string

// Users Usersはemailをkeyに使用してUserの構造体を返す
var Users map[string]User

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/chatroom", chatRoom)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)

	Users = map[string]User{}
	sessions = map[string]session{}
	content = make([]string, 0, 16)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
