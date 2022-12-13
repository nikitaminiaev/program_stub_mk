package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"stubMk/controller"
)

type Client struct {
	conn         net.Conn
	fromServerCh chan string
	toServerCh   chan string
	controller   controller.MkController
}

func NewClient(address string) (Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return Client{}, errors.New("connection is not established")
	}

	_, err = fmt.Fprintf(conn, "go client connect")
	if err != nil {
		return Client{}, errors.New("connection is not established")
	}

	fromServerCh := make(chan string)
	toServerCh := make(chan string)
	newController := controller.NewController(fromServerCh, toServerCh)
	newController.GenSurface(11)

	return Client{
		conn:         conn,
		fromServerCh: fromServerCh,
		toServerCh:   toServerCh,
		controller:   newController,
	}, nil
}

func (c *Client) HandleResponse() {
	message, err := bufio.NewReader(c.conn).ReadString('\n')
	fmt.Println(err)
	fmt.Println(message)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(message)
	go c.controller.ProcessData()
	c.fromServerCh <- message
}

func (c *Client) SendMsgToServer() {
	fmt.Println("SendMsgToServer")
	_, err := fmt.Fprintf(c.conn, <-c.toServerCh)
	if err != nil {
		err := c.conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
}
