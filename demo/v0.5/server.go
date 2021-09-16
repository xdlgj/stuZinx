package main

import (
	"fmt"
	"stuZinx/ziface"
	"stuZinx/znet"
)

//PingRouter ping test自定义路由
type PingRouter struct {
	znet.BaseRequest
}

//func (pr *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping "))
//	if err != nil{
//		fmt.Println("call back before ping error")
//	}
//}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//先读取客户端的数据，再回写
	fmt.Println("receive from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	//回写数据
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

/*
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping "))
	if err != nil{
		fmt.Println("call back after ping error")
	}
}
*/

func main() {
	//1、创建一个server句柄
	s := znet.NewServer()
	//2、添加自定义的router
	s.AddRouter(&PingRouter{})
	//3、启动server
	s.Serve()
}
