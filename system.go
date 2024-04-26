package ecs

import (
	"time"
)

type System interface {
	Run(ecs *ECS, dt time.Duration)
	addEntity(e Entity)
	addEntityId(eId uint64)
}

type BaseSystem struct {
	entities []uint64
}

func (this *BaseSystem) addEntity(e Entity) {
	this.addEntityId(e.Id())
}

func (this *BaseSystem) addEntityId(eId uint64) {
	this.entities = append(this.entities, eId)
}

//type System interface {
//	Run(ecs *ECS, dt time.Duration)
//	addEntity(e *Entity)
//	GetEntities() []*Entity
//}
//
//type HasSystem[T any] struct {
//	components map[reflect.Type]map[uint64]any
//}
//
//func (this *HasSystem[T]) addEntity(e *Entity) {
//	for _, component := range e.components {
//		this.addComponent(e, component)
//	}
//}
//
//func (this *HasSystem[T]) addComponent(e *Entity, c any) {
//	cType := reflect.TypeOf(c).Elem()
//	if _, ok := this.components[cType]; !ok {
//		this.components[cType] = make(map[uint64]any)
//	}
//	this.components[cType][e.Id] = c
//}
//
//func (this *HasSystem[T]) getComponents() map[uint64]T {
//	components := this.components[reflect.TypeFor[T]()]
//	return reflect.ValueOf(components).Interface().(map[uint64]T)
//}

//type System interface {
//	System() *HasSystem
//	Components() []Component
//	AddComponent(c *Component) *HasSystem
//	Run(ecs *ECS, dt time.Duration)
//}
//
//type HasSystem struct {
//	components map[reflect.Type][]*Component
//}
//
//func (this *HasSystem) System() *HasSystem { return this }
//
//func (this *HasSystem) AddComponent(c *Component) *HasSystem {
//	if this.components == nil {
//		this.components = make(map[reflect.Type][]*Component)
//	}
//	cType := reflect.TypeOf(*c)
//	this.components[cType] = append(this.components[cType], c)
//	return this
//}
//
//func (this *HasSystem) GetComponent(c Component) *HasSystem {
//	velocities := this.components[reflect.TypeOf(&c)]
//	for _, velocity := range velocities {
//		log.Println("velocity", *velocity)
//	}
//	return this
//}
