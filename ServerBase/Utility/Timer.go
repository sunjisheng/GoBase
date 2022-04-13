package Utility

import(
	"time"
)

type TimerFunc func(uint32, ...interface{})

type Timer struct {
	timerID uint32
	ticker *time.Ticker
	fun TimerFunc
	stop chan uint8
	args []interface{}
 }

var timers map[uint32]*Timer

func init() {
	timers = make(map[uint32]*Timer)
}

func  SetTimer(timerID uint32, ms uint32, fun TimerFunc, args ...interface{}) bool {
	if timers[timerID] != nil {
		return false
	}
	timer := new (Timer)
	timer.timerID = timerID
	timer.ticker = time.NewTicker(time.Millisecond * time.Duration(ms))
	timer.stop = make(chan uint8)
	timer.fun = fun
	timer.args = args
	timers[timerID] = timer
	go TimerLoop(timer)
	return true
}

func KillTimer(timerID uint32) {
	timer, ok := timers[timerID]
	if ok {
		delete(timers, timerID)
		timer.stop <- 1
	}
}

func TimerLoop(timer *Timer) {
	for {
		select {
		case <-timer.ticker.C:
			timer.fun(timer.timerID, timer.args...)
		case <-timer.stop:
			return
		}
	}
}