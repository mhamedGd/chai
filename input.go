package chai

import (
	"syscall/js"

	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

var inputs_map customtypes.Map[string, ChaiInput]

var current_frame_pressed_inputs map[string]ChaiInput
var prev_frame_pressed_inputs map[string]ChaiInput

type ChaiInput struct {
	Name             string
	CorrespondingKey KeyCode
	ActionStrength   float32
	IsPressed        bool
}

func initInputs() {
	inputs_map = customtypes.NewMap[string, ChaiInput]()

	current_frame_pressed_inputs = make(map[string]ChaiInput)
	prev_frame_pressed_inputs = make(map[string]ChaiInput)
}

func updateInput() {
	for key, val := range inputs_map.AllItems() {
		if val.IsPressed {
			_, curr_ok := current_frame_pressed_inputs[key]
			if !curr_ok {
				current_frame_pressed_inputs[key] = val
			} else {
				_, prev_ok := prev_frame_pressed_inputs[key]
				if !prev_ok {
					prev_frame_pressed_inputs[key] = current_frame_pressed_inputs[key]
				}
			}
		} else {
			_, curr_ok := current_frame_pressed_inputs[key]
			if !curr_ok {
				delete(prev_frame_pressed_inputs, key)
			}
			delete(current_frame_pressed_inputs, key)
		}
	}
	isPreviousMousePressed = isCurrentMousePressed
	previousNumberOfFingersTouching = currentNumberOfFingersTouching
}

func BindInput(_inputName string, _corrKey KeyCode) {
	if inputs_map.Has(_inputName) {
		return
	}

	inputs_map.Set(_inputName, ChaiInput{
		Name:             _inputName,
		CorrespondingKey: _corrKey,
		ActionStrength:   0.0,
		IsPressed:        false,
	})

	addEventListenerWindow(JS_KEYDOWN, func(ae *AppEvent) {
		inp := inputs_map.Get(_inputName)
		if ae.Key == inp.CorrespondingKey {

			inp.ActionStrength = 1.0
			inp.IsPressed = true

			inputs_map.Set(_inputName, inp)
		}
	})

	addEventListenerWindow(JS_KEYUP, func(ae *AppEvent) {
		inp := inputs_map.Get(_inputName)
		if ae.Key == inp.CorrespondingKey {

			inp.ActionStrength = 0.0
			inp.IsPressed = false

			inputs_map.Set(_inputName, inp)
		}
	})

}

func ChangeInputBinding(_inputName string, _newBinding KeyCode) {
	if inputs_map.Has(_inputName) {
		inp := inputs_map.Get(_inputName)
		inp.CorrespondingKey = _newBinding
		inputs_map.Set(_inputName, inp)
	}
}

func GetActionStrength(_inputName string) float32 {
	return inputs_map.Get(_inputName).ActionStrength
}
func IsPressed(_inputName string) bool {
	return inputs_map.Get(_inputName).IsPressed
}

func IsJustPressed(_inputName string) bool {
	_, curr_ok := current_frame_pressed_inputs[_inputName]
	_, prev_ok := prev_frame_pressed_inputs[_inputName]
	return curr_ok && !prev_ok
}

func IsJustReleased(_inputName string) bool {
	_, curr_ok := current_frame_pressed_inputs[_inputName]
	_, prev_ok := prev_frame_pressed_inputs[_inputName]
	return !curr_ok && prev_ok
}

var isCurrentMousePressed bool
var isPreviousMousePressed bool

func IsMousePressed(_mouseButton MouseButton) bool {
	return mousePressed == _mouseButton
}

func IsMouseJustPressed() bool {
	return mousePressed != MouseButtonNull && isCurrentMousePressed != isPreviousMousePressed
}

func IsMouseJustReleased() bool {
	return mousePressed == MouseButtonNull && isCurrentMousePressed != isPreviousMousePressed
}

func onMousePressed() {
	isCurrentMousePressed = true
}
func onMouseReleased() {
	isCurrentMousePressed = false

}

func GetMousePosition(_evt js.Value) Vector2f {
	rect := canvas.Call("getBoundingClientRect")
	return NewVector2f((float32(_evt.Get("clientX").Int())-float32(rect.Get("left").Int()))/float32(rect.Get("width").Int())*float32(canvas.Get("width").Int()),
		float32(canvas.Get("height").Int())-(float32(_evt.Get("clientY").Int())-float32(rect.Get("top").Int()))/float32(rect.Get("height").Int())*float32(canvas.Get("height").Int()))
}

var currentNumberOfFingersTouching uint8
var previousNumberOfFingersTouching uint8

func GetNumberOfFingersTouching() uint8 {
	return numOfFingersTouching
}

func AreFingersTouching(_numOfFingers uint8) bool {
	return numOfFingersTouching >= _numOfFingers
}

func IsJustTouched(_numOfFingers uint8) bool {
	return _numOfFingers == currentNumberOfFingersTouching && currentNumberOfFingersTouching != previousNumberOfFingersTouching
}

func IsJustTouchReleased(_numOfFingers uint8) bool {
	return _numOfFingers != currentNumberOfFingersTouching && currentNumberOfFingersTouching != previousNumberOfFingersTouching
}

func onTouchStart(_numOfFingers uint8) {
	currentNumberOfFingersTouching = _numOfFingers
}

func onTouchEnd(_numOfFingers uint8) {
	currentNumberOfFingersTouching = _numOfFingers
}

/*
##############################################################################
####################									######################
#################### JS EVENTS - JS EVENTS - JS EVENTS	######################
####################									######################
##############################################################################
*/

func addEventListenerWindow(_eventType JsEventType, _callback func(*AppEvent)) {

	eventListener := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_callback(parseJSEvent(args[0]))

		//ae := parseJSEvent(m_Event)
		return nil
	})

	js.Global().Call("addEventListener", js.ValueOf(_eventType), eventListener)

}

type jsEvent = js.Value
type JsEventType = string
type MouseButton = int
type KeyCode = string

// ------------ EVENT TYPE ---------------------
const JS_KEYDOWN JsEventType = "keydown"
const JS_KEYUP JsEventType = "keyup"
const JS_MOUSEDOWN JsEventType = "mousedown"
const JS_MOUSEUP JsEventType = "mouseup"
const JS_MOUSEMOVED JsEventType = "mousemove"
const JS_TOUCHMOVED JsEventType = "touchmove"
const JS_TOUCHSTART JsEventType = "touchstart"
const JS_TOUCHEND JsEventType = "touchend"

// ---------------------------------------------

const LEFT_MOUSE_BUTTON = 0
const MIDLE_MOUSE_BUTTON = 1
const RIGHT_MOUSE_BUTTON = 2

const MouseButtonNull MouseButton = -1
const KeyNull string = ""
const CodeNull string = ""

type Gamepad struct {
	Start, Select                         bool
	South, East, West, North              bool
	DPADDown, DPADRight, DPADLeft, DPADUp bool
	LeftTriggerA, RightTriggerA           bool
	LeftTriggerB, RightTriggerB           float32
	LeftJoystick, RightJoystick           Vector2f
}

var Gamepads customtypes.List[Gamepad]

type AppEvent struct {
	// --------------------------
	m_Event jsEvent
	Type    JsEventType
	// FOR KEYBOARD EVENT
	Code KeyCode
	Key  string
	// --------------------------
	// FOR MOUSE EVENT
	OffsetX int
	OffsetY int
	Button  MouseButton
	// --------------------------
	// FOR MOUSE SCROLLING
	DeltaX, DeltaY, DeltaZ float64
	NUM_FINGERS            uint8
	/* --------------------------
	AltKey   bool
	CtrlKey  bool
	ShiftKey bool
	/ -------------------------- */
}

func (ap *AppEvent) GetJsEvent() js.Value {
	return js.Value(ap.m_Event)
}

func (e *AppEvent) PreventDefault() {
	e.m_Event.Call("preventDefualt")
}
func (e *AppEvent) StopPropagation() {
	e.m_Event.Call("stopPropagation")
}

func parseJSEvent(_event jsEvent) *AppEvent {
	var eventType JsEventType = _event.Get("type").String()
	switch eventType {
	case JS_KEYDOWN, JS_KEYUP:
		return &AppEvent{
			m_Event: _event,
			Type:    eventType,
			Code:    _event.Get("keycode").String(),
			Key:     _event.Get("code").String(),
			OffsetX: 0,
			OffsetY: 0,
			Button:  MouseButtonNull,
		}
	case JS_MOUSEMOVED:
		return &AppEvent{
			m_Event: _event,
			Type:    eventType,
			Code:    CodeNull,
			Key:     KeyNull,
			OffsetX: _event.Get("offsetX").Int(),
			OffsetY: _event.Get("offsetY").Int(),
		}
	case JS_MOUSEDOWN, JS_MOUSEUP:
		return &AppEvent{
			m_Event: _event,
			Type:    eventType,
			Code:    CodeNull,
			Key:     KeyNull,
			OffsetX: 0,
			OffsetY: 0,
			Button:  _event.Get("button").Int(),
		}
	case JS_TOUCHSTART:
		return &AppEvent{
			m_Event:     _event,
			Type:        eventType,
			Code:        CodeNull,
			Key:         KeyNull,
			OffsetX:     0,
			OffsetY:     0,
			Button:      MouseButtonNull,
			NUM_FINGERS: numOfFingersTouching + 1,
		}
	case JS_TOUCHEND:
		return &AppEvent{
			m_Event:     _event,
			Type:        eventType,
			Code:        CodeNull,
			Key:         KeyNull,
			OffsetX:     0,
			OffsetY:     0,
			Button:      MouseButtonNull,
			NUM_FINGERS: numOfFingersTouching - 1,
		}
	case JS_TOUCHMOVED:
		return &AppEvent{
			m_Event: _event,
			Type:    eventType,
			Code:    CodeNull,
			Key:     KeyNull,
			OffsetX: 0,
			OffsetY: 0,
		}
	}

	return &AppEvent{
		m_Event: _event,
		Code:    CodeNull,
		Key:     KeyNull,
		OffsetX: 0,
		OffsetY: 0,
		Button:  MouseButtonNull,
	}
}

/*
################### CUSTOM KEYCODES - CUSTOM KEYCODES ########################
##############################################################################
*/

const (
	KEY_BACKSPACE    KeyCode = "Backspace"
	KEY_TAP          KeyCode = "Tap"
	KEY_ENTER        KeyCode = "Enter"
	KEY_SHIFTLEFT    KeyCode = "ShiftLeft"
	KEY_SHIFTRIGHT   KeyCode = "ShiftRight"
	KEY_CONTROLLEFT  KeyCode = "ControlLeft"
	KEY_CONTROLRIGHT KeyCode = "ControlRight"
	KEY_ALTLEFT      KeyCode = "AltLeft"
	KEY_ALTRIGHT     KeyCode = "AltRight"
	KEY_PAUSE        KeyCode = "Pause"
	KEY_CAPSLOCK     KeyCode = "CapsLock"
	KEY_ESCAPE       KeyCode = "Escape"
	KEY_SPACE        KeyCode = "Space"
	KEY_PAGEUP       KeyCode = "PageUp"
	KEY_PAGEDOWN     KeyCode = "PageDown"
	KEY_END          KeyCode = "End"
	KEY_HOME         KeyCode = "Home"
	KEY_ARROWLEFT    KeyCode = "ArrowLeft"
	KEY_ARROWUP      KeyCode = "ArrowUp"
	KEY_ARROWRIGHT   KeyCode = "ArrowRight"
	KEY_ARROWDOWN    KeyCode = "ArrowDown"
	KEY_PRINTSCREEN  KeyCode = "PrintScreen"
	KEY_INSERT       KeyCode = "Insert"
	KEY_DELETE       KeyCode = "Delete"
	KEY_0            KeyCode = "Digit0"
	KEY_1            KeyCode = "Digit1"
	KEY_2            KeyCode = "Digit2"
	KEY_3            KeyCode = "Digit3"
	KEY_4            KeyCode = "Digit4"
	KEY_5            KeyCode = "Digit5"
	KEY_6            KeyCode = "Digit6"
	KEY_7            KeyCode = "Digit7"
	KEY_8            KeyCode = "Digit8"
	KEY_9            KeyCode = "Digit9"
	KEY_A            KeyCode = "KeyA"
	KEY_B            KeyCode = "KeyB"
	KEY_C            KeyCode = "KeyC"
	KEY_D            KeyCode = "KeyD"
	KEY_E            KeyCode = "KeyE"
	KEY_F            KeyCode = "KeyF"
	KEY_G            KeyCode = "KeyG"
	KEY_H            KeyCode = "KeyH"
	KEY_I            KeyCode = "KeyI"
	KEY_J            KeyCode = "KeyJ"
	KEY_K            KeyCode = "KeyK"
	KEY_L            KeyCode = "KeyL"
	KEY_M            KeyCode = "KeyM"
	KEY_N            KeyCode = "KeyN"
	KEY_O            KeyCode = "KeyO"
	KEY_P            KeyCode = "KeyP"
	KEY_Q            KeyCode = "KeyQ"
	KEY_R            KeyCode = "KeyR"
	KEY_S            KeyCode = "KeyS"
	KEY_T            KeyCode = "KeyT"
	KEY_U            KeyCode = "KeyU"
	KEY_V            KeyCode = "KeyV"
	KEY_W            KeyCode = "KeyW"
	KEY_X            KeyCode = "KeyX"
	KEY_Y            KeyCode = "KeyY"
	KEY_Z            KeyCode = "KeyZ"
	KEY_SUPERLEFT    KeyCode = "MetaLeft"
	KEY_SUPERRIGHT   KeyCode = "MetaRight"
	KEY_SELECT       KeyCode = "ContextMenu"
	KEY_NUM0         KeyCode = "Numpad0"
	KEY_NUM1         KeyCode = "Numpad1"
	KEY_NUM2         KeyCode = "Numpad2"
	KEY_NUM3         KeyCode = "Numpad3"
	KEY_NUM4         KeyCode = "Numpad4"
	KEY_NUM5         KeyCode = "Numpad5"
	KEY_NUM6         KeyCode = "Numpad6"
	KEY_NUM7         KeyCode = "Numpad7"
	KEY_NUM8         KeyCode = "Numpad8"
	KEY_NUM9         KeyCode = "Numpad9"
	KEY_NUMMULTIPLY  KeyCode = "NumpadMultiply"
	KEY_NUMADD       KeyCode = "NumpadAdd"
	KEY_NUMSUBTRACT  KeyCode = "NumpadSubtract"
	KEY_NUMDOT       KeyCode = "NumpadDecimal"
	KEY_NUMDIVIDE    KeyCode = "NumpadDivide"
	KEY_F1           KeyCode = "F1"
	KEY_F2           KeyCode = "F2"
	KEY_F3           KeyCode = "F3"
	KEY_F4           KeyCode = "F4"
	KEY_F5           KeyCode = "F5"
	KEY_F6           KeyCode = "F6"
	KEY_F7           KeyCode = "F7"
	KEY_F8           KeyCode = "F8"
	KEY_F9           KeyCode = "F9"
	KEY_F10          KeyCode = "F10"
	KEY_F11          KeyCode = "F11"
	KEY_F12          KeyCode = "F12"
	KEY_NUMLOCK      KeyCode = "NumLock"
	KEY_SCROLLLOCK   KeyCode = "ScrollLock"
	KEY_SEMICOLOR    KeyCode = "Semicolon"
	KEY_EQUAL        KeyCode = "Equal"
	KEY_COMMA        KeyCode = "Comma"
	KEY_MINUS        KeyCode = "Minus"
	KEY_PERIOD       KeyCode = "Period"
	KEY_SLASH        KeyCode = "Slash"
	KEY_BACKQUOTE    KeyCode = "Backquote"
	KEY_BRACKETLEFT  KeyCode = "BracketLeft"
	KEY_BRACKETRIGHT KeyCode = "BracketRight"
	KEY_BACKSLASH    KeyCode = "Backslash"
	KEY_QUOTE        KeyCode = "Quote"
)

/*
##############################################################################
##############################################################################
*/
