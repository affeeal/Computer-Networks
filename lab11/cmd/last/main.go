package main

import (
	"context"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
	"google.golang.org/api/option"
)

type Message struct {
	Number     string `json:"number"`
	Time       string `json:"time"`
	Difficulty string `json:"difficulty"`
	Hash       string `json:"hash"`
	TxsLen     string `json:"txs_len"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleLast(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	ethClient, err := ethclient.Dial("") // to fill
	if err != nil {
		log.Println(err)
		return
	}
	defer ethClient.Close()

	ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: "", // to fill
	}

	opt := option.WithCredentialsFile("") // to fill
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Println(err)
		return
	}

	dbClient, err := app.Database(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			header, err := ethClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Println(err)
				return
			}

			blockNumber := big.NewInt(header.Number.Int64())
			block, err := ethClient.BlockByNumber(context.Background(), blockNumber)
			if err != nil {
				log.Println(err)
				return
			}

			message := Message{
				Number:     block.Number().String(),
				Time:       strconv.FormatUint(block.Time(), 10),
				Difficulty: block.Difficulty().String(),
				Hash:       block.Hash().String(),
				TxsLen:     strconv.Itoa(len(block.Transactions())),
			}

			if err := dbClient.NewRef("last").Set(ctx, message); err != nil {
				log.Println(err)
				return
			}

			err = c.WriteJSON(message)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/last", handleLast)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
