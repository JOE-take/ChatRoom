package main

import (
	"html/template"
	"net/http"
)

var content []string

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

	content = make([]string, 0, 16)
	mux.HandleFunc("/", index)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/index.html")

	if r.Method == "POST" {
		content = append(content, r.FormValue("message"))
		//同じ場所にリダイレクトすることによって更新でフォームが送られるのを防ぐ
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	t.Execute(w, content)
}
