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

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping "))
	if err != nil{
		fmt.Println("call back before ping error")
	}
}
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping "))
	if err != nil{
		fmt.Println("call back ping...ping...ping error")
	}
}
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping "))
	if err != nil{
		fmt.Println("call back after ping error")
	}
}

func main() {
	//1、创建一个server句柄
	s := znet.NewServer()
	//2、添加自定义的router
	s.AddRouter(&PingRouter{})
	//3、启动server
	s.Serve()
}
