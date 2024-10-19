package chai

import (
	. "github.com/mhamedGd/chai/math"
	"github.com/udhos/goglmath"
)

type Camera2D struct {
	position      Vector2f
	centerOffset  Vector2f
	scale         float32
	projectMatrix goglmath.Matrix4
	viewMatrix    goglmath.Matrix4
	mustUpdate    bool
}

func (cam *Camera2D) Init(_app App) {
	cam.position = Vector2fZero
	cam.scale = 1.0
	cam.projectMatrix = goglmath.Matrix4{}
	cam.projectMatrix = Ortho(0, float32(_app.Width), 0, float32(_app.Height), -5.0, 5.0)

	cam.viewMatrix = goglmath.NewMatrix4Identity()
	cam.mustUpdate = true
	//goglmath.SetOrthoMatrix(&cam.ProjectionMatrix, 0, float64(_app.Width), 0, float64(_app.Height), -1, 1)

}

func (cam *Camera2D) Update(_app App) {
	if !cam.mustUpdate {
		return
	}

	cam.viewMatrix = Translate(cam.projectMatrix, -cam.position.X+cam.centerOffset.X, -cam.position.Y+cam.centerOffset.Y, 0.0, 1.0, cam.scale)
	cam.mustUpdate = false
}

func Ortho(left, right, bottom, top, near, far float32) goglmath.Matrix4 {

	matrix := goglmath.Matrix4{}
	goglmath.SetOrthoMatrix(&matrix, float64(left), float64(right), float64(bottom), float64(top), -1, 1)

	return matrix
}

func Translate(refMatrix goglmath.Matrix4, x, y, z, w float32, scale float32) goglmath.Matrix4 {
	matrix := refMatrix
	matrix.Translate(float64(x), float64(y), float64(z), float64(w))

	m2 := goglmath.NewMatrix4Identity()
	m2.Scale(float64(scale), float64(scale), 0.0, 1.0)
	m2.Multiply(&matrix)

	return m2
}

func (cam *Camera2D) GetPosition() Vector2f {
	return cam.position
}

func (cam *Camera2D) GetScale() float32 {
	return cam.scale
}

func (cam *Camera2D) IsBoxInView(pos, dim Vector2f) bool {
	scaledScreenDimensions := NewVector2f(float32(GetCanvasWidth()), float32(GetCanvasHeigth())).Scale(1 / (cam.scale))

	MIN_DISTANCE_X := dim.X/0.5 + scaledScreenDimensions.X/2.0
	MIN_DISTANCE_Y := dim.Y/0.5 + scaledScreenDimensions.Y/2.0

	distVector := (pos.Add(dim)).Subtract(cam.position)

	xDepth := MIN_DISTANCE_X - AbsFloat32(distVector.X)
	yDepth := MIN_DISTANCE_Y - AbsFloat32(distVector.Y)

	return xDepth > 0 && yDepth > 0
}

func Scale(s float32) goglmath.Matrix4 {
	matrix := goglmath.NewMatrix4Identity()
	matrix.Scale(float64(s), float64(s), 0, 1)

	return matrix
}

const (
	scalePercentage = 0.01
	maxCamScale     = 10000
	minCamScale     = 0.0001
)

// Moves the viewport to a specific position (Bottom-Left)
func ScrollTo(newPosition Vector2f) {
	Cam.position = newPosition
	Cam.mustUpdate = true
}

// Offsets the viewport by a specific value (Bottom-Left)
func ScrollView(offset Vector2f) {
	ScrollTo(Cam.position.Add(offset))
}

func ScaleView(newScale float32) {
	Cam.scale = newScale
	RegulateScale()
}

func IncreaseScale(increment float32) {
	Cam.scale += increment
	RegulateScale()
}

func IncreaseScaleU(increment float32) {
	Cam.scale += (scalePercentage * Cam.scale) * increment
	RegulateScale()
}

func RegulateScale() {
	Cam.scale = ClampFloat32(Cam.scale, minCamScale, maxCamScale)
	Cam.mustUpdate = true
}

func GetMouseWorldPosition() Vector2f {
	screenPoint := MouseCanvasPos
	screenPoint = screenPoint.Subtract(NewVector2f(float32(GetCanvasWidth())/2.0, float32(GetCanvasHeigth())/2.0))
	screenPoint = screenPoint.Scale(1 / Cam.scale)
	screenPoint = screenPoint.Add(Cam.position)
	return screenPoint
}

func GetMouseScreenPosition() Vector2f {
	screenPoint := MouseCanvasPos
	screenPoint = screenPoint.Scale(1 / uiCam.scale)
	screenPoint = screenPoint.Add(uiCam.position)
	return screenPoint
}
