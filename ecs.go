package ecs

import (
	"hash/fnv"
	"log"
	"reflect"
	"sync/atomic"
	"time"
)

type ECS struct {
	// unique atomic counter per ECS
	entityCounter atomic.Uint64

	entities   map[uint64]*Entity
	systems    map[uint64]System
	components map[reflect.Type]map[uint64]any
	context    map[reflect.Type]any
}

func New() (this *ECS) {
	this = new(ECS)

	this.entities = make(map[uint64]*Entity)
	this.systems = make(map[uint64]System)
	this.components = make(map[reflect.Type]map[uint64]any)
	this.context = make(map[reflect.Type]any)

	return this
}

// AddContext attaches any service, map or other interfaces to this ECS (there can only be one per type)
func (this *ECS) AddContext(c any) *ECS {
	this.context[reflect.TypeOf(c).Elem()] = c
	return this
}

// GetContext returns a context from the ECS (there can only be one per type)
func (this *ECS) GetContext(t any) any {
	return this.context[reflect.TypeOf(t).Elem()]
}

// GetContext returns a typed context from the ECS (there can only be one per type)
func GetContext[T any](ecs *ECS) T {
	return ecs.context[reflect.TypeFor[T]()].(T)
}

// CreateEntity scaffolds a new entity with the given components
func (this *ECS) CreateEntity(components ...any) *Entity {
	entity := newEntity(&this.entityCounter)
	this.entities[entity.Id] = entity

	// Add components to systems
	compHash := this.getComponentHash(components)
	for sysHash, s := range this.systems {
		if sysHash == compHash {
			s.AddEntity(entity)
		}
	}

	// Store components globally by type
	for _, c := range components {
		cType := reflect.TypeOf(c).Elem()
		if _, ok := this.components[cType]; !ok {
			this.components[cType] = make(map[uint64]any)
		}
		this.components[cType][entity.Id] = c
	}

	return entity
}

func (this *ECS) AddEntity(e any) *Entity {
	entity := newEntity(&this.entityCounter)

	v := reflect.ValueOf(e)
	t := reflect.TypeOf(e)
	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i)

		switch val.Kind() {
		case reflect.Struct:
			log.Println("struct", v, v.Kind(), reflect.TypeOf(v.Interface()))
		default:
		}
	}

	return entity
}

// GetComponents by given type
func (this *ECS) GetComponents(t any) map[uint64]any {
	return this.components[reflect.TypeOf(t).Elem()]
}

// GetComponents by given type
func GetComponents[T any](ecs *ECS) map[uint64]T {
	components := ecs.components[reflect.TypeFor[T]()]
	return reflect.ValueOf(components).Interface().(map[uint64]T)
}

func (this *ECS) AddSystem(s System, components ...any) *ECS {
	compHash := this.getComponentHash(components)
	this.systems[compHash] = s
	return this
}

func (this *ECS) getComponentHash(components ...any) uint64 {
	h := fnv.New64a()
	for _, component := range components {
		_, _ = h.Write([]byte(reflect.TypeOf(component).Elem().String()))
	}
	return h.Sum64()
}

// Update calls all systems to run and do their stuff
func (this *ECS) Update(dt time.Duration) *ECS {
	for _, system := range this.systems {
		system.Run(this, dt)
	}
	return this
}

//type ECS struct {
//	entities   map[uint64]*Entity
//	components map[reflect.Type]map[uint64]*Component
//	systems    []System
//}
//
//func New() (this *ECS) {
//	this = new(ECS)
//
//	// Mem prep
//	this.entities = make(map[uint64]*Entity)
//	this.components = make(map[reflect.Type]map[uint64]*Component)
//	this.systems = make([]System, 0)
//
//	return this
//}
//
//func (this *ECS) CreateEntity() *Entity {
//	return NewEntity()
//}
//
//func (this *ECS) AddEntity(e *Entity) *ECS {
//	this.entities[e.Id] = e
//	this.addEntityComponents(e)
//	this.addEntityToSystems(e)
//	return this
//}
//
//func (this *ECS) addEntityComponents(e *Entity) *ECS {
//	for _, c := range e.components {
//		// Check if a component of this type is already known
//		cType := reflect.TypeOf(c)
//		if _, ok := this.components[cType]; !ok {
//			// If not, prepare
//			this.components[cType] = make(map[uint64]*Component)
//		}
//
//		// Add component to entity id
//		this.components[cType][e.Id] = &c
//	}
//	return this
//}
//
//func (this *ECS) addEntityToSystems(e *Entity) *ECS {
//	for _, s := range this.systems {
//		for _, sysC := range s.Components() {
//			sysCType := reflect.TypeOf(sysC)
//
//			// Check if this system matches the entity components
//			for _, entityC := range e.components {
//				entityCType := reflect.TypeOf(entityC)
//				log.Println("comparing system components", sysCType, entityCType)
//				if entityCType == sysCType {
//					log.Println("Add component to system:", sysC, entityC)
//					s.AddComponent(&entityC)
//				}
//			}
//		}
//	}
//	return this
//}
//
//func (this *ECS) addComponentsToSystem(s System) *ECS {
//	//for _, sysC := range s.Components() {
//	//	sysCType := reflect.TypeOf(sysC)
//	//
//	//}
//	return this
//}
//
//func (this *ECS) AddSystem(s System) *ECS {
//	this.systems = append(this.systems, s)
//	this.addComponentsToSystem(s)
//	return this
//}
//
//func (this *ECS) Update(dt time.Duration) *ECS {
//	for _, system := range this.systems {
//		system.Run(this, dt)
//	}
//	return this
//}
