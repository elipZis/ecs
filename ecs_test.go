package ecs

import (
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
	BaseSystem
}

func (this *MoveSystem) Run(ecs *ECS, dt time.Duration) {
	pos := ecs.GetComponents(PositionComponent{})
	vel := ecs.GetComponents(VelocityComponent{})

	for _, entityId := range this.entities {
		p := pos[entityId].(*PositionComponent)
		v := vel[entityId].(*VelocityComponent)
		p.X += v.DX
		p.Y += v.DY
	}
}

type CollisionSystem struct {
	BaseSystem
}

func (this *CollisionSystem) Run(ecs *ECS, dt time.Duration) {
	positions := ecs.GetComponents(PositionComponent{})
	bounds := ecs.GetComponents(BoundsComponent{})

	for _, entityId := range this.entities {
		_ = positions[entityId].(*PositionComponent)
		_ = bounds[entityId].(*BoundsComponent)
	}

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

	s string
}

// createPlayer is a helper to create testable player objects
func createPlayer(s string) Player {
	return Player{
		PositionComponent: PositionComponent{X: 1, Y: 1},
		VelocityComponent: VelocityComponent{DX: 2, DY: 2},
		BoundsComponent:   BoundsComponent{Width: 10, Height: 10},
		CommComponent:     CommComponent{},
		s:                 s,
	}
}

func TestECS(t *testing.T) {
	// Create a new world
	ecs := New()

	// Add some interacting systems, working with same and different components
	ecs.AddSystem(&MoveSystem{}, PositionComponent{}, VelocityComponent{})
	ecs.AddSystem(&CollisionSystem{}, PositionComponent{}, BoundsComponent{})

	// Add some entities, which are not captured by any system
	ecs.CreateEntity(&PositionComponent{X: 1, Y: 1})          // 1
	ecs.CreateEntity(&VelocityComponent{DX: 2, DY: 2})        // 2
	ecs.CreateEntity(&BoundsComponent{Width: 10, Height: 10}) // 3

	// Add two fully fledged players, with more than needed components
	player4 := createPlayer("player4")
	ecs.CreateEntity(&player4.PositionComponent, &player4.VelocityComponent, &player4.BoundsComponent) // 4
	player5 := createPlayer("player5")
	ecs.CreateEntity(&player5.PositionComponent, &player5.VelocityComponent, &player5.BoundsComponent) // 5

	// Add reduced players, only fitting one system each
	player6 := createPlayer("player6")
	ecs.CreateEntity(&player6.PositionComponent, &player6.VelocityComponent) // 6
	player7 := createPlayer("player7")
	ecs.CreateEntity(&player7.PositionComponent, &player7.VelocityComponent) // 7
	player8 := createPlayer("player8")
	ecs.CreateEntity(&player8.PositionComponent, &player8.BoundsComponent) // 8
	player9 := createPlayer("player9")
	ecs.CreateEntity(&player9.PositionComponent, &player9.BoundsComponent) // 9

	// Tick for 30 fps to simulate a game update
	var fps = 30
	for i := 0; i < fps; i++ {
		ecs.Update(33 * time.Millisecond)
	}

	// Assertions
	if player4.X != (1+player4.DX*fps) || player4.Y != (1+player4.DY*fps) {
		t.Errorf("player4(%d, %d); expected %d", player4.X, player4.Y, 1+player4.DX*30)
	}
	if player5.X != (1+player5.DX*fps) || player5.Y != (1+player5.DY*fps) {
		t.Errorf("player5(%d, %d); expected %d", player5.X, player5.Y, 1+player5.DX*30)
	}
	if player6.X != (1+player6.DX*fps) || player6.Y != (1+player6.DY*fps) {
		t.Errorf("player6(%d, %d); expected %d", player6.X, player6.Y, 1+player6.DX*30)
	}
	if player7.X != (1+player7.DX*fps) || player7.Y != (1+player7.DY*fps) {
		t.Errorf("player7(%d, %d); expected %d", player7.X, player7.Y, 1+player7.DX*30)
	}
	if player8.X != 1 || player8.Y != 1 {
		t.Errorf("player8(%d, %d); expected %d", player8.X, player8.Y, 1)
	}
	if player9.X != 1 || player9.Y != 1 {
		t.Errorf("player9(%d, %d); expected %d", player9.X, player9.Y, 1)
	}
}

func BenchmarkECS(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Create a new world
		ecs := New()

		// Add some interacting systems, working with same and different components
		ecs.AddSystem(&MoveSystem{}, PositionComponent{}, VelocityComponent{})
		ecs.AddSystem(&CollisionSystem{}, PositionComponent{}, BoundsComponent{})

		// Add some entities, which are not captured by any system
		ecs.CreateEntity(&PositionComponent{X: 1, Y: 1})          // 1
		ecs.CreateEntity(&VelocityComponent{DX: 2, DY: 2})        // 2
		ecs.CreateEntity(&BoundsComponent{Width: 10, Height: 10}) // 3

		// Add two fully fledged players, with more than needed components
		player4 := createPlayer("player4")
		ecs.CreateEntity(&player4.PositionComponent, &player4.VelocityComponent, &player4.BoundsComponent) // 4
		player5 := createPlayer("player5")
		ecs.CreateEntity(&player5.PositionComponent, &player5.VelocityComponent, &player5.BoundsComponent) // 5

		// Add reduced players, only fitting one system each
		player6 := createPlayer("player6")
		ecs.CreateEntity(&player6.PositionComponent, &player6.VelocityComponent) // 6
		player7 := createPlayer("player7")
		ecs.CreateEntity(&player7.PositionComponent, &player7.VelocityComponent) // 7
		player8 := createPlayer("player8")
		ecs.CreateEntity(&player8.PositionComponent, &player8.BoundsComponent) // 8
		player9 := createPlayer("player9")
		ecs.CreateEntity(&player9.PositionComponent, &player9.BoundsComponent) // 9

		// bench one update
		ecs.Update(33 * time.Millisecond)
	}
}
