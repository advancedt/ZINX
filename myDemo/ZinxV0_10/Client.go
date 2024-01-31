package main

import (
	"ZINX/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("Client start")
	time.Sleep(1 * time.Second)
	//1. 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		// 发送封包的msg消息 MsgId == 0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx Client Test Message")))
		if err != nil {
			fmt.Println("Pack error", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		// 服务器应该回复一个msg数据 msgId 1 ping...ping..
		// 先读取流中的head部分 得到msg ID和datalen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error")
			break
		}

		//将二进制Head拆包到msg结构体中
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack message head error")
			break
		}
		// Msg 里有数据
		if msgHead.GetMsgLen() > 0 {
			// 再根据DataLen进行第二次读取，将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("recv Server Msg : ID = ", msg.Id, ", len = ", msg.DataLen, ", data = ", string(msg.Data))
		}

		// CPU 阻塞
		time.Sleep(1 * time.Second)
	}

}
