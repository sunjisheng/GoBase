package Container

type ListNode struct {
	prev *ListNode
	next *ListNode
	value interface{}
}

func (this *ListNode) Value() interface{}{
	return this.value
}

func (this *ListNode) Prev() *ListNode{
	return this.prev
}

func (this *ListNode) Next() *ListNode{
	return this.next
}

type List struct {
	head *ListNode
	tail *ListNode
	count uint32
}

func NewList() *List {
	this := new(List)
	return this
}
//增加到尾部
func (this *List) Push_Back(value interface{}) {
	node := new(ListNode)
	node.value = value
	if  this.tail != nil {
	    this.tail.next = node
		node.prev = this.tail
		this.tail = node
	}	else	{
		this.head = node
		this.tail = node
	}
	this.count++
}

//增加到头部
func (this *List) Push_Front(value interface{}) {
	node := new(ListNode)
	node.value = value
	node.next = this.head
	if  this.head != nil {
		this.head.prev = node
	}	else	{
		this.tail = node
	}
	this.head = node
	this.count++
}

func (this *List) Pop_Front() interface{} {
	if  this.head == nil	{
		 return nil;
	}
	value := this.head.value
	next := this.head.next
	if next != nil {
		next.prev = nil
	}	else	{
		this.tail = nil
	}
	this.head = next
	this.count--
	return value
}

func (this *List) Pop_Back() interface{} {
	if  this.tail == nil	{
		return nil;
	}
	value := this.tail.value
	prev := this.tail.prev
	if prev != nil {
		prev.next = nil
	}	else	{
		this.head = nil
	}
	this.tail = prev
	this.count--
	return value
}

func (this *List) Erase(node *ListNode) *ListNode{
	if node == this.head {
		this.head = node.next
	}
	if node == this.tail {
		this.tail = node .prev
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	this.count--
	return node.next
}

func (this *List) Clear() {
	this.head = nil
	this.tail = nil
	this.count = 0
}

func (this *List) Empty() bool {
	if this.head == nil {
		return true
	} else {
		return false
	}
}

func (this *List) Begin() *ListNode {
	return this.head
}
