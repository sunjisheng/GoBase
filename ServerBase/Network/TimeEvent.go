package Network

type TimerOwner interface {
	OnTimeout()
}

type TimeEvent struct {
	owner TimerOwner
	prev *TimeEvent
	next *TimeEvent
	frame int32
}

func (this *TimeEvent) Init(owner TimerOwner){
	this.owner = owner
	this.frame = -1
	this.prev = nil
	this.next = nil
}

func (this *TimeEvent) IsScheduling()  bool {
	if this.frame >= 0 {
		return true
	} else {
		return false
	}
}

func (this *TimeEvent) Reset() {
	this.frame = -1
	this.prev = nil
	this.next = nil
}