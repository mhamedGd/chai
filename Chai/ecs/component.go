package ecs

import "fmt"

type componentId uint16

type Component interface {
	write(*archEngine, archetypeId, int)
	id() componentId
}

type Box[T any] struct {
	Comp   T
	compId componentId
}

func C[T any](comp T) Box[T] {
	return Box[T]{
		Comp:   comp,
		compId: name(comp),
	}
}

func (c Box[T]) write(engine *archEngine, archId archetypeId, index int) {
	store := getStorageByCompId[T](engine, c.id())
	writeArch[T](engine, archId, index, store, c.Comp)
}
func (c Box[T]) id() componentId {
	if c.compId == inavlidComponentId {
		c.compId = name(c.Comp)
	}

	return c.compId
}

func (c Box[T]) Get() T {
	return c.Comp
}

const maxComponentId = 255

var blankArchMask archetypeMask

type archetypeMask [4]uint64

func buildArchMask(comps ...Component) archetypeMask {
	var mask archetypeMask
	for _, comp := range comps {
		c := comp.id()
		idx := c / 64
		offset := c - (64 * idx)
		mask[idx] |= (1 << offset)
	}

	return mask
}

func buildArchMaskFromAny(comps ...any) archetypeMask {
	var mask archetypeMask
	for _, comp := range comps {
		c := name(comp)
		idx := c / 64
		offset := c - (64 * idx)
		mask[idx] |= (1 << offset)
	}

	return mask
}

func (m archetypeMask) bitwiseOr(a archetypeMask) archetypeMask {
	for i := range m {
		m[i] = m[i] | a[i]
	}

	return m
}

func (m archetypeMask) bitwiseAnd(a archetypeMask) archetypeMask {
	for i := range m {
		m[i] = m[i] & a[i]
	}

	return m
}

type componentRegistry struct {
	archSet     [][]archetypeId
	archMask    map[archetypeMask]archetypeId
	revArchMask map[archetypeId]archetypeMask
}

func newComponentRegistry() *componentRegistry {
	r := &componentRegistry{
		archSet:     make([][]archetypeId, maxComponentId+1),
		archMask:    make(map[archetypeMask]archetypeId),
		revArchMask: make(map[archetypeId]archetypeMask),
	}

	return r
}

func (r *componentRegistry) print() {
	fmt.Println("--- componentRegistry ---")
	fmt.Println("-- archSet --")
	for name, set := range r.archSet {
		fmt.Printf("name(%d): archId: [ ", name)
		for archId := range set {
			fmt.Printf("%d ", archId)
		}
		fmt.Printf("]\n")
	}
}

func (r *componentRegistry) getArchetypeId(engine *archEngine, comps ...Component) archetypeId {
	mask := buildArchMask(comps...)
	archId, ok := r.archMask[mask]
	if !ok {
		archId = engine.newArchetypeId(mask)
		r.archMask[mask] = archId
		r.revArchMask[archId] = mask

		for _, comp := range comps {
			compId := comp.id()
			r.archSet[compId] = append(r.archSet[compId], archId)
		}
	}

	return archId
}

func (r *componentRegistry) archIdOverlapsMask(archId archetypeId, compArchMask archetypeMask) bool {
	archMaskToCheck, ok := r.revArchMask[archId]
	if !ok {
		panic("Bug: Invalid ArchId used")
	}

	resultArchMask := archMaskToCheck.bitwiseAnd(compArchMask)
	if resultArchMask != blankArchMask {
		return true
	}

	return false
}
