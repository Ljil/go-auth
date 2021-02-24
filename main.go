package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"tutorial/posts"
	"tutorial/user"

	_ "github.com/lib/pq"
)

func main() {
	//Инициализация БД
	db, dbConnectionError := sql.Open(
		"postgres",
		"user=postgres dbname=gram password=postgres sslmode=disable")
	if dbConnectionError != nil {
		log.Fatal(dbConnectionError)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Print(err)
		}
	}()

	http.Handle("/login", user.Login(db))
	http.Handle("/posts/", user.LoginRequiredMiddleWare(db, posts.PostList(db)))
	http.Handle("/register", user.CreateNewUser(db))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
