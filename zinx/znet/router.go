package znet

import "ZINX/zinx/ziface"

// 实现Router时，先嵌入BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

/*
	之所以BaseRouter的方法都为空，目的是有的Router不希望有PreHandle或者PostHandle业务
	所以Router全部继承BaseRouter的好处是 不需要实现PreHandle或者PostHandle
*/
// 处理conn业务之前的hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理业务conn的主hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// 处理conn业务之后的hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
