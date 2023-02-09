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
	ChainId  string `json:"chain_id"`
	Hash     string `json:"hash"`
	Value    string `json:"value"`
	Cost     string `json:"cost"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gas_price"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleTxs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	ethClient, err := ethclient.Dial("https://mainnet.infura.io/v3/2b77419575644789b219d5a09a0f0f7b")
	if err != nil {
		log.Println(err)
		return
	}
	defer ethClient.Close()

	ctx := context.Background()
	config := &firebase.Config{
		DatabaseURL: "https://lab11-19f34-default-rtdb.firebaseio.com/",
	}

	opt := option.WithCredentialsFile("") // path to file
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
			blockNumber := big.NewInt(16240027)
			block, err := ethClient.BlockByNumber(context.Background(), blockNumber)
			if err != nil {
				log.Fatal(err)
			}

			messages := make([]Message, 0)
			for _, tx := range block.Transactions() {
				messages = append(messages, Message{
					ChainId:  tx.ChainId().String(),
					Hash:     tx.Hash().String(),
					Value:    tx.Value().String(),
					Cost:     tx.Cost().String(),
					To:       tx.To().String(),
					Gas:      strconv.FormatUint(tx.Gas(), 10),
					GasPrice: tx.GasPrice().String(),
				})
			}

			if err := dbClient.NewRef("txs").Set(ctx, messages); err != nil {
				log.Println(err)
				return
			}

			err = c.WriteJSON(messages)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/txs", handleTxs)
	log.Fatal(http.ListenAndServe("localhost:8083", nil))
}
