# Entity Component System

A very simple, reflect.Type-based, opinionated ECS for quick usage (but without optimization).

## Installation

`go get github.com/elipZis/ecs`

## Usage

Create a new world with `world := ecs.New()`. This is your starting point.

### Components

Components can be any struct you define, for example

```go
type PositionComponent struct {
    X int
    Y int
}

type VelocityComponent struct {
    DX int
    DY int
}
```

You can create components on their own or embedded in other structs

```go
type Player struct {
    PositionComponent
    VelocityComponent
    OtherComponent
    
    name string
}
```

### Systems

A system must embed the `EntitySystem` and implement the `System` interface `Run(ecs *ecs.ECS, dt time.Duration)` function:

```go
type MoveSystem struct {
   EntitySystem
}

func (this *MoveSystem) Run(ecs *ecs.ECS, dt time.Duration) {
   ...
}

```

Add systems to the world via e.g. `world.AddSystem(&MoveSystem{}, &PositionComponent{}, &VelocityComponent{})`

* First is the function to run
* Following are N component types to listen on
  * The system will only be invoked, if the registered component types match those of an entity

**Note: Systems must be registered before entities!**

#### Priority

By default all systems have a priority of 0. 
If you want to order systems, you can implement

```go
Priority() int
```

The higher the number, the earlier the system is called.

#### Components & Entities

Inside a system you can access the ECS itself and you can get all components. 
The entity system itself knows about all their own entity ids and can iterate on these. For example:

```go
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
```

or

```go
func (this *MoveSystem) Run(ecs *ECS, dt time.Duration) {
    positions := ecs.GetComponentsFor[*component.PositionComponent](e)
    velocities := ecs.GetComponents(VelocityComponent{})
    
    for _, entityId := range this.entities {
        p := positions[entityId]
        v := velocities[entityId]
        p.X += v.DX
        p.Y += v.DY
    }
}
```

There are several helper functions to provide access to the different components, context, entities etc.

### Entities

To create and register a new entity, call 

```go
world.CreateEntity(&player.PositionComponent, &player.VelocityComponent, &player.OtherComponent)
```

It returns the entity with id and components, unique to this world.

The entity and their components will be injected into all systems, intersecting the component type combination. More is ok, less does not match!

#### Remove Entity

To remove an entity, call e.g. `ecs.RemoveEntity(id uint64)` on the world or in a system.
If you remove an entity, it will be marked to be removed before the next iteration.

To immediately remove an entity, with consequences for subsequent systems, call `ecs.RemoveEntityNow(id uint64)`.

### Context

Via `world.AddContext(...)` you can add anything as context, available globally to all systems to query for via `world.GetContext(...)`.

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Notes

This package is heavily inspired by

- https://github.com/EngoEngine/ecs
- https://github.com/amethyst/legion

Kudos, shout-out and thanks to them 🙏

## Credits

- [elipZis GmbH](https://elipZis.com)
- [NeA](https://github.com/nea)
- [All Contributors](https://github.com/elipZis/laravel-cacheable-model/contributors)

## Disclaimer

This source and the whole package comes without a warranty. It may or may not harm your computer. Please use with care. 
Any damage cannot be related back to the author. The source has been tested on a virtual environment and scanned for viruses and has passed all tests.
It is not optimized for production, speed or memory consumption but just a personal open-source package.

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.