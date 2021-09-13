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