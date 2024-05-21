package chai

func GetCurrentCanvasPageSize() Vector2i {
	canvasBoundingBox := canvas.Call("getBoundingClientRect")
	cSize := NewVector2i(canvasBoundingBox.Get("right").Int()-canvasBoundingBox.Get("left").Int(), canvasBoundingBox.Get("bottom").Int()-canvasBoundingBox.Get("top").Int())
	return cSize
}
