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
	return Client{
		conn:         conn,
		fromServerCh: fromServerCh,
		toServerCh:   toServerCh,
		controller:   controller.NewController(fromServerCh, toServerCh),
	}, nil
}

func (c *Client) HandleResponse() {
	message, err := bufio.NewReader(c.conn).ReadString('\n')

	if err != nil {
		log.Println(err)
	}

	go c.controller.ProcessData()
	c.fromServerCh <- message
}

func (c *Client) SendMsgToServer() {
	_, err := fmt.Fprintf(c.conn, <-c.fromServerCh)
	if err != nil {
		err := c.conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
}
