package ecs

import (
	"time"
)

type System interface {
	Run(ecs *ECS, dt time.Duration)
	AttachEntity(e Entity)
	DetachEntity(e Entity)
}

type EntitySystem struct {
	entities []uint64
}

func (this *EntitySystem) Entities() []uint64 {
	return this.entities
}

func (this *EntitySystem) AttachEntity(e Entity) {
	this.attachEntityById(e.Id())
}

func (this *EntitySystem) attachEntityById(eId uint64) {
	this.entities = append(this.entities, eId)
}

func (this *EntitySystem) DetachEntity(e Entity) {
	this.detachEntityById(e.Id())
}

func (this *EntitySystem) detachEntityById(eId uint64) {
	for i, e := range this.entities {
		if e == eId {
			this.entities = append(this.entities[:i], this.entities[i+1:]...)
			break
		}
	}
}
