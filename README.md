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
# V0.5
## 定义一个消息的结构体Message
### 属性
1. 消息ID
2. 消息长度
3. 消息内容
### 方法
1. 获取消息ID
2. 获取消息的长度
3. 获取消息的内容
4. 设置消息ID
5. 设置消息长度
6. 设置消息内容
## 定义一个解决TCP粘包问题的封包拆包的模块
### 针对Message进行TLV格式的封装
1. 写入message的长度
2. 写入message的id
3. 写入message的内容
### 针对Message进行TLV格式的拆包
1. 先读取固定长度的head-->消息内容的长度和消息的类型
2. 再根据消息内容的长度，再次进行一次读取，从conn中读取消息的内容
## 将消息封装集成到zinx中
1. 将Message添加到Request属性中
2. 修改连接读取数据的方式
3. 给连接提供一个发包机制
## 使用zinx v0.5开发
# V0.6
我们之前在已经给Zinx配置了路由模式，但是很惨，之前的Zinx好像只能绑定一个路由的处理业务方法。显然这是无法满足基本的服务器需求的，那么现在我们要在之前的基础上，给Zinx添加多路由的方式。
## 消息管理模块(支持多路由业务api调度管理)
### 属性
1. 集合-消息ID对应的router的关系map
### 方法
1. 根据msgID来索引调度路由方法 DoMsgHandler()
2. 添加路由方法到map集合中 AddRouter()
## 将消息管理模块集成到zinx框架中
1. 将server模块中的Router属性 替换成MsgHandler属性
2. 将server之前的AddRouter修改 调用MsgHandler的AddRouter
3. 将connection模块Router属性 替换成MsgHandler，修改初始化Connection方法
4. 修改Connection的之前调度Router的业务 替换换成MsgHandler调度，修改StartReader方法
## 使用zinx v0.6开发
# V0.7
1. 添加一个Reader和Writer之间通信的channel
2. 添加一个Writer Goroutine
3. Reader由之前直接发送给客户端 改成发送给通信Channel
4. 启动Reader和Writer一同工作
# V0.8
## 消息队列及worker工作池实现
### 创建一个消息队列
#### MsgHandler消息管理模块，增加属性
1. 消息队列
2. worker工作池的数量
### 创建多任务worker的工作池并启动
1. 创建一个worker工作池
   1. 根据workerPoolSize的数量去创建worker
   2. 每个worker都应该用go去承载
      1. 阻塞等待与当前worker对应的channel的消息
      2. 一旦有消息到来，worker应该处理当前消息对应的业务，调用MsgHandler
### 修改之前的发送消息，全部改成 把消息发送给消息队列和worker工作池来处理
1. 定义一个方法，将消息发送给消息队列工作池的方法
   1. 保证每个worker所收到的任务是均衡的，让哪个worker处理就将消息发送给对应的taskQueue即可
   2. 将消息发送到对应的channel
## 将消息队列机制集成到zinx框架中
### 开启并调用消息队列及工作池
1. 保证工作池只有一个，需要在启动服务的时候开启
### 将从客户端处理的消息，发送给当前的worker来处理
1. 在已经处理完拆包，得到request请求交给工作池来处理