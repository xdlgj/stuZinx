package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//TestDataPack 负责测试dataPack拆包 封装的单元测试
func TestDataPack(t *testing.T) {
	/*模拟服务器*/
	go func() {
		//1、创建socketTCP
		listener, err := net.Listen("tcp", "127.0.0.1:7777")
		if err != nil{
			fmt.Println("server listen err: ", err)
			return
		}
		//2、从客户端读取数据
		for {
			conn, err := listener.Accept()
			if err != nil{
				fmt.Println("server accept err: ", err)
			}
			go func(conn net.Conn) {
				//0、创建封包拆包对象dp
				dp := NewDataPack()
				for {
					//1、先读出流中的head部分
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) //ReadFull会把headData填充满为止
					if err != nil {
						fmt.Println("read head error")
						break
					}
					//2、将headData字节流 拆包到msg中
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						//3、再次读取data数据
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())
						//4、根据dataLen从IO中读取字节流
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}
				}
			}(conn)
		}

	}()

	/*模拟客户端*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	//创建一个封包对象dp
	dp := NewDataPack()
	//封装两个数据包
	msg1 := &Message{
		Id: 0,
		DataLen: 5,
		Data: []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}
	msg2 := &Message{
		Id: 1,
		DataLen: 7,
		Data: []byte("world!!"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}
	//将sendData1和sendData2拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)
	//向服务器端写入数据
	conn.Write(sendData1)
	//客户端阻塞
	select {}
}
