package chai

import (
	"strings"
	"syscall/js"
)

var app_url string

type App struct {
	Width    int
	Height   int
	Title    string
	OnStart  func()
	OnUpdate func(float32)
	OnDraw   func(float32)
	OnEvent  func(*AppEvent)
}

// Used to make the update function only available in the local App struct, to the whole file
var tempUpdate func(float32)
var tempDraw func(float32)

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
var Shapes ShapeBatch
var Sprites SpriteBatch
var Renderer Renderer2D

var uiCam Camera2D
var UIShapes ShapeBatch
var UISprites SpriteBatch

var started bool = false

var physics_world PhysicsWorld
var RenderQuadTreeContainer StaticQuadTreeContainer[Pair[VisualTransform, RenderObject]]
var DynamicRenderQuadTreeContainer DynamicQuadTreeContainer[RenderObject]

var MouseCanvasPos Vector2f
var TouchCanvasPos [2]Vector2f
var canvasBoundingClientRect js.Value

var mousePressed MouseButton
var numOfFingersTouching uint8
var LeftMouseJustPressed ChaiEvent1[int]

func GetPhysicsWorld() *PhysicsWorld {
	return &physics_world
}

func GetDeltaTime() float32 {
	return deltaTime
}

func (_app *App) fillDefaults() {
	if _app.OnStart == nil {
		_app.OnStart = func() {

		}
	}
	if _app.OnUpdate == nil {
		_app.OnUpdate = func(dt float32) {

		}
	}
	if _app.OnDraw == nil {
		_app.OnDraw = func(dt float32) {

		}
	}
	if _app.OnEvent == nil {
		_app.OnEvent = func(ae *AppEvent) {

		}
	}

	setPhysicsFunctions(PHYSICS_ENGINE_BOX2D)
}

func Run(_app *App) {

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		ErrorF("PANICKED - %v", r)
	// 	}
	// }()

	appRef = _app
	_app.fillDefaults()
	app_url = js.Global().Get("location").Get("href").String()
	if strings.Contains(app_url, "index.html") {
		app_url = strings.ReplaceAll(app_url, "index.html", "")
	}
	LogF("%v", app_url)

	js.Global().Get("document").Set("title", _app.Title)

	canvas = js.Global().Get("document").Call("getElementById", "viewport")

	canvasContext = canvas.Call("getContext", "webgl2")
	Assert(!canvasContext.IsNull(), "CANVAS: Failed to Get Context")

	canvas.Set("width", _app.Width)
	canvas.Set("height", _app.Height)

	canvasContext.Call("blendFunc", canvasContext.Get("ONE"), canvasContext.Get("ONE_MINUS_SRC_ALPHA"), canvasContext.Get("ONE"), canvasContext.Get("ONE"))
	canvasContext.Call("enable", canvasContext.Get("BLEND"))
	canvasContext.Call("enable", canvasContext.Get("DEPTH_TEST"))

	tempUpdate = _app.OnUpdate
	tempDraw = _app.OnDraw

	audioContext = js.Global().Get("AudioContext").New()
	js.Global().Get("document").Call("addEventListener", "visibilitychange", js.FuncOf(func(this js.Value, args []js.Value) any {
		if this.Get("visibilityState").String() == "hidden" {
			SuspendAudioContext()
		} else {
			ResumeAudioContext()
		}

		return 0
	}))

	js.Global().Call("addEventListener", "focus", js.FuncOf(func(this js.Value, args []js.Value) any {
		ResumeAudioContext()

		return 0
	}))
	js.Global().Call("addEventListener", "blur", js.FuncOf(func(this js.Value, args []js.Value) any {
		SuspendAudioContext()

		return 0
	}))

	TouchCanvasPos[0] = NewVector2f(0.0, 0.0)
	TouchCanvasPos[1] = NewVector2f(0.0, 0.0)

	// canvasContext.Call("pixelStorei", canvasContext.Get("UNPACK_ALIGNMENT"))
	initTextures()
	InitInputs()

	physics_world = newPhysicsWorld(NewVector2f(0.0, -98.0))

	js.Global().Set("js_update", js.FuncOf(JSUpdate))
	js.Global().Set("js_draw", js.FuncOf(JSDraw))

	// if I put it above the "js_start" then it would take a lot of time to run
	Cam.Init(*_app)
	Cam.centerOffset = NewVector2f(float32(_app.Width)/2.0, float32(appRef.Height)/2.0)
	Cam.Update(*_app)

	uiCam.Init(*_app)
	uiCam.Update(*_app)

	RenderQuadTreeContainer = NewStaticQuadTreeContainer[Pair[VisualTransform, RenderObject]]()
	DynamicRenderQuadTreeContainer = NewDynamicQuadTreeContainer[RenderObject]()
	RenderQuadTreeContainer.Resize(Rect{Position: Vector2fZero, Size: NewVector2f(1000.0, 1000.0)})
	DynamicRenderQuadTreeContainer.Resize(Rect{Position: Vector2fZero, Size: NewVector2f(10000.0, 10000.0)})

	Shapes.Init()
	Assert(Shapes.Initialized, "Shapes Rendering was not initialized successfully")
	UIShapes.Init()
	Assert(UIShapes.Initialized, "Shapes Rendering was not initialized successfully")

	Sprites.Init("")
	UISprites.Init("")
	Sprites.RenderCam = &Cam
	UISprites.RenderCam = &uiCam
	// Renderer = NewRenderer(_app, 50_000)

	canvasContext.Call("viewport", 0, 0, appRef.Width, appRef.Height)

	mousePressed = MouseButtonNull
	LeftMouseJustPressed.init()

	addEventListenerWindow(JS_KEYUP, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_KEYDOWN, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEDOWN, func(ae *AppEvent) {
		mousePressed = ae.Button
		onMousePressed()
		switch mousePressed {
		case LEFT_MOUSE_BUTTON:
			LeftMouseJustPressed.Invoke(0)
		}
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEUP, func(ae *AppEvent) {
		mousePressed = MouseButtonNull
		onMouseReleased()
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEMOVED, func(ae *AppEvent) {
		canvasBoundingClientRect = canvas.Call("getBoundingClientRect")

		MouseCanvasPos.X = (float32(ae.GetJsEvent().Get("clientX").Int()) - float32(canvasBoundingClientRect.Get("left").Int())) / float32(canvasBoundingClientRect.Get("width").Int()) * float32(canvas.Get("width").Int())
		MouseCanvasPos.Y = float32(canvas.Get("height").Int()) - (float32(ae.GetJsEvent().Get("clientY").Int())-float32(canvasBoundingClientRect.Get("top").Int()))/float32(canvasBoundingClientRect.Get("height").Int())*float32(canvas.Get("height").Int())
		_app.OnEvent(ae)
	})

	addEventListenerWindow(JS_TOUCHSTART, func(ae *AppEvent) {
		numOfFingersTouching = ae.NUM_FINGERS
		onTouchStart(ae.NUM_FINGERS)

		canvasBoundingClientRect = canvas.Call("getBoundingClientRect")

		MouseCanvasPos.X = (float32(ae.GetJsEvent().Get("touches").Index(0).Get("clientX").Int()) - float32(canvasBoundingClientRect.Get("left").Int())) / float32(canvasBoundingClientRect.Get("width").Int()) * float32(canvas.Get("width").Int())
		MouseCanvasPos.Y = float32(canvas.Get("height").Int()) - (float32(ae.GetJsEvent().Get("touches").Index(0).Get("clientY").Int())-float32(canvasBoundingClientRect.Get("top").Int()))/float32(canvasBoundingClientRect.Get("height").Int())*float32(canvas.Get("height").Int())
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_TOUCHEND, func(ae *AppEvent) {
		numOfFingersTouching = ae.NUM_FINGERS
		onTouchEnd(ae.NUM_FINGERS)

		_app.OnEvent(ae)
	})

	addEventListenerWindow(JS_TOUCHMOVED, func(ae *AppEvent) {
		canvasBoundingClientRect = canvas.Call("getBoundingClientRect")

		MouseCanvasPos.X = (float32(ae.GetJsEvent().Get("touches").Index(0).Get("clientX").Int()) - float32(canvasBoundingClientRect.Get("left").Int())) / float32(canvasBoundingClientRect.Get("width").Int()) * float32(canvas.Get("width").Int())
		MouseCanvasPos.Y = float32(canvas.Get("height").Int()) - (float32(ae.GetJsEvent().Get("touches").Index(0).Get("clientY").Int())-float32(canvasBoundingClientRect.Get("top").Int()))/float32(canvasBoundingClientRect.Get("height").Int())*float32(canvas.Get("height").Int())
		_app.OnEvent(ae)
	})
	/*
		////////// FINAL CHECKS ////////////
	*/
	_app.OnStart()
	Assert(current_scene != nil, "Current Scene is none")
	if !started {
		started = true
	}

	select {}
}

var ElapsedTime float32
var deltaTime float32

const CAP_DELTA_TIME float32 = 50.0 / 1000.0

const FIXED_UPDATE_INTERVAL float32 = 1.0 / 60.0
const MAX_FIXED_CYCLES_PER_FRAME = 5

var timeAccumulation float32

func JSUpdate(this js.Value, inputs []js.Value) interface{} {
	if !started {
		return nil
	}

	deltaTime = float32(inputs[0].Float())
	if deltaTime > CAP_DELTA_TIME {
		deltaTime = CAP_DELTA_TIME
	}
	// Renderer.Begin()

	currentWidth = canvas.Get("width").Int()
	currentHeight = canvas.Get("height").Int()
	current_scene.OnUpdate(deltaTime)
	updateInput()
	Cam.Update(*appRef)
	uiCam.Update(*appRef)

	for _, v := range DynamicRenderQuadTreeContainer.allItems.AllItems() {
		t := GetComponentPtr[VisualTransform](current_scene, v.item.entId)
		DynamicRenderQuadTreeContainer.Relocate(&v, Rect{t.Position.Subtract(t.Dimensions.Scale(0.5)), t.Dimensions})
	}
	// LogF("%v", RenderQuadTreeContainer.allItems.Count())

	tempUpdate(deltaTime)
	ElapsedTime += deltaTime

	timeAccumulation += deltaTime
	if timeAccumulation > (MAX_FIXED_CYCLES_PER_FRAME * FIXED_UPDATE_INTERVAL) {
		timeAccumulation = FIXED_UPDATE_INTERVAL
	}

	for timeAccumulation >= FIXED_UPDATE_INTERVAL {
		timeAccumulation -= FIXED_UPDATE_INTERVAL
		// physics_world.cpSpace.Step(float64(1 / 60.0))
		physics_world.box2dWorld.Step(1/60.0, 8, 3)
	}
	return nil
}

func JSDraw(this js.Value, inputs []js.Value) interface{} {
	if !started {
		return nil
	}
	canvasContext.Call("viewport", 0, 0, currentWidth, currentHeight)
	setBackgroundColor(current_scene.Background)
	canvasContext.Call("clear", canvasContext.Get("COLOR_BUFFER_BIT"))
	canvasContext.Call("clear", canvasContext.Get("DEPTH_BUFFER_BIT"))

	current_scene.OnDraw()
	RenderQuadTreeContainer.QuadsCount = 0

	ScreenDims := NewVector2f(float32(appRef.Width), float32(appRef.Height))
	ScreenRect := Rect{Position: ScreenDims.Scale(-0.5).Scale(1 / Cam.GetScale()).Add(Cam.GetPosition()), Size: ScreenDims.Scale(1 / Cam.GetScale())}

	for _, v := range RenderQuadTreeContainer.Search(ScreenRect).Data {
		t := v.First
		v.Second.objectType(&Shapes, &Sprites, t.Position, t.Dimensions, t.UV1, t.UV2, t.Tint, t.Rotation, t.Z, v.Second.texture)
		RenderQuadTreeContainer.QuadsCount++
	}
	dynmaicQuadInView := DynamicRenderQuadTreeContainer.Search(ScreenRect)

	for _, v := range dynmaicQuadInView.AllItems() {
		// chai.Shapes.DrawFillRect(v.First.Position, v.First.Dimensions, v.Second.Tint)
		// rects_count++
		it := v.item
		t := GetComponentPtr[VisualTransform](current_scene, it.entId)
		// t := current_scene.transforms.Get(it.entId)
		// v.objectType(&Shapes, t.Position, t.Dimensions, v.tint, t.Rotation)
		// Shapes.DrawFillRectRotated(t.Position, t.Dimensions, it.tint, t.Rotation)
		v.item.objectType(&Shapes, &Sprites, t.Position, t.Dimensions, t.UV1, t.UV2, t.Tint, t.Rotation, t.Z, v.item.texture)

	}
	tempDraw(deltaTime)
	Sprites.Render()
	Shapes.Render(&Cam)
	// Renderer.End()
	// Renderer.Render()
	UIShapes.Render(&uiCam)
	UISprites.Render()
	return nil
}

func setBackgroundColor(_color RGBA8) {
	canvasContext.Call("clearColor", _color.GetColorRFloat32(), _color.GetColorGFloat32(), _color.GetColorBFloat32(), 1.0)

}

func NumOfQuadsInView() int {
	return RenderQuadTreeContainer.QuadsInViewCount()
}

func GetStaticQuadsInRect(rArea Rect) List[*Pair[VisualTransform, RenderObject]] {
	list := RenderQuadTreeContainer.Search(rArea)
	return list
}

func GetDynamicQuadsInRect(rArea Rect) List[*QuadTreeItem[RenderObject]] {
	list := DynamicRenderQuadTreeContainer.Search(rArea)
	return list
}
