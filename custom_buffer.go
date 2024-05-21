package chai

import (
	"reflect"
	"syscall/js"
	"unsafe"
)

func VertexSliceAsBytes(_verts []Vertex) []byte {
	n := int(vertexByteSize) * len(_verts)

	up := unsafe.Pointer(&(_verts[0]))
	pi := (*[1]byte)(up)
	buf := (*pi)[:]
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Len = n
	sh.Cap = n

	return buf
}

func vertexBufferToJsVertexBuffer(_buffer []Vertex) js.Value {
	jsVerts := js.Global().Get("Uint8Array").New(len(_buffer) * VertexSize)
	var verticesBytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&verticesBytes))
	header.Cap = cap(_buffer) * VertexSize
	header.Len = len(_buffer) * VertexSize
	header.Data = uintptr(unsafe.Pointer(&_buffer[0]))

	js.CopyBytesToJS(jsVerts, verticesBytes)
	return jsVerts
}

func float32BufferToJsFloat32Buffer(_buffer []float32) js.Value {
	jsElements := js.Global().Get("Uint8Array").New(len(_buffer) * 4)
	var elementsBytes []byte
	headerElem := (*reflect.SliceHeader)(unsafe.Pointer(&elementsBytes))
	headerElem.Cap = cap(_buffer) * 4
	headerElem.Len = len(_buffer) * 4
	headerElem.Data = uintptr(unsafe.Pointer(&_buffer[0]))

	js.CopyBytesToJS(jsElements, elementsBytes)
	return jsElements
}

func int32BufferToJsInt32Buffer(_buffer []int32) js.Value {
	jsElements := js.Global().Get("Uint8Array").New(len(_buffer) * 4)
	var elementsBytes []byte
	headerElem := (*reflect.SliceHeader)(unsafe.Pointer(&elementsBytes))
	headerElem.Cap = cap(_buffer) * 4
	headerElem.Len = len(_buffer) * 4
	headerElem.Data = uintptr(unsafe.Pointer(&_buffer[0]))
	js.CopyBytesToJS(jsElements, elementsBytes)

	return jsElements
}

func pixelBufferToJsPixelBubffer(_buffer []Pixel) js.Value {
	jsPixels := js.Global().Get("Uint8Array").New(len(_buffer) * 4)
	var pixelsBytes []byte
	headerPixels := (*reflect.SliceHeader)(unsafe.Pointer(&pixelsBytes))
	headerPixels.Cap = cap(_buffer) * 4
	headerPixels.Len = len(_buffer) * 4
	headerPixels.Data = uintptr(unsafe.Pointer(&_buffer[0]))

	js.CopyBytesToJS(jsPixels, pixelsBytes)
	return jsPixels
}
