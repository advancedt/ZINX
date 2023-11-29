package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是负责测试datapack 拆包和封包的单元测试
func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 1. 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	// 创建一个go 承载负责从客户端处理业务
	go func() {
		// 2. 从客户端读取数据进行拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn2 net.Conn) {
				// 处理客户端请求
				// 拆包
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					// 把headData读到满为止
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err")
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg 是有数据的
						// 第二次读，根据head中的datalen读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("Server unpack err:", err)
							return
						}
						//完整的消息读取完毕
						fmt.Println("Receive MsgID:, ", msg.Id, ", datalen = ", msg.DataLen, ", data = ", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Client dial err:", err)
		return
	}
	// 创建一个封包对象dp
	dp := NewDataPack()
	//封装两个msg一同发送
	// 封装第一个msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}
	// 封装第二个msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data: []byte{
			'h', 'e', 'l', 'l', 'o',
		},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}
	// 将两个包粘在一起，一次性发送给服务端
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
