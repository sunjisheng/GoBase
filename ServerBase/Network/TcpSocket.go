package Network

import (
	"../Container"
	"../Log"
	"../Utility"
	"net"
	"sync"
)

const (
	Socket_Act_Send = 1
	Socket_Act_Close = 2
)

const (
	Msg_Header_Len = 8
	InBuf_MaxSize = 65536
	OutBuf_MaxSize = 65536
)

type ITcpServer interface {
	OnSocketConnect(socketID uint32, ip string, port uint32)
	OnSocketClose(socketID uint32, ip string, port uint32)
	OnSocketMsg(socketID uint32, msgID uint32, msgBody []byte) bool
}

type MsgHeader struct{
	msgID uint32
	len uint32
}

type OutBuf struct {
	msgID uint32
	msgBody []byte
}

func (this *OutBuf) SetMsgID(msgID uint32) {
	this.msgID = msgID
}

func (this *OutBuf) SetMsgBody(msgBody []byte) {
	this.msgBody = msgBody
}

type TcpSocket struct {
	socketID uint32
	conn *net.TCPConn
	ip string
	port uint32
	inBuf *Container.CircleBuf
	outBuf *Container.CircleBuf
	outMutex sync.Mutex
	tcpServer ITcpServer
	actChan chan uint8
}

func  NewTcpSocket(socketID uint32, conn *net.TCPConn, tcpServer ITcpServer) *TcpSocket {
	socket := new(TcpSocket)
	socket.tcpServer = tcpServer
	socket.inBuf = Container.NewCircleQueue(InBuf_MaxSize)
	socket.outBuf = Container.NewCircleQueue(OutBuf_MaxSize)
	socket.actChan = make(chan uint8, 10)
	socket.conn = conn
	conn.SetNoDelay(true)
	conn.SetLinger(1)
	socket.ip, socket.port = Utility.SplitIPPort(conn.RemoteAddr().String())
	socket.socketID = socketID
	return socket
}

func (this *TcpSocket)SocketID() uint32{
	return this.socketID
}

func (this *TcpSocket)Reset() {
	this.port = 0
	this.ip = ""
	this.inBuf.Reset()
	this.outBuf.Reset()
}

func (this *TcpSocket)Start() {
	go this.HandleRecv()
	go this.HandleSend()
}

func  (this *TcpSocket) SendMsg(msgID uint32, msgBody []byte) bool {
	len := uint32(len(msgBody))
	header_buf := make([]byte, Msg_Header_Len)
	Utility.WriteUInt32(header_buf, 0, msgID)
	Utility.WriteUInt32(header_buf, 4, len)

	this.outMutex.Lock()
	space := this.outBuf.Space()
	if space < Msg_Header_Len + len {
		this.outMutex.Unlock()
		Log.WriteLog(Log.Log_Level_Error, "SendMsg outBuf msgtotallen=%d space=%d not enough", Msg_Header_Len + len, space)
		return false
	}
	this.outBuf.Write(header_buf, Msg_Header_Len)
	this.outBuf.Write(msgBody, len)
	this.outMutex.Unlock()
	this.actChan<-Socket_Act_Send
	return true
}

func  (this *TcpSocket)HandleRecv() {
	defer  this.OnClose()
	defer this.conn.Close()
	for {
		writeSlice := this.inBuf.GetWriteSlice()
		len,err := this.conn.Read(writeSlice)
		if err != nil{
			Log.WriteLog(Log.Log_Level_Info, "HandleRecv err:%s", err.Error())
			return
		}
		this.inBuf.OnWrited(uint32(len))
		if !this.ProcessMsg() {
			return
		}
		this.inBuf.Arrange()
	}
}

func  (this *TcpSocket)HandleSend() {
	defer	this.conn.Close()
	for {
		select {
		case act := <-this.actChan:
			if act == Socket_Act_Close {
				return
			} else if act == Socket_Act_Send {
				if this.InnerSend() == false {
					return
				}
			}
		}
	}
}

func  (this *TcpSocket) InnerSend() bool{
	for this.outBuf.Size() > 0 {
		buf := this.outBuf.GetReadSlice()
		_,err := this.conn.Write(buf)
		if err != nil  {
			Log.WriteLog(Log.Log_Level_Info, "InnerSend err:%s", err.Error())
			return false
		}
		this.outBuf.Skip(uint32(len(buf)))
	}
	return true
}

func  (this *TcpSocket)ProcessMsg() bool{
	for {
		header_buf := make([]byte, Msg_Header_Len)
		if !this.inBuf.Peek(header_buf, Msg_Header_Len) {
			break
		}
		msgID := Utility.ReadUInt32(header_buf, 0)
		len := Utility.ReadUInt32(header_buf, 4)
		if this.inBuf.Size() < Msg_Header_Len + len {
			break
		}
		this.inBuf.Skip(Msg_Header_Len)
		msgBody := this.inBuf.GetReadSlice_Len(len)
		if !this.tcpServer.OnSocketMsg(this.socketID, msgID, msgBody) {
			return false
		}
		this.inBuf.Skip(len)
	}
	return true
}

func  (this *TcpSocket)Close() {
	this.actChan <- Socket_Act_Close
}

func  (this *TcpSocket)OnClose() {
	this.tcpServer.OnSocketClose(this.socketID, this.ip, this.port)
}


