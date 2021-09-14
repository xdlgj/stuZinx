# V0.1
基础的server需要如下属性和方法
## 方法
1. 启动服务器
   1. 创建addr
   ```go
    net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
   ```
   2. 创建listener
   ```go
    net.ListenTCP(s.IPVersion, addr)
    ```
   3. 接收客户端数据
    ```go
    listener.AcceptTCP()
    ```
2. 停止服务器
3. 运行服务器
4. 初始化server
## 属性
1. 服务名称
2. tcp版本
3. 监听IP
4. 监听端口
## 测试
```
go run demo/v0.1/server.go  // 启动服务器
go run demo/v0.1/client.go  // 启动客户端
```
# V0.2简单的连接封装和业务绑定
## 连接模块
### 方法
1. 启动连接 Start()
2. 停止连接 Stop()
3. 获取当前连接的conn对象(套接字)  GetTCPConnection() *net.TCPConn
4. 获取连接ID GetConnID() uint32
5. 获取客户端连接的地址和端口 RemoteAddr() net.Addr
6. 发送数据的方法 Send(data []byte) error
7. 连接所绑定的处理业务的函数类型 type HandFunc func(*net.TCPConn, []byte, int) error
### 属性
1. socket TCP套接字 
2. 连接ID
3. 当前连接的状态
4. 与当前连接所绑定的处理业务方法
5. 等待连接被退出的channel