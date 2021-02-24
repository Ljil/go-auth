package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// CreateNewUser Регистрация нового пользователя
func CreateNewUser(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		//получаем и декодируем json с данными нового пользователя
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Println(err)
			return
		}

		//Создаём пользователя
		var dbErr error
		_, dbErr = db.Exec("insert into users (username, password, email, token) values ($1,$2,$3,$4)",
			user.Username, user.Password, user.Email, user.Token)
		if dbErr != nil {
			fmt.Println(dbErr)
		}
	})
}
