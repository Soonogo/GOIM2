package main

import (
	"flag"
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

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置链接ip地址")
	flag.IntVar(&serverPort, "port", 8090, "设置链接port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("=======>链接失败<=======")
		return
	}
	fmt.Println("=======>链接成功<=======")

	select {}
}
