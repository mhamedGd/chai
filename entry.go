package chai

import (
	"strings"
	"syscall/js"
)

var running bool
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
var tempStart func()
var tempUpdate func(float32)
var tempDraw func(float32)

/*
USING THE EventFunc[T] type ------- (1)

	var custom_func EventFunc[string] = func(x ...string) {
		fmt.Printf(x[0] + "\n")
	}

USING THE EventFunc[T] type ------- (1)
*/

// *** Declaring an ChaiEvent[T] *** var event ChaiEvent[int] ------- (2)

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

var uiCam Camera2D
var UIShapes ShapeBatch
var UISprites SpriteBatch

var started bool = false

var physics_world PhysicsWorld

var MouseCanvasPos Vector2f
var canvasBoundingClientRect js.Value

var mousePressed MouseButton
var numOfFingersTouching uint8
var LeftMouseJustPressed ChaiEvent[int]

func GetPhysicsWorld() *PhysicsWorld {
	return &physics_world
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
}

func Run(_app *App) {
	running = true

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

	tempStart = _app.OnStart
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

	initTextures()
	InitInputs()
	physics_world = newPhysicsWorld(NewVector2f(0.0, -200.0))
	// physics_world.box2dWorld.SetContactListener(worldContactListener)

	//js.Global().Set("js_start", js.FuncOf(JSStart))
	js.Global().Set("js_update", js.FuncOf(JSUpdate))
	js.Global().Set("js_draw", js.FuncOf(JSDraw))

	// if I put it above the "js_start" then it would take a lot of time to run
	Cam.Init(*_app)
	Cam.centerOffset = NewVector2f(float32(_app.Width)/2.0, float32(appRef.Height)/2.0)
	Cam.Update(*_app)

	uiCam.Init(*_app)
	//uiCam.position.AddXY(-float32(appRef.Width)/2.0, -float32(appRef.Height)/2.0)
	uiCam.Update(*_app)

	Shapes.Init()
	Assert(Shapes.Initialized, "Shapes Rendering was not initialized successfully")
	UIShapes.Init()
	Assert(UIShapes.Initialized, "Shapes Rendering was not initialized successfully")

	Sprites.Init("")
	UISprites.Init("")

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
		Using the ChaiEvent[T] ------- (2)

			event.AddListener(print_s)
			event.Invoke(1, 20003)
			event.RemoveListener(print_s)
	*/

	/*
		////////// FINAL CHECKS ////////////
	*/

	//custom_func("STRING") ------- (1)
	_app.OnStart()
	Assert(current_scene != nil, "Current Scene is none")
	if !started {
		started = true
	}
	select {}
}

/*
--------- (2)

	func print_s(s ...int) {
		fmt.Println(s[1])
	}

--------- (2)
*/

//	func JSStart(this js.Value, inputs []js.Value) interface{} {
//		tempStart()
//		return nil
//	}

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
	// LogF("Canvas Size: %v", GetCurrentCanvasPageSize())

	deltaTime = float32(inputs[0].Float())
	if deltaTime > CAP_DELTA_TIME {
		deltaTime = CAP_DELTA_TIME
	}

	currentWidth = canvas.Get("width").Int()
	currentHeight = canvas.Get("height").Int()
	current_scene.OnUpdate(deltaTime)
	updateInput()
	Cam.Update(*appRef)
	uiCam.Update(*appRef)
	tempUpdate(deltaTime)
	ElapsedTime += deltaTime

	timeAccumulation += deltaTime
	if timeAccumulation > (MAX_FIXED_CYCLES_PER_FRAME * FIXED_UPDATE_INTERVAL) {
		timeAccumulation = FIXED_UPDATE_INTERVAL
	}

	for timeAccumulation >= FIXED_UPDATE_INTERVAL {
		timeAccumulation -= FIXED_UPDATE_INTERVAL
		// physics_world.box2dWorld.ClearForces()
		physics_world.cpSpace.Step(float64(1 / 60.0))
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

	//Shapes.DrawLine(NewVector2f(0.0, 0.0), NewVector2f(2.5, 0.5), RGBA8{255, 255, 0, 255})
	current_scene.OnDraw()
	tempDraw(deltaTime)
	Sprites.Render(&Cam)
	Shapes.Render(&Cam)
	UISprites.Render(&uiCam)
	UIShapes.Render(&uiCam)
	return nil
}

func setBackgroundColor(_color RGBA8) {
	canvasContext.Call("clearColor", _color.GetColorRFloat32(), _color.GetColorGFloat32(), _color.GetColorBFloat32(), 1.0)

}
