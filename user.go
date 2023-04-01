package main

import (
	"net"
	"strings"
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
			onLineMap := "[" + us.Addr + "]" + us.Name + ":" + "在线\n"
			u.C <- onLineMap
			//u.SendMsg(onLineMap)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("用户名以存在")
		} else {

			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("用户名已更新为：" + newName)
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// format:  to|Jack|hello
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("格式不正确，请输入\"to|Jack|Welcome\"")
		}
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("用户不存在")
			return
		}
		remoteMsg := strings.Split(msg, "|")[2]
		if remoteMsg == "" {
			u.SendMsg("内容不能为空")
		}
		remoteUser.SendMsg(u.Name + "对你说：" + remoteMsg)

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
