package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Name + "]" + ":" + msg

	s.Message <- sendMsg
}

func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Handler(conn net.Conn) {

	fmt.Println("accept successful")

	user := NewUser(conn, s)

	//用户上线，加入onlineMap
	user.OnLine()
	//广播

	s.BroadCast(user, "上线了")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Error:", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)
		}
	}()

	select {}
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

	//启动监听Message的gorutine
	go s.ListenMessager()

	for {

		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		//handlers
		go s.Handler(conn)
	}

}
