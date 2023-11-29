package znet

import (
	"ZINX/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

// 连接模块
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前连接的状态
	isClosed bool

	// 告知当前连接已经退出/停止 channel
	ExitChan chan bool

	//该连接处理的方法
	Router ziface.IRouter
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader GoRoutine is running")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取客户端的数据到buf，最大512byte
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		// 创建一个拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的Msg Head

		// 得到Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		//拆包，得到msgID和msgDataLen放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Unpack error:", err)
			break
		}
		//根据Datalen再次读取data，放在msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("Read msg data error", err)
				break
			}
		}

		msg.SetData(data)

		// 得到当前连接的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		// 从路由中，找到注册绑定的Conn对应的Router调用

	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start(), ConnID =", c.ConnID)

	// 启动从当前连接的读数据的业务
	go c.StartReader()
	// TODO 启动从当前写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop(), ConnID =", c.ConnID)
	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Connection Error", err)
		panic(err)
	}

	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connction closed when send msg")
	}
	// 将data进行封包 MsgDataLen|MsgId|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error, msg id = ", msgId)
		return errors.New("Pack error msg")
	}
	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Wiite msg id:", msgId, "error:", err)
		return errors.New("conn Write error")
	}
	return nil
}
