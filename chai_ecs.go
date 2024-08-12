package chai

import (
	"github.com/mhamedGd/chai/ecs"
)

var current_scene *Scene

type EntId = ecs.Id

type EcsSystem interface {
	Update(dt float32)
}

func Iterate1[A any](f func(EntId, *A)) {
	query1 := ecs.Query1[A](GetCurrentScene().Ecs_World)
	query1.MapId(f)
}
func Iterate2[A, B any](f func(EntId, *A, *B)) {
	query2 := ecs.Query2[A, B](GetCurrentScene().Ecs_World)
	query2.MapId(f)
}
func Iterate3[A, B, C any](f func(EntId, *A, *B, *C)) {
	query3 := ecs.Query3[A, B, C](GetCurrentScene().Ecs_World)
	query3.MapId(f)
}
func Iterate4[A, B, C, D any](f func(EntId, *A, *B, *C, *D)) {
	query4 := ecs.Query4[A, B, C, D](GetCurrentScene().Ecs_World)
	query4.MapId(f)
}
func Iterate5[A, B, C, D, E any](f func(EntId, *A, *B, *C, *D, *E)) {
	query4 := ecs.Query5[A, B, C, D, E](GetCurrentScene().Ecs_World)
	query4.MapId(f)
}

type Scene struct {
	Background     RGBA8
	Ecs_World      *ecs.World
	start_systems  List[func(*Scene)]
	update_systems List[func(*Scene, float32)]
	render_systems List[func(*Scene, float32)]
	tags           Map[string, List[EntId]]
}

func NewScene() Scene {
	return Scene{
		Ecs_World:      ecs.NewWorld(),
		start_systems:  NewList[func(*Scene)](),
		update_systems: NewList[func(*Scene, float32)](),
		render_systems: NewList[func(*Scene, float32)](),
		tags:           NewMap[string, List[EntId]](),
	}
}

func (scene *Scene) AddTag(ent_id EntId, tag_name string) {
	if !scene.tags.Has(tag_name) {
		scene.tags.Insert(tag_name, NewList[EntId]())
	}

	tags := scene.tags.Get(tag_name)
	tags.PushBack(ent_id)
	scene.tags.Set(tag_name, tags)
}

func (scene *Scene) HasTag(ent_id EntId, tag_name string) bool {
	if !scene.tags.Has(tag_name) {
		return false
	}

	for _, id := range scene.tags.Get(tag_name).Data {
		if id == ent_id {
			return true
		}
	}

	return false
}

func ChangeScene(scene *Scene) {
	if current_scene != nil {
		current_scene.terminateScene()
	}

	current_scene = scene
	go func() {
		for _, s := range scene.start_systems.Data {
			s(scene)
		}
	}()
}

func (scene *Scene) terminateScene() {
	Iterate1[RigidBodyComponent](func(i ecs.Id, rbc *RigidBodyComponent) {
		freeRigidbody(rbc)
	})
	// scene.transforms.Clear()
	// scene.update_systems = scene.update_systems[:0]
	// scene.render_systems = scene.render_systems[:0]
	scene.update_systems.Clear()
	scene.render_systems.Clear()

}

func (scene *Scene) NewEntityId() ecs.Id {
	id := scene.Ecs_World.NewId()
	// scene.transforms.Insert(id, t)
	return id
}

// A gateway to the Component type in chai/Ecs
type convComponent = ecs.Component

func ToComponent[T any](comp T) ecs.Box[T] {
	return ecs.C(comp)
}

func (scene *Scene) AddComponents(EntId ecs.Id, comps ...convComponent) {
	ecs.Write(scene.Ecs_World, EntId, comps...)
}

func GetComponent[T any](scene *Scene, entityId EntId) (T, bool) {
	return ecs.Read[T](scene.Ecs_World, entityId)
}

func GetComponentPtr[T any](scene *Scene, entityId EntId) *T {
	return ecs.ReadPtr[T](scene.Ecs_World, entityId)
}

func (scene *Scene) AddEntity(comps ...convComponent) EntId {
	id := scene.NewEntityId()
	scene.AddComponents(id, comps...)

	return id
}

func DestroyComponent[T any](scene *Scene, entityId EntId) {
	comp, _ := GetComponent[T](scene, entityId)
	ecs.DeleteComponent(scene.Ecs_World, entityId, ToComponent(comp))
}

func Destroy(scene *Scene, entityId EntId) {
	ecs.Delete(scene.Ecs_World, entityId)
}

func (scene *Scene) NewStartSystem(sys func(*Scene)) {
	scene.start_systems.PushBack(sys)
}
func (scene *Scene) NewUpdateSystem(sys func(_this_scene *Scene, _dt float32)) {
	scene.update_systems.PushBack(sys)
}
func (scene *Scene) NewRenderSystem(sys func(_this_scene *Scene, _dt float32)) {
	scene.render_systems.PushBack(sys)
}

func (scene *Scene) SetGravity(new_gravity Vector2f) {
	physics_world.cpSpace.SetGravity(cpVector2f(new_gravity))
}

func (scene *Scene) OnUpdate(dt float32) {
	for _, sys := range scene.update_systems.Data {
		sys(scene, dt)
	}
}

func (scene *Scene) OnDraw() {
	for _, sys := range scene.render_systems.Data {
		sys(scene, deltaTime)
	}
}

func GetCurrentScene() *Scene {
	return current_scene
}
