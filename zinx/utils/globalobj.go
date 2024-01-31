package utils

import (
	"ZINX/zinx/ziface"
	"encoding/json"
	"os"
)

/*
存储有关Zinx框架的全局参数，供其他模块进行使用
一些参数可以通过zinx.json 由用户进行配置
*/
type GlobalObj struct {
	// Server 配置
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口‘
	Name      string         // 当前服务器的名称

	// Zinx配置
	Version          string // 当前Zinx版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前允许的数据包的最大值
	WorkerPoolSize   uint32 // 当前业务Worker工作池的GoRoutine的数量
	MaxWorkerTaskLen uint32 // Zinx框架允许用户最多开辟多少个Worker
}

// 定义一个全局对外的GlobalObj
var GlobalObject *GlobalObj

// 当前从zinx.json加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

// 初始化对象
func init() {
	// 如果配置文件没有加载的默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.10",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   // Worker工作池队列的个数
		MaxWorkerTaskLen: 1024, // 每个Worker对应的消息队列的任务的数量的最大值
	}
	// 尝试从conf/zinx.json加载自定义参数
	GlobalObject.Reload()
}
