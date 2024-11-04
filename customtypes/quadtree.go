package customtypes

import (
	"fmt"

	. "github.com/mhamedGd/chai/math"
)

type Pair[A, B any] struct {
	First  A
	Second B
}

const MAX_QUADTREE_DEPTH = 8

type item_key_type = int64

type QuadTreeItemLocation[T any] struct {

	// m_Container *[]Pair[Rect, T]
	m_Container *Map[item_key_type, Pair[Rect, T]]
	m_Index     int64
}

type QuadTreeItem[T any] struct {
	m_Item  T
	m_PItem QuadTreeItemLocation[*QuadTreeItem[T]]
}

func (qtt *QuadTreeItem[T]) GetItem() T {
	return qtt.m_Item
}

type StaticQuadTree[T any] struct {
	m_Depth   int
	m_Rect    Rect
	m_RChild  [4]Rect
	m_PChild  [4]*StaticQuadTree[T]
	m_PItems  List[Pair[Rect, T]]
	m_Counter int64
}

func NewStaticQuadTree[T any](_size Rect, _nDepth int) StaticQuadTree[T] {
	// l := list.New()
	sqt := StaticQuadTree[T]{}
	sqt.m_Depth = _nDepth
	sqt.Resize(_size)

	return sqt
}

func (sqt *StaticQuadTree[T]) Resize(_rArea Rect) {
	sqt.clear()

	sqt.m_Rect = _rArea
	vChildSize := sqt.m_Rect.Size.Scale(0.5)

	sqt.m_RChild = [4]Rect{
		{sqt.m_Rect.Position, vChildSize},
		{NewVector2f(sqt.m_Rect.Position.X+vChildSize.X, sqt.m_Rect.Position.Y), vChildSize},
		{NewVector2f(sqt.m_Rect.Position.X, sqt.m_Rect.Position.Y+vChildSize.Y), vChildSize},
		{sqt.m_Rect.Position.Add(vChildSize), vChildSize},
	}
}

func (sqt *StaticQuadTree[T]) clear() {
	sqt.m_PItems.Clear()

	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			sqt.m_PChild[i].clear()
		}
		sqt.m_PChild[i] = nil
	}
}

func (sqt *StaticQuadTree[T]) size() int {
	nCount := sqt.m_PItems.Count()
	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			nCount += len(sqt.m_PChild)
		}
	}

	return nCount
}

func (sqt *StaticQuadTree[T]) Insert(_item T, _itemsize Rect) {

	for i := 0; i < 4; i++ {
		if sqt.m_RChild[i].ContainsRect(_itemsize) {
			if sqt.m_Depth+1 < MAX_QUADTREE_DEPTH {
				if sqt.m_PChild[i] == nil {
					_tree := NewStaticQuadTree[T](sqt.m_RChild[i], sqt.m_Depth+1)
					sqt.m_PChild[i] = &_tree
				}

				sqt.m_PChild[i].Insert(_item, _itemsize)
				return
			}
		}
	}

	// sqt.m_PItems = append(sqt.m_PItems, Pair[Rect, T]{itemsize, m_Item})
	sqt.m_PItems.PushBack(Pair[Rect, T]{_itemsize, _item})
}

func (sqt *StaticQuadTree[T]) Search(_rArea Rect) List[T] {
	listItems := NewList[T]()

	return sqt.searchThrough(_rArea, listItems)
}

func (sqt *StaticQuadTree[T]) searchThrough(_rArea Rect, _listItems List[T]) List[T] {
	for _, v := range sqt.m_PItems.Data {
		if _rArea.OverlapsRect(v.First) {
			_listItems.PushBack(v.Second)
		}
	}

	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			if _rArea.ContainsRect(sqt.m_RChild[i]) {
				_listItems = sqt.m_PChild[i].Items(_listItems)

			} else if sqt.m_RChild[i].OverlapsRect(_rArea) {
				_listItems = sqt.m_PChild[i].searchThrough(_rArea, _listItems)
			}
		}
	}

	return _listItems
}

func (sqt *StaticQuadTree[T]) Items(_listItems List[T]) List[T] {

	for _, v := range sqt.m_PItems.AllItems() {
		// listItems = append(listItems, v.Second)
		_listItems.PushBack(v.Second)
	}
	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			_listItems = sqt.m_PChild[i].Items(_listItems)
		}
	}
	return _listItems
}

func (sqt *StaticQuadTree[T]) ItemsList() List[T] {
	listItems := NewList[T]()
	return sqt.Items(listItems)
}

func (sqt *StaticQuadTree[T]) Area() Rect {
	return sqt.m_Rect
}

type StaticQuadTreeContainer[T any] struct {
	allItems   List[T]
	root       StaticQuadTree[*T]
	QuadsCount int
	m_Counter  int64
}

func NewStaticQuadTreeContainer[T any]() StaticQuadTreeContainer[T] {
	stQT := NewStaticQuadTree[*T](Rect{NewVector2f(0.0, 0.0), NewVector2f(100.0, 100.0)}, 0)

	return StaticQuadTreeContainer[T]{
		allItems: NewList[T](),
		root:     stQT,
	}
}

func (stQtC *StaticQuadTreeContainer[T]) Resize(_rect Rect) {
	stQtC.root.Resize(_rect)
}

func (stQtc *StaticQuadTreeContainer[T]) Empty() bool {
	return stQtc.allItems.IsEmpty()
}

func (stQtc *StaticQuadTreeContainer[T]) Clear() {
	stQtc.root.clear()
	stQtc.allItems.Clear()
}

// The issue is that stQtc.allItems.PushBack is only adding first row elements for some reason???
func (stQtc *StaticQuadTreeContainer[T]) Insert(_item T, _itemsize Rect) {

	stQtc.allItems.PushBack(_item)
	// newItem.m_PItem = stQtc.root.Insert(&newItem, itemsize)

	stQtc.root.Insert(&_item, _itemsize)
}

func (stQtc *StaticQuadTreeContainer[T]) Search(_rArea Rect) List[*T] {
	listItems := stQtc.root.Search(_rArea)
	return listItems
}

func (stQtc *StaticQuadTreeContainer[T]) QuadsInViewCount() int {
	return stQtc.QuadsCount
}

///////////////////////////////////////////////////
//////////////// DYNAMIC QUADTREE ////////////////

type DynamicQuadTree[T any] struct {
	m_Depth   int
	m_Rect    Rect
	m_RChild  [4]Rect
	m_PChild  [4]*DynamicQuadTree[T]
	m_PItems  Map[item_key_type, Pair[Rect, T]]
	m_Counter int64
}

func NewDynamicQuadTree[T any](_size Rect, _nDepth int) DynamicQuadTree[T] {
	// l := list.New()
	sqt := DynamicQuadTree[T]{}
	sqt.m_Depth = _nDepth
	sqt.Resize(_size)

	return sqt
}

func (sqt *DynamicQuadTree[T]) Resize(_rArea Rect) {
	sqt.clear()

	sqt.m_Rect = _rArea
	vChildSize := sqt.m_Rect.Size.Scale(0.5)

	sqt.m_RChild = [4]Rect{
		{sqt.m_Rect.Position, vChildSize},
		{NewVector2f(sqt.m_Rect.Position.X+vChildSize.X, sqt.m_Rect.Position.Y), vChildSize},
		{NewVector2f(sqt.m_Rect.Position.X, sqt.m_Rect.Position.Y+vChildSize.Y), vChildSize},
		{sqt.m_Rect.Position.Add(vChildSize), vChildSize},
	}
}

func (sqt *DynamicQuadTree[T]) clear() {
	sqt.m_PItems.Clear()

	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			sqt.m_PChild[i].clear()
		}
		sqt.m_PChild[i] = nil
	}
}

func (sqt *DynamicQuadTree[T]) size() int {
	nCount := sqt.m_PItems.Count()
	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			nCount += len(sqt.m_PChild)
		}
	}

	return nCount
}

func (sqt *DynamicQuadTree[T]) Insert(_item T, _itemsize Rect) QuadTreeItemLocation[T] {

	for i := 0; i < 4; i++ {
		if sqt.m_RChild[i].ContainsRect(_itemsize) {
			if sqt.m_Depth+1 < MAX_QUADTREE_DEPTH {
				if sqt.m_PChild[i] == nil {
					_tree := NewDynamicQuadTree[T](sqt.m_RChild[i], sqt.m_Depth+1)
					sqt.m_PChild[i] = &_tree
				}

				return sqt.m_PChild[i].Insert(_item, _itemsize)
			}
		}
	}

	// sqt.m_PItems = append(sqt.m_PItems, Pair[Rect, T]{itemsize, m_Item})
	sqt.m_PItems.Insert(sqt.m_Counter, Pair[Rect, T]{_itemsize, _item})
	sqt.m_Counter += 1
	return QuadTreeItemLocation[T]{
		m_Container: &sqt.m_PItems,
		m_Index:     sqt.m_Counter,
	}
}
func (sqt *DynamicQuadTree[T]) InsertWithIndex(_item T, _itemsize Rect, _index int64) QuadTreeItemLocation[T] {
	for i := 0; i < 4; i++ {
		if sqt.m_RChild[i].ContainsRect(_itemsize) {
			if sqt.m_Depth+1 < MAX_QUADTREE_DEPTH {
				if sqt.m_PChild[i] == nil {
					_tree := NewDynamicQuadTree[T](sqt.m_RChild[i], sqt.m_Depth+1)
					sqt.m_PChild[i] = &_tree
				}

				return sqt.m_PChild[i].InsertWithIndex(_item, _itemsize, _index)
			}
		}
	}

	// sqt.m_PItems = append(sqt.m_PItems, Pair[Rect, T]{itemsize, m_Item})
	sqt.m_PItems.Insert(_index, Pair[Rect, T]{_itemsize, _item})
	return QuadTreeItemLocation[T]{
		m_Container: &sqt.m_PItems,
		m_Index:     _index,
	}
}

// func (sqt *StaticQuadTree[T]) Remove(m_Item T) bool {
// 	it := sqt.m_PItems.FindIf(func(a Pair[Rect, T]) bool {
// 		&m_Item == a.Second
// 	})
// }

func (sqt *DynamicQuadTree[T]) Search(_rArea Rect) List[T] {
	listItems := NewList[T]()

	return sqt.searchThrough(_rArea, listItems)
}

func (sqt *DynamicQuadTree[T]) searchThrough(_rArea Rect, _listItems List[T]) List[T] {
	for _, v := range sqt.m_PItems.AllItems() {
		if _rArea.OverlapsRect(v.First) {
			// listItems = append(listItems, v.Second)
			_listItems.PushBack(v.Second)
		}
	}

	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			if _rArea.ContainsRect(sqt.m_RChild[i]) {
				_listItems = sqt.m_PChild[i].Items(_listItems)

			} else if sqt.m_RChild[i].OverlapsRect(_rArea) {
				_listItems = sqt.m_PChild[i].searchThrough(_rArea, _listItems)
			}
		}
	}

	return _listItems
}

func (sqt *DynamicQuadTree[T]) Items(_listItems List[T]) List[T] {

	for _, v := range sqt.m_PItems.AllItems() {
		// listItems = append(listItems, v.Second)
		_listItems.PushBack(v.Second)
	}
	for i := 0; i < 4; i++ {
		if sqt.m_PChild[i] != nil {
			_listItems = sqt.m_PChild[i].Items(_listItems)
		}
	}
	return _listItems
}

func (sqt *DynamicQuadTree[T]) ItemsList() List[T] {
	listItems := NewList[T]()
	return sqt.Items(listItems)
}

func (sqt *DynamicQuadTree[T]) Area() Rect {
	return sqt.m_Rect
}

type DynamicQuadTreeContainer[T any] struct {
	allItems   Map[item_key_type, QuadTreeItem[T]]
	root       DynamicQuadTree[*QuadTreeItem[T]]
	QuadsCount int
	m_Counter  int64
}

func (dqtc *DynamicQuadTreeContainer[T]) AllItems() *Map[item_key_type, QuadTreeItem[T]] {
	return &dqtc.allItems
}

func NewDynamicQuadTreeContainer[T any]() DynamicQuadTreeContainer[T] {
	stQT := NewDynamicQuadTree[*QuadTreeItem[T]](Rect{NewVector2f(0.0, 0.0), NewVector2f(100.0, 100.0)}, 0)

	return DynamicQuadTreeContainer[T]{
		allItems: NewMap[item_key_type, QuadTreeItem[T]](),
		root:     stQT,
	}
}

func (stQtC *DynamicQuadTreeContainer[T]) Resize(_rect Rect) {
	stQtC.root.Resize(_rect)
}

func (stQtc *DynamicQuadTreeContainer[T]) Empty() bool {
	return stQtc.allItems.IsEmpty()
}

func (stQtc *DynamicQuadTreeContainer[T]) Clear() {
	stQtc.root.clear()
	stQtc.allItems.Clear()
}

// The issue is that stQtc.allItems.PushBack is only adding first row elements for some reason???
func (stQtc *DynamicQuadTreeContainer[T]) Insert(_item T, _itemSize Rect) int64 {

	var newItem QuadTreeItem[T]
	newItem.m_Item = _item
	// stQtc.allItems.PushBack(-1, newItem)

	// stQtc.root.Insert(&m_Item, itemsize)
	// stQtc.allItems.LastAddedElement().m_PItem = stQtc.root.Insert(stQtc.allItems.LastAddedElement(), itemsize)
	stQtc.m_Counter += 1
	stQtc.allItems.Insert(stQtc.m_Counter, newItem)
	// newItem.m_PItem = stQtc.root.Insert(&newItem, itemsize)

	tempItem := stQtc.allItems.LastAddedElement()
	tempItem.m_PItem = stQtc.root.InsertWithIndex(&tempItem, _itemSize, stQtc.m_Counter)

	stQtc.allItems.SetLatestElement(tempItem)

	return stQtc.m_Counter
}

func (stQtc *DynamicQuadTreeContainer[T]) InsertWithIndex(_item T, _itemSize Rect, _index int64) {

	var newItem QuadTreeItem[T]
	newItem.m_Item = _item
	// stQtc.allItems.PushBack(-1, newItem)

	// stQtc.root.Insert(&m_Item, itemsize)
	// stQtc.allItems.LastAddedElement().m_PItem = stQtc.root.Insert(stQtc.allItems.LastAddedElement(), itemsize)

	stQtc.allItems.Insert(_index, newItem)
	// newItem.m_PItem = stQtc.root.Insert(&newItem, itemsize)

	// tempItem := stQtc.allItems.LastAddedElement()
	tempItem := stQtc.allItems.Get(_index)
	tempItem.m_PItem = stQtc.root.InsertWithIndex(&tempItem, _itemSize, _index)

	// stQtc.allItems.SetLatestElement(tempItem)
	stQtc.allItems.Set(_index, tempItem)

}

func (stQtc *DynamicQuadTreeContainer[T]) Remove(_item *QuadTreeItem[T]) {
	_item.m_PItem.m_Container.Erase(_item.m_PItem.m_Index)

	stQtc.allItems.Erase(_item.m_PItem.m_Index)
}

func (stQtc *DynamicQuadTreeContainer[T]) RemoveWithIndex(_itemIndex int64) {

	searchedItem := stQtc.searchIndex(_itemIndex)
	if searchedItem == nil {
		fmt.Printf("Didn't find object\n")
		return
	}

	searchedItem.m_PItem.m_Container.Erase(_itemIndex)
	stQtc.allItems.Erase(_itemIndex)
}

func (stQtc *DynamicQuadTreeContainer[T]) Relocate(_item *QuadTreeItem[T], _itemSize Rect) {
	// m_Item.m_PItem.m_Container.Erase(m_Item.m_PItem.m_Index)
	// m_Item.m_PItem = stQtc.root.InsertWithIndex(m_Item, itemsize, m_Item.m_PItem.m_Index)
	stQtc.Remove(_item)
	stQtc.InsertWithIndex(_item.m_Item, _itemSize, _item.m_PItem.m_Index)
}

func (stQtc *DynamicQuadTreeContainer[T]) Search(_rArea Rect) List[*QuadTreeItem[T]] {
	listItems := stQtc.root.Search(_rArea)
	return listItems
}

func (dQtc *DynamicQuadTreeContainer[T]) searchIndex(_index int64) *QuadTreeItem[T] {
	for _, v := range dQtc.allItems.AllItems() {
		if v.m_PItem.m_Index == _index {
			return &v
		}
	}
	return nil
}

func (stQtc *DynamicQuadTreeContainer[T]) QuadsInViewCount() int {
	return stQtc.QuadsCount
}
