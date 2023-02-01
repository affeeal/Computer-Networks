package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type Message struct {
	Out string `json:"out"`
	Err string `json:"err"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var config = &ssh.ClientConfig{
	User: "test",
	Auth: []ssh.AuthMethod{
		ssh.Password("SDHBCXdsedfs222"),
	},
	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
}

func handleSSH(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	client, err := ssh.Dial("tcp", "151.248.113.144:443", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			session, err := client.NewSession()
			if err != nil {
				log.Println(err)
				return
			}

			var outBuff, errBuff bytes.Buffer
			session.Stdout = &outBuff
			session.Stderr = &errBuff
			session.Run("cat achtung.txt") // ошибку игнорирую

			session.Close()
			message := Message{
				Out: outBuff.String(),
				Err: errBuff.String(),
			}

			if err = c.WriteJSON(message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ssh", handleSSH)
	log.Fatal(http.ListenAndServe("localhost:8282", nil))
}
