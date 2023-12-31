package ziface

/*
路由的抽象接口
路由里的数据都是irequest
*/
type IRouter interface {
	// 处理conn业务之前的hook
	PreHandle(request IRequest)
	// 处理业务conn的主hook
	Handle(request IRequest)
	// 处理conn业务之后的hook
	PostHandle(request IRequest)
}
