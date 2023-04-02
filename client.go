package main

import (
	"fmt"
	"net"
)

type Client struct {
	ClientIP   string
	ClientPort int
	Name       string
	Conn       net.Conn
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ClientIP:   serverIP,
		ClientPort: serverPort,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net Dial error:", err)
		return nil
	}
	client.Conn = conn

	return client
}

func main() {
	client := NewClient("127.0.0.1", 8090)
	if client == nil {
		fmt.Println("=======>链接失败<=======")
		return
	}
	fmt.Println("=======>链接成功<=======")

	select {}
}
