package ziface

//IDataPack 封包数据和拆包数据
// 直接面向TCP连接中的数据流，为传输数据添加头部信息，用户处理TCP粘包问题
type IDataPack interface {
	GetHeadLen() uint32                // 获取包头的长度
	Pack(msg IMessage) ([]byte, error) // 封装方法
	Unpack([]byte) (IMessage, error)   // 拆包方法
}

