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
	request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {

}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {

}
func main() {
	// 创建一个server句柄,使用Zinx的api
	s := znet.NewServer("[zinx V0.3]")
	// 启动server
	s.Serve()
}
