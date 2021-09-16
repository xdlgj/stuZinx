package ziface

type IMessage interface {
	GetDataLen() uint32 //获取消息长度
	GetMsgId() uint32   //获取消息ID
	GetData() []byte    //获取消息内容
	SetMsgId(uint32)    //设置消息ID
	SetData([]byte)     // 设置消息内容
	SetDataLen(uint32)  //设置消息数据长度
}
