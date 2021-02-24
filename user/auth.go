package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Login выдача токена на верный логин+пароль
func Login(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		type LoginUser struct {
			UserName string `json:"user_name"`
			Password string `json:"password"`
		}

		var user LoginUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Println(err)
		}

		type Token struct {
			Token string `json:"token"`
		}
		var token Token
		err = db.QueryRow("select token from users where username = $1 and password = $2",
			user.UserName, user.Password,
		).Scan(&token.Token)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if jsonToken, err := json.Marshal(token); err == nil {
			w.Header().Add("Content-type", "application/json")
			_, err = w.Write(jsonToken)
		} else {
			fmt.Println(err)
		}
	})
}

// LoginRequiredMiddleWare asd
func LoginRequiredMiddleWare(db *sql.DB, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var id int
		err := db.QueryRow("select user_id from users where token = $1", r.Header.Get("token")).Scan(&id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
