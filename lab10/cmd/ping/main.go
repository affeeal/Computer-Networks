package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Status string `json:"status"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cmd := exec.Command("ping", "-c", "1", "-W", "5", "bmstu.ru")
			message := Message{
				Status: "BAUMAN is ok",
			}

			if err = cmd.Run(); err != nil {
				message.Status = "BAUMAN is not available"
			}

			if err = c.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ping", handlePing)
	log.Fatal(http.ListenAndServe("localhost:8383", nil))
}
