package Orm

import (
	"database/sql"
)

type KVDBInterface interface {
	NewStoreObj(key interface{}) StoreObj
    LoadDB(key interface{}) (StoreObj,error)
	InsertDB(obj StoreObj) bool
	UpdateDB(obj StoreObj) bool
	DeleteDB(key interface{}) bool
}

type KVTable struct {
	KVDBInterface
	KVCache
	TableName string
	DB *sql.DB
	insertChan  chan StoreObj
	updateChan  chan StoreObj
	deleteChan  chan interface{}
	stopChan chan int
}

func(this *KVTable) Init(dbInterface KVDBInterface, db *sql.DB, tableName string, max int) {
	this.KVCache.Init(max)
	this.KVDBInterface = dbInterface
	this.DB = db
	this.TableName = tableName
	this.insertChan = make(chan StoreObj, 512)
	this.updateChan = make(chan StoreObj, 512)
	this.deleteChan = make(chan interface{}, 128)
	go this.Run()
}

func (this *KVTable)Stop() {
	this.stopChan <- 1
}

func(this *KVTable) Get(key interface{}) StoreObj {
	obj, ok := this.KVCache.Get(key)
	if ok {
		return obj
	}
	obj, _= this.LoadDB(key)
	if obj == nil {
		obj = this.NewStoreObj(key)
		this.Insert(key, obj)
	}
	this.KVCache.Set(key, obj)
	return obj
}

func(this *KVTable) Insert(key interface{}, obj StoreObj) bool{
	this.KVCache.Set(key, obj)
	obj.AddDirty(Dirty_New)
	this.insertChan <- obj
	return true
}

func(this *KVTable) Update(key interface{}, obj StoreObj) {
	this.KVCache.Set(key, obj)
	obj.AddDirty(Dirty_Update)
	this.updateChan <- obj
}

func(this *KVTable) Delete(key interface{}) {
	this.KVCache.Delete(key)
	this.deleteChan <- key
}

func(this *KVTable) Run() {
	var stop bool = false
	for !stop {
		select {
		case obj := <-this.insertChan:
			if obj != nil {
				this.InsertDB(obj)
			}
		case obj := <-this.updateChan:
			if obj != nil {
				this.UpdateDB(obj)
			}
		case key := <-this.deleteChan:
			this.DeleteDB(key)
		case <-this.stopChan:
			stop = true
		}
	}
}

