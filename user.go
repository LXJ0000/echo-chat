package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	Chan   chan string
	Conn   net.Conn
	Server *Server
}

// NewUser 创建用户实例
func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		Chan:   make(chan string),
		Conn:   conn,
		Server: server,
	}

	go user.ListenMsg()

	return user
}

func (u *User) Online() {
	u.Server.MapLock.Lock()
	u.Server.OnlineMap[u.Name] = u
	u.Server.MapLock.Unlock()

	u.Server.BroadCast(u, "User OnLine")
}

func (u *User) OffLine() {
	u.Server.BroadCast(u, "User OffLine")
}

func (u *User) DoMsg(msg string) {
	u.Server.BroadCast(u, msg)
}

// ListenMsg 监听服务器发来的消息
func (u *User) ListenMsg() {
	for {
		msg := <-u.Chan
		_, _ = u.Conn.Write([]byte(strings.Join([]string{msg, "\n"}, "")))
	}
}
