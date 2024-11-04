package customtypes

type Map[A comparable, B any] struct {
	m_Data          map[A]B
	m_LatestAddedId A
}

func NewMap[A comparable, B any]() Map[A, B] {
	return Map[A, B]{
		m_Data: make(map[A]B),
	}
}

func (m *Map[A, B]) Insert(_id A, _item B) {
	m.m_Data[_id] = _item
	m.m_LatestAddedId = _id
}

func (m *Map[A, B]) Erase(_id A) {
	delete(m.m_Data, _id)
}

func (m *Map[A, B]) IsEmpty() bool {
	return len(m.m_Data) == 0
}

func (m *Map[A, B]) Clear() {
	m.m_Data = make(map[A]B)
}

func (m *Map[A, B]) Count() int {
	return len(m.m_Data)
}

func (m *Map[A, B]) AllItems() map[A]B {
	return m.m_Data
}

func (m *Map[A, B]) LastAddedElement() B {
	return m.m_Data[m.m_LatestAddedId]
}

func (m *Map[A, B]) SetLatestElement(_item B) {
	m.m_Data[m.m_LatestAddedId] = _item
}

func (m *Map[A, B]) Get(_key A) B {
	return m.m_Data[_key]
}
func (m *Map[A, B]) GetOk(_key A) (B, bool) {
	e, ok := m.m_Data[_key]
	return e, ok
}

func (m *Map[A, B]) Has(_key A) bool {
	_, ok := m.m_Data[_key]
	return ok
}

func (m *Map[A, B]) Set(_key A, _value B) {
	m.m_Data[_key] = _value
}
