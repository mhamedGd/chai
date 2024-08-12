package chai

import (
	"fmt"
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

/*
###################################################################################
############## CLAMP - CLAMP - CLAMP ##############################################
*/
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

/*
############## CLAMP - CLAMP - CLAMP ##############################################
###################################################################################
*/

/*
###################################################################################
############## LERP - LERP - LERP #################################################
*/

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

/*
############## LERP - LERP - LERP #################################################
###################################################################################
*/

/*
###################################################################################
############## CONVERSION - CONVERSION - CONVERSION ###############################
*/

func BoolToFloat64(_boolean bool) float64 {
	return float64(*(*byte)(unsafe.Pointer(&_boolean)))
}

func BoolToFloat32(_boolean bool) float32 {
	return float32(*(*byte)(unsafe.Pointer(&_boolean)))
}

func BoolToInt(_boolean bool) int {
	return int(*(*byte)(unsafe.Pointer(&_boolean)))
}

/*
############## CONVERSION - CONVERSION - CONVERSION ###############################
###################################################################################
*/

/*
###################################################################################
############## MIN MAX - MIN MAX - MIN MAX ########################################
*/

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

// ################################################################################

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

/*
############## MIN MAX - MIN MAX - MIN MAX ########################################
###################################################################################
*/

/*
###################################################################################
############## ABS - ABS - ABS - ABS ##############################################
*/

func AbsFloat32(_v float32) float32 {
	return SignFloat32(_v) * _v
}

func AbsInt(_v int) int {
	return SignInt(_v) * _v
}

/*
############## ABS - ABS - ABS - ABS ##############################################
###################################################################################
*/

/*
###################################################################################
############## SIGN - SIGN - SIGN #################################################
*/

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

/*
############## ABS - ABS - ABS - ABS ##############################################
###################################################################################
*/

/*
###################################################################################
############## MAP - MAP - MAP - MAP ##############################################
*/

func MapFloat32(value, low1, high1, low2, high2 float32) float32 {
	return low2 + (value-low1)*(high2-low2)/(high1-low1)
}

/*
############## MAP - MAP - MAP - MAP ##############################################
###################################################################################
*/

/*
###################################################################################
############## VECTOR2  - VECTOR2 #################################################
*/

type Vector2f struct {
	X float32
	Y float32
}

func NewVector2f(x, y float32) Vector2f {
	return Vector2f{
		x, y,
	}
}

var (
	Vector2fZero  Vector2f = Vector2f{0.0, 0.0}
	Vector2fOne   Vector2f = Vector2f{1.0, 1.0}
	Vector2fRight Vector2f = Vector2f{1.0, 0.0}
	Vector2fLeft  Vector2f = Vector2f{-1.0, 0.0}
	Vector2fUp    Vector2f = Vector2f{0.0, 1.0}
	Vector2fDown  Vector2f = Vector2f{0.0, -1.0}
)

func (v1 Vector2f) Equal(v2 Vector2f) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

func (v1 Vector2f) NearlyEqual(v2 Vector2f) bool {
	var factor float32 = 0.001

	diff := v1.Subtract(v2)
	diff = AbsVector2f(diff)

	if diff.X <= factor && diff.Y <= factor {
		return true
	}

	return false
}

func (v1 Vector2f) Add(v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v1 Vector2f) AddXY(x, y float32) Vector2f {
	return Vector2f{X: v1.X + x, Y: v1.Y + y}
}
func (v1 Vector2f) Subtract(v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v1 Vector2f) SubtractXY(x, y float32) Vector2f {
	return Vector2f{X: v1.X - x, Y: v1.Y - y}
}

func (v1 Vector2f) Multp(v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X * v2.X, Y: v1.Y * v2.Y}
}

func (v1 Vector2f) MultpXY(x, y float32) Vector2f {
	return Vector2f{X: v1.X * x, Y: v1.Y * y}
}

func (v1 Vector2f) Div(v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X / v2.X, Y: v1.Y / v2.Y}
}

func (v Vector2f) Scale(_value float32) Vector2f {
	return Vector2f{X: v.X * _value, Y: v.Y * _value}
}

func AbsVector2f(_v Vector2f) Vector2f {
	return Vector2f{
		AbsFloat32(_v.X), AbsFloat32(_v.Y),
	}
}

func (v *Vector2f) Length() float32 {
	return float32(math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y))))
}

func (v *Vector2f) LengthSquared() float32 {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v Vector2f) Normalize() Vector2f {
	leng := v.Length()
	return Vector2f{v.X / leng, v.Y / leng}
}

func (v1 Vector2f) Distance(v2 Vector2f) float32 {
	return float32(math.Sqrt(math.Pow(float64(v2.X-v1.X), 2) + math.Pow(float64(v2.Y-v1.Y), 2)))
}

func DotProduct(v1, v2 Vector2f) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v Vector2f) Perpendicular() Vector2f {
	return Vector2f{-v.Y, v.X}
}

func (v Vector2f) Angle() float32 {
	ang := float32(math.Atan2(float64(v.Y), float64(v.X)))
	return ang
}

func (v Vector2f) Rotate(_angle float32, _pivot Vector2f) Vector2f {
	anglePolar := _angle * math.Pi / 180.0
	x := v.X
	y := v.Y

	v.X = (x-_pivot.X)*float32(math.Cos(float64(anglePolar))) - (y-_pivot.Y)*float32(math.Sin(float64(anglePolar))) + _pivot.X
	v.Y = (x-_pivot.X)*float32(math.Sin(float64(anglePolar))) + (y-_pivot.Y)*float32(math.Cos(float64(anglePolar))) + _pivot.Y
	return Vector2f{
		v.X, v.Y,
	}
}

func (v Vector2f) RotateCenter(_angle float32) Vector2f {
	anglePolar := _angle * math.Pi / 180.0
	x := v.X
	y := v.Y

	v.X = (x)*float32(math.Cos(float64(anglePolar))) - (y)*float32(math.Sin(float64(anglePolar)))
	v.Y = (x)*float32(math.Sin(float64(anglePolar))) + (y)*float32(math.Cos(float64(anglePolar)))
	return Vector2f{
		v.X, v.Y,
	}
}

func Vector2fMidpoint(v1, v2 Vector2f) Vector2f {
	return v1.Add(v2).Scale(0.5)
}

func (v Vector2f) ToString() string {
	return fmt.Sprint(v.X, v.Y)
}

func RandVector2f() Vector2f {
	v := NewVector2f(0.0, 0.0)
	v.X = rand.Float32()*(1.0 - -1.0) + -1.0
	v.Y = rand.Float32()*(1.0 - -1.0) + -1.0

	return v
}

func RandPosVector2f() Vector2f {
	return NewVector2f(rand.Float32(), rand.Float32())
}

type Vector2i struct {
	X, Y int
}

func NewVector2i(x, y int) Vector2i {
	return Vector2i{
		x, y,
	}
}

func (v1 *Vector2i) Equal(v2 *Vector2i) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

func (v1 Vector2i) Add(v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v1 Vector2i) AddXY(x, y int) Vector2i {
	return Vector2i{X: v1.X + x, Y: v1.Y + y}
}
func (v1 Vector2i) Subtract(v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v1 Vector2i) SubtractXY(x, y int) Vector2i {
	return Vector2i{X: v1.X - x, Y: v1.Y - y}
}

func (v1 Vector2i) Multp(v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X * v2.X, Y: v1.Y * v2.Y}
}

func (v1 Vector2i) MultpXY(x, y int) Vector2i {
	return Vector2i{X: v1.X * x, Y: v1.Y * y}
}

func (v1 Vector2i) Div(v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X / v2.X, Y: v1.Y / v2.Y}
}

func (v Vector2i) Scale(_value int) Vector2i {
	return Vector2i{X: int(v.X * _value), Y: int(v.Y * _value)}
}

func AbsVector2i(_v *Vector2i) Vector2i {
	return Vector2i{
		AbsInt(_v.X), AbsInt(_v.Y),
	}
}

func (v *Vector2i) Length() int {
	return int(math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y))))
}

func (v *Vector2i) LengthSquared() int {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v Vector2i) Normalize() Vector2i {
	leng := v.Length()
	return Vector2i{v.X / leng, v.Y / leng}
}

func DotProductInt(v1, v2 Vector2i) int {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v *Vector2i) Perpendicular() Vector2i {
	return Vector2i{-v.Y, v.X}
}

func (v *Vector2i) Angle() int {
	return int(math.Atan2(float64(v.Y), float64(v.X)))
}

func (v Vector2i) Rotate(_angle float32, _pivot Vector2i) Vector2i {
	anglePolar := _angle * math.Pi / 180.0
	x := v.X
	y := v.Y

	v.X = (x-_pivot.X)*int(math.Cos(float64(anglePolar))) - (y-_pivot.Y)*int(math.Sin(float64(anglePolar))) + _pivot.X
	v.Y = (x-_pivot.X)*int(math.Sin(float64(anglePolar))) + (y-_pivot.Y)*int(math.Cos(float64(anglePolar))) + _pivot.Y
	return Vector2i{
		v.X, v.Y,
	}
}

func (v Vector2i) RotateCenter(_angle float32) Vector2i {
	anglePolar := _angle * math.Pi / 180.0
	x := v.X
	y := v.Y

	v.X = (x)*int(math.Cos(float64(anglePolar))) - (y)*int(math.Sin(float64(anglePolar)))
	v.Y = (x)*int(math.Sin(float64(anglePolar))) + (y)*int(math.Cos(float64(anglePolar)))
	return Vector2i{
		v.X, v.Y,
	}
}

func Vector2fMidpointInt(v1, v2 Vector2i) Vector2i {
	return v1.Add(v2).Scale(1 / 2)
}

func (v Vector2i) ToString() string {
	return fmt.Sprint(v.X, v.Y)
}

/*
############## VECTOR2  - VECTOR2 #################################################
###################################################################################
*/

func PointVsRect(_point Vector2f, lower_left, upper_right Vector2f) bool {
	return (_point.X >= lower_left.X && _point.Y >= lower_left.Y && _point.X <= upper_right.X && _point.Y <= upper_right.Y)
}
