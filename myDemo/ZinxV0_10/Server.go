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

type HelloZinxRouter struct {
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
	err := request.GetConnection().SendMsg(200, []byte("ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// hello zinx test 自定路由
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写ping...ping...
	fmt.Println("recv from client: msgID = ", request.GetMsgId(), ", data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("hello zinx"))
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

// 创建连接之后执行的钩子函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("====> DoConnBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnBegin")); err != nil {
		fmt.Println(err)
	}
	// 给当前连接设置一些属性
	fmt.Println("Set Conn Name ...")
	conn.SetProperty("Name", "Dehua")
	conn.SetProperty("Home", "https://google.com")
	conn.SetProperty("Career", "OD")
}

// 连接断开前需要执行的函数
func DoConnLost(connection ziface.IConnection) {
	fmt.Println("====>DoConnLost is Called ...")
	fmt.Println("ConnID =", connection.GetConnID(), "is Lost")

	// 获取连接属性
	if name, err := connection.GetProperty("Name"); err == nil {
		fmt.Println("Name =", name)
	}

	if home, err := connection.GetProperty("Home"); err == nil {
		fmt.Println("Home =", home)
	}

	if career, err := connection.GetProperty("Career"); err == nil {
		fmt.Println("Name =", career)
	}
}

func main() {
	// 创建一个server句柄,使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")

	// 注册连接的Hook钩子函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	// 给当前框架添加一个自定义的Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

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
