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
	update_systems []EcsSystem
	render_systems []EcsSystem
	OnSceneStart   func(thisScene *Scene)
	OnSceneUpdate  func(dt float32, thisScene *Scene)
}

func NewScene() Scene {
	return Scene{
		Ecs_World:      ecs.NewWorld(),
		update_systems: make([]EcsSystem, 0),
		render_systems: make([]EcsSystem, 0),
		OnSceneStart:   func(thisScene *Scene) {},
		OnSceneUpdate:  func(dt float32, thisScene *Scene) {},
	}
}

func ChangeScene(scene *Scene) {
	if current_scene != nil {
		current_scene.terminateScene()
	}

	current_scene = scene
	go func() {
		current_scene.OnSceneStart(current_scene)
	}()
}

func (scene *Scene) terminateScene() {
	Iterate1[RigidBodyComponent](func(i ecs.Id, rbc *RigidBodyComponent) {
		freeRigidbody(rbc)
	})
	// scene.transforms.Clear()
	scene.update_systems = scene.update_systems[:0]
	scene.render_systems = scene.render_systems[:0]

}

func (scene *Scene) NewEntityId() ecs.Id {
	id := scene.Ecs_World.NewId()
	// scene.transforms.Insert(id, t)
	return id
}

// func GetTransform(entId EntId) Transform {
// 	return current_scene.transforms.Get(entId)
// }

// func GetPosition(entId EntId) Vector2f {
// 	return current_scene.transforms.Get(entId).Position
// }

// func SetPosition(entId EntId, newPosition Vector2f) {
// 	_t := current_scene.transforms.data[entId]
// 	_t.Position = newPosition
// 	current_scene.transforms.data[entId] = _t
// }

type Component = ecs.Component

func ToComponent[T any](comp T) ecs.Box[T] {
	return ecs.C(comp)
}

func (scene *Scene) AddComponents(EntId ecs.Id, comps ...Component) {
	ecs.Write(scene.Ecs_World, EntId, comps...)
}

func GetComponent[T any](scene *Scene, entityId EntId) (T, bool) {
	return ecs.Read[T](scene.Ecs_World, entityId)
}

func GetComponentPtr[T any](scene *Scene, entityId EntId) *T {
	return ecs.ReadPtr[T](scene.Ecs_World, entityId)
}

func (scene *Scene) NewUpdateSystem(sys EcsSystem) {
	scene.update_systems = append(scene.update_systems, sys)
}

func (scene *Scene) NewRenderSystem(sys EcsSystem) {
	scene.render_systems = append(scene.render_systems, sys)
}

func (scene *Scene) SetGravity(new_gravity Vector2f) {
	physics_world.cpSpace.SetGravity(cpVector2f(new_gravity))
}

func (scene *Scene) OnUpdate(dt float32) {
	for _, sys := range scene.update_systems {
		sys.Update(dt)
	}
}

func (scene *Scene) OnDraw() {
	for _, sys := range scene.render_systems {
		sys.Update(deltaTime)
	}
}

func GetCurrentScene() *Scene {
	return current_scene
}

type Transform struct {
	Position   Vector2f
	Dimensions Vector2f
	Rotation   float32
	Scale      float32
}

type SpriteComponent struct {
	Texture Texture2D
	Tint    RGBA8
}

type SpriteRenderSystem struct {
	EcsSystem
	Sprites *SpriteBatch
	Offset  Vector2f
	Scale   float32
}

func (_render *SpriteRenderSystem) Update(dt float32) {
	query2 := ecs.Query2[Transform, SpriteComponent](GetCurrentScene().Ecs_World)
	query2.MapId(func(id ecs.Id, t *Transform, s *SpriteComponent) {
		newOffset := _render.Offset.Rotate(t.Rotation, Vector2fZero)
		halfDim := NewVector2f(newOffset.X*float32(s.Texture.Width)/2.0, newOffset.Y*float32(s.Texture.Height)/2.0)
		_render.Sprites.DrawSpriteOriginScaledRotated(t.Position.Add(halfDim), Vector2fZero, Vector2fOne, _render.Scale, &s.Texture, s.Tint, t.Rotation)
	})
}

type ShapesDrawingSystem struct {
	EcsSystem
	Shapes *ShapeBatch
}

func (sds *ShapesDrawingSystem) Update(dt float32) {
	// lineQuery := ecs.Query1[LineRenderComponent](GetCurrentScene().Ecs_World)
	// lineQuery.MapId(func(id ecs.Id, line *LineRenderComponent) {
	// 	if Cam.IsBoxInView(line.FromPoint, AbsVector2f(line.ToPoint.Subtract(line.FromPoint))) {
	// 		sds.Shapes.DrawLine(line.FromPoint, line.ToPoint, line.Tint)
	// 	}
	// })

	queryTri := ecs.Query2[Transform, TriangleRenderComponent](GetCurrentScene().Ecs_World)
	queryTri.MapId(func(id ecs.Id, t *Transform, tri *TriangleRenderComponent) {
		if Cam.IsBoxInView(t.Position, tri.Dimensions.Scale(t.Scale)) {
			sds.Shapes.DrawTriangleRotated(t.Position, tri.Dimensions.Scale(t.Scale), tri.Tint, t.Rotation)
		}
	})

	// queryFillTri := ecs.Query2[Transform, FillTriangleRenderComponent](GetCurrentScene().Ecs_World)
	// queryFillTri.MapId(func(id ecs.Id, t *Transform, tri *FillTriangleRenderComponent) {
	// 	if Cam.IsBoxInView(t.Position, t.Dimensions.Scale(t.Scale)) {
	// 		sds.Shapes.DrawFillTriangleRotated(t.Position, t.Dimensions.Scale(t.Scale), tri.Tint, t.Rotation)
	// 	}
	// })

	queryRect := ecs.Query2[Transform, RectRenderComponent](GetCurrentScene().Ecs_World)
	queryRect.MapId(func(id ecs.Id, t *Transform, rect *RectRenderComponent) {
		if Cam.IsBoxInView(t.Position, rect.Dimensions.Scale(t.Scale)) {
			sds.Shapes.DrawRectRotated(t.Position, rect.Dimensions.Scale(t.Scale), rect.Tint, t.Rotation)
		}
	})
	// go func() {
	// 	queryFillRect := ecs.Query2[Transform, FillRectRenderComponent](GetCurrentScene().Ecs_World)
	// 	queryFillRect.MapId(func(id ecs.Id, t *Transform, rect *FillRectRenderComponent) {
	// 		if Cam.IsBoxInView(t.Position, t.Dimensions.Scale(t.Scale)) {
	// 			sds.Shapes.DrawFillRectRotated(t.Position, t.Dimensions.Scale(t.Scale), rect.Tint, t.Rotation)
	// 		}
	// 	})
	// }()

	queryFillRectBottom := ecs.Query2[Transform, FillRectBottomRenderComponent](GetCurrentScene().Ecs_World)
	queryFillRectBottom.MapId(func(id ecs.Id, t *Transform, rect *FillRectBottomRenderComponent) {
		rectDims := rect.Dimensions.Scale(t.Scale)
		if Cam.IsBoxInView(t.Position.Subtract(rectDims.Scale(0.5)), rectDims) {
			sds.Shapes.DrawFillRectBottomRotated(t.Position, rectDims, rect.Tint, t.Rotation)
		}
	})

	// queryCircle := ecs.Query2[Transform, CircleRenderComponent](GetCurrentScene().Ecs_World)
	// queryCircle.MapId(func(id ecs.Id, t *Transform, circ *CircleRenderComponent) {
	// 	if Cam.IsBoxInView(t.Position, NewVector2f(circ.Radius*t.Scale, circ.Radius*t.Scale)) {
	// 		sds.Shapes.DrawCircle(t.Position, circ.Radius*t.Scale, circ.Tint)
	// 	}
	// })
}

type LineRenderComponent struct {
	Tint RGBA8
}

type TriangleRenderComponent struct {
	Dimensions  Vector2f
	OffsetPivot Vector2f
	Tint        RGBA8
}

type FillTriangleRenderComponent struct {
	Tint RGBA8
}

type RectRenderComponent struct {
	Dimensions Vector2f
	Tint       RGBA8
}

type FillRectRenderComponent struct {
	Tint RGBA8
}

type FillRectBottomRenderComponent struct {
	Dimensions Vector2f
	Tint       RGBA8
}

type CircleRenderComponent struct {
	Tint RGBA8
}

/*
###################################################################
###################################################################
*/

type FontRenderSystem struct {
	EcsSystem
	fontbatch_atlas FontBatchAtlas
	FontSettings    FontBatchSettings
}

func (frs *FontRenderSystem) SetFont(fontPath string) {
	frs.fontbatch_atlas = LoadFontToAtlas(fontPath, &frs.FontSettings)
}

func (frs *FontRenderSystem) SetFontBatchRenderer(sb *SpriteBatch) {
	frs.fontbatch_atlas.sPatch = sb
}

// func (frs *FontRenderSystem) SetFont(_font *FontBatchAtlas) {
// 	frs.fontbatch_atlas = _font
// }

type FontRenderComponent struct {
	Text   string
	Scale  float32
	Offset Vector2f
	Tint   RGBA8
}

func (frs *FontRenderSystem) Update(dt float32) {
	Iterate2[Transform, FontRenderComponent](func(i ecs.Id, t *Transform, frc *FontRenderComponent) {
		frs.fontbatch_atlas.DrawString(frc.Text, t.Position.Add(frc.Offset), frc.Scale, frc.Tint)
	})
}
