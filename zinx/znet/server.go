package znet

import (
	"ZINX/zinx/utils"
	"ZINX/zinx/ziface"
	"fmt"
	"net"
)

// iServer接口的实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的IP版本
	IPVesion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	// 当前Server 消息管理模块，绑定msgId和对应处理业务API关系
	MsgHandler ziface.IMsgHandle
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, Listener at IP: %s, Port: %d is start\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn:%d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		// 获取一个tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IPVesion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error :", err)
		}

		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVesion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVesion, "err", err)
			return
		}
		fmt.Println("Start Zinx server success", s.Name, "success, Listening")
		var cid uint32
		cid = 0
		// 阻塞的等待客户端进行连接，处理客户端的业务（Write && Read）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 将处理新连接的业务方法和conn进行绑定，得到连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	// 启动server服务功能
	s.Start()

	// 启动服务器后做一些额外的业务

	// 阻塞状态
	select {}
}

// 给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success")
}

// 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVesion:   "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
