package main

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func dbInit() error {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/chatroom")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	Db = db
	return nil
}

func getUserFromDb(email string) (User, error) {
	query := fmt.Sprintf("select * from user where email = \"%s\"", email)
	row, err := Db.Query(query)
	if err != nil {
		return User{}, err
	}

	user := User{}
	for row.Next() {
		if err := row.Scan(&user.UserName, &user.Email, &user.Password); err != nil {
			return User{}, err
		}
	}

	return user, nil
}

func startSession(w http.ResponseWriter, u *User) {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(60 * time.Minute)
	userName := u.UserName

	sessions[sessionToken] = session{
		username: u.UserName,
		expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   userName + "|" + sessionToken,
		Expires: expiresAt,
	})
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// passwordEncrypt パスワードを受け取り、ハッシュとエラーを返す
func passwordEncrypt(password string) ([]byte, error) {
	if hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		return hash, nil
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
func (formData *User) Validate() (ok bool, result string) {
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

// セッションが存在するか、セッションがタイムアウトしていないかを確認する
func authenticateSession(cookie *http.Cookie) bool {

	// セッションが存在するかどうかの確認
	v := strings.Split(cookie.Value, "|")
	sessionToken := v[1]
	userSession, ok := sessions[sessionToken]
	if !ok {
		return false
	}

	// セッションがタイムアウトしていないかの確認
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		return false
	}

	// 何もなければtrueを返す
	return true
}
