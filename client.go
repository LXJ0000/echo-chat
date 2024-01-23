package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Conn       net.Conn
	Name       string
}

func NewClient(serverIP string, serverPort int) *Client {
	conn, err := net.Dial("tcp", strings.Join([]string{serverIP, ":", strconv.Itoa(serverPort)}, ""))
	if err != nil {
		fmt.Println("net Dial error: ", err)
		return nil
	}
	return &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		Conn:       conn,
	}
}

var (
	serverIP   string
	serverPort int
)

func init() {
	// 命令行解析
	// ./client -p 127.0.0.1 -port 8888
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "服务器IP，默认127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口，默认8888")
}

func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("==========> 连接服务器失败...")
		return
	}

	fmt.Println("==========> 连接服务器成功...")
	select {}
}
