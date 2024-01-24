package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//server := NewServer("127.0.0.1", 8888)
	//server.Run()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入字符串：")
	scanner.Scan()
	line1 := scanner.Text()
	fmt.Println("请输入字符串：")
	scanner.Scan()
	line2 := scanner.Text()

	fmt.Println("1:", line1)
	fmt.Println("2:", line2)

}
