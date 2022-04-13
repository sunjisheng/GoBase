package Network

import (
	"../Log"
	"../Routine"
	"fmt"
	"net"
	"sync"
)

const (
	InvalidSocketID = 0
	Max_Socket_Count = 65535
)

type TcpManager struct {
	Routine.Routine
	sockets map[uint32]*TcpSocket
	socketsMutex sync.Mutex
	nextSocketID uint32
	listener *net.TCPListener
	tcpServer ITcpServer
}

func (this *TcpManager) Init() {
	this.nextSocketID = 1
	this.sockets = make(map[uint32]*TcpSocket)
	this.Routine.Init("TcpManager",this.Loop)
	Log.WriteLog(Log.Log_Level_Info, "TcpManager Init OK")
}

func (this *TcpManager) NextSocketID() uint32 {
	this.socketsMutex.Lock()
	var socketID uint32
	for {
		socketID = this.nextSocketID
		_, ok := this.sockets[socketID]
		this.nextSocketID++
		if this.nextSocketID == InvalidSocketID {
			this.nextSocketID = 1
		}
		if !ok {
			break
		}
	}
	this.socketsMutex.Unlock()
	return socketID
}

func (this *TcpManager) AddSocket(socket *TcpSocket) {
	this.socketsMutex.Lock()
	this.sockets[socket.socketID] = socket
	this.socketsMutex.Unlock()
}

func (this *TcpManager) GetSocket(socketID uint32) *TcpSocket{
	this.socketsMutex.Lock()
	socket, ok := this.sockets[socketID]
	this.socketsMutex.Unlock()
	if ok {
		return socket
	}
	return nil
}

func (this *TcpManager) DelSocket(socketID uint32) *TcpSocket{
	this.socketsMutex.Lock()
	socket, ok := this.sockets[socketID]
	defer this.socketsMutex.Unlock()
	if ok {
		delete(this.sockets, socketID)
	}
	return socket
}

func (this *TcpManager) Listen(port uint32) {
	go this.ListenLoop(port)
	Log.WriteLog(Log.Log_Level_Info, "TcpManager Listen Port %d", port)
}

func (this *TcpManager) Loop() {
	var stop bool = false
	for !stop {
		select {
		case <-this.StopChan:
			this.socketsMutex.Lock()
			for _, socket := range this.sockets {
				socket.Close()
			}
			this.socketsMutex.Unlock()
			stop = true
		}
	}
}

func (this *TcpManager) ListenLoop(port uint32) {
	addrstr := fmt.Sprintf(":%d",port)
	addr, err := net.ResolveTCPAddr("tcp", addrstr) //创建 tcpAddr数据
	if err != nil {
		return
	}
	this.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "TcpManager Listen Error:%s", err.Error())
		return
	}
	go this.HandleAccept()
}

func (this *TcpManager) HandleAccept() {
	for {
		conn, err := this.listener.AcceptTCP()
		if err != nil {
			Log.WriteLog(Log.Log_Level_Error, "Accept Error:", err.Error())
		}
		socketID := this.NextSocketID()
		socket := NewTcpSocket(socketID, conn, this.tcpServer)
		this.AddSocket(socket)
		this.tcpServer.OnSocketConnect(socket.socketID, socket.ip, socket.port)
		socket.Start()
	}
}

func (this *TcpManager) Connect(ip string, port uint32) {
	go this.InnerConnect(ip, port)
}

func (this *TcpManager) InnerConnect(ip string, port uint32) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "TcpManager Connect Error ip=%s port=%d errorMsg=%s", ip, port, err.Error())
		this.tcpServer.OnSocketClose(InvalidSocketID, ip, port)
		return
	}
	socketID := this.NextSocketID()
	socket := NewTcpSocket(socketID, conn, this.tcpServer)
	this.AddSocket(socket)
	this.tcpServer.OnSocketConnect(socket.socketID, socket.ip, socket.port)
	socket.Start()
}

func (this *TcpManager) SendMsg(socketID uint32, msgID uint32, msgBody []byte) bool {
	socket := this.GetSocket(socketID)
	if socket != nil {
		return socket.SendMsg(msgID, msgBody)
	}
	return false
}

func (this *TcpManager) Close(socketID uint32) {
	socket := this.DelSocket(socketID)
	if socket != nil {
		socket.Close()
	}
}

