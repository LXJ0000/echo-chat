package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户
	OnlineMap map[string]*User
	MapLock   sync.RWMutex
	//接收消息用于广播
	MessageChan chan string
}

// NewServer Server 创建
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:          ip,
		Port:        port,
		OnlineMap:   make(map[string]*User),
		MessageChan: make(chan string)}
}

// Handler ... 处理当前连接的业务
func (s *Server) Handler(conn net.Conn) {

	user := NewUser(conn, s)

	user.Online()

	//
	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)

			// 客户端主动关闭连接
			if n == 0 {
				user.OffLine()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read error: ", err)
				return
			}

			// get user's msg and remove '\n'
			msg := string(buf[:n-1])

			user.DoMsg(msg)

			isLive <- true
		}
	}()

	//阻塞,避免断开连接
	for {
		select {
		// 重置定时器
		case <-isLive:
		// 定时器
		case <-time.After(time.Second * 10):
			user.SendMsg("超时下线")
			// 资源清理
			close(user.Chan)
			conn.Close()
			return
		}
	}

}

// BroadCast 接收客户端消息
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := strings.Join([]string{"[", user.Addr, "]", user.Name, ": ", msg}, "")
	s.MessageChan <- sendMsg
}

func (s *Server) ListenMsg() {
	for {
		msg := <-s.MessageChan

		s.MapLock.Lock()
		for _, user := range s.OnlineMap {
			user.Chan <- msg
		}
		s.MapLock.Unlock()
	}
}

// Run Server 启动
func (s *Server) Run() {
	//	1.socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
	}

	//启动监听Msg广播
	go s.ListenMsg()

	//	2. defer close
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	for {
		//	3. accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept error: ", err)
			continue
		}
		//	4. do handler
		go s.Handler(conn)
	}
}
