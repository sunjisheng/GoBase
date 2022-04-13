package Container

type CircleBuf struct {
	maxSize uint32
	buffer []byte
	head uint32
	tail uint32
}

func NewCircleQueue(maxSize uint32) *CircleBuf {
	this := new(CircleBuf)
	this.buffer = make([]byte, maxSize, maxSize)
	this.maxSize = maxSize
	this.head = 0
	this.tail = 0
	return this
}

func (this *CircleBuf) Reset() {
	this.head = 0
	this.tail = 0
}

func (this *CircleBuf) IsEmpty() bool {
	return this.head == this.tail
}

func (this *CircleBuf) Size() uint32 {
	return (this.tail + this.maxSize - this.head) % this.maxSize
}

func (this *CircleBuf) Head() uint32 {
	return this.head
}

func (this *CircleBuf) Tail() uint32 {
	return this.tail
}

func (this *CircleBuf) Space() uint32 {
	return (this.head + this.maxSize - this.tail - 1) % this.maxSize
}

func (this *CircleBuf) Write(buf []byte, len uint32) bool {
	if this.Space() < len {
		return false
	}
	toEndSize := this.maxSize - this.tail
	if toEndSize >= len {
		copy(this.buffer[this.tail:this.tail+len], buf[0:len])
	} else {
		copy(this.buffer[this.tail:this.tail+toEndSize], buf[0:toEndSize])
		copy(this.buffer[0:len-toEndSize], buf[toEndSize:len])
	}
	this.tail = (this.tail + len) % this.maxSize
	return true
}

func (this *CircleBuf) Peek(buf []byte, len uint32) bool {
	if this.Size() < len	{
		return false
	}
	toEndSize := this.maxSize - this.head
	if toEndSize >= len	{
		copy(buf[0:len], this.buffer[this.head:this.head+len])
	} else	{
		copy(buf[0:toEndSize], this.buffer[this.head:this.head+toEndSize])
		copy(buf[toEndSize:len], this.buffer)
	}
	return true
}

func (this *CircleBuf) Read(buf []byte, len uint32) bool {
	if !this.Peek(buf, len)	{
		return false
	}
	this.head = (this.head + len) % this.maxSize
	return true
}

func (this *CircleBuf) Skip(len uint32) {
 	this.head = (this.head + len) % this.maxSize;
}

func (this *CircleBuf) OnWrited(len uint32) {
	this.tail = (this.tail + len) % this.maxSize;
}

func (this *CircleBuf) GetReadSlice_Len(len uint32) []byte {
	return this.buffer[this.head:this.head+len]
}

func (this *CircleBuf) GetReadSlice() []byte {
	if(this.head < this.tail) {
		return this.buffer[this.head:this.tail]
	} else {
		return this.buffer[this.head:]
	}
}

func (this *CircleBuf) GetWriteSlice() []byte {
	toEndSize := this.maxSize - this.tail
	space := (this.head + this.maxSize - this.tail - 1) % this.maxSize
	if toEndSize > space {
		toEndSize = space
	}
	return this.buffer[this.tail:this.tail + toEndSize]
}

func (this *CircleBuf) Arrange() {
	len := this.tail - this.head
	copy(this.buffer[0:len],this.buffer[this.head:this.tail])
	this.head = 0
	this.tail = this.head + len
}
