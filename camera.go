package chai

import (
	. "github.com/mhamedGd/chai/math"
	"github.com/udhos/goglmath"
)

type Camera2D struct {
	m_Position      Vector2f
	m_CenterOffset  Vector2f
	m_Scale         float32
	m_ProjectMatrix goglmath.Matrix4
	m_ViewMatrix    goglmath.Matrix4
	m_MustUpdate    bool
}

func (cam *Camera2D) Init(_app App) {
	cam.m_Position = Vector2fZero
	cam.m_Scale = 1.0
	cam.m_ProjectMatrix = goglmath.Matrix4{}
	cam.m_ProjectMatrix = ortho(0, float32(_app.Width), 0, float32(_app.Height), -5.0, 5.0)

	cam.m_ViewMatrix = goglmath.NewMatrix4Identity()
	cam.m_MustUpdate = true
	//goglmath.SetOrthoMatrix(&cam.ProjectionMatrix, 0, float64(_app.Width), 0, float64(_app.Height), -1, 1)

}

func (cam *Camera2D) Update(_app App) {
	if !cam.m_MustUpdate {
		return
	}

	cam.m_ViewMatrix = translate(cam.m_ProjectMatrix, -cam.m_Position.X+cam.m_CenterOffset.X, -cam.m_Position.Y+cam.m_CenterOffset.Y, 0.0, 1.0, cam.m_Scale)
	cam.m_MustUpdate = false
}

func ortho(_left, _right, _bottom, _top, _near, _far float32) goglmath.Matrix4 {

	matrix := goglmath.Matrix4{}
	goglmath.SetOrthoMatrix(&matrix, float64(_left), float64(_right), float64(_bottom), float64(_top), -1, 1)

	return matrix
}

func translate(_refMatrix goglmath.Matrix4, _x, _y, _z, _w float32, _scale float32) goglmath.Matrix4 {
	matrix := _refMatrix
	matrix.Translate(float64(_x), float64(_y), float64(_z), float64(_w))

	m2 := goglmath.NewMatrix4Identity()
	m2.Scale(float64(_scale), float64(_scale), 0.0, 1.0)
	m2.Multiply(&matrix)

	return m2
}

func (cam *Camera2D) GetPosition() Vector2f {
	return cam.m_Position
}

func (cam *Camera2D) GetScale() float32 {
	return cam.m_Scale
}

func (cam *Camera2D) IsBoxInView(_pos, _dim Vector2f) bool {
	scaledScreenDimensions := NewVector2f(float32(GetCanvasWidth()), float32(GetCanvasHeigth())).Scale(1 / (cam.m_Scale))

	MIN_DISTANCE_X := _dim.X/0.5 + scaledScreenDimensions.X/2.0
	MIN_DISTANCE_Y := _dim.Y/0.5 + scaledScreenDimensions.Y/2.0

	distVector := (_pos.Add(_dim)).Subtract(cam.m_Position)

	xDepth := MIN_DISTANCE_X - AbsFloat32(distVector.X)
	yDepth := MIN_DISTANCE_Y - AbsFloat32(distVector.Y)

	return xDepth > 0 && yDepth > 0
}

func Scale(_s float32) goglmath.Matrix4 {
	matrix := goglmath.NewMatrix4Identity()
	matrix.Scale(float64(_s), float64(_s), 0, 1)

	return matrix
}

const (
	scalePercentage = 0.01
	maxCamScale     = 10000
	minCamScale     = 0.0001
)

// Moves the viewport to a specific m_Position (Bottom-Left)
func ScrollTo(_newPosition Vector2f) {
	Cam.m_Position = _newPosition
	Cam.m_MustUpdate = true
}

// Offsets the viewport by a specific value (Bottom-Left)
func ScrollView(_offset Vector2f) {
	ScrollTo(Cam.m_Position.Add(_offset))
}

func ScaleView(_newScale float32) {
	Cam.m_Scale = _newScale
	RegulateScale()
}

func IncreaseScale(_increment float32) {
	Cam.m_Scale += _increment
	RegulateScale()
}

func IncreaseScaleU(_increment float32) {
	Cam.m_Scale += (scalePercentage * Cam.m_Scale) * _increment
	RegulateScale()
}

func RegulateScale() {
	Cam.m_Scale = ClampFloat32(Cam.m_Scale, minCamScale, maxCamScale)
	Cam.m_MustUpdate = true
}

func GetMouseWorldPosition() Vector2f {
	screenPoint := MouseCanvasPos
	screenPoint = screenPoint.Subtract(NewVector2f(float32(GetCanvasWidth())/2.0, float32(GetCanvasHeigth())/2.0))
	screenPoint = screenPoint.Scale(1 / Cam.m_Scale)
	screenPoint = screenPoint.Add(Cam.m_Position)
	return screenPoint
}

func GetMouseScreenPosition() Vector2f {
	screenPoint := MouseCanvasPos
	screenPoint = screenPoint.Scale(1 / uiCam.m_Scale)
	screenPoint = screenPoint.Add(uiCam.m_Position)
	return screenPoint
}
