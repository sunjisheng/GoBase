package Routine

import (
	"time"
)

type RoutinePool struct{
	inch    chan ITask
	count   uint32
	stop    chan struct{}
	routines []*Routine
}

func (this *RoutinePool) Init(count uint32){
	this.inch = make(chan ITask, 256)
	this.stop = make(chan struct{})
	this.count = count
	this.routines = make([]*Routine, count)
	for i := uint32(0); i < count; i++ {
		routine := new(Routine)
		routine.InitWithArg("RoutinePool", this.Run, routine)
		this.routines[i] = routine
		routine.Gid = g_RoutineMonitor.Add2Monitor(routine)
	}
}

func (this *RoutinePool) Stop(){
	//禁止重复调用
	close(this.stop)
}

func (this *RoutinePool) AddTask (task ITask){
	this.inch <- task
}

func (this *RoutinePool) Run(arg interface{}){
	ticker := time.NewTicker(time.Second * 60)
	var stop bool = false
	for !stop {
		select {
		case task := <- this.inch:
			if task != nil {
				task.Execute()
			}
		case <-ticker.C:
			arg.(*Routine).Heartbeat(uint32(time.Now().Unix()))
		case <-this.stop:
			stop = true
		}
	}
}

