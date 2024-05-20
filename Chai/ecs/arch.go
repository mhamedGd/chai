package ecs

import "fmt"

type Id uint32

type archetypeId uint32

type componentSlice[T any] struct {
	comp []T
}

func (s *componentSlice[T]) Write(index int, val T) {
	if index == len(s.comp) {
		s.comp = append(s.comp, val)
	} else {
		s.comp[index] = val
	}
}

type lookupList struct {
	index *internalMap[Id, int]
	id    []Id
	holes []int
	mask  archetypeMask
}

func (l *lookupList) Len() int {
	return l.index.Len()
}

func (l *lookupList) addToEasiestHole(id Id) int {
	if len(l.holes) > 0 {
		lastHoleIndex := len(l.holes) - 1

		index := l.holes[lastHoleIndex]
		l.id[index] = id
		l.index.Put(id, index)

		l.holes = l.holes[:lastHoleIndex]

		return index
	} else {
		l.id = append(l.id, id)
		index := len(l.id) - 1
		l.index.Put(id, index)
		return index
	}
}

type storage interface {
	ReadToEntity(*Entity, archetypeId, int) bool
	ReadToRawEntity(*RawEntity, archetypeId, int) bool
	Delete(archetypeId, int)
	print(int)
}

type componentSliceStorage[T any] struct {
	slice map[archetypeId]*componentSlice[T]
}

func (ss *componentSliceStorage[T]) ReadToEntity(entity *Entity, archId archetypeId, index int) bool {
	cSlice, ok := ss.slice[archId]
	if !ok {
		return false
	}

	entity.Add(C(cSlice.comp[index]))

	return true
}

func (ss *componentSliceStorage[T]) ReadToRawEntity(entity *RawEntity, archId archetypeId, index int) bool {
	cSlice, ok := ss.slice[archId]
	if !ok {
		return false
	}

	entity.Add(&cSlice.comp[index])
	return true
}

func (ss *componentSliceStorage[T]) Delete(archId archetypeId, index int) {
	cSlice, ok := ss.slice[archId]

	if !ok {
		return
	}

	lastVal := cSlice.comp[len(cSlice.comp)-1]
	cSlice.comp[index] = lastVal
	cSlice.comp = cSlice.comp[:len(cSlice.comp)-1]
}

func (s *componentSliceStorage[T]) print(amount int) {
	for archId, compSlice := range s.slice {
		fmt.Printf("archId(%d) - %v\n", archId, *compSlice)
	}
}

type archEngine struct {
	generation int

	lookup           []*lookupList
	compSliceStorage []storage
	dcr              *componentRegistry

	archCount map[archetypeId]int
}

func newArchEngine() *archEngine {
	return &archEngine{
		generation:       1,
		lookup:           make([]*lookupList, 0, DefaultAllocation),
		compSliceStorage: make([]storage, maxComponentId+1),
		dcr:              newComponentRegistry(),
		archCount:        make(map[archetypeId]int),
	}
}

func (e *archEngine) newArchetypeId(archMask archetypeMask) archetypeId {
	e.generation++

	archId := archetypeId(len(e.lookup))

	e.lookup = append(e.lookup, &lookupList{
		index: newInternalMap[Id, int](0),
		id:    make([]Id, 0, DefaultAllocation),
		holes: make([]int, 0, DefaultAllocation),
		mask:  archMask,
	})

	return archId
}

func (e *archEngine) getGeneration() int {
	return e.generation
}

func (e *archEngine) count(anything ...any) int {
	comps := make([]componentId, len(anything))

	for i, c := range anything {
		comps[i] = name(c)
	}

	archIds := make([]archetypeId, 0)
	archIds = e.FilterList(archIds, comps)

	total := 0
	for _, archId := range archIds {
		lookup := e.lookup[archId]
		if lookup == nil {
			panic(fmt.Sprintf("Couldn't find archId in archEngine lookup table: %d", archId))
		}

		total = total + len(lookup.id) - len(lookup.holes)
	}

	return total
}

func (e *archEngine) getArchetypeId(comp ...Component) archetypeId {
	return e.dcr.getArchetypeId(e, comp...)
}

func (e *archEngine) FilterList(archIds []archetypeId, comp []componentId) []archetypeId {
	for k := range e.archCount {
		delete(e.archCount, k)
	}

	for _, compId := range comp {
		for _, archId := range e.dcr.archSet[compId] {
			e.archCount[archId] = e.archCount[archId] + 1
		}
	}

	numComponents := len(comp)

	archIds = archIds[:0]
	for archId, count := range e.archCount {
		if count >= numComponents {
			archIds = append(archIds, archId)
		}
	}

	return archIds
}

func getStorage[T any](e *archEngine) *componentSliceStorage[T] {
	var val T
	n := name(val)
	return getStorageByCompId[T](e, n)
}

func getStorageByCompId[T any](e *archEngine, compId componentId) *componentSliceStorage[T] {
	ss := e.compSliceStorage[compId]
	if ss == nil {
		ss = &componentSliceStorage[T]{
			slice: make(map[archetypeId]*componentSlice[T], DefaultAllocation),
		}
		e.compSliceStorage[compId] = ss
	}
	storage := ss.(*componentSliceStorage[T])

	return storage
}

func (e *archEngine) getOrAddLookupIndex(archId archetypeId, id Id) int {
	lookup := e.lookup[archId]

	if len(lookup.holes) >= 1024 {
		e.CleanupHoles(archId)
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		index = lookup.addToEasiestHole(id)
	}

	return index
}

func (e *archEngine) write(archId archetypeId, id Id, comp ...Component) {
	index := e.getOrAddLookupIndex(archId, id)

	for i := range comp {
		comp[i].write(e, archId, index)
	}
}

func writeArch[T any](e *archEngine, archId archetypeId, index int, store *componentSliceStorage[T], val T) {
	cSlice, ok := store.slice[archId]
	if !ok {
		cSlice = &componentSlice[T]{
			comp: make([]T, 0, DefaultAllocation),
		}
		store.slice[archId] = cSlice
	}

	cSlice.Write(index, val)
}

func readArch[T any](e *archEngine, archId archetypeId, id Id) (T, bool) {
	var ret T
	lookup := e.lookup[archId]
	if lookup == nil {
		return ret, false
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		return ret, false
	}

	n := name(ret)
	ss := e.compSliceStorage[n]

	if ss == nil {
		return ret, false
	}

	storage, ok := ss.(*componentSliceStorage[T])
	if !ok {
		panic(fmt.Sprintf("Wrong componentSliceStorage[T] type: %d != %d", name(ss), name(ret)))
	}

	cSlice, ok := storage.slice[archId]
	if !ok {
		return ret, false
	}

	return cSlice.comp[index], true
}

func readPtrArch[T any](e *archEngine, archId archetypeId, id Id) *T {
	var ret T

	lookup := e.lookup[archId]
	if lookup == nil {
		return nil
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		return nil
	}

	n := name(ret)
	ss := e.compSliceStorage[n]
	if ss == nil {
		return nil
	}

	storage, ok := ss.(*componentSliceStorage[T])
	if !ok {
		panic(fmt.Sprintf("Wrong componentSliceStorage[T] type: %d != %d", name(ss), name(ret)))
	}

	cSlice, ok := storage.slice[archId]
	if !ok {
		return nil
	}

	return &cSlice.comp[index]
}

func (e *archEngine) rewriteArch(archId archetypeId, id Id, comp ...Component) archetypeId {
	lookup := e.lookup[archId]
	oldMask := lookup.mask
	addMask := buildArchMask(comp...)
	newMask := oldMask.bitwiseOr(addMask)

	if oldMask == newMask {
		e.write(archId, id, comp...)

		return archId
	} else {
		ent := e.ReadEntity(archId, id)
		ent.Add(comp...)
		combinedComps := ent.Comps()

		newArchId := e.getArchetypeId(combinedComps...)

		e.TagForDeleltion(archId, id)

		e.write(newArchId, id, combinedComps...)
		return newArchId
	}
}

func (e *archEngine) ReadEntity(archId archetypeId, id Id) *Entity {
	lookup := e.lookup[archId]
	if lookup == nil {
		panic("Archetype doesn't have lookup list")
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		panic("Archetype doesn't contain ID")
	}

	ent := NewEntity()
	for n := range e.compSliceStorage {
		if e.compSliceStorage[n] != nil {
			e.compSliceStorage[n].ReadToEntity(ent, archId, index)
		}
	}

	return ent
}

func (e *archEngine) ReadToRawEntity(archId archetypeId, id Id) *RawEntity {
	lookup := e.lookup[archId]
	if lookup == nil {
		panic("Archetype doesn't have lookup list")
	}

	index, ok := lookup.index.Get(id)
	if !ok {
		panic("Archetype doesn't contain ID")
	}

	ent := NewRawEntity()
	for n := range e.compSliceStorage {
		if e.compSliceStorage[n] != nil {
			e.compSliceStorage[n].ReadToRawEntity(ent, archId, index)
		}
	}

	return ent
}

func (e *archEngine) TagForDeleltion(archId archetypeId, id Id) {
	lookup := e.lookup[archId]
	if lookup == nil {
		panic("archetype doesn't have lookup list")
	}

	index, ok := lookup.index.Get(id)

	if !ok {
		panic("Archetype doesn't contain ID")
	}

	lookup.id[index] = InvalidEntity
	lookup.index.Delete(id)

	lookup.holes = append(lookup.holes, index)
}

func (e *archEngine) CleanupHoles(archId archetypeId) {
	lookup := e.lookup[archId]

	if lookup == nil {
		panic("Archetype doesn't have lookup list")
	}

	for _, index := range lookup.holes {
		for {
			lastIndex := len(lookup.id) - 1
			if lastIndex < 0 {
				break
			}

			lastId := lookup.id[lastIndex]
			if lastId == InvalidEntity {
				lookup.id = lookup.id[:lastIndex]
				for n := range e.compSliceStorage {
					if e.compSliceStorage[n] != nil {
						e.compSliceStorage[n].Delete(archId, lastIndex)
					}
				}

				continue
			}

			break
		}

		if index >= len(lookup.id) {
			continue
		}

		lastIndex := len(lookup.id) - 1

		lastId := lookup.id[lastIndex]
		if lastId == InvalidEntity {
			panic("Bug: This shouldn't happen")
		}

		lookup.id[index] = lastId
		lookup.id = lookup.id[:lastIndex]
		lookup.index.Put(lastId, index)
		for n := range e.compSliceStorage {
			if e.compSliceStorage[n] != nil {
				e.compSliceStorage[n].Delete(archId, index)
			}
		}
	}

	lookup.holes = lookup.holes[:0]
}
