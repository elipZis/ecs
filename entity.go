package ecs

import "sync/atomic"

type Entity struct {
	Id         uint64
	components []any
}

func newEntity(counter *atomic.Uint64) (this *Entity) {
	this = new(Entity)
	this.Id = counter.Add(1)
	return this
}

func (this *Entity) addComponent(component any) {
	this.components = append(this.components, component)
}

//var entityCounter atomic.Uint64
//
//type Entity struct {
//	Id         uint64
//	components []Component
//}
//
//func NewEntity() (this *Entity) {
//	this = new(Entity)
//	this.Id = entityCounter.Add(1)
//	this.components = make([]Component, 0)
//	return this
//}
//
//func (this *Entity) AddComponent(c Component) *Entity {
//	this.components = append(this.components, c)
//	return this
//}
