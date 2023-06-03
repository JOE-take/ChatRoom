package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var content []string

var Db *sql.DB

func main() {
	// マルチプレクサにmuxを使う
	mux := http.NewServeMux()

	// 静的ファイルの配信
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

	// Dbの設定
	err := dbInit()
	defer Db.Close()
	if err != nil {
		panic(err)
	}

	// ハンドラの設定
	mux.HandleFunc("/", index)
	mux.HandleFunc("/chatroom", chatRoom)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)

	sessions = map[string]session{}
	content = make([]string, 0, 16)

	// サーバの設定、起動
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
