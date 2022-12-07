package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var address string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ip, exists := os.LookupEnv("IP")
	if !exists {
		return
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		return
	}

	address = fmt.Sprintf("%s:%s", ip, port)
}

func main() {
	defer handlePanic()

	client, err := NewClient(address)

	if err != nil {
		panic(err)
	}

	for {
		format := "go clietn connect"
		client.SendMsgToServer(format)

		client.HandleResponse()
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}
