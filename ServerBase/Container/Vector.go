package Container

type Vector struct {
	values []interface{}
}

func NewVector(cap uint32) *Vector {
	this := new(Vector)
	this.values = make([]interface{}, 0, cap)
	return this
}

func (this *Vector) Capacity() int  {
	return cap(this.values)
}

func (this *Vector) Size() int {
	return len(this.values)
}

func (this *Vector) Push(value interface{}) {
	this.values = append(this.values, value)
}

func (this *Vector) Pop() interface{} {
	len := len(this.values)
	if len <= 0 {
		return nil
	}
	value := this.values[len - 1]
	this.values[len - 1] = nil
	this.values = this.values[:len - 1]
	return value
}

func (this *Vector) Clear()  {
	for i := 0; i < len(this.values); i++ {
		this.values[i] = nil
	}
	this.values = this.values[0:0]
}

func (this *Vector) Get(index uint32) interface{} {
	if index < 0 || index >= uint32(len(this.values)) {
		return nil
	}
	return this.values[index]
}

func (this *Vector) Set(index uint32, value interface{}) bool {
	if index < 0 || index >= uint32(len(this.values)) {
		return false
	}
	this.values[index] = value
	return true
}