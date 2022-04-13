package Routine

import (
	"../Log"
	"sync"
	"time"
)

type RoutineMonitor struct {
	Routine
	routines map[uint64]*Routine
	currentRoutineID uint64
	mutex sync.Mutex
}

var g_RoutineMonitor RoutineMonitor

func RoutineMonitor_Instance() *RoutineMonitor{
	return &g_RoutineMonitor
}

func (this *RoutineMonitor) Init() {
	this.currentRoutineID = 1
	this.routines = make(map[uint64]*Routine)
	this.Routine.Init("RoutineMonitor", this.Run)
}

func (this *RoutineMonitor) Add2Monitor(routine *Routine) uint64 {
	this.mutex.Lock()
	gid := this.currentRoutineID
	this.currentRoutineID++
	this.routines[gid] = routine
	this.mutex.Unlock()
	return gid
}

func (this *RoutineMonitor) DelFromMonitor(gid uint64) {
	this.mutex.Lock()
	delete(this.routines, gid)
	this.mutex.Unlock()
}

func (this *RoutineMonitor) Run() {
	ticker := time.NewTicker(time.Minute)
	var stop bool = false
	for !stop {
		select {
		case <-ticker.C:
			this.Tick()
			break
		case <-this.StopChan:
			stop = true
		}
	}
}

func (this *RoutineMonitor) Tick() {
	this.mutex.Lock()
	for gid, routine := range this.routines {
		if routine.IsDead() {
			Log.WriteLog(Log.Log_Level_Error, "Routine %s Is Dead gid=%d ", routine.Name, gid)
		}
	}
	this.mutex.Unlock()
}
