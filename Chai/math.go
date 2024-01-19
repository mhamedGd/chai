package chai

import (
	"fmt"
	"math"
)

const (
	PI = math.Pi
)

func Deg2Rad(_degrees float32) float32 {
	return _degrees * PI / 180.0
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

func LerpFloat64(_a, _b, _t float64) float64 {
	return _a + (_b-_a)*_t
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
	if _boolean {
		return 1.0
	}

	return 0.0
}

func BoolToFloat32(_boolean bool) float32 {
	if _boolean {
		return 1.0
	}

	return 0.0
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

/*
############## ABS - ABS - ABS - ABS ##############################################
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

var Vector2fZero Vector2f = Vector2f{0.0, 0.0}
var Vector2fOne Vector2f = Vector2f{1.0, 1.0}
var Vector2fRight Vector2f = Vector2f{1.0, 0.0}
var Vector2fUp Vector2f = Vector2f{0.0, 1.0}

func (v1 Vector2f) Equal(v2 *Vector2f) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

func (v1 Vector2f) NearlyEqual(v2 Vector2f) bool {
	var factor float32 = 0.001

	diff := v1.Subtract(v2)
	diff = AbsVector2f(&diff)

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
func (v1 Vector2f) Div(v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X / v2.X, Y: v1.Y / v2.Y}
}

func (v Vector2f) Scale(_value float32) Vector2f {
	return Vector2f{X: v.X * _value, Y: v.Y * _value}
}

func AbsVector2f(_v *Vector2f) Vector2f {
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

func DotProduct(v1, v2 Vector2f) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v Vector2f) Perpendicular() Vector2f {
	return Vector2f{-v.Y, v.X}
}

func (v *Vector2f) Angle() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
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
