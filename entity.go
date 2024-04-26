package ecs

import "sync/atomic"

type Entity interface {
	Id() uint64
	GetComponents() []any
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

func (this *BaseEntity) AddComponent(component any) {
	this.components = append(this.components, component)
}

func (this *BaseEntity) AddComponents(component ...any) {
	for _, c := range component {
		this.AddComponent(c)
	}
}

func (this *BaseEntity) GetComponents() []any {
	return this.components
}

func (this *BaseEntity) Id() uint64 {
	return this.id
}
