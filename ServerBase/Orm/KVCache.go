package Orm

import (
	"container/list"
	"sync"
)

type KVNode struct {
	key interface{}
	obj StoreObj
}

type KVCache struct {
	max   int
	list     *list.List
	elements map[interface{}]*list.Element
	mutex sync.Mutex
}

func(this *KVCache) Init(max int) {
	this.max = max
	if this.max == 0 {
		this.max = 10000
	}
	this.list = list.New()
	this.elements = make(map[interface{}]*list.Element)
}

func(this *KVCache) Get(key interface{}) (StoreObj,bool) {
	if this.elements == nil {
		return nil,false
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ele, ok := this.elements[key]
	if  ok {
		this.list.MoveToFront(ele)
		return ele.Value.(*KVNode).obj, true
	}
	return nil, false
}

func(this *KVCache) Set(key interface{}, obj StoreObj) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	e, ok := this.elements[key];
	if ok {
		e.Value.(*KVNode).obj = obj
		this.list.MoveToFront(e)
		return
	}
	ele := this.list.PushFront(&KVNode{key: key, obj: obj})
	this.elements[key] = ele
	if this.max != 0 && this.list.Len() > this.max {
		oldest := this.list.Back()
		if oldest != nil {
			this.list.Remove(oldest)
			node := e.Value.(*KVNode)
			delete(this.elements, node.key)
		}
	}
}

func(this *KVCache) Delete(key interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ele, ok := this.elements[key]
	if ok {
		delete(this.elements, ele)
		this.list.Remove(ele)
	}
}
