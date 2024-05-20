package ecs

type View1[A any] struct {
	world    *World
	filter   filterList
	storageA *componentSliceStorage[A]
}

func (v *View1[A]) initialize(world *World) any {
	return Query1[A](world)
}

func Query1[A any](world *World, filters ...Filter) *View1[A] {
	storageA := getStorage[A](world.engine)

	var AA A

	comps := []componentId{
		name(AA),
	}

	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View1[A]{
		world:    world,
		filter:   filterList,
		storageA: storageA,
	}

	return v
}

func (v *View1[A]) Read(id Id) *A {
	if id == InvalidEntity {
		return nil
	}

	archId, ok := v.world.arch.Get(id)

	if !ok {
		return nil
	}

	lookup := v.world.engine.lookup[archId]
	if lookup == nil {
		panic("LookupList is missing!")
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		return nil
	}

	var retA *A

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}

	return retA
}

func (v *View1[A]) Count() int {
	v.filter.regenerate(v.world)

	total := 0
	for _, archId := range v.filter.archIds {
		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		total += lookup.Len()
	}

	return total
}

func (v *View1[A]) MapId(lambda func(id Id, a *A)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	for _, archId := range v.filter.archIds {
		sliceA, _ = v.storageA.slice[archId]

		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		ids := lookup.id

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}

		retA = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			}

			if compA != nil {
				retA = &compA[idx]
			}
			lambda(ids[idx], retA)
		}
	}
}

// - Query 2
// ------------------------------------------------------

type View2[A, B any] struct {
	world    *World
	filter   filterList
	storageA *componentSliceStorage[A]
	storageB *componentSliceStorage[B]
}

func (v *View2[A, B]) initialize(world *World) any {
	return Query2[A, B](world)
}

func Query2[A, B any](world *World, filters ...Filter) *View2[A, B] {
	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)

	var AA A
	var BB B

	comps := []componentId{
		name(AA),
		name(BB),
	}

	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View2[A, B]{
		world:    world,
		filter:   filterList,
		storageA: storageA,
		storageB: storageB,
	}

	return v
}

func (v *View2[A, B]) Read(id Id) (*A, *B) {
	if id == InvalidEntity {
		return nil, nil
	}

	archId, ok := v.world.arch.Get(id)

	if !ok {
		return nil, nil
	}

	lookup := v.world.engine.lookup[archId]
	if lookup == nil {
		panic("LookupList is missing!")
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		return nil, nil
	}

	var retA *A
	var retB *B

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}

	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}

	return retA, retB
}

func (v *View2[A, B]) Count() int {
	v.filter.regenerate(v.world)

	total := 0
	for _, archId := range v.filter.archIds {
		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		total += lookup.Len()
	}

	return total
}

func (v *View2[A, B]) MapId(lambda func(id Id, a *A, b *B)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	for _, archId := range v.filter.archIds {
		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]

		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		ids := lookup.id

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}

		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}

		retA = nil
		retB = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			}

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}

			lambda(ids[idx], retA, retB)
		}
	}
}

// - Query 3
// ------------------------------------------------------
// ------------------------------------------------------

type View3[A, B, C any] struct {
	world    *World
	filter   filterList
	storageA *componentSliceStorage[A]
	storageB *componentSliceStorage[B]
	storageC *componentSliceStorage[C]
}

func (v *View3[A, B, C]) initialize(world *World) any {
	return Query3[A, B, C](world)
}

func Query3[A, B, C any](world *World, filters ...Filter) *View3[A, B, C] {
	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)

	var AA A
	var BB B
	var CC C

	comps := []componentId{
		name(AA),
		name(BB),
		name(CC),
	}

	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View3[A, B, C]{
		world:    world,
		filter:   filterList,
		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
	}

	return v
}

func (v *View3[A, B, C]) Read(id Id) (*A, *B, *C) {
	if id == InvalidEntity {
		return nil, nil, nil
	}

	archId, ok := v.world.arch.Get(id)

	if !ok {
		return nil, nil, nil
	}

	lookup := v.world.engine.lookup[archId]
	if lookup == nil {
		panic("LookupList is missing!")
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		return nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}

	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}

	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}

	return retA, retB, retC
}

func (v *View3[A, B, C]) Count() int {
	v.filter.regenerate(v.world)

	total := 0
	for _, archId := range v.filter.archIds {
		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		total += lookup.Len()
	}

	return total
}

func (v *View3[A, B, C]) MapId(lambda func(id Id, a *A, b *B, c *C)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	for _, archId := range v.filter.archIds {
		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]

		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		ids := lookup.id

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}

		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}

		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}

		retA = nil
		retB = nil
		retC = nil
		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			}

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}

			if compB != nil {
				retC = &compC[idx]
			}

			lambda(ids[idx], retA, retB, retC)
		}
	}
}

// - Query 4
// ------------------------------------------------------
// ------------------------------------------------------

type View4[A, B, C, D any] struct {
	world    *World
	filter   filterList
	storageA *componentSliceStorage[A]
	storageB *componentSliceStorage[B]
	storageC *componentSliceStorage[C]
	storageD *componentSliceStorage[D]
}

func (v *View4[A, B, C, D]) initialize(world *World) any {
	return Query4[A, B, C, D](world)
}

func Query4[A, B, C, D any](world *World, filters ...Filter) *View4[A, B, C, D] {
	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D

	comps := []componentId{
		name(AA),
		name(BB),
		name(CC),
		name(DD),
	}

	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View4[A, B, C, D]{
		world:    world,
		filter:   filterList,
		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
	}

	return v
}

func (v *View4[A, B, C, D]) Read(id Id) (*A, *B, *C, *D) {
	if id == InvalidEntity {
		return nil, nil, nil, nil
	}

	archId, ok := v.world.arch.Get(id)

	if !ok {
		return nil, nil, nil, nil
	}

	lookup := v.world.engine.lookup[archId]
	if lookup == nil {
		panic("LookupList is missing!")
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		return nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}

	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}

	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}

	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}

	return retA, retB, retC, retD
}

func (v *View4[A, B, C, D]) Count() int {
	v.filter.regenerate(v.world)

	total := 0
	for _, archId := range v.filter.archIds {
		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		total += lookup.Len()
	}

	return total
}

func (v *View4[A, B, C, D]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	for _, archId := range v.filter.archIds {
		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]

		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		ids := lookup.id

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}

		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}

		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}

		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil

		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			}

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}

			lambda(ids[idx], retA, retB, retC, retD)
		}
	}
}

// - Query 5
// ------------------------------------------------------
// ------------------------------------------------------

type View5[A, B, C, D, E any] struct {
	world    *World
	filter   filterList
	storageA *componentSliceStorage[A]
	storageB *componentSliceStorage[B]
	storageC *componentSliceStorage[C]
	storageD *componentSliceStorage[D]
	storageE *componentSliceStorage[E]
}

func (v *View5[A, B, C, D, E]) initialize(world *World) any {
	return Query5[A, B, C, D, E](world)
}

func Query5[A, B, C, D, E any](world *World, filters ...Filter) *View5[A, B, C, D, E] {
	storageA := getStorage[A](world.engine)
	storageB := getStorage[B](world.engine)
	storageC := getStorage[C](world.engine)
	storageD := getStorage[D](world.engine)
	storageE := getStorage[E](world.engine)

	var AA A
	var BB B
	var CC C
	var DD D
	var EE E

	comps := []componentId{
		name(AA),
		name(BB),
		name(CC),
		name(DD),
		name(EE),
	}

	filterList := newFilterList(comps, filters...)
	filterList.regenerate(world)

	v := &View5[A, B, C, D, E]{
		world:    world,
		filter:   filterList,
		storageA: storageA,
		storageB: storageB,
		storageC: storageC,
		storageD: storageD,
		storageE: storageE,
	}

	return v
}

func (v *View5[A, B, C, D, E]) Read(id Id) (*A, *B, *C, *D, *E) {
	if id == InvalidEntity {
		return nil, nil, nil, nil, nil
	}

	archId, ok := v.world.arch.Get(id)

	if !ok {
		return nil, nil, nil, nil, nil
	}

	lookup := v.world.engine.lookup[archId]
	if lookup == nil {
		panic("LookupList is missing!")
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		return nil, nil, nil, nil, nil
	}

	var retA *A
	var retB *B
	var retC *C
	var retD *D
	var retE *E

	sliceA, ok := v.storageA.slice[archId]
	if ok {
		retA = &sliceA.comp[index]
	}

	sliceB, ok := v.storageB.slice[archId]
	if ok {
		retB = &sliceB.comp[index]
	}

	sliceC, ok := v.storageC.slice[archId]
	if ok {
		retC = &sliceC.comp[index]
	}

	sliceD, ok := v.storageD.slice[archId]
	if ok {
		retD = &sliceD.comp[index]
	}

	sliceE, ok := v.storageE.slice[archId]
	if ok {
		retE = &sliceE.comp[index]
	}

	return retA, retB, retC, retD, retE
}

func (v *View5[A, B, C, D, E]) Count() int {
	v.filter.regenerate(v.world)

	total := 0
	for _, archId := range v.filter.archIds {
		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		total += lookup.Len()
	}

	return total
}

func (v *View5[A, B, C, D, E]) MapId(lambda func(id Id, a *A, b *B, c *C, d *D, e *E)) {
	v.filter.regenerate(v.world)

	var sliceA *componentSlice[A]
	var compA []A
	var retA *A

	var sliceB *componentSlice[B]
	var compB []B
	var retB *B

	var sliceC *componentSlice[C]
	var compC []C
	var retC *C

	var sliceD *componentSlice[D]
	var compD []D
	var retD *D

	var sliceE *componentSlice[E]
	var compE []E
	var retE *E

	for _, archId := range v.filter.archIds {
		sliceA, _ = v.storageA.slice[archId]
		sliceB, _ = v.storageB.slice[archId]
		sliceC, _ = v.storageC.slice[archId]
		sliceD, _ = v.storageD.slice[archId]
		sliceE, _ = v.storageE.slice[archId]

		lookup := v.world.engine.lookup[archId]
		if lookup == nil {
			panic("LookupList is missing!")
		}

		ids := lookup.id

		compA = nil
		if sliceA != nil {
			compA = sliceA.comp
		}

		compB = nil
		if sliceB != nil {
			compB = sliceB.comp
		}

		compC = nil
		if sliceC != nil {
			compC = sliceC.comp
		}

		compD = nil
		if sliceD != nil {
			compD = sliceD.comp
		}

		compE = nil
		if sliceE != nil {
			compE = sliceE.comp
		}

		retA = nil
		retB = nil
		retC = nil
		retD = nil
		retE = nil

		for idx := range ids {
			if ids[idx] == InvalidEntity {
				continue
			}

			if compA != nil {
				retA = &compA[idx]
			}
			if compB != nil {
				retB = &compB[idx]
			}
			if compC != nil {
				retC = &compC[idx]
			}
			if compD != nil {
				retD = &compD[idx]
			}
			if compE != nil {
				retE = &compE[idx]
			}

			lambda(ids[idx], retA, retB, retC, retD, retE)
		}
	}
}
