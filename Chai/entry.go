package chai

import (
	"syscall/js"
)

var app_url string

type App struct {
	Width    int
	Height   int
	Title    string
	OnStart  func()
	OnUpdate func(float32)
	OnDraw   func()
	OnEvent  func(*AppEvent)
}

// Used to make the update function only available in the local App struct, to the whole file
var tempStart func()
var tempUpdate func(float32)
var tempDraw func()

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
var debug_console js.Value
var appRef *App

var Cam Camera2D
var Shapes ShapeBatch
var Sprites SpriteBatch

var started bool = false

func (_app *App) fillDefaults() {
	if _app.OnStart == nil {
		_app.OnStart = func() {

		}
	}
	if _app.OnUpdate == nil {
		_app.OnUpdate = func(f float32) {

		}
	}
	if _app.OnDraw == nil {
		_app.OnDraw = func() {

		}
	}
	if _app.OnEvent == nil {
		_app.OnEvent = func(ae *AppEvent) {

		}
	}
}

func Run(_app *App) {
	defer func() {
		if r := recover(); r != nil {
			ErrorF("PANICKED - %v", r)
		}
	}()

	appRef = _app
	_app.fillDefaults()
	app_url = js.Global().Get("location").Get("href").String()

	js.Global().Get("document").Set("title", _app.Title)

	canvas = js.Global().Get("document").Call("getElementById", "viewport")
	debug_console = js.Global().Get("document").Call("querySelector", ".debug-console")

	canvasContext = canvas.Call("getContext", "webgl2")
	Assert(!canvasContext.IsNull(), "CANVAS: Failed to Get Context")

	canvas.Set("width", _app.Width)
	canvas.Set("height", _app.Height)

	canvasContext.Call("blendFuncSeparate", canvasContext.Get("SRC_ALPHA"), canvasContext.Get("ONE_MINUS_SRC_ALPHA"), canvasContext.Get("ONE"), canvasContext.Get("ONE"))
	canvasContext.Call("enable", canvasContext.Get("BLEND"))

	tempStart = _app.OnStart
	tempUpdate = _app.OnUpdate
	tempDraw = _app.OnDraw

	InitInputs()

	//js.Global().Set("js_start", js.FuncOf(JSStart))
	js.Global().Set("js_update", js.FuncOf(JSUpdate))
	js.Global().Set("js_draw", js.FuncOf(JSDraw))

	js.Global().Set("js_dpad_up", js.FuncOf(JSDpadUp))
	js.Global().Set("js_dpad_down", js.FuncOf(JSDpadDown))
	js.Global().Set("js_dpad_left", js.FuncOf(JSDpadLeft))
	js.Global().Set("js_dpad_right", js.FuncOf(JSDpadRight))
	js.Global().Set("js_main_button", js.FuncOf(JSMainButton))
	js.Global().Set("js_side_button", js.FuncOf(JSSideButton))

	// if I put it above the "js_start" then it would take a lot of time to run
	Cam.Init(*_app)
	Cam.Update(*_app)

	Shapes.Init()
	Assert(Shapes.Initialized, "Shapes Rendering was not initialized successfully")

	Sprites.Init("")
	canvasContext.Call("viewport", 0, 0, appRef.Width, appRef.Height)

	addEventListenerWindow(JS_KEYUP, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_KEYDOWN, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEDOWN, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEUP, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEMOVED, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})

	/*
		Using the ChaiEvent[T] ------- (2)

			event.AddListener(print_s)
			event.Invoke(1, 20003)
			event.RemoveListener(print_s)
	*/

	//UnuseShader()
	//custom_func("STRING") ------- (1)
	_app.OnStart()
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
var deltaTime float32

const CAP_DELTA_TIME float32 = 50.0 / 1000.0

func JSUpdate(this js.Value, inputs []js.Value) interface{} {
	if !started {
		return nil
	}
	deltaTime = float32(inputs[0].Float())
	if deltaTime > CAP_DELTA_TIME {
		deltaTime = CAP_DELTA_TIME
	}
	currentWidth = canvas.Get("width").Int()
	currentHeight = canvas.Get("height").Int()
	tempUpdate(deltaTime)
	updateInput()
	Cam.Update(*appRef)
	return nil
}

func JSDraw(this js.Value, inputs []js.Value) interface{} {
	if !started {
		return nil
	}
	canvasContext.Call("viewport", 0, 0, currentWidth, currentHeight)
	canvasContext.Call("clearColor", 0.1, 0.0, 0.2, 1.0)
	canvasContext.Call("clear", canvasContext.Get("COLOR_BUFFER_BIT"))

	//Shapes.DrawLine(NewVector2f(0.0, 0.0), NewVector2f(2.5, 0.5), RGBA8{255, 255, 0, 255})
	tempDraw()
	Sprites.Render(&Cam)
	Shapes.Render(&Cam)
	return nil
}
