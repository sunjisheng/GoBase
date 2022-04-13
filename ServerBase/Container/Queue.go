package Container

type Queue struct {
	values []interface{}
}

func NewQueue(cap uint32) *Queue {
	this := new(Queue)
	this.values = make([]interface{}, 0, cap)
	return this
}

func (this *Queue) Capacity() int  {
	return cap(this.values)
}

func (this *Queue) Size() int {
	return len(this.values)
}

func (this *Queue) Push(value interface{}) {
	this.values = append(this.values, value)
}

func (this *Queue) Pop() interface{} {
	len := len(this.values)
	if len <= 0 {
		return nil
	}
	value := this.values[0]
	this.values[0] = nil
	this.values = this.values[1:len]
	return value
}

func (this *Queue) Clear()  {
	for i := 0; i < len(this.values); i++ {
		this.values[i] = nil
	}
	this.values = this.values[0:0]
}
