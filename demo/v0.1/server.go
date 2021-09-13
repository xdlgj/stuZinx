package main

import "stuZinx/znet"

func main() {
	//1、创建一个server句柄
	s := znet.NewServer("[tcp server v1.0]")
	//2、启动server
	s.Serve()
}
