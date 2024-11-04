package chai

import (
	"github.com/mhamedGd/chai/customtypes"
	"github.com/mhamedGd/chai/ecs"
	. "github.com/mhamedGd/chai/math"
)

var current_scene *Scene

type EntId = ecs.Id

type EcsSystem interface {
	Update(dt float32)
}

func Iterate1[A any](_f func(EntId, *A)) {
	query1 := ecs.Query1[A](GetCurrentScene().m_EcsWorld)
	query1.MapId(_f)
}
func Iterate2[A, B any](_f func(EntId, *A, *B)) {
	query2 := ecs.Query2[A, B](GetCurrentScene().m_EcsWorld)
	query2.MapId(_f)
}
func Iterate3[A, B, C any](_f func(EntId, *A, *B, *C)) {
	query3 := ecs.Query3[A, B, C](GetCurrentScene().m_EcsWorld)
	query3.MapId(_f)
}
func Iterate4[A, B, C, D any](_f func(EntId, *A, *B, *C, *D)) {
	query4 := ecs.Query4[A, B, C, D](GetCurrentScene().m_EcsWorld)
	query4.MapId(_f)
}
func Iterate5[A, B, C, D, E any](_f func(EntId, *A, *B, *C, *D, *E)) {
	query4 := ecs.Query5[A, B, C, D, E](GetCurrentScene().m_EcsWorld)
	query4.MapId(_f)
}

type Scene struct {
	Background      RGBA8
	m_EcsWorld      *ecs.World
	m_StartSystems  customtypes.List[func(*Scene)]
	m_UpdateSystems customtypes.List[func(*Scene, float32)]
	m_RenderSystems customtypes.List[func(*Scene, float32)]
	m_Tags          customtypes.Map[string, customtypes.List[EntId]]
}

func NewScene() Scene {
	return Scene{
		m_EcsWorld:      ecs.NewWorld(),
		m_StartSystems:  customtypes.NewList[func(*Scene)](),
		m_UpdateSystems: customtypes.NewList[func(*Scene, float32)](),
		m_RenderSystems: customtypes.NewList[func(*Scene, float32)](),
		m_Tags:          customtypes.NewMap[string, customtypes.List[EntId]](),
	}
}

func (scene *Scene) AddTag(_entId EntId, _tagName string) {
	if !scene.m_Tags.Has(_tagName) {
		scene.m_Tags.Insert(_tagName, customtypes.NewList[EntId]())
	}

	m_Tags := scene.m_Tags.Get(_tagName)
	m_Tags.PushBack(_entId)
	scene.m_Tags.Set(_tagName, m_Tags)
}

func (scene *Scene) HasTag(_entId EntId, _tagName string) bool {
	if !scene.m_Tags.Has(_tagName) {
		return false
	}

	for _, id := range scene.m_Tags.Get(_tagName).Data {
		if id == _entId {
			return true
		}
	}

	return false
}

func ChangeScene(_scene *Scene) {
	if current_scene != nil {
		current_scene.terminateScene()
	}

	current_scene = _scene
	go func() {
		for _, s := range _scene.m_StartSystems.Data {
			s(_scene)
		}
	}()
}

func (scene *Scene) terminateScene() {
	// Iterate1[RigidBodyComponent](func(i ecs.Id, rbc *RigidBodyComponent) {
	// 	freeRigidbody(rbc)
	// })
	Iterate1[DynamicBodyComponent](func(i ecs.Id, dbc *DynamicBodyComponent) {
		// freeRigidbody(rbc)
		physics_world.box2dWorld.DestroyBody(dbc.m_B2Body)
	})
	Iterate1[StaticBodyComponent](func(i ecs.Id, sbc *StaticBodyComponent) {
		// freeRigidbody(rbc)
		physics_world.box2dWorld.DestroyBody(sbc.m_B2Body)
	})
	Iterate1[KinematicBodyComponent](func(i ecs.Id, kbc *KinematicBodyComponent) {
		// freeRigidbody(rbc)
		physics_world.box2dWorld.DestroyBody(kbc.m_B2Body)
	})
	physics_world = newPhysicsWorld(NewVector2f(0.0, -98.0))

	DynamicRenderQuadTreeContainer.Clear()
	RenderQuadTreeContainer.Clear()
	scene.m_EcsWorld = ecs.NewWorld()
	// scene.transforms.Clear()
	// scene.m_UpdateSystems = scene.m_UpdateSystems[:0]
	// scene.m_RenderSystems = scene.m_RenderSystems[:0]
	scene.m_UpdateSystems.Clear()
	scene.m_RenderSystems.Clear()

}

func (scene *Scene) NewEntityId() ecs.Id {
	id := scene.m_EcsWorld.NewId()
	// scene.transforms.Insert(id, t)
	return id
}

// A gateway to the Component type in chai/Ecs
type convComponent = ecs.Component

func ToComponent[T any](_comp T) ecs.Box[T] {
	return ecs.C(_comp)
}

func (scene *Scene) AddComponents(_entId ecs.Id, _comps ...convComponent) {
	ecs.Write(scene.m_EcsWorld, _entId, _comps...)
}

func GetComponent[T any](_scene *Scene, _entityId EntId) (T, bool) {
	return ecs.Read[T](_scene.m_EcsWorld, _entityId)
}

func GetComponentPtr[T any](_scene *Scene, _entityId EntId) *T {
	return ecs.ReadPtr[T](_scene.m_EcsWorld, _entityId)
}

func (scene *Scene) AddEntity(_comps ...convComponent) EntId {
	id := scene.NewEntityId()
	scene.AddComponents(id, _comps...)

	return id
}

func DestroyComponent[T any](_scene *Scene, _entityId EntId) {
	comp, _ := GetComponent[T](_scene, _entityId)
	ecs.DeleteComponent(_scene.m_EcsWorld, _entityId, ToComponent(comp))
}

func Destroy(_scene *Scene, _entityId EntId) {
	ecs.Delete(_scene.m_EcsWorld, _entityId)
}

func (scene *Scene) NewStartSystem(_sys func(*Scene)) {
	scene.m_StartSystems.PushBack(_sys)
}
func (scene *Scene) NewUpdateSystem(_sys func(_this_scene *Scene, _dt float32)) {
	scene.m_UpdateSystems.PushBack(_sys)
}
func (scene *Scene) NewRenderSystem(_sys func(_this_scene *Scene, _dt float32)) {
	scene.m_RenderSystems.PushBack(_sys)
}

func (scene *Scene) SetGravity(_newGravity Vector2f) {
	// physics_world.cpSpace.SetGravity(cpVector2f(new_gravity))
	physics_world.box2dWorld.SetGravity(vec2fToB2Vec(_newGravity))
}

func (scene *Scene) OnUpdate(_dt float32) {
	for _, sys := range scene.m_UpdateSystems.Data {
		sys(scene, _dt)
	}
}

func (scene *Scene) OnDraw() {
	for _, sys := range scene.m_RenderSystems.Data {
		sys(scene, deltaTime)
	}
}

func GetCurrentScene() *Scene {
	return current_scene
}
