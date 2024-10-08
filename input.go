package chai

import "syscall/js"

var inputs_map Map[string, ChaiInput]

var current_frame_pressed_inputs map[string]ChaiInput
var prev_frame_pressed_inputs map[string]ChaiInput

type ChaiInput struct {
	Name             string
	CorrespondingKey KeyCode
	ActionStrength   float32
	IsPressed        bool
}

func InitInputs() {
	// inputs_map = make(map[string]ChaiInput)
	inputs_map = NewMap[string, ChaiInput]()

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

func BindInput(_input_name string, _corr_key KeyCode) {
	if inputs_map.Has(_input_name) {
		return
	}

	inputs_map.Set(_input_name, ChaiInput{
		Name:             _input_name,
		CorrespondingKey: _corr_key,
		ActionStrength:   0.0,
		IsPressed:        false,
	})

	addEventListenerWindow(JS_KEYDOWN, func(ae *AppEvent) {
		inp := inputs_map.Get(_input_name)
		if ae.Key == inp.CorrespondingKey {

			inp.ActionStrength = 1.0
			inp.IsPressed = true

			// inputs_map[_input_name] = inp
			inputs_map.Set(_input_name, inp)
		}
	})

	addEventListenerWindow(JS_KEYUP, func(ae *AppEvent) {
		// inp := inputs_map[_input_name]
		// if ae.Key == inp.CorrespondingKey {
		// 	inp.ActionStrength = 0.0
		// 	inp.IsPressed = false
		// 	inputs_map[_input_name] = inp
		// }
		inp := inputs_map.Get(_input_name)
		if ae.Key == inp.CorrespondingKey {

			inp.ActionStrength = 0.0
			inp.IsPressed = false

			// inputs_map[_input_name] = inp
			inputs_map.Set(_input_name, inp)
		}
	})

}

func ChangeInputBinding(_input_name string, _new_binding KeyCode) {
	if inputs_map.Has(_input_name) {
		inp := inputs_map.Get(_input_name)
		inp.CorrespondingKey = _new_binding
		inputs_map.Set(_input_name, inp)
	}
}

func GetActionStrength(_input_name string) float32 {
	return inputs_map.Get(_input_name).ActionStrength
}
func IsPressed(_input_name string) bool {
	return inputs_map.Get(_input_name).IsPressed
}

func IsJustPressed(_input_name string) bool {
	_, curr_ok := current_frame_pressed_inputs[_input_name]
	_, prev_ok := prev_frame_pressed_inputs[_input_name]
	return curr_ok && !prev_ok
}

func IsJustReleased(_input_name string) bool {
	_, curr_ok := current_frame_pressed_inputs[_input_name]
	_, prev_ok := prev_frame_pressed_inputs[_input_name]
	return !curr_ok && prev_ok
}

var isCurrentMousePressed bool
var isPreviousMousePressed bool

func IsMousePressed(mouseButton MouseButton) bool {
	return mousePressed == mouseButton
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

func GetMousePosition(evt js.Value) Vector2f {
	rect := canvas.Call("getBoundingClientRect")
	return NewVector2f((float32(evt.Get("clientX").Int())-float32(rect.Get("left").Int()))/float32(rect.Get("width").Int())*float32(canvas.Get("width").Int()),
		float32(canvas.Get("height").Int())-(float32(evt.Get("clientY").Int())-float32(rect.Get("top").Int()))/float32(rect.Get("height").Int())*float32(canvas.Get("height").Int()))
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
