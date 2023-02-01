package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/websocket"

	"proto"
)

const G = 6.674e-11

var (
	addr     = flag.String("addr", "localhost:8080", "http service address")
	upgrader = websocket.Upgrader{}
)

func calculateForce(b1, b2 proto.Body) float64 {
	r := math.Sqrt(math.Pow(b1.X-b2.X, 2) + math.Pow(b1.Y-b2.Y, 2))
	return G * b1.Mass * b2.Mass * r / 2
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		bodies := make([]proto.Body, 3)
		for i := 0; i < 3; i++ {
			_, p, err := c.ReadMessage()
			if err != nil {
				log.Println("read :", err)
				return
			}
			err = json.Unmarshal(p, &bodies[i])
			if err != nil {
				log.Println("unmarshal: ", err)
				return
			}
		}

		results := []float64{
			calculateForce(bodies[0], bodies[1]),
			calculateForce(bodies[0], bodies[2]),
			calculateForce(bodies[1], bodies[2]),
		}

		for _, result := range results {
			json_result, err := json.Marshal(result)
			if err != nil {
				log.Println("marshal: ", err)
				return
			}
			err = c.WriteMessage(websocket.TextMessage, json_result)
			if err != nil {
				log.Println("write message: ", err)
				return
			}
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", defaultHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
