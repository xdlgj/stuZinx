package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"stuZinx/utils"
	"stuZinx/ziface"
)

//DataPack 封包拆包类实例，暂时不需要成员
type DataPack struct{}

//NewDataPack 封包拆包实例初始化
func NewDataPack() *DataPack {
	return &DataPack{}
}

//GetHeadLen 获取包头的长度
func (dp *DataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}

//Pack 封装方法(压缩数据)
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//1、创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//2、写入dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//3、写入msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//4、写入data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//Unpack 拆包方法(解压数据)
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//1、创建一个从输入二进制数据ioReader
	dataBuff := bytes.NewReader(binaryData)
	//2、只解压head的信息，得到dataLen和msgID
	msg := &Message{}
	//3、读取dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//4、读取msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//5、判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}
	//6、这里只需把head的数据拆包出来就可以了，然后在通过head的长度，再从conn读取一次数据
	return msg, nil
}
