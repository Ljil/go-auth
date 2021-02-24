package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	// Не забывать писать имена полей с заглавной буквы, чтобы они появлялись в json
	Username  string `json:"username"`
	PostID    int    `json:"post_id"`
	PostTitle string `json:"post_title"`
	PostText  string `json:"post_text"`
}

func PostList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Posts struct {
			Posts []Post `json:"posts"`
		}
		posts := &Posts{}
		var queryString string
		// получить url вида /posts/{%username%}
		if len(r.URL.Path[len("/posts/"):]) > 0 {
			queryString = fmt.Sprintf(
				`select u.username, p.post_id, p.post_title, p.post_text 
				from users u join posts p on u.user_id = p.post_author 
				where u.username = '%v' order by p.post_id desc;`,
				r.URL.Path[len("/posts/"):])
		} else {
			//просмотр всех постов подряд
			queryString = `select u.username, p.post_id, p.post_title, p.post_text 
			from users u join posts p on u.user_id = p.post_author 
			group by u.username, p.post_id 
			order by p.post_id desc;`
		}
		//для просмотра постов конкретного авторизованного пользователя (через токен)
		// rows, err := db.Query("select u.username, p.post_id, p.post_title, p.post_text from users u join posts p on u.user_id = p.post_author where u.token = $1 group by u.username, p.post_id order by p.post_id desc;", r.Header.Get("token"))

		rows, err := db.Query(queryString)
		if err != nil {
			fmt.Println(err)
			return
		}

		for rows.Next() {
			post := Post{}
			err := rows.Scan(
				&post.Username,
				&post.PostID,
				&post.PostTitle,
				&post.PostText,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			posts.Posts = append(posts.Posts, post)
		}

		if posts, err := json.Marshal(posts); err == nil {
			w.Header().Add("Content-type", "application/json")
			_, err = w.Write(posts)
		} else {
			fmt.Println(err)
		}

	}
}
