package ecs

import "slices"

type Filter interface {
	Filter([]componentId) []componentId
}

type without struct {
	mask archetypeMask
}

func Without(comps ...any) without {
	return without{
		mask: buildArchMaskFromAny(comps...),
	}
}

func (w without) Filter(list []componentId) []componentId {
	return list
}

type with struct {
	copms []componentId
}

func With(comps ...any) with {
	ids := make([]componentId, len(comps))
	for i := range comps {
		ids[i] = name(comps[i])
	}

	return with{
		copms: ids,
	}
}

func (w with) Filter(list []componentId) []componentId {
	return append(list, w.copms...)
}

type optional struct {
	comps []componentId
}

func Optional(comps ...any) optional {
	ids := make([]componentId, len(comps))
	for i := range comps {
		ids[i] = name(comps[i])
	}

	return optional{
		comps: ids,
	}
}

func (f optional) Filter(list []componentId) []componentId {
	for i := 0; i < len(list); i++ {
		for j := range f.comps {
			if list[i] == f.comps[j] {
				list[i] = list[len(list)-1]
				list = list[:len(list)-1]

				i--
				break
			}
		}
	}

	return list
}

type filterList struct {
	comps                     []componentId
	withoutArchMask           archetypeMask
	cachedArchetypeGeneration int
	archIds                   []archetypeId
}

func newFilterList(comps []componentId, filters ...Filter) filterList {
	var withoutArchMask archetypeMask
	for _, f := range filters {
		withoutFilter, isWithout := f.(without)
		if isWithout {
			withoutArchMask = withoutFilter.mask
		} else {
			comps = f.Filter(comps)
		}
	}

	return filterList{
		comps:           comps,
		withoutArchMask: withoutArchMask,
		archIds:         make([]archetypeId, 0),
	}
}

func (f *filterList) regenerate(world *World) {
	if world.engine.getGeneration() != f.cachedArchetypeGeneration {
		f.archIds = world.engine.FilterList(f.archIds, f.comps)

		if f.withoutArchMask != blankArchMask {
			f.archIds = slices.DeleteFunc(f.archIds, func(archId archetypeId) bool {
				return world.engine.dcr.archIdOverlapsMask(archId, f.withoutArchMask)
			})
		}

		f.cachedArchetypeGeneration = world.engine.getGeneration()
	}
}
