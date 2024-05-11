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
	// group systems without overlapping types to parallelize
	parallel        bool
	parallelSystems [][]System
}

func NewSystemStorage(ecs *ECS, parallel bool) (this *SystemStorage) {
	this = new(SystemStorage)
	this.ecs = ecs
	this.parallel = parallel
	this.systemTypes = make(map[System][]reflect.Type)
	return
}

// NewParallelSystemStorage creates a parallel-systems tracking storage (costs, do not use if you don't need it)
func NewParallelSystemStorage(ecs *ECS) (this *SystemStorage) {
	return NewSystemStorage(ecs, true)
}

// Clear nils all systems
func (this *SystemStorage) Clear() {
	this.systems = nil
	this.systemTypes = nil
	this.parallelSystems = nil
}

// All returns all systems
func (this *SystemStorage) All() []System {
	return this.systems
}

// AllParallel returns all systems grouped by parallelity
func (this *SystemStorage) AllParallel() [][]System {
	return this.parallelSystems
}

// AddSystem stores the given system under every type to this storage
func (this *SystemStorage) AddSystem(system System, types ...any) []reflect.Type {
	// add to slice
	this.systems = append(this.systems, system)
	this.systemTypes[system] = make([]reflect.Type, len(types))
	for i, t := range types {
		// add to types
		this.systemTypes[system][i] = reflect.TypeOf(t) //this.ecs.getPlainType(t)
	}

	// Sort
	this.sort()
	this.parallelize()

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

	// delete types
	delete(this.systemTypes, system)

	// Sort
	this.sort()
	this.parallelize()
}

// sort systems by priority (higher = better)
func (this *SystemStorage) sort() []System {
	slices.SortStableFunc(this.systems, func(a, b System) int {
		return cmp.Compare(b.Priority(), a.Priority())
	})
	return this.systems
}

// parallelize compares all systems against each other to build a 2-dim slice of parallel runnable systems
func (this *SystemStorage) parallelize() [][]System {
	if !this.parallel {
		return this.parallelSystems
	}

	// Check whether systems can run in parallel anew
	this.parallelSystems = make([][]System, 0)
	systems := append([]System(nil), this.systems...)
	var pSystems []System
	for len(systems) > 0 {
		pSystems, systems = this.testSystemsOverlap(systems[0], systems)
		this.parallelSystems = append(this.parallelSystems, pSystems)
	}
	return this.parallelSystems
}

// testSystemsOverlap compares systemA vs all other systems to find commonalities
func (this *SystemStorage) testSystemsOverlap(a System, systems []System) (pSystems []System, otherSystems []System) {
	// The parallel systems to run with a
	pSystems = append(pSystems, a)

	// Check against all others
	for i := 1; i < len(systems); i++ {
		b := systems[i]

		// In case of no type equality, we can run in parallel (no overlap)
		if !this.testTypesOverlap(this.systemTypes[a], this.systemTypes[b]) {
			pSystems = append(pSystems, b)
		} else {
			otherSystems = append(otherSystems, b)
		}
	}
	return pSystems, otherSystems
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

// testTypesOverlap compares the type slices for any overlap (not a single same type)
func (this *SystemStorage) testTypesOverlap(a, b []reflect.Type) bool {
	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				return true
			}
		}
	}
	return false
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
