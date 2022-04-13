package Network

import "../Log"
type TimeFrame struct {
	head *TimeEvent
	tail *TimeEvent
	frame int32
}

func (this *TimeFrame) Insert(event *TimeEvent) {
	event.frame = this.frame
	if this.head == nil {
		this.head = event
		this.tail = event
	} else {
		this.tail.next = event
		event.prev = this.tail
		this.tail = event
	}
}

func (this *TimeFrame) Remove(event *TimeEvent) {
	if event.frame != this.frame {
		return
	}
	if event.prev != nil {
		event.prev.next = event.next
	} else {
		this.head = event.next
	}

	if event.next != nil {
		event.next.prev = event.prev
	} else {
		this.tail = event.prev
	}
	event.Reset()
}

func (this *TimeFrame) Fire() {
	event := this.head
	var next *TimeEvent = nil
	for event != nil {
		if event.owner != nil  {
			event.owner.OnTimeout()
			Log.WriteLog(Log.Log_Level_Info, "Peer OnTimeout this.frame=%d event.frame=%d", this.frame, event.frame)
		}
		next = event.next
		event.Reset()
		event = next
	}
	this.head = nil
	this.tail = nil
}