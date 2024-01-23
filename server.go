package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer Server 创建
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

func (s *Server) Handler(conn net.Conn) {
	//	... 当前连接的业务
	fmt.Println("连接建立成功")
}

// Run Server 启动
func (s *Server) Run() {
	//	1.socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
	}
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
