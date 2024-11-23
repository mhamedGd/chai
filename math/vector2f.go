package math

import (
	"fmt"
	"math"
	"math/rand"
)

type Vector2f struct {
	X float32
	Y float32
}

func NewVector2f(_x, _y float32) Vector2f {
	return Vector2f{
		_x, _y,
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

func (v1 Vector2f) Equal(_v2 Vector2f) bool {
	return v1.X == _v2.X && v1.Y == _v2.Y
}

func (v1 Vector2f) NearlyEqual(_v2 Vector2f) bool {
	var factor float32 = 0.001

	diff := v1.Subtract(_v2)
	diff = AbsVector2f(diff)

	if diff.X <= factor && diff.Y <= factor {
		return true
	}

	return false
}

func (v1 Vector2f) Add(_v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X + _v2.X, Y: v1.Y + _v2.Y}
}

func (v1 Vector2f) AddXY(_x, _y float32) Vector2f {
	return Vector2f{X: v1.X + _x, Y: v1.Y + _y}
}
func (v1 Vector2f) Subtract(_v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X - _v2.X, Y: v1.Y - _v2.Y}
}

func (v1 Vector2f) SubtractXY(_x, _y float32) Vector2f {
	return Vector2f{X: v1.X - _x, Y: v1.Y - _y}
}

func (v1 Vector2f) Multp(_v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X * _v2.X, Y: v1.Y * _v2.Y}
}

func (v1 Vector2f) MultpXY(_x, _y float32) Vector2f {
	return Vector2f{X: v1.X * _x, Y: v1.Y * _y}
}

func (v1 Vector2f) Div(_v2 Vector2f) Vector2f {
	return Vector2f{X: v1.X / _v2.X, Y: v1.Y / _v2.Y}
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

func (v1 Vector2f) Distance(_v2 Vector2f) float32 {
	return float32(math.Sqrt(math.Pow(float64(_v2.X-v1.X), 2) + math.Pow(float64(_v2.Y-v1.Y), 2)))
}

func DotProduct(_v1, _v2 Vector2f) float32 {
	return _v1.X*_v2.X + _v1.Y*_v2.Y
}

func (v Vector2f) Perpendicular() Vector2f {
	return Vector2f{-v.Y, v.X}
}

func (v Vector2f) Angle() float32 {
	ang := float32(math.Atan2(float64(v.Y), float64(v.X)))
	return ang
}

func (v Vector2f) Rotate(_rotation float32, _pivot Vector2f) Vector2f {
	anglePolar := _rotation * Deg2Rad
	_x := v.X
	_y := v.Y

	v.X = (_x-_pivot.X)*float32(math.Cos(float64(anglePolar))) - (_y-_pivot.Y)*float32(math.Sin(float64(anglePolar))) + _pivot.X
	v.Y = (_x-_pivot.X)*float32(math.Sin(float64(anglePolar))) + (_y-_pivot.Y)*float32(math.Cos(float64(anglePolar))) + _pivot.Y
	return Vector2f{
		v.X, v.Y,
	}
}

func (v Vector2f) RotateCenter(_rotation float32) Vector2f {
	anglePolar := _rotation * math.Pi / 180.0
	_x := v.X
	_y := v.Y

	v.X = (_x)*float32(math.Cos(float64(anglePolar))) - (_y)*float32(math.Sin(float64(anglePolar)))
	v.Y = (_x)*float32(math.Sin(float64(anglePolar))) + (_y)*float32(math.Cos(float64(anglePolar)))
	return Vector2f{
		v.X, v.Y,
	}
}

func LerpVector2f(_a, _b Vector2f, _t float32) Vector2f {
	return Vector2f{
		X: _a.X + (_b.X-_a.X)*_t,
		Y: _a.Y + (_b.Y-_a.Y)*_t,
	}
}

func Vector2fMidpoint(_v1, _v2 Vector2f) Vector2f {
	return _v1.Add(_v2).Scale(0.5)
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
