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
	Router ziface.IRouter //当前Server由用户绑定的回调router，也就是Server注册的连接对应的处理业务
}

//NewServer 初始化Server模块
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router: nil,
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
		//已经监听成功
		fmt.Println("start tcp server success, ", s.Name, "success listening...")
		// todo 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0
		//3、启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//3.2 todo Server.Start()  设置服务器最大连接控制，如果超过最大连接，那么则关闭此新的连接

			//3.3 处理该连接请求的 业务方法，此时应该有handle 和conn是绑定的
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			//3.4 启动当前连接的处理业务
			go dealConn.Start()
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

//AddRouter 给当前服务注册一个路由业务方法，供客户端连接处理使用
func (s *Server)AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success! " )
}
