package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

var sessions map[string]session

type session struct {
	username string
	expiry   time.Time
}

type User struct {
	UserName string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/index.html")
	t.Execute(w, "")
}

func chatRoom(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/chatRoom.html")

	c, err := r.Cookie("session_token")
	if err != nil && err != http.ErrNoCookie {
		panic(err)
	} else if c == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if !authenticateSession(c) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	v := strings.Split(c.Value, "|")
	userName := v[0]

	if v := r.FormValue("message"); r.Method == "POST" && v != "" {
		post := userName + ": " + v
		content = append(content, post)
		http.Redirect(w, r, "/chatroom", http.StatusSeeOther)
	}

	t.Execute(w, content)
}

func signup(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/signup.html")

	if r.Method == "POST" {
		user := &User{
			UserName: r.FormValue("userName"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		//フォームのValidation
		if ok, errorMessage := user.Validate(); !ok {
			t.Execute(w, errorMessage)
			return
		}

		//パスワードの暗号化
		hashedPassword, _ := passwordEncrypt(user.Password)
		user.Password = string(hashedPassword)

		//データベースにユーザ情報を登録
		insert, err := Db.Prepare("insert into user values(?, ?, ?)")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to insert data to database\n")
			return
		}
		insert.Exec(user.UserName, user.Email, user.Password)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/login.html")

	if r.Method == "POST" {
		//フォームのデータを受け取る
		user := &User{
			UserName: "noname",
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		//フォームのValidation
		if ok, errorMessage := user.Validate(); !ok {
			t.Execute(w, errorMessage)
			return
		}

		registeredUser, err := getUserFromDb(user.Email)
		if err != nil {
			fmt.Println(err)
			return
		} else if registeredUser.Email == "" {
			//userが空なら登録されていない
			errorMessage := "enter a correct information"
			t.Execute(w, errorMessage)
			return
		} else {

			//パスワードが正しいかどうかのチェック
			if checkPassword([]byte(registeredUser.Password), user.Password) {
				//パスワード一致
				user.UserName = registeredUser.UserName

				//Sessionの開始(Cookieの設定)
				startSession(w, user)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				//パスワード不一致
				t.Execute(w, "enter a correct information")
				return
			}
		}
	}
	t.Execute(w, nil)
}

//func logout(w http.ResponseWriter)
