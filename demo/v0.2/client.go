package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(time.Second)
	//1、直接连接远程服务器，得到一个conn
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
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
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
		}
		fmt.Printf("server call back: %s, cnt=%d\n", buf, cnt)
		time.Sleep(time.Second*3)

		}
}


