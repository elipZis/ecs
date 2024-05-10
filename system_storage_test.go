package ecs

import (
	"reflect"
	"testing"
)

func Test_IntersectSystems(t *testing.T) {
	ecs := New()
	storage := NewSystemStorage(ecs)

	s1 := MoveSystem{}
	s2 := CollisionSystem{}
	a := []System{
		&s1,
		&s2,
	}
	b := []System{
		&s1,
		&s2,
	}
	intersect := storage.intersectSystems(a, b)

	// Assertions
	if len(intersect) != 2 {
		t.Errorf("intersect %d; expected %d", len(intersect), 2)
	}

	a = []System{
		&s1,
		&s2,
	}
	b = []System{
		&s1,
		&MoveWithoutPtrSystem{},
	}
	intersect = storage.intersectSystems(a, b)

	// Assertions
	if len(intersect) != 1 {
		t.Errorf("intersect %d; expected %d", len(intersect), 1)
	}
}

func Test_NoIntersectSystems(t *testing.T) {
	ecs := New()
	storage := NewSystemStorage(ecs)

	a := []System{
		&MoveSystem{},
		&CollisionSystem{},
	}
	b := []System{
		&MoveWithoutPtrSystem{},
	}
	intersect := storage.intersectSystems(a, b)

	// Assertions
	if len(intersect) != 0 {
		t.Errorf("intersect %d; expected %d", len(intersect), 0)
	}
}

func Test_TypesOverlap(t *testing.T) {
	ecs := New()
	storage := NewSystemStorage(ecs)

	// Completely equal
	a := []reflect.Type{
		reflect.TypeOf(BoundsComponent{}),
		reflect.TypeOf(PositionComponent{}),
	}
	b := []reflect.Type{
		reflect.TypeOf(PositionComponent{}),
		reflect.TypeOf(BoundsComponent{}),
	}
	overlap := storage.testTypesOverlap(a, b)

	// Assertions
	if !overlap {
		t.Errorf("overlap %v; expected %v", overlap, true)
	}

	// One common component
	a = []reflect.Type{
		reflect.TypeOf(BoundsComponent{}),
		reflect.TypeOf(PositionComponent{}),
	}
	b = []reflect.Type{
		reflect.TypeOf(PositionComponent{}),
		reflect.TypeOf(CommComponent{}),
	}
	overlap = storage.testTypesOverlap(a, b)

	// Assertions
	if !overlap {
		t.Errorf("overlap %v; expected %v", overlap, true)
	}

	// No common component
	a = []reflect.Type{
		reflect.TypeOf(BoundsComponent{}),
		reflect.TypeOf(PositionComponent{}),
	}
	b = []reflect.Type{
		reflect.TypeOf(CommComponent{}),
	}
	overlap = storage.testTypesOverlap(a, b)

	// Assertions
	if overlap {
		t.Errorf("overlap %v; expected %v", overlap, false)
	}
}
