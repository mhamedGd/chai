package chai

type Pair[A, B any] struct {
	First  A
	Second B
}

const MAX_QUADTREE_DEPTH = 8

type item_key_type = int64

type QuadTreeItemLocation[T any] struct {

	// container *[]Pair[Rect, T]
	container *Map[item_key_type, Pair[Rect, T]]
	index     int64
}

type QuadTreeItem[T any] struct {
	item  T
	pItem QuadTreeItemLocation[*QuadTreeItem[T]]
}

type StaticQuadTree[T any] struct {
	depth   int
	rect    Rect
	rChild  [4]Rect
	pChild  [4]*StaticQuadTree[T]
	pItems  Map[item_key_type, Pair[Rect, T]]
	counter int64
}

func NewStaticQuadTree[T any](size Rect, _nDepth int) StaticQuadTree[T] {
	// l := list.New()
	sqt := StaticQuadTree[T]{}
	sqt.depth = _nDepth
	sqt.Resize(size)

	return sqt
}

func (sqt *StaticQuadTree[T]) Resize(rArea Rect) {
	sqt.clear()

	sqt.rect = rArea
	vChildSize := sqt.rect.Size.Scale(0.5)

	sqt.rChild = [4]Rect{
		{sqt.rect.Position, vChildSize},
		{NewVector2f(sqt.rect.Position.X+vChildSize.X, sqt.rect.Position.Y), vChildSize},
		{NewVector2f(sqt.rect.Position.X, sqt.rect.Position.Y+vChildSize.Y), vChildSize},
		{sqt.rect.Position.Add(vChildSize), vChildSize},
	}
}

func (sqt *StaticQuadTree[T]) clear() {
	sqt.pItems.Clear()

	for i := 0; i < 4; i++ {
		if sqt.pChild[i] != nil {
			sqt.pChild[i].clear()
		}
		sqt.pChild[i] = nil
	}
}

func (sqt *StaticQuadTree[T]) size() int {
	nCount := sqt.pItems.Count()
	for i := 0; i < 4; i++ {
		if sqt.pChild[i] != nil {
			nCount += len(sqt.pChild)
		}
	}

	return nCount
}

func (sqt *StaticQuadTree[T]) Insert(item T, itemsize Rect) QuadTreeItemLocation[T] {

	for i := 0; i < 4; i++ {
		if sqt.rChild[i].ContainsRect(itemsize) {
			if sqt.depth+1 < MAX_QUADTREE_DEPTH {
				if sqt.pChild[i] == nil {
					_tree := NewStaticQuadTree[T](sqt.rChild[i], sqt.depth+1)
					sqt.pChild[i] = &_tree
				}

				return sqt.pChild[i].Insert(item, itemsize)
			}
		}
	}

	// sqt.pItems = append(sqt.pItems, Pair[Rect, T]{itemsize, item})
	sqt.pItems.Insert(sqt.counter, Pair[Rect, T]{itemsize, item})
	sqt.counter += 1
	return QuadTreeItemLocation[T]{
		container: &sqt.pItems,
		index:     sqt.counter,
	}
}
func (sqt *StaticQuadTree[T]) InsertWithIndex(item T, itemsize Rect, index int64) QuadTreeItemLocation[T] {
	for i := 0; i < 4; i++ {
		if sqt.rChild[i].ContainsRect(itemsize) {
			if sqt.depth+1 < MAX_QUADTREE_DEPTH {
				if sqt.pChild[i] == nil {
					_tree := NewStaticQuadTree[T](sqt.rChild[i], sqt.depth+1)
					sqt.pChild[i] = &_tree
				}

				return sqt.pChild[i].InsertWithIndex(item, itemsize, index)
			}
		}
	}

	// sqt.pItems = append(sqt.pItems, Pair[Rect, T]{itemsize, item})
	sqt.pItems.Insert(index, Pair[Rect, T]{itemsize, item})
	return QuadTreeItemLocation[T]{
		container: &sqt.pItems,
		index:     index,
	}
}

// func (sqt *StaticQuadTree[T]) Remove(item T) bool {
// 	it := sqt.pItems.FindIf(func(a Pair[Rect, T]) bool {
// 		&item == a.Second
// 	})
// }

func (sqt *StaticQuadTree[T]) Search(rArea Rect) List[T] {
	listItems := NewList[T]()

	return sqt.searchThrough(rArea, listItems)
}

func (sqt *StaticQuadTree[T]) searchThrough(rArea Rect, listItems List[T]) List[T] {
	for _, v := range sqt.pItems.AllItems() {
		if rArea.OverlapsRect(v.First) {
			// listItems = append(listItems, v.Second)
			listItems.PushBack(v.Second)
		}
	}

	for i := 0; i < 4; i++ {
		if sqt.pChild[i] != nil {
			if rArea.ContainsRect(sqt.rChild[i]) {
				listItems = sqt.pChild[i].Items(listItems)

			} else if sqt.rChild[i].OverlapsRect(rArea) {
				listItems = sqt.pChild[i].searchThrough(rArea, listItems)
			}
		}
	}

	return listItems
}

func (sqt *StaticQuadTree[T]) Items(listItems List[T]) List[T] {

	for _, v := range sqt.pItems.AllItems() {
		// listItems = append(listItems, v.Second)
		listItems.PushBack(v.Second)
	}
	for i := 0; i < 4; i++ {
		if sqt.pChild[i] != nil {
			listItems = sqt.pChild[i].Items(listItems)
		}
	}
	return listItems
}

func (sqt *StaticQuadTree[T]) ItemsList() List[T] {
	listItems := NewList[T]()
	return sqt.Items(listItems)
}

func (sqt *StaticQuadTree[T]) Area() Rect {
	return sqt.rect
}

type StaticQuadTreeContainer[T any] struct {
	allItems   Map[item_key_type, QuadTreeItem[T]]
	root       StaticQuadTree[*QuadTreeItem[T]]
	QuadsCount int
	counter    int64
}

func NewStaticQuadTreeContainer[T any]() StaticQuadTreeContainer[T] {
	stQT := NewStaticQuadTree[*QuadTreeItem[T]](Rect{NewVector2f(0.0, 0.0), NewVector2f(100.0, 100.0)}, 0)

	return StaticQuadTreeContainer[T]{
		allItems: NewMap[item_key_type, QuadTreeItem[T]](),
		root:     stQT,
	}
}

func (stQtC *StaticQuadTreeContainer[T]) Resize(rect Rect) {
	stQtC.root.Resize(rect)
}

func (stQtc *StaticQuadTreeContainer[T]) Empty() bool {
	return stQtc.allItems.IsEmpty()
}

func (stQtc *StaticQuadTreeContainer[T]) Clear() {
	stQtc.root.clear()
	stQtc.allItems.Clear()
}

// The issue is that stQtc.allItems.PushBack is only adding first row elements for some reason???
func (stQtc *StaticQuadTreeContainer[T]) Insert(item T, itemsize Rect) {

	var newItem QuadTreeItem[T]
	newItem.item = item
	// stQtc.allItems.PushBack(-1, newItem)

	// stQtc.root.Insert(&item, itemsize)
	// stQtc.allItems.LastAddedElement().pItem = stQtc.root.Insert(stQtc.allItems.LastAddedElement(), itemsize)
	stQtc.counter += 1
	stQtc.allItems.Insert(stQtc.counter, newItem)
	// newItem.pItem = stQtc.root.Insert(&newItem, itemsize)

	tempItem := stQtc.allItems.LastAddedElement()
	tempItem.pItem = stQtc.root.InsertWithIndex(&tempItem, itemsize, stQtc.counter)

	stQtc.allItems.SetLatestElement(tempItem)
}

func (stQtc *StaticQuadTreeContainer[T]) InsertWithIndex(item T, itemsize Rect, index int64) {

	var newItem QuadTreeItem[T]
	newItem.item = item
	// stQtc.allItems.PushBack(-1, newItem)

	// stQtc.root.Insert(&item, itemsize)
	// stQtc.allItems.LastAddedElement().pItem = stQtc.root.Insert(stQtc.allItems.LastAddedElement(), itemsize)

	stQtc.allItems.Insert(index, newItem)
	// newItem.pItem = stQtc.root.Insert(&newItem, itemsize)

	// tempItem := stQtc.allItems.LastAddedElement()
	tempItem := stQtc.allItems.Get(index)
	tempItem.pItem = stQtc.root.InsertWithIndex(&tempItem, itemsize, index)

	// stQtc.allItems.SetLatestElement(tempItem)
	stQtc.allItems.Set(index, tempItem)

	LogF("Len: %v", index)
}

func (stQtc *StaticQuadTreeContainer[T]) Remove(item *QuadTreeItem[T]) {
	item.pItem.container.Erase(item.pItem.index)

	stQtc.allItems.Erase(item.pItem.index)
}

func (stQtc *StaticQuadTreeContainer[T]) Relocate(item *QuadTreeItem[T], itemsize Rect) {
	// item.pItem.container.Erase(item.pItem.index)
	// item.pItem = stQtc.root.InsertWithIndex(item, itemsize, item.pItem.index)
	stQtc.Remove(item)
	stQtc.Insert(item.item, itemsize)
}

func (stQtc *StaticQuadTreeContainer[T]) Search(rArea Rect) List[*QuadTreeItem[T]] {
	listItems := stQtc.root.Search(rArea)
	return listItems
}

func (stQtc *StaticQuadTreeContainer[T]) QuadsInViewCount() int {
	return stQtc.QuadsCount
}
