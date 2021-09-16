package ziface

//IMsgHandle 消息管路抽象层
type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 马上以阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
}
