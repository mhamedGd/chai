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

func (s *BasicStorage) delete(entity *EcsEntity, val interface{}) {
	delete(s.list, entity)
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

func DeleteComponent(e *EcsEngine, entity *EcsEntity, val interface{}) {
	storage := GetStorage(e, val)
	storage.delete(entity, val)
}

func EachEntity(val interface{}, f func(entity *EcsEntity, a interface{})) {
	storage := GetStorage(&current_scene.Ecs_engine, val)
	for entity, a := range storage.list {
		f(entity, a)
	}
}

// If change anything in the entity then call WriteToEntity(index, new entity)
func EachEntityAll(engine *EcsEngine, f func(entity *EcsEntity, entity_index int)) {
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
	id         Id
	Pos        Vector2f
	Rot        float32
	Dimensions Vector2f
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
	Background     RGBA8
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

	for b := GetPhysicsWorld().box2dWorld.GetBodyList(); b != nil; b = b.GetNext() {
		GetPhysicsWorld().box2dWorld.DestroyBody(b)
	}
}

func (scene *Scene) NewEntity(pos Vector2f, dim Vector2f, rot float32) *EcsEntity {
	ent := scene.Ecs_engine.NewEntity()
	ent.Pos = pos
	ent.Rot = rot
	ent.Dimensions = dim
	scene.entities = append(scene.entities, ent)
	return &ent
}

func (scene *Scene) GetLastEntity() *EcsEntity {
	return &scene.entities[len(scene.entities)-1]
}

func (scene *Scene) WriteComponentToLastEntity(component interface{}) {
	WriteComponent(&scene.Ecs_engine, &scene.entities[len(scene.entities)-1], component)
}

func (scene *Scene) DeleteComponentFromLastEntity(component interface{}) {
	DeleteComponent(&scene.Ecs_engine, &scene.entities[len(scene.entities)-1], component)
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

type SpriteComponent struct {
	Component
	Texture Texture2D
	Tint    RGBA8
}

func (t *SpriteComponent) ComponentSet(val interface{}) { *t = val.(SpriteComponent) }

type SpriteRenderOriginSystem struct {
	EcsSystemImpl
	Sprites *SpriteBatch
	Offset  Vector2f
	Scale   float32
}

func (_render *SpriteRenderOriginSystem) Update(dt float32) {
	EachEntity(SpriteComponent{}, func(entity *EcsEntity, a interface{}) {
		sprite := a.(SpriteComponent)
		halfDim := NewVector2f(_render.Offset.X*float32(sprite.Texture.Width)/2.0, _render.Offset.Y*float32(sprite.Texture.Height)/2.0)
		_render.Sprites.DrawSpriteOriginScaledRotated(entity.Pos.Add(halfDim), Vector2fZero, Vector2fOne, _render.Scale, &sprite.Texture, sprite.Tint, entity.Rot)
	})

}

type ShapesDrawingSystem struct {
	EcsSystemImpl
	Shapes *ShapeBatch
}

func (sds *ShapesDrawingSystem) Update(dt float32) {
	EachEntity(LineRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		lineComp := a.(LineRenderComponent)
		sds.Shapes.DrawLine(lineComp.FromPoint, lineComp.ToPoint, WHITE)
	})
	EachEntity(TriangleRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		lineComp := a.(TriangleRenderComponent)
		sds.Shapes.DrawTriangleRotated(entity.Pos, lineComp.Dimensions, WHITE, float32(entity.Rot))
	})
	EachEntity(RectRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		rectComp := a.(RectRenderComponent)
		sds.Shapes.DrawRectRotated(entity.Pos, entity.Dimensions, rectComp.Tint, entity.Rot)
	})
	EachEntity(FillRectRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		rectComp := a.(FillRectRenderComponent)
		sds.Shapes.DrawFillRectRotated(entity.Pos, entity.Dimensions, rectComp.Tint, entity.Rot)
	})
	EachEntity(CircleRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		circle := a.(CircleRenderComponent)
		sds.Shapes.DrawCircle(entity.Pos, entity.Dimensions.X, circle.Tint)
	})
}

type LineRenderComponent struct {
	Component
	FromPoint Vector2f
	ToPoint   Vector2f
}

func (t *LineRenderComponent) ComponentSet(val interface{}) { *t = val.(LineRenderComponent) }

type TriangleRenderComponent struct {
	Component
	Dimensions Vector2f
}

func (t *TriangleRenderComponent) ComponentSet(val interface{}) { *t = val.(TriangleRenderComponent) }

type RectRenderComponent struct {
	Component
	Tint RGBA8
}

func (t *RectRenderComponent) ComponentSet(val interface{}) { *t = val.(RectRenderComponent) }

type FillRectRenderComponent struct {
	Component
	Tint RGBA8
}

func (t *FillRectRenderComponent) ComponentSet(val interface{}) { *t = val.(FillRectRenderComponent) }

type CircleRenderComponent struct {
	Component
	Tint RGBA8
}

func (t *CircleRenderComponent) ComponentSet(val interface{}) { *t = val.(CircleRenderComponent) }

/*
###################################################################
################# UI UI UI UI UI UI UI UI UI UI UI ################
###################################################################
*/

type UIShapesDrawingSystem struct {
	EcsSystemImpl
	Shapes *ShapeBatch
}

func (sds *UIShapesDrawingSystem) Update(dt float32) {
	EachEntity(LineUIRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		lineComp := a.(LineUIRenderComponent)
		sds.Shapes.DrawLine(lineComp.FromPoint, lineComp.ToPoint, WHITE)
	})
	EachEntity(TriangleUIRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		triComp := a.(TriangleUIRenderComponent)
		sds.Shapes.DrawTriangleRotated(entity.Pos, triComp.Dimensions, triComp.Tint, float32(entity.Rot))
	})
	EachEntity(RectUIRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		rectComp := a.(RectUIRenderComponent)
		sds.Shapes.DrawRectRotated(entity.Pos, rectComp.Dimensions, rectComp.Tint, entity.Rot)
	})
	EachEntity(FillRectUIRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		rectComp := a.(FillRectUIRenderComponent)
		sds.Shapes.DrawFillRectRotated(entity.Pos, entity.Dimensions, rectComp.Tint, entity.Rot)
	})
	EachEntity(CircleUIRenderComponent{}, func(entity *EcsEntity, a interface{}) {
		circle := a.(CircleUIRenderComponent)
		sds.Shapes.DrawCircle(entity.Pos, entity.Dimensions.X, circle.Tint)
	})
}

type LineUIRenderComponent struct {
	Component
	FromPoint Vector2f
	ToPoint   Vector2f
}

func (t *LineUIRenderComponent) ComponentSet(val interface{}) { *t = val.(LineUIRenderComponent) }

type TriangleUIRenderComponent struct {
	Component
	Dimensions Vector2f
	Tint       RGBA8
}

func (t *TriangleUIRenderComponent) ComponentSet(val interface{}) {
	*t = val.(TriangleUIRenderComponent)
}

type RectUIRenderComponent struct {
	Component
	Dimensions Vector2f
	Tint       RGBA8
}

func (t *RectUIRenderComponent) ComponentSet(val interface{}) { *t = val.(RectUIRenderComponent) }

type FillRectUIRenderComponent struct {
	Component
	Dimensions Vector2f
	Tint       RGBA8
}

func (t *FillRectUIRenderComponent) ComponentSet(val interface{}) {
	*t = val.(FillRectUIRenderComponent)
}

type CircleUIRenderComponent struct {
	Component
	Tint RGBA8
}

func (t *CircleUIRenderComponent) ComponentSet(val interface{}) { *t = val.(CircleUIRenderComponent) }

type CircleUIRenderSystem struct {
	EcsSystemImpl
	Shapes *ShapeBatch
}

/*
###################################################################
###################################################################
*/
