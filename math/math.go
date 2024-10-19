package math

import (
	"math"
	"math/rand"
	"unsafe"
)

const (
	PI      = math.Pi
	Rad2Deg = 180.0 / PI
	Deg2Rad = 1 / Rad2Deg
)

func RandomRangeFloat32(_min, _max float32) float32 {
	return rand.Float32()*(_max-_min) + _min
}

func ClampFloat32(_value float32, _min float32, _max float32) float32 {
	if _value > _max {
		_value = _max
	} else if _value < _min {
		_value = _min
	}

	return _value
}

func ClampFloat64(_value float64, _min float64, _max float64) float64 {
	if _value > _max {
		_value = _max
	} else if _value < _min {
		_value = _min
	}

	return _value
}

func ClampInt(_value, _min, _max int) int {
	_value = ((_value - _min) & -(BoolToInt(_value >= _min))) + _min
	return ((_value - _max) & -(BoolToInt(_value <= _max))) + _max
}

func RemapFloat32(_original, _origial_lowest, _original_highest, _new_lowest, _new_highest float32) float32 {
	return _new_lowest + (_new_highest-_new_lowest)*(_original/(_original_highest-_origial_lowest))
}

func LerpFloat32(_a, _b, _t float32) float32 {
	return _a + (_b-_a)*_t
}

func LerpRot(_a, _b, _t float32) float32 {
	_a = NormalizeAngle(_a)
	_b = NormalizeAngle(_b)
	delta := _b - _a

	if delta > 180.0 {
		delta -= 360
	} else if delta < -180.0 {
		delta += 360
	}
	return _a + delta*_t
}

func NormalizeAngle(_a float32) float32 {
	// Normalize angle to [-180, 180]
	// _a = _a % 360
	_a = float32(math.Mod(float64(_a), 360))
	if _a > 180 {
		_a -= 360
	} else if _a < -180 {
		_a += 360
	}
	return _a
}

func LerpFloat64(_a, _b, _t float64) float64 {
	return _a + (_b-_a)*_t
}

func LerpInt(_a, _b int, _t float32) int {
	return _a + (_b-_a)*int(_t*float32(_b-_a))/(_b-_a)
}

func LerpUint8(_a, _b uint8, _t float32) uint8 {
	return _a + (_b-_a)*uint8(_t*255)
}

func Float1ToUint8_255(_inp float32) uint8 {
	output := (_inp * 255.0)
	if output >= 255 {
		output = 255
	} else if output <= 0.0 {
		output = 0
	}
	return uint8(output)
}

func Uint8ToFloat1(_inp uint8) float32 {
	return float32(_inp) / 255.0
}

func LerpVector2f(_a, _b Vector2f, _t float32) Vector2f {
	return Vector2f{
		X: _a.X + (_b.X-_a.X)*_t,
		Y: _a.Y + (_b.Y-_a.Y)*_t,
	}
}

func BoolToFloat64(_boolean bool) float64 {
	return float64(*(*byte)(unsafe.Pointer(&_boolean)))
}

func BoolToFloat32(_boolean bool) float32 {
	return float32(*(*byte)(unsafe.Pointer(&_boolean)))
}

func BoolToInt(_boolean bool) int {
	return int(*(*byte)(unsafe.Pointer(&_boolean)))
}

func MinFloat32(_v1, _v2 float32) float32 {
	if _v1 <= _v2 {
		return _v1
	}

	return _v2
}

func MinFloat64(_v1, _v2 float64) float64 {
	if _v1 <= _v2 {
		return _v1
	}

	return _v2
}

func MinInt(_v1, _v2 int) int {
	if _v1 <= _v2 {
		return _v1
	}

	return _v2
}

func MaxFloat32(_v1, _v2 float32) float32 {
	if _v1 >= _v2 {
		return _v1
	}

	return _v2
}

func MaxFloat64(_v1, _v2 float64) float64 {
	if _v1 >= _v2 {
		return _v1
	}

	return _v2
}

func MaxInt(_v1, _v2 int) int {
	if _v1 >= _v2 {
		return _v1
	}

	return _v2
}

func MaxUint8(_v1, _v2 uint8) uint8 {
	if _v1 >= _v2 {
		return _v1
	}
	return _v2
}

func AbsFloat32(_v float32) float32 {
	return SignFloat32(_v) * _v
}

func AbsInt(_v int) int {
	return SignInt(_v) * _v
}

func SignFloat32(_v float32) float32 {
	switch {
	case _v < 0.0:
		return -1.0
	case _v > 0.0:
		return 1.0
	}

	return 0.0
}

func SignInt(_v int) int {
	switch {
	case _v < 0:
		return -1
	case _v > 0:
		return 1
	}

	return 0
}

func MapFloat32(value, low1, high1, low2, high2 float32) float32 {
	return low2 + (value-low1)*(high2-low2)/(high1-low1)
}
