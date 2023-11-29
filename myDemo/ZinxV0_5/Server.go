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
	znet.BaseRouter
}

// Test PreHandle
//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("call back before ping error")
//	}
//}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写ping...ping...
	fmt.Println("recv from client: msgID = ", request.GetMsgId(), ", data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// Test PostHandle
//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
//	if err != nil {
//		fmt.Println("call after before ping error")
//	}
//}

func main() {
	// 创建一个server句柄,使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")
	// 给当前框架添加一个自定义的Router
	s.AddRouter(&PingRouter{})
	// 启动server
	s.Serve()
}

/*
只能给一个Server 添加一个 Router
*/
/*
1. 创建server句柄
2. 添加一个自定义的router
3. 启动server
需要继承BaseRouter
实现PreHandle， Handle，PostHandle
*/
