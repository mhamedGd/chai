package chai

import "syscall/js"

var inputs_map map[string]ChaiInput

var current_frame_pressed_inputs map[string]ChaiInput
var prev_frame_pressed_inputs map[string]ChaiInput

type ChaiInput struct {
	Name             string
	CorrespondingKey KeyCode
	ActionStrength   float32
	IsPressed        bool
}

func InitInputs() {
	inputs_map = make(map[string]ChaiInput)

	current_frame_pressed_inputs = make(map[string]ChaiInput)
	prev_frame_pressed_inputs = make(map[string]ChaiInput)
}

func updateInput() {
	for key, val := range inputs_map {
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
}

func BindInput(_input_name string, _corr_key KeyCode) {
	_, ok := inputs_map[_input_name]
	if ok {
		return
	}

	inputs_map[_input_name] = ChaiInput{
		Name:             _input_name,
		CorrespondingKey: _corr_key,
		ActionStrength:   0.0,
		IsPressed:        false,
	}

	addEventListenerWindow(JS_KEYDOWN, func(ae *AppEvent) {
		inp := inputs_map[_input_name]
		if ae.Key == inp.CorrespondingKey {

			inp.ActionStrength = 1.0
			inp.IsPressed = true

			inputs_map[_input_name] = inp
		}
	})

	addEventListenerWindow(JS_KEYUP, func(ae *AppEvent) {
		inp := inputs_map[_input_name]
		if ae.Key == inp.CorrespondingKey {
			inp.ActionStrength = 0.0
			inp.IsPressed = false
			inputs_map[_input_name] = inp
		}
	})

}

func ChangeInputBinding(_input_name string, _new_binding KeyCode) {
	inp, ok := inputs_map[_input_name]
	if ok {
		inp.CorrespondingKey = _new_binding
		inputs_map[_input_name] = inp
	}
}

func GetActionStrength(_input_name string) float32 {
	inp := inputs_map[_input_name]
	return inp.ActionStrength
}
func IsPressed(_input_name string) bool {

	return inputs_map[_input_name].IsPressed
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

func GetMousePosition(evt js.Value) Vector2f {
	rect := canvas.Call("getBoundingClientRect")
	return NewVector2f((float32(evt.Get("clientX").Int())-float32(rect.Get("left").Int()))/float32(rect.Get("width").Int())*float32(canvas.Get("width").Int()),
		float32(canvas.Get("height").Int())-(float32(evt.Get("clientY").Int())-float32(rect.Get("top").Int()))/float32(rect.Get("height").Int())*float32(canvas.Get("height").Int()))
}
