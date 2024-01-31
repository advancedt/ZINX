package znet

import (
	"ZINX/zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

/*
连接管理模块
*/

type ConnManager struct {
	// 管理的连接信息集合
	connections map[uint32]ziface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("ConnID = ", conn.GetConnID(), "add to ConnManager successfully: conn num = ", connMgr.Len())
}

// 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("ConnID = ", conn.GetConnID(), "remove from ConnManager successfully: conn num = ", connMgr.Len())
}

// 根据ConnID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.Unlock()

	if conn, ok := connMgr.connections[connID]; ok {
		// 找到
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 得到当前的连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有的连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear all connections successfully, conn num = ", connMgr.Len())
}
