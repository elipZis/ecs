package ecs

import (
	"cmp"
	"reflect"
	"slices"
)

type SystemStorage struct {
	ecs *ECS

	// n systems could be registered
	systems []System
	// per system, n types it requires
	systemTypes map[System][]reflect.Type
}

func NewSystemStorage(ecs *ECS) (this *SystemStorage) {
	this = new(SystemStorage)
	this.ecs = ecs
	this.systemTypes = make(map[System][]reflect.Type)
	return
}

func (this *SystemStorage) Clear() {
	this.systems = nil
	this.systemTypes = nil
}

// All returns all systems
func (this *SystemStorage) All() []System {
	return this.systems
}

// AddSystem stores the given system under every type to this storage
func (this *SystemStorage) AddSystem(system System, types ...any) []reflect.Type {
	// add to slice
	this.systems = append(this.systems, system)
	this.sort()

	this.systemTypes[system] = make([]reflect.Type, len(types))
	for i, t := range types {
		// add to types
		this.systemTypes[system][i] = reflect.TypeOf(t) //this.ecs.getPlainType(t)
	}
	return this.systemTypes[system]
}

// RemoveSystem slices the given system out of every type from this storage
func (this *SystemStorage) RemoveSystem(system System) {
	// delete from slice
	for i, s := range this.systems {
		if s == system {
			this.systems = append(this.systems[:i], this.systems[i+1:]...)
		}
	}
	this.sort()

	// delete types
	delete(this.systemTypes, system)
}

// sort systems by priority (higher = better)
func (this *SystemStorage) sort() []System {
	slices.SortStableFunc(this.systems, func(a, b System) int {
		return cmp.Compare(b.Priority(), a.Priority())
	})
	return this.systems
}

// QuerySystems returns all systems matching all given types connotations
func (this *SystemStorage) QuerySystems(types ...any) []System {
	systems := make([]System, 0)

	// Check the types of the given
	reflectTypes := make([]reflect.Type, len(types))
	for i, t := range types {
		reflectTypes[i] = reflect.TypeOf(t) //this.ecs.getPlainType(t))
	}

	for system, systemTypes := range this.systemTypes {
		if this.testTypesSubset(systemTypes, reflectTypes) {
			systems = append(systems, system)
		}
	}

	return systems
}

// testTypesSubset checks if the needle is fully contained in the haystack
func (this *SystemStorage) testTypesSubset(needle, haystack []reflect.Type) bool {
	set := make(map[reflect.Type]int, len(haystack))
	for _, value := range haystack {
		set[value] += 1
	}

	for _, value := range needle {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}

// testTypesEq compares the type slices for equal contents and returns false if any difference occurs
func (this *SystemStorage) testTypesEq(a, b []reflect.Type) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// intersectSystems returns the intersection of the given systems
func (this *SystemStorage) intersectSystems(a, b []System) []System {
	intersection := make([]System, 0)
	for _, s1 := range a {
		for _, s2 := range b {
			if s1 == s2 {
				intersection = append(intersection, s2)
			}
		}
	}
	return intersection
}
