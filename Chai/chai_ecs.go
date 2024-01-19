package chai

import (
	"reflect"
)

type Id uint32

type BasicStorage struct {
	list map[*EcsEntity]interface{}
}

// Copy Paste for new types
// type TYPE struct {
// }
// func (t *TYPE) ComponentSet(val interface{}) { *t = val.(TYPE) }

type Component interface {
	ComponentSet(interface{})
}

func NewBasicStorage() *BasicStorage {
	return &BasicStorage{
		list: make(map[*EcsEntity]interface{}),
	}
}

func (s *BasicStorage) read(entity *EcsEntity) (interface{}, bool) {
	val, ok := s.list[entity]

	return val, ok
}

func (s *BasicStorage) write(entity *EcsEntity, val interface{}) {
	s.list[entity] = val
}

type EcsEngine struct {
	reg       map[string]*BasicStorage
	entities  []*EcsEntity
	idCounter Id
}

func NewEcsEngine() *EcsEngine {
	return &EcsEngine{
		reg:       make(map[string]*BasicStorage),
		entities:  make([]*EcsEntity, 0),
		idCounter: 0,
	}
}

func (e *EcsEngine) NewEntity() EcsEntity {
	id := e.idCounter
	e.idCounter++
	ent := EcsEntity{id: id}
	e.entities = append(e.entities, &ent)
	return *e.entities[len(e.entities)-1]
}

func name(t interface{}) string {
	name := reflect.TypeOf(t).String()
	if name[0] == '*' {
		return name[1:]
	}

	return name
}

func GetStorage(e *EcsEngine, t interface{}) *BasicStorage {
	name := name(t)

	storage, ok := e.reg[name]
	if !ok {
		e.reg[name] = NewBasicStorage()
		storage, _ = e.reg[name]
	}
	return storage
}

func Read(e *EcsEngine, entity *EcsEntity, val Component) bool {
	storage := GetStorage(e, val)
	newVal, ok := storage.read(entity)
	if ok {
		val.ComponentSet(newVal)
	}
	return ok
}

func Write(e *EcsEngine, entity *EcsEntity, val interface{}) {
	storage := GetStorage(e, val)
	storage.write(entity, val)
}

func Each(engine *EcsEngine, val interface{}, f func(entity *EcsEntity, a interface{})) {
	storage := GetStorage(engine, val)
	for entity, a := range storage.list {
		f(entity, a)
		LogF("%v", *entity)

	}
}

func EachAll(engine *EcsEngine, f func(entity *EcsEntity)) {
	// Terrible Solution, Try to connect engine.entities to engine.reg
	// for _, ent := range engine.entities {
	// 	f(ent)

	// }
	for _, storage := range engine.reg {
		for entity, _ := range storage.list {
			f(entity)
		}
	}
}

type EcsEntity struct {
	id  Id
	Pos Vector2f
	Rot float32
}

type EcsSystem interface {
	Update(dt float32) func(float32)
}

type Scene struct {
	systems []EcsSystem
}
