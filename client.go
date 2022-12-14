package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
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
	newController.SurfaceGenerator.GenSurface(11)

	return Client{
		conn:         conn,
		fromServerCh: fromServerCh,
		toServerCh:   toServerCh,
		controller:   newController,
	}, nil
}

func (c *Client) HandleResponse() {
	message, err := bufio.NewReader(c.conn).ReadString('\n')
	if strings.HasPrefix(message, "connected") {
		fmt.Println(message)
		message = string([]rune(message)[9:])
		return
	}

	if err != nil {
		log.Println(err)
	}

	go c.controller.ProcessData()
	c.fromServerCh <- message
}

func (c *Client) SendMsgToServer() {
	msg := <-c.toServerCh
	_, err := fmt.Fprintf(c.conn, msg)
	if err != nil {
		err := c.conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
}
