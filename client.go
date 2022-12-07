package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
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

	return Client{
		conn: conn,
	}, nil
}

func (c *Client) HandleResponse() {
	message, _ := bufio.NewReader(c.conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}

func (c *Client) SendMsgToServer(msg string) {
	_, err := fmt.Fprintf(c.conn, msg)
	if err != nil {
		err := c.conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
}
