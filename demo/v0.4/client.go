package main

import (
	"fmt"
	"net"
	"stuZinx/utils"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(time.Second)
	//1、直接连接远程服务器，得到一个conn
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TcpPort))
	if err != nil {
		fmt.Println("client start error, exit!!!!!")
		return
	}
	//2、连接调用write写数据
	for {
		if _, err := conn.Write([]byte("Hello World")); err != nil {
			fmt.Println("write conn err", err)
			return
		}
		buf := make([]byte, utils.GlobalObject.MaxPacketSize)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
		}
		fmt.Printf("server call back: %s, cnt=%d\n", buf, cnt)
		time.Sleep(time.Second*3)

		}
}


