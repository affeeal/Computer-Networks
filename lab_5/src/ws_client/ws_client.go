package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/websocket"

	"proto"
)

const COMMANDS = `/calculate - вычислить силы притяжения;
/quit      - закрыть соединение.`

var (
	addr    = flag.String("addr", "localhost:8080", "http service address")
	scanner = bufio.NewScanner(os.Stdin)
)

func scanProp(prop string) (float64, error) {
	fmt.Printf("%s: ", prop)
	scanner.Scan()
	return strconv.ParseFloat(scanner.Text(), 64)
}

func calculate(c *websocket.Conn) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Введите характеристики %d-го тела\n", i+1)
		mass, err := scanProp("Масса")
		if err != nil {
			log.Println("mass: ", err)
			return
		}
		x, err := scanProp("Координата X")
		if err != nil {
			log.Println("x: ", err)
			return
		}
		y, err := scanProp("Координата Y")
		if err != nil {
			log.Println("y: ", err)
			return
		}
		json_body, err := json.Marshal(proto.Body{Mass: mass, X: x, Y: y})
		if err != nil {
			log.Println("marshal: ", err)
			return
		}
		err = c.WriteMessage(websocket.TextMessage, json_body)
		if err != nil {
			log.Println("write message: ", err)
			return
		}
	}

	results := make([]float64, 3)
	for i := 0; i < 3; i++ {
		_, p, err := c.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			return
		}
		err = json.Unmarshal(p, &results[i])
		if err != nil {
			log.Println("unmarshal: ", err)
			return
		}
	}
	fmt.Println("\nРЕЗУЛЬТАТЫ:")
	fmt.Printf("F(1, 2) ≈ %.3e Н.\n", results[0])
	fmt.Printf("F(1, 3) ≈ %.3e Н.\n", results[1])
	fmt.Printf("F(2, 3) ≈ %.3e Н.\n", results[2])
}

func quit(c *websocket.Conn) {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatal("write close:", err)
	}
	fmt.Println()
	os.Exit(0)
}

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		<-interrupt
		quit(c)
	}()

	fmt.Print("> ")
	for scanner.Scan() {
		cmd := scanner.Text()
		switch cmd {
		case "/calculate":
			calculate(c)
		case "/quit":
			quit(c)
		default:
			fmt.Println(COMMANDS)
		}
		fmt.Print("> ")
	}
}
