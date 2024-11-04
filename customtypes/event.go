package customtypes

import (
	"reflect"
)

type EventFunc[T any] func(...T)

type EventFunc1[A any] func(A)
type EventFunc2[A, B any] func(A, B)
type EventFunc3[A, B, C any] func(A, B, C)
type EventFunc4[A, B, C, D any] func(A, B, C, D)
type EventFunc5[A, B, C, D, E any] func(A, B, C, D, E)

// type ChaiEvent[T any] struct {
// 	m_Listeners []EventFunc[T]
// }

type ChaiEvent1[A any] struct {
	m_Listeners List[EventFunc1[A]]
}
type ChaiEvent2[A, B any] struct {
	m_Listeners List[EventFunc2[A, B]]
}
type ChaiEvent3[A, B, C any] struct {
	m_Listeners List[EventFunc3[A, B, C]]
}
type ChaiEvent4[A, B, C, D any] struct {
	m_Listeners List[EventFunc4[A, B, C, D]]
}
type ChaiEvent5[A, B, C, D, E any] struct {
	m_Listeners List[EventFunc5[A, B, C, D, E]]
}

// func NewChaiEvent[T any]() ChaiEvent[T] {
// 	var c ChaiEvent[T]
// 	c.init()
// 	return c
// }

// func (e *ChaiEvent[T]) init() {
// 	e.m_Listeners = make([]EventFunc[T], 0)
// }

func NewChaiEvent1[A any]() ChaiEvent1[A] {
	return ChaiEvent1[A]{
		m_Listeners: NewList[EventFunc1[A]](),
	}
}
func NewChaiEvent2[A, B any]() ChaiEvent2[A, B] {
	return ChaiEvent2[A, B]{
		m_Listeners: NewList[EventFunc2[A, B]](),
	}
}
func NewChaiEvent3[A, B, C any]() ChaiEvent3[A, B, C] {
	return ChaiEvent3[A, B, C]{
		m_Listeners: NewList[EventFunc3[A, B, C]](),
	}
}
func NewChaiEvent4[A, B, C, D any]() ChaiEvent4[A, B, C, D] {
	return ChaiEvent4[A, B, C, D]{
		m_Listeners: NewList[EventFunc4[A, B, C, D]](),
	}
}
func NewChaiEvent5[A, B, C, D, E any]() ChaiEvent5[A, B, C, D, E] {
	return ChaiEvent5[A, B, C, D, E]{
		m_Listeners: NewList[EventFunc5[A, B, C, D, E]](),
	}
}

func (e *ChaiEvent1[A]) init() {
	e.m_Listeners = NewList[EventFunc1[A]]()
}
func (e *ChaiEvent2[A, B]) init() {
	e.m_Listeners = NewList[EventFunc2[A, B]]()
}
func (e *ChaiEvent3[A, B, C]) init() {
	e.m_Listeners = NewList[EventFunc3[A, B, C]]()
}
func (e *ChaiEvent4[A, B, C, D]) init() {
	e.m_Listeners = NewList[EventFunc4[A, B, C, D]]()
}
func (e *ChaiEvent5[A, B, C, D, E]) init() {
	e.m_Listeners = NewList[EventFunc5[A, B, C, D, E]]()
}

// func (e *ChaiEvent[T]) AddListener(ef EventFunc[T]) {
// 	e.m_Listeners = append(e.m_Listeners, ef)
// }

func (ev *ChaiEvent1[A]) AddListener(_eventFunc EventFunc1[A]) {
	ev.m_Listeners.PushBack(_eventFunc)
}
func (ev *ChaiEvent2[A, B]) AddListener(_eventFunc EventFunc2[A, B]) {
	ev.m_Listeners.PushBack(_eventFunc)
}
func (ev *ChaiEvent3[A, B, C]) AddListener(_eventFunc EventFunc3[A, B, C]) {
	ev.m_Listeners.PushBack(_eventFunc)
}
func (ev *ChaiEvent4[A, B, C, D]) AddListener(_eventFunc EventFunc4[A, B, C, D]) {
	ev.m_Listeners.PushBack(_eventFunc)
}
func (ev *ChaiEvent5[A, B, C, D, E]) AddListener(_eventFunc EventFunc5[A, B, C, D, E]) {
	ev.m_Listeners.PushBack(_eventFunc)
}

// func (e *ChaiEvent[T]) RemoveListener(ef EventFunc[T]) {
// 	for i, fn := range e.m_Listeners {
// 		f1 := reflect.ValueOf(fn)
// 		f2 := reflect.ValueOf(ef)
// 		if f1 == f2 {
// 			e.m_Listeners = append(e.m_Listeners[:i], e.m_Listeners[i+1:]...)
// 		}
// 	}
// }

func (ev *ChaiEvent1[A]) RemoveListener(_eventFunc EventFunc1[A]) {
	for i, fn := range ev.m_Listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(_eventFunc)
		if f1 == f2 {
			ev.m_Listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent2[A, B]) RemoveListener(_eventFunc EventFunc2[A, B]) {
	for i, fn := range ev.m_Listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(_eventFunc)
		if f1 == f2 {
			ev.m_Listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent3[A, B, C]) RemoveListener(_eventFunc EventFunc3[A, B, C]) {
	for i, fn := range ev.m_Listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(_eventFunc)
		if f1 == f2 {
			ev.m_Listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent4[A, B, C, D]) RemoveListener(_eventFunc EventFunc4[A, B, C, D]) {
	for i, fn := range ev.m_Listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(_eventFunc)
		if f1 == f2 {
			ev.m_Listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent5[A, B, C, D, E]) RemoveListener(_eventFunc EventFunc5[A, B, C, D, E]) {
	for i, fn := range ev.m_Listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(_eventFunc)
		if f1 == f2 {
			ev.m_Listeners.Erase(i)
		}
	}
}

// func (e *ChaiEvent[T]) Invoke(x ...T) {
// 	for _, f := range e.m_Listeners {
// 		//fmt.Println(index)
// 		f(x...)
// 	}
// }

func (e *ChaiEvent1[A]) Invoke(_a A) {
	for _, f := range e.m_Listeners.Data {
		//fmt.Println(index)
		f(_a)
	}
}
func (e *ChaiEvent2[A, B]) Invoke(_a A, _b B) {
	for _, f := range e.m_Listeners.Data {
		//fmt.Println(index)
		f(_a, _b)
	}
}
func (e *ChaiEvent3[A, B, C]) Invoke(_a A, _b B, _c C) {
	for _, f := range e.m_Listeners.Data {
		//fmt.Println(index)
		f(_a, _b, _c)
	}
}
func (e *ChaiEvent4[A, B, C, D]) Invoke(_a A, _b B, _c C, _d D) {
	for _, f := range e.m_Listeners.Data {
		//fmt.Println(index)
		f(_a, _b, _c, _d)
	}
}
func (ev *ChaiEvent5[A, B, C, D, E]) Invoke(_a A, _b B, _c C, _d D, _e E) {
	for _, f := range ev.m_Listeners.Data {
		//fmt.Println(index)
		f(_a, _b, _c, _d, _e)
	}
}
