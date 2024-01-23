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
	u.Conn.Write([]byte(strings.Join([]string{msg, "\n"}, "")))
}

func (u *User) DoMsg(msg string) {
	if msg == "who" {
		u.Server.MapLock.Lock()
		for _, user := range u.Server.OnlineMap {
			onlineMsg := strings.Join([]string{"[", user.Addr, "]", user.Name, ": 在线 ..."}, "")
			u.SendMsg(onlineMsg)
		}
		u.Server.MapLock.Unlock()
	} else if len(msg) > 3 && msg[:3] == "to|" {
		// to|jannan|content
		// 1. get recv content
		split := strings.Split(msg, "|")
		if len(split) < 3 || split[1] == "" || split[2] == "" {
			u.SendMsg("消息格式有误，应该为 \"to|recv|content\" 格式。")
			return
		}
		recv, content := split[1], split[2]
		// 2. get user
		recvUser, ok := u.Server.OnlineMap[recv]
		if !ok {
			u.SendMsg("用户名不存在")
			return
		}
		// 3. send msg
		recvUser.SendMsg(strings.Join([]string{u.Name, "对你说：", content}, ""))
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// rename|jannan
		newName := msg[7:]
		u.Server.MapLock.Lock()
		if _, ok := u.Server.OnlineMap[newName]; ok {
			u.SendMsg("用户名已存在")
		} else {
			delete(u.Server.OnlineMap, u.Name)
			u.Server.OnlineMap[newName] = u
			u.Name = newName
			u.SendMsg(strings.Join([]string{"用户名修改成功:", u.Name}, ""))
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
