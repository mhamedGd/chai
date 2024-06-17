package chai

type List[T any] struct {
	Data []T
}

func NewList[T any]() List[T] {
	return List[T]{
		Data: make([]T, 0),
	}
}

func (l *List[T]) PushBack(item T) {
	l.Data = append(l.Data, item)
}

func (l *List[T]) Erase(index int) {
	if index < 0 || index >= l.Count() {
		return
	}
	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

func (l *List[T]) EraseByPointer(p *T) {

}

func (l *List[T]) Clear() {
	l.Data = l.Data[:0]
}

func (l *List[T]) Count() int {
	return len(l.Data)
}

func (l *List[T]) IsEmpty() bool {
	return len(l.Data) == 0
}

func (l *List[T]) Front() *T {
	if !l.IsEmpty() {
		return &l.Data[0]
	}

	return nil
}

func (l *List[T]) Back() *T {
	if !l.IsEmpty() {
		return &l.Data[len(l.Data)-1]
	}

	return nil
}

func (l *List[T]) Iterate(action func(_index int, _item *T)) {
	for i := 0; i < l.Count(); i++ {
		action(i, &l.Data[i])
	}
}

func (l *List[T]) FindIf(return_action func(a T) bool) int {
	for i := 0; i < l.Count(); i++ {
		if return_action(l.Data[i]) {
			return i
		}
	}
	return -1
}

func (l *List[T]) AllItems() []T {
	return l.Data
}

func (l *List[T]) GetItemIndex(item *T) int {
	temp := int(-1)
	l.Iterate(func(_index int, _item *T) {
		if _item == item {
			temp = _index
			return
		}
	})

	return temp
}

// -------------------------------------------------------------------------------

type Map[A comparable, B any] struct {
	data          map[A]B
	latestAddedId A
}

func NewMap[A comparable, B any]() Map[A, B] {
	return Map[A, B]{
		data: make(map[A]B),
	}
}

func (m *Map[A, B]) Insert(id A, item B) {
	m.data[id] = item
	m.latestAddedId = id
}

func (m *Map[A, B]) Erase(id A) {
	delete(m.data, id)
}

func (m *Map[A, B]) IsEmpty() bool {
	return len(m.data) == 0
}

func (m *Map[A, B]) Clear() {
	m.data = make(map[A]B)
}

func (m *Map[A, B]) Count() int {
	return len(m.data)
}

func (m *Map[A, B]) AllItems() map[A]B {
	return m.data
}

func (m *Map[A, B]) LastAddedElement() B {
	return m.data[m.latestAddedId]
}

func (m *Map[A, B]) SetLatestElement(item B) {
	m.data[m.latestAddedId] = item
}

func (m *Map[A, B]) Get(key A) B {
	return m.data[key]
}

func (m *Map[A, B]) Set(key A, value B) {
	m.data[key] = value
}
