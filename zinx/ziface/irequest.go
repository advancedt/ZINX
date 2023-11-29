package ziface

/*
irequest接口实际上把客户端请求的连接信息 和 请求的数据 包装到了一个Request中
*/
type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
	// 得到请求的消息ID
	GetMsgId() uint32
}
