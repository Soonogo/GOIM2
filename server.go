package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string `"json:ip"`
	Port int    `"json:port"`
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}
func (s *Server) Handler(conn *net.Conn) {
	fmt.Println(*conn)
	fmt.Println("accept successful")
}
func (s *Server) Start() {
	// socket listening
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("listen error:", err)
	}
	fmt.Println("开始监听：", fmt.Sprintf("%s:%d", s.Ip, s.Port))

	//closes
	defer listener.Close()

	for {

		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		//handlers
		go s.Handler(&conn)
	}

}
