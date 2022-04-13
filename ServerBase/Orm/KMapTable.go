package Orm

import (
	"database/sql"
)

type KMapDBInterface interface {
	LoadDB(key interface{}) map[interface{}]StoreObj
	InsertDB(obj StoreObj) bool
	UpdateDB(obj StoreObj) bool
	DeleteDB(key interface{}, field interface{}) bool
	DeleteDB_Collection(key interface{}) bool
}

type KeyField struct {
	key interface{}
	field interface{}
}

type KMapTable struct {
	KMapDBInterface
	KMapCache
	DB *sql.DB
	TableName string
	insertChan  chan StoreObj
	updateChan  chan StoreObj
	deleteChan  chan KeyField
	deleteAllChan  chan interface{}
	stopChan chan int
}

func(this *KMapTable) Init(dbInterface KMapDBInterface, db *sql.DB, tableName string, max int) {
	this.KMapCache.Init(max)
	this.KMapDBInterface = dbInterface
	this.DB = db
	this.TableName = tableName
	this.insertChan = make(chan StoreObj, 1024)
	this.updateChan = make(chan StoreObj, 1024)
	this.deleteChan = make(chan KeyField, 128)
	this.deleteAllChan = make(chan interface{}, 128)
	go this.Run()
}

func (this *KMapTable)Stop() {
	this.stopChan <- 1
}

func(this *KMapTable) GetCollection(key uint64) (map[interface{}]StoreObj) {
	objs, ok := this.KMapCache.GetCollection(key)
	if ok {
		return objs
	}
    objs = this.LoadDB(key)
    return objs
}

func(this *KMapTable) DeleteCollection(key interface{}) {
	this.KMapCache.DeleteCollection(key)
	this.deleteAllChan <- key
}

func(this *KMapTable) Insert(key interface{}, field interface{}, obj StoreObj) bool{
	this.KMapCache.Set(key, field, obj)
	obj.AddDirty(Dirty_New)
	this.insertChan <- obj
	return true
}

func(this *KMapTable) Update(key interface{},  field interface{}, obj StoreObj) {
	this.KMapCache.Set(key, field, obj)
	obj.AddDirty(Dirty_Update)
	this.updateChan <- obj
}

func(this *KMapTable) Delete(key interface{}, field interface{}) {
	this.KMapCache.Delete(key, field)
	this.deleteChan <- KeyField{key: key, field: field}
}

func(this *KMapTable) Run() {
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
		case keyfield := <-this.deleteChan:
			this.DeleteDB(keyfield.key, keyfield.field)
		case key := <-this.deleteAllChan:
			this.DeleteDB_Collection(key)
		case <-this.stopChan:
			stop = true
		}
	}
}