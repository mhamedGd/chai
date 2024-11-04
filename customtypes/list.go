package customtypes

type List[T any] struct {
	Data []T
}

func NewList[T any]() List[T] {
	return List[T]{
		Data: make([]T, 0),
	}
}
func NewListSized[T any](_size int) List[T] {
	return List[T]{
		Data: make([]T, _size),
	}
}

func ListFromSlice[T any](_slice []T) List[T] {
	return List[T]{
		Data: _slice,
	}
}

func (l *List[T]) PushBack(item T) {
	l.Data = append(l.Data, item)
}

func (l *List[T]) PushbackList(pl List[T]) {
	l.Data = append(l.Data, pl.Data...)
}

func (l *List[T]) PushBackArray(arr []T) {
	l.Data = append(l.Data, arr...)
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

func (l *List[T]) FindIf(_returnAction func(a T) bool) int {
	for i := 0; i < l.Count(); i++ {
		if _returnAction(l.Data[i]) {
			return i
		}
	}
	return -1
}

func (l *List[T]) AllItems() []T {
	return l.Data
}

func (l *List[T]) GetItemIndex(_item *T) int {
	temp := int(-1)
	l.Iterate(func(_index int, _item_it *T) {
		if _item_it == _item {
			temp = _index
			return
		}
	})

	return temp
}
