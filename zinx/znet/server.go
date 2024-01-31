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
	// Server的连接管理器
	ConnMgr ziface.IConnManager
	// Server 创建连接之后自动调用的Hook函数--OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// Server 销毁连接之后自动调用的Hook函数--OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, Listener at IP: %s, Port: %d is start\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn:%d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		// 0. 开启消息队列以及WorkerPool
		s.MsgHandler.StartWorkerPool()
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

			//设置最大连接个数的判断，如果超过最大连接的数量，则关闭此新连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too many Connections, MaxConn =", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定，得到连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将服务器的资源，状态以及开辟的连接进行回收
	fmt.Println("[STOP] Zinx Server name =", s.Name)
	s.ConnMgr.ClearConn()
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

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMgr
}

// 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVesion:   "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 注册OnConnStart钩子函数方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop钩子函数方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart钩子函数方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----->Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop钩子函数方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----->Call OnConnStop()")
		s.OnConnStop(conn)
	}
}
