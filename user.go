package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	server *Server
	conn   net.Conn
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		server: server,
		conn:   conn,
	}
	//启动userMessage 监听进程
	go user.ListenMessage()

	return user
}

func (u *User) OnLine() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播

	u.server.BroadCast(u, "上线了")
}
func (u *User) OffLine() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播

	u.server.BroadCast(u, "下线了")
}

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, us := range u.server.OnlineMap {
			onLineMap := "[" + us.Addr + "]" + ":" + "在线\n"
			u.C <- onLineMap
			//u.SendMsg(onLineMap)
		}
		u.server.mapLock.Unlock()
	} else {

		u.server.BroadCast(u, msg)
	}

}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
