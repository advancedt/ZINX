package znet

import (
	"ZINX/zinx/utils"
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
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的Worker数量
	WorkerPoolSize uint32
}

// 初始化

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
		// 从全局配置中获取，也可以在配置文件中让用户进行设置
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// 启动一个Worker工作池(只能发生一次)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据WorkerPoolSize分别开启Worker，每个Worker用Go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个Worker被启动
		//1. 给当前的Worker对应的channel消息队列开辟空间
		// 第i个worker用第i个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前的Worker，阻塞等待消息从channel中进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) startOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID =", workerId, "is started...")
	// 不断阻塞对应的队列的消息
	for {
		select {
		//如果有消息过来，出来的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueqe,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1. 将消息平均分配给不同的Worker
	// 根据客户端建立的ConnID进行分配
	// 平均分配轮询法则
	WorkerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnId =", request.GetConnection().GetConnID(),
		"request MsgID =", request.GetMsgId(),
		"to WorkerID = ", WorkerId)
	//2. 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[WorkerId] <- request
}
