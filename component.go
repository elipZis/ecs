package ecs

type Component[T any] struct {
	components map[uint64]T
}

func newComponent[T any]() (this *Component[T]) {
	this = new(Component[T])
	this.components = make(map[uint64]T)
	return this
}

func (this *Component[T]) add(entityId uint64, c T) {
	this.components[entityId] = c
}

func (this *Component[T]) remove(entityId uint64) {
	delete(this.components, entityId)
}

func (this *Component[T]) get() map[uint64]T {
	return this.components
}

func (this *Component[T]) query(entityIds []uint64) map[uint64]T {
	components := make(map[uint64]T)
	for _, entityId := range entityIds {
		components[entityId] = this.components[entityId]
	}
	return components
}
