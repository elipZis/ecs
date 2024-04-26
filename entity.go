package ecs

import "sync/atomic"

type Entity interface {
	Id() uint64
}

type BaseEntity struct {
	id         uint64
	components []any
}

func NewEntity(counter *atomic.Uint64) (this *BaseEntity) {
	this = new(BaseEntity)
	this.id = counter.Add(1)
	return this
}

func (this *BaseEntity) addComponent(component any) {
	this.components = append(this.components, component)
}

func (this *BaseEntity) Id() uint64 {
	return this.id
}
