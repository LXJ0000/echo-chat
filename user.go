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

func (u *User) SendMsg(msg string) {
	u.Conn.Write([]byte(msg))
}

func (u *User) DoMsg(msg string) {
	if msg == "who" {
		u.Server.MapLock.Lock()
		for _, user := range u.Server.OnlineMap {
			onlineMsg := strings.Join([]string{"[", user.Addr, "]", user.Name, ": 在线 ...\n"}, "")
			u.SendMsg(onlineMsg)
		}
		u.Server.MapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := msg[7:]
		u.Server.MapLock.Lock()
		if _, ok := u.Server.OnlineMap[newName]; ok {
			u.SendMsg("用户名已存在")
		} else {
			delete(u.Server.OnlineMap, u.Name)
			u.Server.OnlineMap[newName] = u
			u.Name = newName
			u.SendMsg(strings.Join([]string{"用户名修改成功:", u.Name, "\n"}, ""))
		}
		u.Server.MapLock.Unlock()
	} else {
		u.Server.BroadCast(u, msg)
	}

}

// ListenMsg 监听服务器发来的消息
func (u *User) ListenMsg() {
	for {
		msg := <-u.Chan
		_, _ = u.Conn.Write([]byte(strings.Join([]string{msg, "\n"}, "")))
	}
}
