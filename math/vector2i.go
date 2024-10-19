package math

import (
	"fmt"
	"math"
)

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

func (v Vector2i) ToString() string {
	return fmt.Sprint(v.X, v.Y)
}
