package Orm

const(
	Dirty_New = 1
	Dirty_Update = 2
	Dirty_Delete = 4
)

type StoreObj interface {
	GetDirty() int
	AddDirty(dirty int)
}

type  StoreObjBase struct {
	dirty int
}

func(this *StoreObjBase) GetDirty() int {
	return this.dirty
}

func(this *StoreObjBase) AddDirty(dirty int){
	this.dirty += dirty
}


