package Network

import (
	"../BaseProtocol"
	"../Database"
	"../Log"
	"../Redis"
	"../Utility"
	"../Routine"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	PeerType_Unknow = 0
	Max_ServerType = 16
	PeerType_Client = Max_ServerType
	Max_ServerName = 32
	Min_ServerID = 1
	Max_ServerID = 16
)

type ILogicServer interface {
	//逻辑服务的网络接口
	OnConnect(peerType uint32, peerID uint64)
	OnClose(peerType uint32, peerID uint64)
	CreateMsg(msgID uint32) proto.Message
    OnSocketMsg(socketID uint32, msgID uint32, msgBody []byte) bool
	OnServerStop()
}

var instance *TcpServer

func TcpServer_Instance() *TcpServer{
	return instance
}

type TcpServer struct {
	serverType       uint32                              //服务器类型
	serverID         uint64                              //服务器ID
	serverName       string                              //服务名
	port             uint32                              //端口
	tcpManager       *TcpManager                         //TcpManager
	timeSchedule     *TimeSchedule                       //超时定时器
	logicServer      ILogicServer                        //逻辑服务
    globalIni        *Utility.IniFile                    //全局配置
	localIni         *Utility.IniFile                    //本地配置
	serverPeers      [Max_ServerType][Max_ServerID]*Peer //服务器节点
	maxServerIDs     [Max_ServerType]uint64              //服务类型节点数量
	addr2PeerMap     map[uint64]*Peer                    //本服务主动连接的地址
	clientPeers      map[uint64]*Peer                    //客户端节点
	clientMutex      sync.Mutex                          //客户端节点锁
	socketID2PeerMap map[uint32]*Peer                    //SocketID->Peer
	socketID2PeerMutex sync.Mutex                        //SocketID->Peer锁
	pendingPeers map[uint32]*Peer                        //尚未Auth的Peer
	pendingMutex sync.Mutex                              //未Auth的Peer锁
}

func (this *TcpServer) SetServerName(serverName string) {
	this.serverName = serverName
}

func (this *TcpServer) GetServerName() string {
	return this.serverName
}

func (this *TcpServer) SetServerType(serverType uint32) {
	this.serverType = serverType
}

func (this *TcpServer) GetServerType() uint32 {
	return this.serverType
}

func (this *TcpServer) SetServerID(serverID uint64) {
    this.serverID = serverID
}

func (this *TcpServer) GetServerID() uint64 {
	return this.serverID
}

func (this *TcpServer) GetPort() uint32 {
	return this.port
}

func (this *TcpServer) SetLogicServer(logicServer ILogicServer) {
	this.logicServer = logicServer
}

func (this *TcpServer) GetLocalIni() *Utility.IniFile {
	return this.localIni
}

func (this *TcpServer) GetGlobalIni() *Utility.IniFile {
	return this.globalIni
}

func (this *TcpServer) Init() bool {
	instance = this
	rand.Seed(time.Now().Unix())
	Log.StartLog(this.serverName, Log.Log_Level_Debug)	//开启日志
	Log.WriteLog(Log.Log_Level_Info, "%s Start...", this.serverName)
	//加载全局配置
	if !this.LoadGlobalConfig() {
		return false
	}
	//加载本地配置
	if !this.LoadLocalConfig() {
		return false
	}
	Routine.RoutineMonitor_Instance().Init()
	this.addr2PeerMap = make(map[uint64]*Peer)
	this.clientPeers = make(map[uint64]*Peer)
	this.socketID2PeerMap = make(map[uint32]*Peer)
	this.pendingPeers = make(map[uint32]*Peer)
	this.tcpManager = new(TcpManager)
	this.tcpManager.tcpServer = this
	this.tcpManager.Init()
	this.timeSchedule = new(TimeSchedule)
	this.timeSchedule.Init(120)
	if this.port > 0 {
		this.tcpManager.Listen(this.port)
	}
	return true
}

func (this *TcpServer) Join() {
	c := make(chan os.Signal)
	signal.Notify(c)
	var stop = false
	var signal os.Signal
	for !stop {
		signal = <-c //阻塞直至有信号传入
		switch signal {
		case syscall.SIGHUP:
		case syscall.SIGINT:
		case syscall.SIGILL:
		case syscall.SIGABRT:
		case syscall.SIGTERM:
		case syscall.SIGQUIT:
			stop = true
			break;
		default:
		}
	}
	this.timeSchedule.Stop()
	this.logicServer.OnServerStop()
	Log.WriteLog(Log.Log_Level_Info, "Server %s receive signal %d, Exit ...", this.serverName, signal)
	time.Sleep(time.Second * 2)
}

func (this *TcpServer) LoadGlobalConfig() bool {

	fullPath := fmt.Sprintf("%s/Global.ini", Utility.GetCurrentPath())
	this.globalIni = new(Utility.IniFile)
	err :=this.globalIni.LoadFile(fullPath)
	if err !=nil {
		Log.WriteLog(Log.Log_Level_Error, "LoadGlobalConfig Error")
		return false
	}
	return true
}

func (this *TcpServer) LoadLocalConfig() bool {
	fullPath := fmt.Sprintf("%s/%s.ini", Utility.GetCurrentPath(), this.serverName)
	this.localIni = new(Utility.IniFile)
	err :=this.localIni.LoadFile(fullPath)
	if err !=nil {
		Log.WriteLog(Log.Log_Level_Error, "LoadLocalConfig %s Error", fullPath)
		return false
	}
	this.serverID = uint64(this.localIni.GetInt("Settings", "ServerID"))
	this.port = uint32(this.localIni.GetInt("Settings", "Port"))
	return true
}

func (this *TcpServer) InitMySql(dbIndex uint32, section string) bool{
	var iniFile *Utility.IniFile = nil
	iniFile = this.GetLocalIni()
	addr := iniFile.GetString(section, "Addr")
	dbName := iniFile.GetString(section, "DBName")
	userName := iniFile.GetString(section, "UserName")
	pwd := iniFile.GetString(section, "Pwd")
	if !Database.MySql_Instance(dbIndex).Open(addr, dbName, userName, pwd) {
		return false
	}
	Log.WriteLog(Log.Log_Level_Info, "MySql Open(%s %s) success", addr, dbName)
	return true
}

func (this *TcpServer) InitRedis(rdbIndex uint32, section string) bool{
	var iniFile *Utility.IniFile = nil
	if rdbIndex == 0 {
		iniFile = this.GetGlobalIni()
	} else {
		iniFile = this.GetLocalIni()
	}

	addr := iniFile.GetString(section, "Addr")
	pwd := iniFile.GetString(section, "Pwd")
	if !Redis.Redis_Init(rdbIndex, addr, pwd) {
		return false
	}
	Log.WriteLog(Log.Log_Level_Info, "Redis_Init(%s) success", addr)
	return true
}

func (this *TcpServer) Connect(serverType uint32, serverName string) {
	if this.serverType == PeerType_Unknow {
		Log.WriteLog(Log.Log_Level_Error, "serverType Is Not Set!");
		return;
	}
	if serverType <= PeerType_Unknow || serverType >= Max_ServerType{
		Log.WriteLog(Log.Log_Level_Error, "serverType Is Not Set!");
		return;
	}
    //相同类型服务连接，serverID
	for serverID := uint64(Min_ServerID); serverID <= uint64(Max_ServerID); serverID++ {
		if this.serverType == serverType && this.serverID == serverID {
			this.SetMaxServerID(serverType, serverID);
			continue
		}
	     if !this.AddPeer(serverType, serverID, serverName) {
		     break;
		 }
	}
	for serverID := uint64(Min_ServerID); serverID <= uint64(Max_ServerID); serverID++ {
		peer := this.serverPeers[serverType][serverID - 1]
		if peer != nil {
			this.tcpManager.Connect(peer.ip, peer.port)
		}
	}
}

func (this *TcpServer) AddPeer(serverType uint32, serverID uint64, serverName string) bool {
	key := strconv.FormatUint(serverID, 10)
	addr := this.globalIni.GetString(serverName, key)
	if len(addr) == 0 {
		return false
	}
	peer := new (Peer)
	peer.Init(serverType, serverID, addr)
	this.serverPeers[serverType][serverID - 1] = peer
	this.SetMaxServerID(serverType, serverID);
	addrKey := Utility.HashStr(addr)
	this.addr2PeerMap[addrKey] = peer;
	return true;
}

func (this *TcpServer)  SendMsg(socketID uint32, msgID uint32, message proto.Message) {
	msgBody ,err := proto.Marshal(message)
	if err != nil {
		return
	}
	this.tcpManager.SendMsg(socketID, msgID, msgBody)
}

func (this *TcpServer) Send2Server(serverType uint32, serverID uint64 , msgID uint32, msg proto.Message) {
	if serverType <= PeerType_Unknow || serverType >= Max_ServerType {
		return
	}
	if serverID < Min_ServerID || serverID > Max_ServerID {
		return
	}
	peer := this.serverPeers[serverType][serverID - 1]
	if peer !=nil {
		this.SendMsg(peer.socketID, msgID, msg)
	}
}

func (this *TcpServer) SendData2Server(serverType uint32, serverID uint64 , msgID uint32, msgBody []byte) {
	if serverType <= PeerType_Unknow || serverType >= Max_ServerType {
		return
	}
	if serverID < Min_ServerID || serverID > Max_ServerID {
		return
	}
	peer := this.serverPeers[serverType][serverID - 1]
	if peer !=nil {
		this.tcpManager.SendMsg(peer.socketID, msgID, msgBody)
	}
}

func (this *TcpServer) Broadcast2Server(serverType uint32, msgID uint32, msg proto.Message) {
	if serverType <= PeerType_Unknow || serverType >= Max_ServerType {
		return
	}
	for serverID := uint64(Min_ServerID); serverID <= this.maxServerIDs[serverType]; serverID++ {
		peer:= this.serverPeers[serverType][serverID - 1]
		if peer != nil {
			this.SendMsg(peer.socketID, msgID, msg)
		}
	}
}

func (this *TcpServer)  Send2Client(peerID uint64, msgID uint32, msg proto.Message) {
	peer := this.GetClientPeer(peerID)
	if peer != nil {
		this.SendMsg(peer.socketID, msgID, msg)
	}
}

func (this *TcpServer)  SendData2Client(peerID uint64, msgID uint32, msgBody []byte) {
	peer := this.GetClientPeer(peerID)
	if peer != nil {
		this.tcpManager.SendMsg(peer.socketID, msgID, msgBody)
	}
}

func (this *TcpServer)  CloseServer(serverType uint32, serverID uint64) {
	peer:= this.serverPeers[serverType][serverID - 1]
	if peer != nil {
		this.tcpManager.Close(peer.socketID)
	}
}

func (this *TcpServer)  CloseClient(peerID uint64) {
	peer := this.GetClientPeer(peerID)
	if peer != nil {
		this.tcpManager.Close(peer.socketID)
	}
}

func (this *TcpServer)  CloseAndDelClientPeer(peerID uint64) {
	peer := this.GetClientPeer(peerID)
	if peer != nil {
		this.tcpManager.Close(peer.socketID)
		peer.KillTimeout()
		peer.ResetPeerInfo()
		this.DelClientPeer(peerID)
	}
}

func (this *TcpServer) AddSocketID2Peer(socketID uint32, peer *Peer){
	this.socketID2PeerMutex.Lock()
	this.socketID2PeerMap[socketID] = peer
	this.socketID2PeerMutex.Unlock()
}

func (this *TcpServer) DelSocketID2Peer(socketID uint32) *Peer{
	this.socketID2PeerMutex.Lock()
	peer,ok := this.socketID2PeerMap[socketID]
	defer this.socketID2PeerMutex.Unlock()
	if ok {
		delete(this.socketID2PeerMap, socketID)
		return peer
	}
	return nil
}

func (this *TcpServer) GetPeerBySocketID(socketID uint32) *Peer {
	this.socketID2PeerMutex.Lock()
	peer, ok := this.socketID2PeerMap[socketID]
	this.socketID2PeerMutex.Unlock()
	if ok {
		return peer
	}
	return nil
}

func (this *TcpServer)  SetPeerInfo(socketID uint32, peerType uint32, peerID uint64, peerData interface{}) {
	this.socketID2PeerMutex.Lock()
	peer, ok := this.socketID2PeerMap[socketID]
	this.socketID2PeerMutex.Unlock()
	if ok {
		peer.peerType = peerType
		peer.peerID = peerID
		peer.peerData = peerData
		if peerType == PeerType_Client {
			this.AddClientPeer(peerID, peer)
			Log.WriteLog(Log.Log_Level_Info,"AddClientPeer PeerID=%d SocketID=%d", peerID, socketID)
		}
		this.RemovePendingPeer(socketID)
	}
}

func (this *TcpServer)  SetMaxServerID(peerType uint32, peerID uint64) {
	if peerID > this.maxServerIDs[peerType] {
		this.maxServerIDs[peerType] = peerID
	}
}

func (this *TcpServer)  GetServerID_Mod(peerType uint32, key uint64) uint64 {
	if this.maxServerIDs[peerType] == 0 {
		return 0
	}
	return (key % this.maxServerIDs[peerType]) + 1
}

func (this *TcpServer)  GetServerID_Rand(peerType uint32) uint64 {
	if this.maxServerIDs[peerType] == 0 {
		return 0
	}
	return uint64(rand.Intn(int(this.maxServerIDs[peerType])) + 1)
}

func (this *TcpServer)  GetServerID_Valid(serverType uint32) uint64 {
	if this.maxServerIDs[serverType] == 0 {
		return 0
	}
	for serverID := uint64(Min_ServerID); serverID <= this.maxServerIDs[serverType]; serverID++ {
		if this.IsServerConnected(serverType, serverID) {
			return serverID
		}
	}
	return 0
}

func (this *TcpServer) AddPendingPeer(peer *Peer) {
	this.pendingMutex.Lock()
	this.pendingPeers[peer.socketID] = peer
	this.pendingMutex.Unlock()
}

func (this *TcpServer)RemovePendingPeer(socketID uint32) {
	this.pendingMutex.Lock()
	delete(this.pendingPeers, socketID)
	this.pendingMutex.Unlock()
}

func (this *TcpServer)ClosePendingPeer(socketID uint32) {
	this.pendingMutex.Lock()
	_, ok :=this.pendingPeers[socketID]
	this.pendingMutex.Unlock()
	if ok {
		this.tcpManager.Close(socketID)
	}
}

func (this *TcpServer)AddClientPeer(peerID uint64, peer *Peer) {
	this.clientMutex.Lock()
	this.clientPeers[peerID] = peer
	this.clientMutex.Unlock()
}

func (this *TcpServer)GetClientPeer(peerID uint64)*Peer {
	this.clientMutex.Lock()
	peer,ok := this.clientPeers[peerID]
	this.clientMutex.Unlock()
	if ok {
		return peer
	}
	return nil
}

func (this *TcpServer)DelClientPeer(peerID uint64) {
	this.clientMutex.Lock()
	delete(this.clientPeers, peerID)
	this.clientMutex.Unlock()
}

func (this *TcpServer) OnSocketConnect(socketID uint32, ip string, port uint32) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	addrKey := Utility.HashStr(addr)
	peer, ok := this.addr2PeerMap[addrKey]
	if ok {  //主动连接
		peer.socketID = socketID
		msg := new(BaseProtocol.ServerConnectMsg)
		msg.ServerType = this.serverType
		msg.ServerID = this.serverID
		this.SendMsg(socketID, uint32(BaseProtocol.MsgID_ServerConnect), msg)
		//上层处理
		this.logicServer.OnConnect(peer.peerType, peer.peerID)
	} else {
		peer = new(Peer)
		peer.InitPending(socketID, addr)
		this.AddPendingPeer(peer)
		//超过一定时间仍未认证，断开连接
		time.AfterFunc(time.Second * 30, func() {
			this.ClosePendingPeer(socketID)
		})
	}
	this.AddSocketID2Peer(socketID, peer)
}

func (this *TcpServer) OnSocketClose(socketID uint32, ip string ,port uint32){
	peer := this.DelSocketID2Peer(socketID)
	if peer != nil {
		//上层处理
		if peer.IsValidPeerInfo() {
			this.logicServer.OnClose(peer.peerType, peer.peerID)
		}
		if peer.peerType > PeerType_Unknow && peer.peerType < Max_ServerType {
			//重连
			addr := fmt.Sprintf("%s:%d", ip, port)
			addrKey := Utility.HashStr(addr)
			_, ok := this.addr2PeerMap[addrKey]
			if ok {
				this.tcpManager.Connect(ip, port)
			} else {
				this.serverPeers[peer.peerType][peer.peerID - 1] = nil;
			}
		}	else if peer.peerType == PeerType_Client {
			peer.KillTimeout()
			this.DelClientPeer(peer.peerID)
		}
	}	else	{
		this.RemovePendingPeer(socketID)
		Log.WriteLog(Log.Log_Level_Error, "TcpServer %s:%d OnClose...", ip, port)
		time.AfterFunc(time.Second * 5, func() {
			Log.WriteLog(Log.Log_Level_Error, "TcpServer %s:%d Reconnect...", ip, port)
			this.tcpManager.Connect(ip, port)
		})
	}
}

func (this *TcpServer) OnSocketMsg(socketID uint32, msgID uint32, msgBody []byte) bool {
	if (this.OnBaseMsg(socketID, msgID, msgBody)) {
		return true
	} else {
		return this.logicServer.OnSocketMsg(socketID, msgID, msgBody)
	}
}

func (this *TcpServer) OnBaseMsg(socketID uint32, msgID uint32, msgBody []byte) bool {
	if msgID == uint32(BaseProtocol.MsgID_ServerHeartbeat) {
		//心跳消息
		msg := new(BaseProtocol.ServerHeartBeatMsg)
		err := proto.Unmarshal(msgBody, msg)
		if err != nil {
			Log.WriteLog(Log.Log_Level_Error, "TcpServer ProcessBaseMsg return FALSE, msgID=%d", msgID)
			return false
		}
		if msg.Heartbeat > uint32(0) {
			msg.Heartbeat--
		}
		this.SendMsg(socketID, uint32(BaseProtocol.MsgID_ServerHeartbeat), msg)
		return true
	} else if msgID == uint32(BaseProtocol.MsgID_ServerConnect) {
		//建立连接
		msg := new(BaseProtocol.ServerConnectMsg)
		err := proto.Unmarshal(msgBody, msg)
		if err != nil {
			Log.WriteLog(Log.Log_Level_Error, "TcpServer ProcessBaseMsg return FALSE, msgID=%d", msgID)
			return false
		}
		oldPeer := this.serverPeers[msg.ServerType][msg.ServerID - 1]
		if oldPeer != nil { //禁止新连接
			Log.WriteLog(Log.Log_Level_Info, "Conflict Peer ServerType=%d ServerID=%lld", msg.ServerType, msg.ServerID);
			this.tcpManager.Close(socketID)
		} else {
			peer := this.GetPeerBySocketID(socketID)
			if peer == nil {
				Log.WriteLog(Log.Log_Level_Error, "TcpServer ProcessBaseMsg SocketID=%d Not Found", socketID)
				return false
			}
			peer.peerType = msg.ServerType
			peer.peerID = msg.ServerID
			this.serverPeers[peer.peerType][peer.peerID-1] = peer
			this.SetMaxServerID(peer.peerType, peer.peerID)
			this.RemovePendingPeer(socketID)
			//上层处理
			this.logicServer.OnConnect(peer.peerType, peer.peerID)
		}
		return true
	}
	return false
}

func (this *TcpServer) IsServerConnected(peerType uint32, serverID uint64) bool{
	if peerType <= PeerType_Unknow || peerType >= Max_ServerType {
		return false
	}
	if serverID < Min_ServerID || serverID > Max_ServerID {
		return false
	}
	peer := this.serverPeers[peerType][serverID - 1]
	if peer == nil || peer.socketID <= 0 {
		return false
	}
	return true
}