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
# V0.3
## 基础router模块
### Request请求封装
将连接和数据绑定在一起
#### 属性
1. 连接 conn ziface.IConnection
2. 请求数据 data []byte
#### 方法
1. 得到当前连接 GetConnection() IConnection
2. 得到当前数据 GetData() []byte
### Router模块
router实际上的作用就是，服务端应用可以给Zinx框架配置当前连接的处理业务方法，之前的Zinx-V0.2我们的Zinx框架处理连接请求的方法是固定的，现在是可以自定义，并且有3种接口可以重写。
* Handle: 处理当前连接的主业务函数
* PreHandle: 如果需要在主业务函数之前有前置业务，可以重写这个方法
* PostHandle: 如果需要在主业务函数之后有后置业务，可以重写这个方法
#### 抽象的Router
```go
type IRouter interface {
	PreHandle(request IRequest) // 在处理conn业务之前的钩子方法
	Handle(request IRequest)//在处理conn业务的方法
	PostHandle(request IRequest)//在处理conn业务之后的钩子方法
}
```
#### 具体的Router
```go
//BaseRequest 实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRequest struct {}
// PreHandle 这里之所以BaseRouter的方法都为空，
//是因为有的Router不希望有PreHandle或PostHandle
//所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化
func (br *BaseRequest) PreHandle(request ziface.IRequest)  {} // 在处理conn业务之前的钩子方法
func (br *BaseRequest) Handle(request ziface.IRequest)     {} //在处理conn业务的方法
func (br *BaseRequest) PostHandle(request ziface.IRequest) {} //在处理conn业务之后的钩子方法

```
### 集成router模块
1. IServer增加路由添加功能
2. Server类增加Router成员
3. Connection类绑定一个Router成员
4. 在Connection中调用已经注册的Router处理业务
## 测试
1. 创建一个server句柄，使用zinxapi
2. 给当前的zinx框架添加一个自定义的router，需要继承BaseRouter，实现相应的方法
3. 启动server
# V0.4
添加全局配置
1. 服务器应用/conf/zinx.json(用户进行填写)
2. 创建一个Zinx的去全局配置模块utils/globalobj.go
   1. 读取用户配置好的zinx.json文件 --->globalobj对象中
   2. 提供一个全局的GlobalObject对象
3. 将zinx框架中全部的硬代码，用globalobj里面的参数进行替换
4. 使用zinxv0.4开发