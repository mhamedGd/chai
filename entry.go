package chai

import (
	"strings"
	"syscall/js"

	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

var app_url string

// Used to make the update function only available in the local App struct, to the whole file
var appUpdate func(float32)
var appDraw func(float32)

var currentWidth, currentHeight int
var canvas js.Value
var appRef *App

func GetCanvasWidth() int {
	return canvas.Get("width").Int()
}
func GetCanvasHeigth() int {
	return canvas.Get("height").Int()
}

var Cam Camera2D
var Renderer Renderer2D

var uiCam Camera2D

var started bool = false

var physics_world PhysicsWorld
var RenderQuadTreeContainer customtypes.StaticQuadTreeContainer[customtypes.Pair[VisualTransform, RenderObject]]
var DynamicRenderQuadTreeContainer customtypes.DynamicQuadTreeContainer[RenderObject]

var MouseCanvasPos Vector2f
var TouchCanvasPos [2]Vector2f
var canvasBoundingClientRect js.Value

var mousePressed MouseButton
var numOfFingersTouching uint8
var LeftMouseJustPressed customtypes.ChaiEvent1[int]

var collisionTiles customtypes.List[customtypes.Pair[Vector2f, Vector2f]]

func GetPhysicsWorld() *PhysicsWorld {
	return &physics_world
}

func GetDeltaTime() float32 {
	return deltaTime
}

func Run(_app *App) {
	// App URL Value Assignment
	// /////////////////////////////
	appRef = _app
	_app.fillDefaults()
	app_url = js.Global().Get("location").Get("href").String()
	if strings.Contains(app_url, "index.html") {
		app_url = strings.ReplaceAll(app_url, "index.html", "")
	}
	LogF("%v", app_url)
	// /////////////////////////////
	// Canvas Context Getter
	// /////////////////////////////
	js.Global().Get("document").Set("title", _app.Title)

	canvas = js.Global().Get("document").Call("getElementById", "viewport")

	canvasContext = canvas.Call("getContext", "webgl2")
	Assert(!canvasContext.IsNull(), "CANVAS: Failed to Get Context")

	canvas.Set("width", _app.Width)
	canvas.Set("height", _app.Height)

	canvasContext.Call("blendFunc", canvasContext.Get("ONE"), canvasContext.Get("ONE_MINUS_SRC_ALPHA"), canvasContext.Get("ONE"), canvasContext.Get("ONE"))
	canvasContext.Call("enable", canvasContext.Get("BLEND"))
	// canvasContext.Call("enable", canvasContext.Get("DEPTH_TEST"))
	// /////////////////////////////

	appUpdate = _app.OnUpdate
	appDraw = _app.OnDraw

	appPresets(_app)

	modulesInitialization(_app)

	canvasContext.Call("viewport", 0, 0, appRef.Width, appRef.Height)

	eventsInitialization(_app)

	_app.OnStart()
	Assert(current_scene != nil, "Current Scene is none")
	started = true

	select {}
}

var ElapsedTime float32
var deltaTime float32

const CAP_DELTA_TIME float32 = 60.0 / 1000.0

const FIXED_UPDATE_INTERVAL float32 = 1.0 / 50.0
const MAX_FIXED_CYCLES_PER_FRAME = 5

var timeAccumulation float32

func jSUpdate(_this js.Value, _inputs []js.Value) interface{} {
	if !started {
		return nil
	}

	deltaTime = float32(_inputs[0].Float())
	if deltaTime > CAP_DELTA_TIME {
		deltaTime = CAP_DELTA_TIME
	}

	currentWidth = canvas.Get("width").Int()
	currentHeight = canvas.Get("height").Int()
	updateInput()
	current_scene.OnUpdate(deltaTime)

	Cam.Update(*appRef)
	uiCam.Update(*appRef)

	for _, v := range DynamicRenderQuadTreeContainer.AllItems().AllItems() {
		t := GetComponentPtr[VisualTransform](current_scene, v.GetItem().m_EntId)
		DynamicRenderQuadTreeContainer.Relocate(&v, Rect{t.Position.Subtract(t.Dimensions.Scale(0.5)), t.Dimensions})
	}
	// LogF("%v", RenderQuadTreeContainer.allItems.Count())

	appUpdate(deltaTime)
	ElapsedTime += deltaTime

	timeAccumulation += deltaTime
	if timeAccumulation > (MAX_FIXED_CYCLES_PER_FRAME * FIXED_UPDATE_INTERVAL) {
		timeAccumulation = FIXED_UPDATE_INTERVAL
	}

	for timeAccumulation >= FIXED_UPDATE_INTERVAL {
		timeAccumulation -= FIXED_UPDATE_INTERVAL
		current_scene.OnPhysics(FIXED_UPDATE_INTERVAL)
		// physics_world.cpSpace.Step(float64(1 / 60.0))
		physics_world.box2dWorld.Step(float64(FIXED_UPDATE_INTERVAL), 16, 6)
	}
	return nil
}

func jSDraw(_this js.Value, _inputs []js.Value) interface{} {
	if !started {
		return nil
	}

	canvasContext.Call("viewport", 0, 0, currentWidth, currentHeight)
	setBackgroundColor(current_scene.Background)
	canvasContext.Call("clear", canvasContext.Get("COLOR_BUFFER_BIT"))
	canvasContext.Call("clear", canvasContext.Get("DEPTH_BUFFER_BIT"))

	Renderer.Begin()
	RenderQuadTreeContainer.QuadsCount = 0

	ScreenDims := NewVector2f(float32(appRef.Width), float32(appRef.Height))
	ScreenRect := Rect{Position: ScreenDims.Scale(-0.5).Scale(1 / Cam.GetScale()).Add(Cam.GetPosition()), Size: ScreenDims.Scale(1 / Cam.GetScale())}

	for _, v := range RenderQuadTreeContainer.Search(ScreenRect).Data {
		t := v.First
		v.Second.m_ObjectType(t.Position, t.Dimensions, t.UV1, t.UV2, t.Tint, t.Rotation, t.Z, v.Second.m_Texture)
		RenderQuadTreeContainer.QuadsCount++
	}
	dynmaicQuadInView := DynamicRenderQuadTreeContainer.Search(ScreenRect)

	for _, v := range dynmaicQuadInView.AllItems() {
		// chai.Shapes.DrawFillRect(v.First.Position, v.First.Dimensions, v.Second.Tint)
		// rects_count++
		it := v.GetItem()
		t := GetComponentPtr[VisualTransform](current_scene, it.m_EntId)
		// t := current_scene.transforms.Get(it.entId)
		// v.objectType(&Shapes, t.Position, t.Dimensions, v.tint, t.Rotation)
		// Shapes.DrawFillRectRotated(t.Position, t.Dimensions, it.tint, t.Rotation)
		v.GetItem().m_ObjectType(t.Position, t.Dimensions, t.UV1, t.UV2, t.Tint, t.Rotation, t.Z, v.GetItem().m_Texture)

	}
	current_scene.OnDraw()
	appDraw(deltaTime)
	Renderer.End()
	Renderer.Render()
	return nil
}

func setBackgroundColor(_color RGBA8) {
	canvasContext.Call("clearColor", _color.GetColorRFloat32(), _color.GetColorGFloat32(), _color.GetColorBFloat32(), 1.0)
}

func NumOfQuadsInView() int {
	return RenderQuadTreeContainer.QuadsInViewCount()
}

func GetStaticQuadsInRect(rArea Rect) customtypes.List[*customtypes.Pair[VisualTransform, RenderObject]] {
	list := RenderQuadTreeContainer.Search(rArea)
	return list
}

func GetDynamicQuadsInRect(rArea Rect) customtypes.List[*customtypes.QuadTreeItem[RenderObject]] {
	list := DynamicRenderQuadTreeContainer.Search(rArea)
	return list
}
