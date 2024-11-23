package chai

import "syscall/js"

func addJsEventListener(_eventName string, _eventFunc func(this js.Value, args []js.Value) any) {
	js.Global().Call("addEventListener", _eventName, js.FuncOf(func(this js.Value, args []js.Value) any {
		return _eventFunc(this, args)
	}))
}
