package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Conn       net.Conn
	Name       string
	Flag       int
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
		Flag:       -1,
	}
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.Flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}

}

func (c *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (c *Client) Run() {
	for c.Flag != 0 {
		for c.menu() == false {
		}

		switch c.Flag {
		case 1:
			fmt.Println("1:todo")
		case 2:
			fmt.Println("2:todo")
		case 3:
			c.UpdateName()
		}
	}
}

// 处理server回应的消息， 直接显示到标准输出即可
func (client *Client) DealResponse() {
	//一旦client.conn有数据，就直接copy到stdout标准输出上, 永久阻塞监听
	io.Copy(os.Stdout, client.Conn)
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

	go client.DealResponse()

	fmt.Println("==========> 连接服务器成功...")
	client.Run()
}
