# echo chat
回顾 Golang 知识点

## 启动
```bash
git clone https://github.com/LXJ0000/echo-chat.git
cd echo-chat
go mod tidy
go build -o server server.go main.go user.go
go build -o client client.go
./server
./client
```