package chai

import (
	"reflect"
	"syscall/js"
)

type EventFunc[T any] func(...T)

type ChaiEvent[T any] struct {
	listeners []EventFunc[T]
}

func (e *ChaiEvent[T]) init() {
	e.listeners = make([]EventFunc[T], 0)
}

func (e *ChaiEvent[T]) AddListener(ef EventFunc[T]) {
	e.listeners = append(e.listeners, ef)
}

func (e *ChaiEvent[T]) RemoveListener(ef EventFunc[T]) {
	for i, fn := range e.listeners {
		f1 := reflect.ValueOf(fn)
		f2 := reflect.ValueOf(ef)
		if f1 == f2 {
			e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
		}
	}
}

func (e *ChaiEvent[T]) Invoke(x ...T) {
	for _, f := range e.listeners {
		//fmt.Println(index)
		f(x...)
	}
}

/*
##############################################################################
####################									######################
#################### JS EVENTS - JS EVENTS - JS EVENTS	######################
####################									######################
##############################################################################
*/

func addEventListenerWindow(eventType JsEventType, callback func(*AppEvent)) {

	eventListener := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(parseJSEvent(args[0]))

		//ae := parseJSEvent(event)
		return nil
	})

	js.Global().Call("addEventListener", js.ValueOf(eventType), eventListener)

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

// ---------------------------------------------

const LEFT_MOUSE_BUTTON = 0
const MIDLE_MOUSE_BUTTON = 1
const RIGHT_MOUSE_BUTTON = 2

const MouseButtonNull MouseButton = -1
const KeyNull string = ""
const CodeNull string = ""

type AppEvent struct {
	// --------------------------
	event jsEvent
	Type  JsEventType
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
	/* --------------------------
	AltKey   bool
	CtrlKey  bool
	ShiftKey bool
	/ -------------------------- */
}

func (ap *AppEvent) GetJsEvent() js.Value {
	return js.Value(ap.event)
}

func (e *AppEvent) PreventDefault() {
	e.event.Call("preventDefualt")
}
func (e *AppEvent) StopPropagation() {
	e.event.Call("stopPropagation")
}

func parseJSEvent(event jsEvent) *AppEvent {
	var eventType JsEventType = event.Get("type").String()
	switch eventType {
	case JS_KEYDOWN, JS_KEYUP:
		return &AppEvent{
			event:   event,
			Type:    eventType,
			Code:    event.Get("keycode").String(),
			Key:     event.Get("code").String(),
			OffsetX: 0,
			OffsetY: 0,
			Button:  MouseButtonNull,
		}
	case JS_MOUSEMOVED:
		return &AppEvent{
			event:   event,
			Type:    eventType,
			Code:    CodeNull,
			Key:     KeyNull,
			OffsetX: event.Get("offsetX").Int(),
			OffsetY: event.Get("offsetY").Int(),
		}
	case JS_MOUSEDOWN, JS_MOUSEUP:
		return &AppEvent{
			event:   event,
			Type:    eventType,
			Code:    CodeNull,
			Key:     KeyNull,
			OffsetX: 0,
			OffsetY: 0,
			Button:  event.Get("button").Int(),
		}
	}

	return &AppEvent{
		event:   event,
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
	KEY_SPCAE        KeyCode = "Space"
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
