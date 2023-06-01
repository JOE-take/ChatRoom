package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var content []string

// Users Usersはemailをkeyに使用してUserの構造体を返す
var Users map[string]User

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

	content = make([]string, 0, 16)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)

	Users = map[string]User{}

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/index.html")

	if v := r.FormValue("message"); r.Method == "POST" && v != "" {
		content = append(content, v)
		fmt.Printf("[%s]\n", v)
		//同じ場所にリダイレクトすることによって更新でPOSTが送られるのを防ぐ
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	t.Execute(w, content)
}
