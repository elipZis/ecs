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
	toRemove   []uint64
	systems    *SystemStorage
	components *ComponentStorage
	context    map[reflect.Type]any
}

// New creates a new ecs world
func New() (this *ECS) {
	this = new(ECS)

	this.entities = make(map[uint64]Entity)
	this.systems = NewSystemStorage(this)
	this.components = NewComponentStorage(this)
	this.context = make(map[reflect.Type]any)

	return this
}

// Clear nils all entities from this world
func (this *ECS) Clear() {
	this.entities = nil
	this.context = nil
	if this.systems != nil {
		this.systems.Clear()
	}
	if this.components != nil {
		this.components.Clear()
	}
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
	// Add components to entity as reference
	entity.AddComponents(components...)

	// Add components to systems
	systems := this.systems.QuerySystems(components...)
	for _, system := range systems {
		system.AttachEntity(entity)
	}

	// Store components globally by type
	this.components.AddComponent(entity, components...)

	return entity
}

// AddEntity via reflection of embedded structs
func (this *ECS) AddEntity(e any) Entity {
	// Get structs as components
	var components []any
	v := reflect.ValueOf(e)
	t := reflect.TypeOf(e)
	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i)

		switch val.Kind() {
		case reflect.Struct:
			if val.CanInterface() {
				// Add all structs as components of this entity as reference
				components = append(components, val.Interface())
			}

		case reflect.Pointer:
			if val.CanInterface() && val.Elem().Kind() == reflect.Struct {
				// Add all structs as components of this entity as reference
				components = append(components, val.Interface())
			}

		default:
		}
	}

	return this.CreateEntity(components...)
}

// RemoveEntity marks an entity for deletion in the next iteration, to not affect the current run
func (this *ECS) RemoveEntity(id uint64) {
	this.toRemove = append(this.toRemove, id)
}

// removeEntities detaches the entity & components from all systems and globally
func (this *ECS) removeEntities() {
	for _, entityId := range this.toRemove {
		this.RemoveEntityNow(entityId)
	}
	this.toRemove = make([]uint64, 0)
}

// RemoveEntityNow detaches the entity now, no matter of more systems are running
func (this *ECS) RemoveEntityNow(id uint64) {
	entity := this.entities[id]

	if entity != nil {
		// Remove from systems
		systems := this.systems.QuerySystems(entity.GetComponents()...)
		for _, system := range systems {
			system.DetachEntity(entity)
		}

		// Remove from global components
		this.components.RemoveComponent(entity, entity.GetComponents()...)

		// delete entities
		delete(this.entities, id)
	}
}

// GetEntity returns the id referenced entity of this ECS
func (this *ECS) GetEntity(id uint64) Entity {
	return this.entities[id]
}

// GetComponents by given type
func (this *ECS) GetComponents(componentType any) map[uint64]interface{} {
	return this.components.GetComponents(componentType)
}

// AddSystem attaches the given system to this ECS under the given types
func (this *ECS) AddSystem(s System, types ...any) *ECS {
	this.systems.AddSystem(s, types...)
	// TODO: Check whether existing entities should be added to this new system
	return this
}

// RemoveSystem deletes the given system from this ECS
func (this *ECS) RemoveSystem(s System) *ECS {
	this.systems.RemoveSystem(s)
	return this
}

// Update calls all systems to run and do their stuff
func (this *ECS) Update(dt time.Duration) *ECS {
	// Clear all marked entities
	this.removeEntities()

	// Iterate on the systems
	systems := this.systems.All()
	for _, s := range systems {
		s.Run(this, dt)
	}
	return this
}

// getPlainType returns a non-pointer type from any given
func (this *ECS) getPlainType(t any) reflect.Type {
	var typ reflect.Type
	if _, ok := t.(reflect.Type); !ok {
		typ = reflect.TypeOf(t)
	} else {
		typ = t.(reflect.Type)
	}
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}
