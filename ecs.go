package ecs

import (
	"reflect"
	"sync/atomic"
	"time"
)

type ECS struct {
	// unique atomic counter per ECS
	entityCounter atomic.Uint64

	entities   map[uint64]Entity
	systems    *SystemStorage
	components map[reflect.Type]map[uint64]any
	context    map[reflect.Type]any
}

func New() (this *ECS) {
	this = new(ECS)

	this.entities = make(map[uint64]Entity)
	this.systems = NewSystemStorage(this)
	this.components = make(map[reflect.Type]map[uint64]any)
	this.context = make(map[reflect.Type]any)

	return this
}

// AddContext attaches any service, map or other interfaces to this ECS (there can only be one per type)
func (this *ECS) AddContext(c any) *ECS {
	this.context[this.getPlainType(c)] = c
	return this
}

// GetContext returns a context from the ECS (there can only be one per type)
func (this *ECS) GetContext(c any) any {
	return this.context[this.getPlainType(c)]
}

// CreateEntity scaffolds a new entity with the given components
func (this *ECS) CreateEntity(components ...any) Entity {
	entity := NewEntity(&this.entityCounter)

	// Store entities
	this.entities[entity.Id()] = entity

	// Add components to systems
	systems := this.systems.QuerySystems(components...)
	for _, system := range systems {
		system.addEntity(entity)
	}

	// Store components globally by type
	for _, c := range components {
		cType := this.getPlainType(c)
		if _, ok := this.components[cType]; !ok {
			this.components[cType] = make(map[uint64]any)
		}
		this.components[cType][entity.Id()] = c
	}

	return entity
}

func (this *ECS) GetEntity(id uint64) Entity {
	return this.entities[id]
}

//func (this *ECS) AddEntity(e any) Entity {
//	entity := newEntity(&this.entityCounter)
//
//	v := reflect.ValueOf(e)
//	t := reflect.TypeOf(e)
//	for i := 0; i < t.NumField(); i++ {
//		val := v.Field(i)
//
//		switch val.Kind() {
//		case reflect.Struct:
//			log.Println("struct", v, v.Kind(), reflect.TypeOf(v.Interface()))
//		default:
//		}
//	}
//
//	return entity
//}

// GetComponents by given type
func (this *ECS) GetComponents(componentType any) map[uint64]interface{} {
	return this.components[this.getPlainType(componentType)]
}

func (this *ECS) AddSystem(s System, types ...any) *ECS {
	this.systems.AddSystem(s, types...)
	return this
}

// Update calls all systems to run and do their stuff
func (this *ECS) Update(dt time.Duration) *ECS {
	systems := this.systems.All()
	for _, s := range systems {
		s.Run(this, dt)
	}
	return this
}

// getPlainType returns a non-pointer type from any given
func (this *ECS) getPlainType(t any) reflect.Type {
	typ := reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		return typ.Elem()
	}
	return typ
}
