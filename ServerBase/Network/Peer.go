package Network

import (
	"../Log"
	"../Utility"
	"time"
)

type Peer struct {
	peerType uint32										 //类型
	peerID uint64                                          //ID
	ip string
	port uint32
	socketID uint32
	createTime time.Time
	timeout *TimeEvent
	peerData interface{}
}

func (this *Peer) Init(peerType uint32, peerID uint64, addr string)  {
	this.peerType = peerType
	this.peerID = peerID
	this.ip, this.port = Utility.SplitIPPort(addr)
	this.createTime = time.Now()
	this.timeout = new(TimeEvent)
	this.timeout.Init(this)
}

func (this *Peer) ResetPeerInfo() {
	this.peerType = PeerType_Unknow
	this.peerID = 0
	this.peerData = nil
}

func (this *Peer) IsValidPeerInfo() bool{
	if this.peerType != PeerType_Unknow && this.peerID > 0 {
		return true
	}
	return false
}

func (this *Peer) InitPending(socketID uint32, addr string)  {
	this.socketID = socketID
	this.ip, this.port = Utility.SplitIPPort(addr)
	this.createTime = time.Now()
	this.timeout = new(TimeEvent)
	this.timeout.Init(this)
}

func (this *Peer) IsAuthed() bool {
	if this.peerType == PeerType_Unknow {
		return false
	}
	return true
}

func (this *Peer) GetPeerType() uint32 {
	return this.peerType
}

func (this *Peer) GetPeerID() uint64 {
	return this.peerID
}

func (this *Peer) GetSocketID() uint32 {
	return this.socketID
}

func (this *Peer) GetPeerData() interface{} {
	return this.peerData
}

func (this *Peer) SetPeerData(data interface{}) {
	this.peerData = data
}

func (this *Peer)SetTimeout(frame int32) {
	TcpServer_Instance().timeSchedule.AddTimer(frame, this.timeout)
}

func (this *Peer)KillTimeout() {
	if this.timeout.IsScheduling() {
		TcpServer_Instance().timeSchedule.DelTimer(this.timeout)
	}
}

func (this *Peer)OnTimeout() {
	Log.WriteLog(Log.Log_Level_Error, "Peer OnTimeout peerType=%d peerID=%d", this.peerType, this.peerID)
	TcpServer_Instance().tcpManager.Close(this.socketID)
}