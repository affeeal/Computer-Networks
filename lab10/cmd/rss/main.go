package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

type Message struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRSS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rows, err := db.Query("SELECT title, link FROM iu9afanasyev WHERE id <= 20")
			if err != nil {
				log.Println(err)
				return
			}

			var messages []Message
			for rows.Next() {
				var message Message

				err := rows.Scan(&message.Title, &message.Link)
				if err != nil {
					log.Println(err)
					return
				}

				messages = append(messages, message)
			}

			err = rows.Err()
			if err != nil {
				log.Println(err)
				return
			}

			rows.Close()
			if err = c.WriteJSON(messages); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/rss", handleRSS)
	log.Fatal(http.ListenAndServe("localhost:8181", nil))
}
