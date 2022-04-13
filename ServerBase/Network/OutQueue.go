package Network

import "sync"

type OutQueue struct {
	total uint32
	head uint32
	tail uint32
	outBufs []*OutBuf
	mutex sync.Mutex
}

func NewOutQueue(total uint32) *OutQueue {
	this := new(OutQueue)
	this.total = total
	this.head = 0
	this.tail = 0
	this.outBufs = make([]*OutBuf, this.total, this.total)
	return this
}

func (this *OutQueue) Reset() {
	this.head = 0
	this.tail = 0
}

func (this *OutQueue) IsEmpty() bool {
	this. mutex.Lock()
	defer this. mutex.Unlock()
	return this.head == this.tail
}

func (this *OutQueue) Count() uint32 {
	this. mutex.Lock()
	defer this. mutex.Unlock()
	return (this.tail + this.total - this.head) % this.total
}

func (this *OutQueue) Space() uint32 {
	this. mutex.Lock()
	defer this. mutex.Unlock()
	return (this.head + this.total - this.tail - 1) % this.total
}

func (this *OutQueue) Push(msg *OutBuf) bool {
	this. mutex.Lock()
	defer this. mutex.Unlock()
	if  (this.head + this.total - this.tail - 1) % this.total < 1 {
		return false
	}
	this.outBufs[this.tail] = msg
	this.tail = (this.tail + 1) % this.total
	return true
}

func (this *OutQueue) Pop() *OutBuf {
	if this.Count() == 0 {
		return nil
	}
	msg := this.outBufs[this.head]
	this.head = (this.head + 1) % this.total
	return msg
}