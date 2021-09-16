package utils

import (
	"encoding/json"
	"io/ioutil"
	"stuZinx/ziface"
)

// GlobalObj 存储一切有关zinx框架的全局参数，供其他模块使用一些参数也可以通过 用户根据zinx.json来配置
type GlobalObj struct {
	TcpServer     ziface.IServer //当前全局Server对象
	Host          string         //当前服务器主机IP
	TcpPort       int            //当前服务器监听端口
	Name          string         //当前服务器名称
	Version       string         //当前zinx版本号
	MaxPacketSize uint32         //数据的的最大长度
	MaxConn       int            //当前服务器主机允许的最大连接个数
}

//GlobalObject 定义一个全局的对象
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("demo/v0.4/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj {
		Name:    "ZinxServerApp",
		Version: "V0.4",
		TcpPort: 7777,
		Host:    "0.0.0.0",
		MaxConn: 12000,
		MaxPacketSize:4096,
	}
	//从配置文件中加载一些用户配置的参数
	//GlobalObject.Reload()
}

