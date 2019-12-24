package main

import (
	"fmt"
	"io"
	"math/rand"
	"mmo-go/zinx/znet"
	"net"
	"time"
)

type Message struct {
	Len   uint32
	MsgId uint32
	Data  []byte
}

//type TcpClient struct {
//	conn     net.Conn
//	X        float32
//	Y        float32
//	Z        float32
//	V        float32
//	Pid      int32
//	isOnline chan bool
//}
//
//func (this *TcpClient) Unpack(headdata []byte) (head *Message, err error) {
//	headbuf := bytes.NewReader(headdata)
//
//	head = &Message{}
//
//	// 读取Len
//	if err = binary.Read(headbuf, binary.LittleEndian, &head.Len); err != nil {
//		return nil, err
//	}
//
//	// 读取MsgId
//	if err = binary.Read(headbuf, binary.LittleEndian, &head.MsgId); err != nil {
//		return nil, err
//	}
//
//	// 封包太大
//	//if head.Len > MaxPacketSize {
//	//	return nil, packageTooBig
//	//}
//
//	return head, nil
//}
//
//func (this *TcpClient) Pack(msgId uint32, dataBytes []byte) (out []byte, err error) {
//	outbuff := bytes.NewBuffer([]byte{})
//	// 写Len
//	if err = binary.Write(outbuff, binary.LittleEndian, uint32(len(dataBytes))); err != nil {
//		return
//	}
//	// 写MsgId
//	if err = binary.Write(outbuff, binary.LittleEndian, msgId); err != nil {
//		return
//	}
//
//	//all pkg data
//	if err = binary.Write(outbuff, binary.LittleEndian, dataBytes); err != nil {
//		return
//	}
//
//	out = outbuff.Bytes()
//
//	return
//}
//
//func (this *TcpClient) SendMsg(msgID uint32, data proto.Message) {
//
//	// 进行编码
//	binary_data, err := proto.Marshal(data)
//	if err != nil {
//		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
//		return
//	}
//
//	sendData, err := this.Pack(msgID, binary_data)
//	if err == nil {
//		this.conn.Write(sendData)
//	} else {
//		fmt.Println(err)
//	}
//
//	return
//}

func main() {
	fmt.Println("a client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("client start err:", err)
		return
	}
	rand.Seed(time.Now().UnixNano())

	for {
		// 发送封包的message消息
		dp := znet.NewDataPack()
		msg, err := dp.Pack(znet.NewMsgPackage(uint32(rand.Intn(2)), []byte("Zinx 0.5 client test mesage")))
		if err != nil {
			fmt.Println("pack error:", err)
			break
		}
		if _, err = conn.Write(msg); err != nil {
			fmt.Println("write error;", err)
			break
		}

		// 服务器回复一个message的数据，msgID：1 pingping
		// 1.先读取流中的head部分，得到id和 len
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}
		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack mesHead error", err)
			break
		}

		if msgHead.GetDataLen() > 0 {
			// 2.在根据len 第二次读取 data
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("client read msg error", err)
				break
			}
			fmt.Println("------> Recv Server MsgID = ", msg.GetMsgId(), "data = ", string(msg.Data))

		}

		//	cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
