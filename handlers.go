package main

import (
	"fmt"
	"html/template"
	"net/http"
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
	authenticate(w, r)
	c, _ := r.Cookie("session_token")
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
			fmt.Println("validator: ", errorMessage)
			t.Execute(w, errorMessage)
			return
		}

		//パスワードの暗号化
		hashedPassword, _ := passwordEncrypt(user.Password)
		user.Password = string(hashedPassword)

		//User情報を追加
		Users[user.Email] = *user

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
			fmt.Println("validator: ", errorMessage)
			t.Execute(w, errorMessage)
			return
		}

		//Userが存在するかどうか
		if registeredUser, ok := Users[user.Email]; !ok {
			fmt.Println("user not found")
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

func authenticate(w http.ResponseWriter, r *http.Request) {

	//Cookieが存在するかどうかの確認
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		fmt.Println(err)
		return
	}

	//セッションが存在するかどうかの確認
	v := strings.Split(cookie.Value, "|")
	sessionToken := v[1]
	userSession, ok := sessions[sessionToken]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	//セッションがタイムアウトしていないかの確認
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

}
