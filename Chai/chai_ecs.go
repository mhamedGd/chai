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

func NewEcsEngine() EcsEngine {
	return EcsEngine{
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
	return ent
}

func (e *EcsEngine) WriteToEntity(index int, ent EcsEntity) {
	e.entities[index] = &ent
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

func ReadComponent(e *EcsEngine, entity *EcsEntity, val Component) bool {
	storage := GetStorage(e, val)
	newVal, ok := storage.read(entity)
	if ok {
		val.ComponentSet(newVal)
	}
	return ok
}

func WriteComponent(e *EcsEngine, entity *EcsEntity, val interface{}) {
	storage := GetStorage(e, val)
	storage.write(entity, val)
}

func Each(engine *EcsEngine, val interface{}, f func(entity *EcsEntity, a interface{})) {
	storage := GetStorage(engine, val)
	for entity, a := range storage.list {
		f(entity, a)
	}
}

// If change anything in the entity then call WriteToEntity(index, new entity)
func EachAll(engine *EcsEngine, f func(entity *EcsEntity, entity_index int)) {
	// Terrible Solution, Try to connect engine.entities to engine.reg
	for index, ent := range engine.entities {
		f(ent, index)
	}
	// for _, storage := range engine.reg {
	// 	for entity := range storage.list {
	// 		f(entity)
	// 	}
	// }
}

var current_scene *Scene

type EcsEntity struct {
	id  Id
	Pos Vector2f
	Rot float32
}

type EcsSystem interface {
	Update(dt float32)
	GetEcsEngine() *EcsEngine
}

type EcsSystemImpl struct {
	EcsSystem
}

func (sys *EcsSystemImpl) GetEcsEngine() *EcsEngine {
	return &current_scene.Ecs_engine
}

type Scene struct {
	Ecs_engine     EcsEngine
	entities       []EcsEntity
	update_systems []EcsSystem
	render_systems []EcsSystem
	OnSceneStart   func()
}

func (scene *Scene) GetNumberOfEntities() int {
	return len(scene.entities)
}

func NewScene() Scene {
	return Scene{
		Ecs_engine:     NewEcsEngine(),
		entities:       make([]EcsEntity, 0),
		update_systems: make([]EcsSystem, 0),
		render_systems: make([]EcsSystem, 0),
	}
}

func ChangeScene(scene *Scene) {
	if current_scene != nil {
		current_scene.terminateScene()
	}
	current_scene = scene
	current_scene.OnSceneStart()
}

func (scene *Scene) terminateScene() {
	scene.entities = scene.entities[:0]
	scene.update_systems = scene.update_systems[:0]
	scene.render_systems = scene.render_systems[:0]
}

func (scene *Scene) NewEntity(pos Vector2f, rot float32) *EcsEntity {
	ent := scene.Ecs_engine.NewEntity()
	ent.Pos = pos
	ent.Rot = rot
	scene.entities = append(scene.entities, ent)
	return &ent
}

func (scene *Scene) GetLastEntity() *EcsEntity {
	return &scene.entities[len(scene.entities)-1]
}

func (scene *Scene) WriteComponentToLastEntity(component interface{}) {
	WriteComponent(&scene.Ecs_engine, &scene.entities[len(scene.entities)-1], component)
}

func (scene *Scene) NewUpdateSystem(sys EcsSystem) {
	scene.update_systems = append(scene.update_systems, sys)
}

func (scene *Scene) NewRenderSystem(sys EcsSystem) {
	scene.render_systems = append(scene.render_systems, sys)
}

func (scene *Scene) OnUpdate(dt float32) {
	for _, sys := range scene.update_systems {
		sys.Update(dt)
	}
}

func (scene *Scene) OnDraw() {
	for _, sys := range scene.render_systems {
		sys.Update(0.0)
	}
}

func GetCurrentScene() *Scene {
	return current_scene
}
