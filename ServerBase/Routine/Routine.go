package Routine

import (
	"../Utility"
	"time"
)

type Routine struct {
	StopChan chan int
	Name string
	Gid uint64
	LastHeartbeat uint32
}

func (this *Routine)Init(name string, fn func()) {
	this.StopChan = make(chan int, 1)
	this.Name = name
	go Utility.StartRoutine(fn)
}

func (this *Routine)InitWithArg(name string, fn func(arg interface{}), arg interface{}) {
	this.StopChan = make(chan int, 1)
	this.Name = name
	go Utility.StartRoutine_Arg(fn, arg)
}

func (this *Routine)Stop() {
	this.StopChan <- 1
}

func (this *Routine)Heartbeat(now uint32) {
	this.LastHeartbeat = now
}

func (this *Routine)IsDead() bool {
	if uint32(time.Now().Unix()) - this.LastHeartbeat > 180 && this.LastHeartbeat > 0{
		return true
	} else {
		return false
	}
}

