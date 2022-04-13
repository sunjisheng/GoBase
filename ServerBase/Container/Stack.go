package Container

type Stack struct {
	values []interface{}
}

func NewStack(cap uint32) *Stack {
	this := new(Stack)
	this.values = make([]interface{}, 0, cap)
	return this
}

func (this *Stack) Capacity() int  {
	return cap(this.values)
}

func (this *Stack) Size() int {
	return len(this.values)
}

func (this *Stack) Push(value interface{}) {
	this.values = append(this.values, value)
}

func (this *Stack) Pop() interface{} {
	len := len(this.values)
	if len <= 0 {
		return nil
	}
	value := this.values[len - 1]
	this.values[len - 1] = nil
	this.values = this.values[:len - 1]
	return value
}

func (this *Stack) Clear()  {
	for i := 0; i < len(this.values); i++ {
		this.values[i] = nil
	}
	this.values = this.values[0:0]
}
