package znet

import (
	"fmt"
	"net"
	"stuZinx/ziface"
)

//Server IServer的接口实现，定义一个Server服务器模块
type Server struct {
	Name      string //服务器名称
	IPVersion string //tcp4 or other
	IP        string //服务绑定的IP地址
	Port      int    //服务绑定的端口
}

//NewServer 初始化Server模块
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

//Start 启动服务器方法
func (s *Server) Start() {
	fmt.Printf("[Start] Server name: %s listener at IP: %s, Port %d, is starting\n", s.Name, s.IP, s.Port)
	//开启一个协程去做服务端listener业务
	go func() {
		//1、获取一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error", err)
			return
		}
		//2、监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " error", err)
			return
		}
		fmt.Println("start tcp server success, ", s.Name, "success listening...")
		//3、阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞回返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 已经与和客户端建立连接，做一些业务，做一个基本的最大512字节长度的回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buf err", err)
						continue
					}
					fmt.Printf("receive client buf %s, cnt %d\n", buf, cnt)
					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf error", err)
						continue
					}
				}
			}()
		}
	}()
}

//Stop 停止服务器方法
func (s Server) Stop() {
	// todo 将一些服务器的资源、状态或者一些已经开启的连接信息进行停止或回收
}

//Serve 开启业务服务方法
func (s Server) Serve() {
	s.Start() // 启动server的服务功能
	// todo 做一些启动服务器之后的额外业务
	// 阻塞，否则主go退出，listener的go将会退出
	select {}
}
