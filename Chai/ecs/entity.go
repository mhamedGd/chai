package ecs

type Entity struct {
	comp []Component
}

func NewEntity(components ...Component) *Entity {
	return &Entity{
		comp: components,
	}
}

func (e *Entity) findIndex(compId componentId) int {
	for i := range e.comp {
		if compId == e.comp[i].id() {
			return i
		}
	}

	return -1
}

func (e *Entity) Add(components ...Component) {
	for i := range components {
		idx := e.findIndex(components[i].id())
		if idx < 0 {
			e.comp = append(e.comp, components[i])
		} else {
			e.comp[idx] = components[i]
		}
	}
}

func (e *Entity) Merge(e2 *Entity) {
	e.Add(e2.comp...)
}

func (e *Entity) Comps() []Component {
	return e.comp
}

func ReadFromEntity[T any](ent *Entity) (T, bool) {
	var t T
	n := name(t)
	idx := ent.findIndex(n)
	if idx < 0 {
		return t, false
	}

	icomp := ent.comp[idx]
	return icomp.(Box[T]).Comp, true
}

func (ent *Entity) Write(world *World, id Id) {
	world.Write(id, ent.comp...)
}

func ReadEntity(world *World, id Id) *Entity {
	archId, ok := world.arch.Get(id)
	if !ok {
		return nil
	}

	return world.engine.ReadEntity(archId, id)
}

func (e *Entity) Delete(c Component) {
	compId := c.id()
	idx := e.findIndex(compId)

	if idx < 0 {
		return
	}

	e.comp[idx] = e.comp[len(e.comp)-1]
	e.comp = e.comp[:len(e.comp)-1]
}

func (e *Entity) Clear() {
	e.comp = e.comp[:0]
}

type RawEntity struct {
	comp map[componentId]any
}

func NewRawEntity(components ...any) *RawEntity {
	c := make(map[componentId]any)
	for i := range components {
		c[name(components[i])] = components[i]
	}
	return &RawEntity{
		comp: c,
	}
}

func (e *RawEntity) Add(components ...any) {
	for i := range components {
		e.comp[name(components[i])] = components[i]
	}
}
