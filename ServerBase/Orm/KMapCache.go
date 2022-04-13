package Orm

import (
	"container/list"
	"sync"
)

type KMapNode struct {
	key interface{}
	objs map[interface{}]StoreObj
}

type KMapCache struct {
	max   int
	list     *list.List
	elements map[interface{}]*list.Element
	mutex sync.Mutex
}

func(this *KMapCache) Init(max int) {
	this.max = max
	if this.max == 0 {
		this.max = 10000
	}
	this.list = list.New()
	this.elements = make(map[interface{}]*list.Element)
}

func(this *KMapCache) GetCollection(key interface{}) (map[interface{}]StoreObj,bool) {
	if this.elements == nil {
		return nil,false
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ele, ok := this.elements[key]
	if  ok {
		this.list.MoveToFront(ele)
		return ele.Value.(*KMapNode).objs, true
	}
	return nil, false
}

func(this *KMapCache) Get(key interface{}, field interface{}) (StoreObj,bool) {
	if this.elements == nil {
		return nil,false
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ele, ok := this.elements[key]
	if  ok {
		this.list.MoveToFront(ele)
		obj, ok := ele.Value.(*KMapNode).objs[field]
		if ok {
			return obj, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

func(this *KMapCache) Set(key interface{}, field interface{}, obj StoreObj) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	e, ok := this.elements[key];
	if ok {
		e.Value.(*KMapNode).objs[field] = obj
		this.list.MoveToFront(e)
		return
	}
	objs := make(map[interface{}]StoreObj)
	objs[field] = obj
	ele := this.list.PushFront(&KMapNode{key: key, objs: objs})
	this.elements[key] = ele
	if this.max != 0 && this.list.Len() > this.max {
		oldest := this.list.Back()
		if oldest != nil {
			this.list.Remove(oldest)
			node := e.Value.(*KMapNode)
			delete(this.elements, node.key)
		}
	}
}

func(this *KMapCache) DeleteCollection(key interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ele, ok := this.elements[key]
	if ok {
		delete(this.elements, ele)
		this.list.Remove(ele)
	}
}

func(this *KMapCache) Delete(key interface{}, field interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ele, ok := this.elements[key]
	if ok {
		delete(ele.Value.(*KMapNode).objs, field)
		if len(ele.Value.(*KMapNode).objs) == 0 {
			this.list.Remove(ele)
		}
	}
}

