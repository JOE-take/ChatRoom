package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type form struct {
	UserName string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type User struct {
	UserName       string
	Email          string
	hashedPassword []byte
}

func signup(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/signup.html")

	if r.Method == "POST" {
		userName := r.FormValue("userName")
		email := r.FormValue("email")
		password := r.FormValue("password")
		hash, _ := passwordEncrypt(password)

		//フォームのValidation
		formData := form{userName, email, password}
		if ok, errorMessage := formData.Validate(); !ok {
			t.Execute(w, errorMessage)
			fmt.Println("validator: ", errorMessage)
			return
		}

		Users[email] = User{
			UserName:       userName,
			Email:          email,
			hashedPassword: hash,
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/login.html")

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		formData := form{"login", email, password}
		if ok, errorMessage := formData.Validate(); !ok {
			t.Execute(w, errorMessage)
			return
		}

		//Userが存在するかどうか
		user, ok := Users[email]
		if !ok {
			fmt.Println("user can't find")
		}

		if checkPassword(user.hashedPassword, password) {
			fmt.Println("login")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			fmt.Println("can't login")
			t.Execute(w, "enter a correct information")
			return
		}
	}
	t.Execute(w, nil)
}

// passwordEncrypt パスワードを受け取り、ハッシュとエラーを返す
// httpsでないので伝送路上ではパスワードが平文で流れてくる
// dbを使っていないので暗号化する意味は今はないかも
func passwordEncrypt(password string) ([]byte, error) {
	if hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		return hash, err
	}
}

// ハッシュとパスワードが正しいかチェックする
func checkPassword(hash []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// Validate Validator
func (formData *form) Validate() (ok bool, result string) {
	err := validator.New().Struct(*formData)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		if len(errors) != 0 {
			for i := range errors {
				switch errors[i].StructField() {
				case "UserName":
					switch errors[i].Tag() {
					case "required":
						// += だと遅いらしい
						result = result + "UserName required\n"
					}
				case "Email":
					switch errors[i].Tag() {
					case "required":
						result = result + "Email required\n"
					case "email":
						result = result + "Enter a valid email address\n"
					}

				case "Password":
					switch errors[i].Tag() {
					case "required":
						result = result + "Password required"
					}
				}
			}
		}
		return false, result
	}
	return true, result
}
