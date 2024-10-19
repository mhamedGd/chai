package customtypes

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
func (m *Map[A, B]) GetOk(key A) (B, bool) {
	e, ok := m.data[key]
	return e, ok
}

func (m *Map[A, B]) Has(key A) bool {
	_, ok := m.data[key]
	return ok
}

func (m *Map[A, B]) Set(key A, value B) {
	m.data[key] = value
}
