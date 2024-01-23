package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	Chan chan string
	Conn net.Conn
}

// NewUser 创建用户实例
func NewUser(conn net.Conn) *User {
	user := &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		Chan: make(chan string),
		Conn: conn,
	}

	go user.ListenMsg()

	return user
}

// ListenMsg 监听服务器发来的消息
func (u *User) ListenMsg() {
	for {
		msg := <-u.Chan
		_, _ = u.Conn.Write([]byte(strings.Join([]string{msg, "\n"}, "")))
	}
}
