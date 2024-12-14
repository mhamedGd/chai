package chai

import (
	"syscall/js"

	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

type App struct {
	Width          int
	Height         int
	Title          string
	PixelsPerMeter int
	OnStart        func()
	OnUpdate       func(float32)
	OnDraw         func(float32)
	OnEvent        func(*AppEvent)
}

var pixelsPerMeter int
var pixelsPerMeterDimensions Vector2f

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

	if _app.PixelsPerMeter == 0 {
		_app.PixelsPerMeter = 1
	}

	setPhysicsFunctions(PHYSICS_ENGINE_BOX2D)
}

func appPresets(_app *App) {
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

	pixelsPerMeter = _app.PixelsPerMeter
	pixelsPerMeterDimensions = NewVector2f(float32(_app.PixelsPerMeter), float32(_app.PixelsPerMeter))

	initTextures()
	initInputs()

	js.Global().Set("js_update", js.FuncOf(jSUpdate))
	js.Global().Set("js_draw", js.FuncOf(jSDraw))
}

func modulesInitialization(_app *App) {
	physics_world = newPhysicsWorld(NewVector2f(0.0, -98.0))

	Cam.Init(*_app)
	Cam.m_CenterOffset = NewVector2f(float32(_app.Width)/2.0, float32(appRef.Height)/2.0)
	Cam.Update(*_app)

	uiCam.Init(*_app)
	uiCam.Update(*_app)

	RenderQuadTreeContainer = customtypes.NewStaticQuadTreeContainer[customtypes.Pair[VisualTransform, RenderObject]]()
	DynamicRenderQuadTreeContainer = customtypes.NewDynamicQuadTreeContainer[RenderObject]()
	RenderQuadTreeContainer.Resize(Rect{Position: Vector2fZero, Size: NewVector2f(1000.0, 1000.0)})
	DynamicRenderQuadTreeContainer.Resize(Rect{Position: Vector2fZero, Size: NewVector2f(10000.0, 10000.0)})

	Renderer = NewRenderer2D(5_000, _app, &Cam)
}

func eventsInitialization(_app *App) {
	mousePressed = MouseButtonNull
	LeftMouseJustPressed = customtypes.NewChaiEvent1[int]()

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

	addJsEventListener("gamepadconnected", func(this js.Value, args []js.Value) any {
		// event := args[0]
		Gamepads.PushBack(Gamepad{
			// LeftJoystick:  NewVector2f(float32(event.Get("axes").Index(0).Float()), float32(event.Get("axes").Index(1).Float())),
			// RightJoystick: NewVector2f(float32(event.Get("axes").Index(2).Float()), float32(event.Get("axes").Index(3).Float())),
		})
		return 0
	})

	// Axes: Left (0, 1), Right (2, 3)
	// JoystickLeft (10), JoystickRight(11)
	// Buttons: South(0), East(1), West(2), North(3)
	// Triggers: LeftTriggerA(4), RightTriggerA(5), LeftTriggerB(6), RightTriggerB(7)
	// Select(8), Start(9)
	// DPad: Up(12), Down(13), Left(14), Right(15)

}
