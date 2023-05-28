package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var content []string

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

	content = make([]string, 0, 16)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)

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

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/login.html")

	if r.Method == "POST" {
		userName := r.FormValue("userName")
		email := r.FormValue("email")
		//パスワードの扱いは変更が必要
		password := r.FormValue("password")
		fmt.Println("userName: ", userName)
		fmt.Println("email: ", email)
		fmt.Println("password: ", password)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}

	t.Execute(w, "")
}
