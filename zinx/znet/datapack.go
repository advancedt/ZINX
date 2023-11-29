package znet

import (
	"ZINX/zinx/utils"
	"ZINX/zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

/*
	封包/拆包模块

直接面向TCP连接中的数据流，用于处理TCP粘包问题
*/
type DataPack struct {
}

// 拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 (4字节) + ID uint32 (4字节)
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节流的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 将DataLen写进Databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgID写进Databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据写入Databuf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// 只需要将包的Head信息读出来即可，之后根据head信息里的data长度在进行一次读
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压Head信息得到dataLen和MsgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断DataLen是否已经超出了允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv")
	}

	return msg, nil
}
