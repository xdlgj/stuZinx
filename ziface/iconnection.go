package ziface

import "net"

// IConnection 定义连接接口
type IConnection interface {
	Start()                         //启动连接，让当前连接开始工作
	Stop()                          //停止连接，结束当前连接状态
	GetTCPConnection() *net.TCPConn //从当前连接获取原始的socket TCPConn
	GetConnID() uint32              //获取当前连接ID
	RemoteAddr() net.Addr           //获取远程客户端地址信息
	SendMsg(msgId uint32, data []byte) error         //发送数据，将数据发送给远程客户端
}

// HandFunc 这个是所有conn链接在处理业务的函数接口，
//第一参数是socket原生链接，
//第二个参数是客户端请求的数据，
//第三个参数是客户端请求的数据长度。
//这样，如果我们想要指定一个conn的处理业务，只要定义一个HandFunc类型的函数，然后和该链接绑定就可以了。
type HandFunc func(*net.TCPConn, []byte, int) error
