package ziface

//IServer 定义一个服务器接口
type IServer interface {
	Start() // 启动服务器方法
	Stop()  // 停止服务器方法
	Serve() // 开启业务服务方法
	AddRouter(msgId uint32, router IRouter) // 给当前服务注册一个路由业务方法，供客户端连接处理使用
}
