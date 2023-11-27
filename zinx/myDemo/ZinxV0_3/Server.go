package main

import (
	"ZINX/zinx/ziface"
	"ZINX/zinx/znet"
	"fmt"
)

/*
	基于Zinx框架开打的服务器应用端程序
*/

// ping test 自定义路由
type PingRouter struct {
}

// Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping...\n"))
	if err != nil {
		fmt.Println("call ping ping error")
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call after before ping error")
	}
}
func main() {
	// 创建一个server句柄,使用Zinx的api
	s := znet.NewServer("[zinx V0.3]")
	// 给当前框架添加一个自定义的Router
	s.AddRouter(&PingRouter{})
	// 启动server
	s.Serve()
}
