package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"stuZinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn      //当前连接的socket TCP套接字
	ConnID       uint32            //当前连接的ID 也可以称作为SessionID，ID全局唯一
	isClosed     bool              //当前连接的关闭状态
	MsgHandler   ziface.IMsgHandle //消息管理MsgId和对应处理方法的消息管理模块
	ExitBuffChan chan bool         //告知该连接已经提出/停止的channel
	msgChan      chan []byte       //无缓冲管道，用户读、写两个goroutine之间的消息通信
}

//NewConnection 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte), //msgChan初始化
	}
	return c
}

//StartReader 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		//读取我们最大的数据到buf中
		/*buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buf err ", err)
			c.ExitBuffChan <- true
			continue
		}*/

		// 创建拆包对象
		dp := NewDataPack()
		//读取客户端的MsgHead
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}
		//拆包，得到msgId和dataLen放到msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}
		//根据dataLen读取data放到msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			data: msg,
		}
		//从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go c.MsgHandler.DoMsgHandler(&req)

	}
}

/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWriter() {

	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:
			//conn已经关闭
			return
		}
	}
}

//Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//1、开启从客户端读取数据的Goroutine
	go c.StartReader()
	//2、开启用户向客户端发送数据的Goroutine
	go c.StartWriter()

	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

//Stop 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	//1、如果当前连接已经关闭
	if c.isClosed {
		return
	}
	// todo 如果用户注册了该连接的关闭回调业务，那么在此刻应该显示调用
	//2、关闭socket连接
	c.isClosed = true
	c.Conn.Close()
	//3、通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true
	//4、关闭该连接全部管道
	close(c.ExitBuffChan)
}

//GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgChan <- msg   //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取

	return nil
}
