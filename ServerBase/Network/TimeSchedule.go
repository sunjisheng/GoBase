package Network

import "time"
import "../Routine"

type TimerInfo struct {
	frame int32
	event *TimeEvent
}

type TimeSchedule struct {
	Routine.Routine
	frames []TimeFrame
	frameCount int32
	curFrame int32
	timerChan  chan *TimerInfo
}

func (this *TimeSchedule) Init(frameCount int32) {
	this.frameCount = frameCount
	this.curFrame = 0
	this.frames = make([]TimeFrame, frameCount)
	this.timerChan = make(chan *TimerInfo, 256)
	for frame:= int32(0); frame < frameCount; frame++ {
		this.frames[frame].frame = frame
	}
	this.Routine.Init("TimeSchedule", this.Run)
}

func (this *TimeSchedule) Run() {
	ticker := time.NewTicker(time.Second)
	var stop bool = false
	for !stop {
		select {
		case info := <-this.timerChan:
			if info != nil {
				if info.event.IsScheduling() {
					//Log.WriteLog(Log.Log_Level_Info, "KillTimer frame=%d event=%d", info.event.frame, info.event)
					this.frames[info.event.frame].Remove(info.event)
				}
				if info.frame >= 0 {
					info.frame = (this.curFrame + info.frame) % this.frameCount
					this.frames[info.frame].Insert(info.event)
					//Log.WriteLog(Log.Log_Level_Info, "AddTimer frame=%d event=%d", info.frame, info.event)
				}
			}
		case <-ticker.C:
			this.Tick()
		case <-this.StopChan:
			stop = true
		}
	}
}

func (this *TimeSchedule) AddTimer(frame int32, event *TimeEvent) {
	this.timerChan <- &TimerInfo{frame:frame, event:event}
}

func (this *TimeSchedule) DelTimer(event *TimeEvent) {
	this.timerChan <- &TimerInfo{frame:-1, event:event}
}

func (this *TimeSchedule) Tick() {
	this.frames[this.curFrame].Fire()
	this.curFrame = (this.curFrame + 1) % this.frameCount
}

