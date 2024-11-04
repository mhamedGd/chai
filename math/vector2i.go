package math

import (
	"fmt"
	"math"
)

type Vector2i struct {
	X, Y int
}

func NewVector2i(_x, _y int) Vector2i {
	return Vector2i{
		_x, _y,
	}
}

func (v1 *Vector2i) Equal(_v2 *Vector2i) bool {
	return v1.X == _v2.X && v1.Y == _v2.Y
}

func (v1 Vector2i) Add(_v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X + _v2.X, Y: v1.Y + _v2.Y}
}

func (v1 Vector2i) AddXY(_x, _y int) Vector2i {
	return Vector2i{X: v1.X + _x, Y: v1.Y + _y}
}
func (v1 Vector2i) Subtract(_v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X - _v2.X, Y: v1.Y - _v2.Y}
}

func (v1 Vector2i) SubtractXY(_x, _y int) Vector2i {
	return Vector2i{X: v1.X - _x, Y: v1.Y - _y}
}

func (v1 Vector2i) Multp(_v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X * _v2.X, Y: v1.Y * _v2.Y}
}

func (v1 Vector2i) MultpXY(_x, _y int) Vector2i {
	return Vector2i{X: v1.X * _x, Y: v1.Y * _y}
}

func (v1 Vector2i) Div(_v2 Vector2i) Vector2i {
	return Vector2i{X: v1.X / _v2.X, Y: v1.Y / _v2.Y}
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

func (v *Vector2i) Perpendicular() Vector2i {
	return Vector2i{-v.Y, v.X}
}

func (v *Vector2i) Angle() int {
	return int(math.Atan2(float64(v.Y), float64(v.X)))
}

func (v Vector2i) Rotate(_angle float32, _pivot Vector2i) Vector2i {
	anglePolar := _angle * math.Pi / 180.0
	_x := v.X
	_y := v.Y

	v.X = (_x-_pivot.X)*int(math.Cos(float64(anglePolar))) - (_y-_pivot.Y)*int(math.Sin(float64(anglePolar))) + _pivot.X
	v.Y = (_x-_pivot.X)*int(math.Sin(float64(anglePolar))) + (_y-_pivot.Y)*int(math.Cos(float64(anglePolar))) + _pivot.Y
	return Vector2i{
		v.X, v.Y,
	}
}

func (v Vector2i) RotateCenter(_angle float32) Vector2i {
	anglePolar := _angle * math.Pi / 180.0
	_x := v.X
	_y := v.Y

	v.X = (_x)*int(math.Cos(float64(anglePolar))) - (_y)*int(math.Sin(float64(anglePolar)))
	v.Y = (_x)*int(math.Sin(float64(anglePolar))) + (_y)*int(math.Cos(float64(anglePolar)))
	return Vector2i{
		v.X, v.Y,
	}
}

func (v Vector2i) ToString() string {
	return fmt.Sprint(v.X, v.Y)
}
