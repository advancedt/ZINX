package znet

import (
	"ZINX/zinx/ziface"
	"fmt"
	"strconv"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 从request中找到MsgID
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not Found! Need register")
	}
	// 根据msgID调度对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// id 已经注册
		panic("repeat API, msgID = " + strconv.Itoa(int(msgId)))
	}
	// 添加msg与API的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api MsgID = ", msgId, " succ!")
}
