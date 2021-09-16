package main

import (
	"fmt"
	"io"
	"net"
	"stuZinx/utils"
	"stuZinx/znet"
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
		//封装发送的数据
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx V0.6 Client1 Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
