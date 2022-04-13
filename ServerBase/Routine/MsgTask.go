package Routine

import (
	"google.golang.org/protobuf/proto"
)

type IMsgHandler interface {
	OnMsg(peerType uint32, peerID uint64, msgID uint32, msg proto.Message) bool
}

type MsgTask struct {
	PeerType    uint32
	PeerID      uint64
	MsgID       uint32
	Msg         proto.Message
	MsgHandler IMsgHandler
}

func (this *MsgTask)Execute() {
	this.MsgHandler.OnMsg(this.PeerType, this.PeerID, this.MsgID, this.Msg)
}
