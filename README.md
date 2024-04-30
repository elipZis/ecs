# Entity Component System

A simple, opinionated ECS for quick usage (but without optimization).

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

### Entities

To create and register a new entity, call `world.CreateEntity(&player.PositionComponent, &player.VelocityComponent, &player.OtherComponent)`

It returns the entity with id and components, unique to this world.

The entity and their components will be injected into all systems, intersecting the component type combination. More is ok, less does not match!

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Credits

- [elipZis GmbH](https://elipZis.com)
- [NeA](https://github.com/nea)
- [All Contributors](https://github.com/elipZis/laravel-cacheable-model/contributors)

Heavily inspired by
- https://github.com/EngoEngine/ecs
- https://github.com/amethyst/legion

Shout-out to them 🙏

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.