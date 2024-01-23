package main

import (
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

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("==========> 连接服务器失败...")
		return
	}

	fmt.Println("==========> 连接服务器成功...")
	select {}
}
