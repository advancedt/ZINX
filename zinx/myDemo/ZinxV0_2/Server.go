package main

import "ZINX/zinx/znet"

/*
	基于Zinx框架开打的服务器应用端程序
*/

func main() {
	// 创建一个server句柄,使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	// 启动server
	s.Serve()
}
