package ecs

import (
	"reflect"
)

type ComponentStorage struct {
	ecs *ECS

	// n components per type could be registered
	components map[reflect.Type]map[uint64]any
}

func NewComponentStorage(ecs *ECS) (this *ComponentStorage) {
	this = new(ComponentStorage)
	this.ecs = ecs
	this.components = make(map[reflect.Type]map[uint64]any)
	return
}

// AddComponent stores the given components
func (this *ComponentStorage) AddComponent(e Entity, components ...any) {
	for _, c := range components {
		cType := this.ecs.getPlainType(c)
		if _, ok := this.components[cType]; !ok {
			this.components[cType] = make(map[uint64]any)
		}
		this.components[cType][e.Id()] = c
	}
}

// RemoveComponent deletes the given components from their respective types and entity
func (this *ComponentStorage) RemoveComponent(e Entity, components ...any) {
	for _, c := range components {
		cType := this.ecs.getPlainType(c)
		delete(this.components[cType], e.Id())
	}
}

// GetComponents by given type
func (this *ComponentStorage) GetComponents(componentType any) map[uint64]interface{} {
	return this.components[this.ecs.getPlainType(componentType)]
}

// GetEntityComponent is a typed helper to get a cast entity component from the ECS
func GetEntityComponent[T any](ecs *ECS, eId uint64) T {
	vals := ecs.GetComponents(reflect.TypeFor[T]())
	val := vals[eId]
	return val.(T)
}

// GetComponent is a casting helper to return a typed component by entity id
func GetComponent[T any](components map[uint64]any, eId uint64) T {
	return components[eId].(T)
}

// GetComponentsFor creates a typed map of the components
func GetComponentsFor[T any](ecs *ECS) (retVals map[uint64]T) {
	comps := ecs.GetComponents(reflect.TypeFor[T]())
	for i, comp := range comps {
		retVals[i] = comp.(T)
	}
	return retVals
}
