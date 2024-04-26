package ecs

import (
	"time"
)

type System interface {
	Run(ecs *ECS, dt time.Duration)
	AttachEntity(e Entity)
	DetachEntity(e Entity)
}

type BaseSystem struct {
	entities []uint64
}

func (this *BaseSystem) AttachEntity(e Entity) {
	this.attachEntityById(e.Id())
}

func (this *BaseSystem) attachEntityById(eId uint64) {
	this.entities = append(this.entities, eId)
}

func (this *BaseSystem) DetachEntity(e Entity) {
	this.detachEntityById(e.Id())
}

func (this *BaseSystem) detachEntityById(eId uint64) {
	for i, e := range this.entities {
		if e == eId {
			this.entities = append(this.entities[:i], this.entities[i+1:]...)
			break
		}
	}
}
