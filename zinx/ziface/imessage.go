package ziface

/*
将请求的消息封装到message中，定义一个抽象层模块
*/
type IMessage interface {
	// 获取消息id
	GetMsgId() uint32
	//获取消息的长度
	GetMsgLen() uint32
	// 获取消息的内容
	GetData() []byte
	// 设置消息的ID
	SetMsgId(uint32)
	// 设置消息的长度
	SetMsgLen(uint32)
	// 设置消息的内容
	SetData([]byte)
}
