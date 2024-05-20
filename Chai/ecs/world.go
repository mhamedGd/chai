package ecs

import (
	"math"
	"reflect"
	"sync/atomic"
)

var (
	DefaultAllocation = 0
)

const (
	InvalidEntity Id = 0
	firstEntity   Id = 1
	MaxEntity     Id = math.MaxUint32
)

type World struct {
	idCounter    atomic.Uint64
	nextId       Id
	minId, maxId Id
	arch         *internalMap[Id, archetypeId]
	engine       *archEngine
	resources    map[reflect.Type]any
}

func NewWorld() *World {
	return &World{
		nextId:    firstEntity + 1,
		minId:     firstEntity + 1,
		maxId:     MaxEntity,
		arch:      newInternalMap[Id, archetypeId](DefaultAllocation),
		engine:    newArchEngine(),
		resources: make(map[reflect.Type]any),
	}
}

func (w *World) SetIdRange(min, max Id) {
	if min <= firstEntity {
		panic("min must be greater than 1")
	}
	if max <= firstEntity {
		panic("max must be greater than 1")
	}
	if max > MaxEntity {
		panic("max must be lest than max")
	}

	w.minId = min
	w.maxId = max
}

func (w *World) NewId() Id {
	for {
		val := w.idCounter.Load()
		if w.idCounter.CompareAndSwap(val, val+1) {
			return (Id(val) % (w.maxId - w.minId)) + w.minId
		}
	}
}

func Write(world *World, id Id, comp ...Component) {
	world.Write(id, comp...)
}

func (world *World) Write(id Id, comp ...Component) {
	if len(comp) <= 0 {
		return
	}

	archId, ok := world.arch.Get(id)

	if ok {
		newarchetypeId := world.engine.rewriteArch(archId, id, comp...)
		world.arch.Put(id, newarchetypeId)
	} else {
		archId = world.engine.getArchetypeId(comp...)
		world.arch.Put(id, archId)
		world.engine.write(archId, id, comp...)
	}
}

func Read[T any](world *World, id Id) (T, bool) {
	var ret T
	archId, ok := world.arch.Get(id)
	if !ok {
		return ret, false
	}

	return readArch[T](world.engine, archId, id)
}

func ReadPtr[T any](world *World, id Id) *T {
	archId, ok := world.arch.Get(id)
	if !ok {
		return nil
	}

	return readPtrArch[T](world.engine, archId, id)
}

func Delete(world *World, id Id) bool {
	archId, ok := world.arch.Get(id)

	if !ok {
		return false
	}

	world.arch.Delete(id)

	world.engine.TagForDeleltion(archId, id)

	return true
}

func DeleteComponent(world *World, id Id, comp ...Component) {
	archId, ok := world.arch.Get(id)

	if !ok {
		return
	}

	ent := world.engine.ReadEntity(archId, id)

	for i := range comp {
		ent.Delete(comp[i])
	}

	world.arch.Delete(id)
	world.engine.TagForDeleltion(archId, id)

	if len(ent.comp) > 0 {
		world.Write(id, ent.comp...)
	}
}

func (world *World) Exists(id Id) bool {
	return world.arch.Has(id)
}

// -------------------------------------------------
// - Resources
// -------------------------------------------------

func resourceName(t any) reflect.Type {
	return reflect.TypeOf(t)
}

func PutResources[T any](world *World, resource *T) {
	name := resourceName(resource)
	world.resources[name] = resource
}

func GetResource[T any](world *World) *T {
	var t T
	name := resourceName(&t)
	anyVal, ok := world.resources[name]
	if !ok {
		return nil
	}

	return anyVal.(*T)
}
