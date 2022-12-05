package main

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
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
	conn, _ := net.Dial("tcp", address)

	_, err := fmt.Fprintf(conn, "go client connect")
	if err != nil {
		return
	}

	defer handlePanic()

	for {
		funcName(conn)

		handleResponse(conn)
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

func handleResponse(conn net.Conn) {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}

func funcName(conn net.Conn) {
	format := "go clietn connect"
	_, err := fmt.Fprintf(conn, format)
	if err != nil {
		panic(err)
	}
}
