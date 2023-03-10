package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var address string

func init() {
	setLogFile()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}

	ip, exists := os.LookupEnv("IP")
	if !exists {
		log.Println("Not found IP")
		return
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("Not found PORT")
		return
	}

	address = fmt.Sprintf("%s:%s", ip, port)
}

func setLogFile() {
	logFile, err := os.OpenFile("stubMk_log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(logFile)
}

func main() {
	defer handlePanic()

	client, err := NewClient(address)

	if err != nil {
		panic(err)
	}

	for {
		go client.SendMsgToServer()

		client.HandleResponse()
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Println(r)
	}
}
