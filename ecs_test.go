package ecs

import (
	"log"
	"reflect"
	"testing"
	"time"
)

type PositionComponent struct {
	X int
	Y int
}

type VelocityComponent struct {
	DX int
	DY int
}

type BoundsComponent struct {
	Width  int
	Height int
}

type CommComponent struct {
}

type MoveSystem struct {
}

func (this *MoveSystem) Run(ecs *ECS, dt time.Duration) {

}

type CollisionSystem struct {
	pos    HasSystem[PositionComponent]
	bounds HasSystem[BoundsComponent]
}

func (this *CollisionSystem) Run(ecs *ECS, dt time.Duration) {
	positions := GetComponents[PositionComponent](ecs)
	bounds := GetComponents[BoundsComponent](ecs)
	log.Println(positions)
	log.Println(bounds)

	sPositions := this.pos.GetComponents()
	log.Println(sPositions)

	//this.players.Range(func(key *uuid.UUID, player *Player) bool {
	//	// Reset collision
	//	player.GetBounds().Colliding = false
	//
	//	for _, mapObj := range this.Map.Objects {
	//		// Only for non-player objects as players do not collide with each other
	//		if _, ok := mapObj.(*Player); !ok {
	//			// TODO: if colliding
	//			if this.detectIntersection(player.GetPosition(), mapObj.GetPosition(), player.GetBounds(), mapObj.GetBounds()) {
	//				// TODO: something
	//				player.GetBounds().Colliding = true
	//			}
	//		}
	//	}
	//
	//	return true
	//})
}

type Player struct {
	PositionComponent
	VelocityComponent
	BoundsComponent
	CommComponent

	something string
}

func TestECS(t *testing.T) {
	player := Player{
		PositionComponent: PositionComponent{},
		VelocityComponent: VelocityComponent{},
		BoundsComponent:   BoundsComponent{},
		CommComponent:     CommComponent{},
	}

	//fields := make([]reflect.Value, 0)
	ifv := reflect.ValueOf(player)
	ift := reflect.TypeOf(player)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)

		switch v.Kind() {
		case reflect.Struct:
			log.Println("struct", v, v.Kind(), reflect.TypeOf(v.Interface()))
		default:
			log.Println("not", v, v.Kind(), reflect.TypeOf(v))
		}
	}

	// Create a new world
	//ecs := New()
	//ecs.AddSystem()

	//// Tick for 30 fps to simulate a game update
	//done := make(chan bool, 1)
	//ticker := time.NewTicker(33 * time.Millisecond)
	//go func() {
	//	var dt time.Duration
	//	lastTimestamp := time.Now()
	//
	//	for {
	//		select {
	//		case <-done:
	//			return
	//		case tick := <-ticker.C:
	//			// Save the diff, time and tick
	//			dt = tick.Sub(lastTimestamp)
	//			lastTimestamp = tick
	//
	//			// ECS update
	//			ecs.Update(dt)
	//		}
	//	}
	//}()
	//
	//time.Sleep(1000 * time.Millisecond)
	//<-done
}

//type Position struct {
//	HasComponent
//	X int
//	Y int
//}
//
//type Velocity struct {
//	HasComponent
//	DX int
//	DY int
//}
//
//type Moveable struct {
//	HasSystem
//}
//
//func (this *Moveable) Components() []Component {
//	return []Component{
//		&Position{},
//		&Velocity{},
//	}
//}
//
//func (this *Moveable) Run(ecs *ECS, dt time.Duration) {
//	log.Println(this.components[reflect.TypeOf(&Velocity{})])
//
//	velocities := this.components[reflect.TypeOf(&Velocity{})]
//	for _, velocity := range velocities {
//		log.Println("velocity", *velocity)
//	}
//
//	for t, c := range this.components {
//		log.Println(t, c)
//		//switch t {
//		//case reflect.TypeOf(&Velocity{}):
//		//	log.Println("Velocity", c)
//		//	for _, component := range c {
//		//		m := reflect.ValueOf(component).Elem()
//		//		log.Println(m)
//		//		var v *Velocity = (*component).(*Velocity)
//		//		log.Println(v, reflect.ValueOf(v).Kind())
//		//	}
//		//case reflect.TypeOf(&Position{}):
//		//	log.Println("Position")
//		//}
//	}
//}
//
//type I interface {
//	M()
//}
//
//type V struct {
//	X int
//	Y int
//}
//
//func (this *V) M() {}
//
//func GetVs[T any](t I) T {
//	return t.(T)
//}
//
//func TestECS(t *testing.T) {
//	var is map[reflect.Type][]I
//	is = make(map[reflect.Type][]I)
//
//	v1 := new(V)
//	//v2 := new(V)
//	log.Println(reflect.TypeFor[V]())
//
//	is[reflect.TypeOf(v1).Elem()] = make([]I, 0)
//	is[reflect.TypeOf(v1).Elem()] = append(is[reflect.TypeOf(v1).Elem()], v1)
//
//	log.Println(is)
//	log.Println(is[reflect.TypeFor[V]()])
//
//	vs := is[reflect.TypeFor[V]()]
//	for _, v := range vs {
//		log.Println(v)
//	}
//
//	// Create a new world
//	ecs := New()
//
//	// Add at least one system to act
//	ecs.AddSystem(&Moveable{})
//
//	// Create a new entity with components
//	entity := ecs.CreateEntity().AddComponent(&Position{
//		X: 2,
//		Y: 2,
//	}).AddComponent(&Velocity{
//		DX: 1,
//		DY: 1,
//	})
//	t.Logf("created entity %d", entity.Id)
//	ecs.AddEntity(entity)
//
//	entity = ecs.CreateEntity().AddComponent(&Position{
//		X: 5,
//		Y: 5,
//	}).AddComponent(&Velocity{
//		DX: 2,
//		DY: 2,
//	})
//	t.Logf("created entity %d", entity.Id)
//	ecs.AddEntity(entity)
//
//	// Update the world
//	for i := range 30 {
//		// Simulate one second at 30 fps
//		ecs.Update(time.Duration(float32(i)*33.33) * time.Millisecond)
//	}
//
//	// Validate that the entity components changed based on the system
//
//}
