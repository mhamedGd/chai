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
// 	listeners []EventFunc[T]
// }

type ChaiEvent1[A any] struct {
	listeners List[EventFunc1[A]]
}
type ChaiEvent2[A, B any] struct {
	listeners List[EventFunc2[A, B]]
}
type ChaiEvent3[A, B, C any] struct {
	listeners List[EventFunc3[A, B, C]]
}
type ChaiEvent4[A, B, C, D any] struct {
	listeners List[EventFunc4[A, B, C, D]]
}
type ChaiEvent5[A, B, C, D, E any] struct {
	listeners List[EventFunc5[A, B, C, D, E]]
}

// func NewChaiEvent[T any]() ChaiEvent[T] {
// 	var c ChaiEvent[T]
// 	c.init()
// 	return c
// }

// func (e *ChaiEvent[T]) init() {
// 	e.listeners = make([]EventFunc[T], 0)
// }

func NewChaiEvent1[A any]() ChaiEvent1[A] {
	return ChaiEvent1[A]{
		listeners: NewList[EventFunc1[A]](),
	}
}
func NewChaiEvent2[A, B any]() ChaiEvent2[A, B] {
	return ChaiEvent2[A, B]{
		listeners: NewList[EventFunc2[A, B]](),
	}
}
func NewChaiEvent3[A, B, C any]() ChaiEvent3[A, B, C] {
	return ChaiEvent3[A, B, C]{
		listeners: NewList[EventFunc3[A, B, C]](),
	}
}
func NewChaiEvent4[A, B, C, D any]() ChaiEvent4[A, B, C, D] {
	return ChaiEvent4[A, B, C, D]{
		listeners: NewList[EventFunc4[A, B, C, D]](),
	}
}
func NewChaiEvent5[A, B, C, D, E any]() ChaiEvent5[A, B, C, D, E] {
	return ChaiEvent5[A, B, C, D, E]{
		listeners: NewList[EventFunc5[A, B, C, D, E]](),
	}
}

func (e *ChaiEvent1[A]) init() {
	e.listeners = NewList[EventFunc1[A]]()
}
func (e *ChaiEvent2[A, B]) init() {
	e.listeners = NewList[EventFunc2[A, B]]()
}
func (e *ChaiEvent3[A, B, C]) init() {
	e.listeners = NewList[EventFunc3[A, B, C]]()
}
func (e *ChaiEvent4[A, B, C, D]) init() {
	e.listeners = NewList[EventFunc4[A, B, C, D]]()
}
func (e *ChaiEvent5[A, B, C, D, E]) init() {
	e.listeners = NewList[EventFunc5[A, B, C, D, E]]()
}

// func (e *ChaiEvent[T]) AddListener(ef EventFunc[T]) {
// 	e.listeners = append(e.listeners, ef)
// }

func (ev *ChaiEvent1[A]) AddListener(ef EventFunc1[A]) {
	ev.listeners.PushBack(ef)
}
func (ev *ChaiEvent2[A, B]) AddListener(ef EventFunc2[A, B]) {
	ev.listeners.PushBack(ef)
}
func (ev *ChaiEvent3[A, B, C]) AddListener(ef EventFunc3[A, B, C]) {
	ev.listeners.PushBack(ef)
}
func (ev *ChaiEvent4[A, B, C, D]) AddListener(ef EventFunc4[A, B, C, D]) {
	ev.listeners.PushBack(ef)
}
func (ev *ChaiEvent5[A, B, C, D, E]) AddListener(ef EventFunc5[A, B, C, D, E]) {
	ev.listeners.PushBack(ef)
}

// func (e *ChaiEvent[T]) RemoveListener(ef EventFunc[T]) {
// 	for i, fn := range e.listeners {
// 		f1 := reflect.ValueOf(fn)
// 		f2 := reflect.ValueOf(ef)
// 		if f1 == f2 {
// 			e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
// 		}
// 	}
// }

func (ev *ChaiEvent1[A]) RemoveListener(ef EventFunc1[A]) {
	for i, fn := range ev.listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(ef)
		if f1 == f2 {
			ev.listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent2[A, B]) RemoveListener(ef EventFunc2[A, B]) {
	for i, fn := range ev.listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(ef)
		if f1 == f2 {
			ev.listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent3[A, B, C]) RemoveListener(ef EventFunc3[A, B, C]) {
	for i, fn := range ev.listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(ef)
		if f1 == f2 {
			ev.listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent4[A, B, C, D]) RemoveListener(ef EventFunc4[A, B, C, D]) {
	for i, fn := range ev.listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(ef)
		if f1 == f2 {
			ev.listeners.Erase(i)
		}
	}
}
func (ev *ChaiEvent5[A, B, C, D, E]) RemoveListener(ef EventFunc5[A, B, C, D, E]) {
	for i, fn := range ev.listeners.Data {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(ef)
		if f1 == f2 {
			ev.listeners.Erase(i)
		}
	}
}

// func (e *ChaiEvent[T]) Invoke(x ...T) {
// 	for _, f := range e.listeners {
// 		//fmt.Println(index)
// 		f(x...)
// 	}
// }

func (e *ChaiEvent1[A]) Invoke(a A) {
	for _, f := range e.listeners.Data {
		//fmt.Println(index)
		f(a)
	}
}
func (e *ChaiEvent2[A, B]) Invoke(a A, b B) {
	for _, f := range e.listeners.Data {
		//fmt.Println(index)
		f(a, b)
	}
}
func (e *ChaiEvent3[A, B, C]) Invoke(a A, b B, c C) {
	for _, f := range e.listeners.Data {
		//fmt.Println(index)
		f(a, b, c)
	}
}
func (e *ChaiEvent4[A, B, C, D]) Invoke(a A, b B, c C, d D) {
	for _, f := range e.listeners.Data {
		//fmt.Println(index)
		f(a, b, c, d)
	}
}
func (ev *ChaiEvent5[A, B, C, D, E]) Invoke(a A, b B, c C, d D, e E) {
	for _, f := range ev.listeners.Data {
		//fmt.Println(index)
		f(a, b, c, d, e)
	}
}
