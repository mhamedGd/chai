package chai

import (
	"math"
	"math/rand"
	"syscall/js"

	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

const SHAPES_SHADER_VERTEX = `#version 300 es

precision mediump float;

in vec2 coordinates;
in float z;
in vec4 colors;

out vec4 vertex_FragColor;

uniform mat4 projection_matrix;
uniform mat4 view_matrix;

void main(void) {
	vec4 global_position = vec4(0.0);
	global_position = view_matrix * vec4(coordinates, 0.0, 1.0);
	global_position.z = z;
	global_position.w = 1.0;		
	gl_Position = global_position;
	
	vertex_FragColor = colors;
}
`
const SHAPES_SHADER_FRAGMENT = `#version 300 es

precision mediump float;

in vec4 vertex_FragColor;

out vec4 fragColor;
void main(void) {
	vec4 finalColor = vertex_FragColor;
	finalColor.r *= vertex_FragColor.a;
	finalColor.g *= vertex_FragColor.a;
	finalColor.b *= vertex_FragColor.a;
    fragColor = finalColor;
}
`

///////// OLD RENDERER //////////////////////////////////////////////
// const VertexSize = 24

// type Vertex struct {
// 	Coordinates Vector2f
// 	Z           float32
// 	Color       RGBA8
// 	UV          Vector2f
// }

// func NewVertex(_coords Vector2f, _z float32, _uv Vector2f, _color RGBA8) Vertex {
// 	return Vertex{_coords, _z, _color, _uv}
// }

// const vertexByteSize uintptr = unsafe.Sizeof(Vertex{})

// type ShapeBatch struct {
// 	m_Vao             js.Value
// 	m_Vbo, m_Ibo      js.Value
// 	Vertices          customtypes.List[Vertex]
// 	Indices           customtypes.List[int32]
// 	NumberOfElements  int
// 	NumberOfInstances int
// 	Shader            ShaderProgram
// 	LineWidth         float32
// 	Initialized       bool
// }

// func (_shapesB *ShapeBatch) Init() {
// 	_shapesB.LineWidth = 50

// 	_shapesB.Vertices = customtypes.NewList[Vertex]()
// 	_shapesB.Indices = customtypes.NewList[int32]()

// 	_shapesB.m_Vao = canvasContext.Call("createVertexArray")
// 	canvasContext.Call("bindVertexArray", _shapesB.m_Vao)

// 	_shapesB.m_Vbo = canvasContext.Call("createBuffer")
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _shapesB.m_Vbo)

// 	canvasContext.Call("enableVertexAttribArray", 0)
// 	canvasContext.Call("enableVertexAttribArray", 1)
// 	canvasContext.Call("enableVertexAttribArray", 2)
// 	canvasContext.Call("enableVertexAttribArray", 3)

// 	canvasContext.Call("vertexAttribPointer", 0, 2, canvasContext.Get("FLOAT"), false, VertexSize, 0)
// 	canvasContext.Call("vertexAttribPointer", 1, 1, canvasContext.Get("FLOAT"), false, VertexSize, 8)
// 	canvasContext.Call("vertexAttribPointer", 2, 4, canvasContext.Get("UNSIGNED_BYTE"), true, VertexSize, 12)
// 	canvasContext.Call("vertexAttribPointer", 3, 2, canvasContext.Get("FLOAT"), false, VertexSize, 16)

// 	_shapesB.m_Ibo = canvasContext.Call("createBuffer")
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), _shapesB.m_Ibo)

// 	canvasContext.Call("bindVertexArray", js.Null())
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), js.Null())
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), js.Null())

// 	canvasContext.Call("disableVertexAttribArray", 0)
// 	canvasContext.Call("disableVertexAttribArray", 1)
// 	canvasContext.Call("disableVertexAttribArray", 2)
// 	canvasContext.Call("disableVertexAttribArray", 3)

// 	_shapesB.Shader.ParseShader(SHAPES_SHADER_VERTEX, SHAPES_SHADER_FRAGMENT)
// 	//_shapesB.Shader.ParseShaderFromFile("shapes.m_Shader")
// 	_shapesB.Shader.CreateShaderProgram()
// 	_shapesB.Shader.AddAttribute("coordinates")
// 	_shapesB.Shader.AddAttribute("z")
// 	_shapesB.Shader.AddAttribute("colors")
// 	_shapesB.Initialized = true
// }

// func (_sp *ShapeBatch) DrawLine(_from, _to Vector2f, _z float32, _color RGBA8) {
// 	//var m_Offset Vector2f = NewVector2f(_sp.LineWidth/2.0, 0.0)
// 	var m_Offset Vector2f

// 	m_Offset = _to.Subtract(_from)
// 	m_Offset = m_Offset.Perpendicular()
// 	m_Offset = m_Offset.Normalize()
// 	m_Offset = m_Offset.Scale(_sp.LineWidth)

// 	vertsSize := _sp.Vertices.Count()
// 	_new_z := remapZForView(_z)
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _from.Subtract(m_Offset), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _to.Subtract(m_Offset), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _from.Add(m_Offset), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _to.Add(m_Offset), Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 3))

// 	_sp.NumberOfInstances++
// }

// func (_sp *ShapeBatch) DrawRect(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8) {
// 	offsetX := _dimensions.Scale(0.5).X
// 	offsetY := _dimensions.Scale(0.5).Y

// 	lineDifference := _sp.LineWidth

// 	_sp.DrawLine(NewVector2f(-offsetX-lineDifference, offsetY).Add(_center), NewVector2f(offsetX+lineDifference, offsetY).Add(_center), _z, _color)
// 	_sp.DrawLine(NewVector2f(-offsetX-lineDifference, -offsetY).Add(_center), NewVector2f(offsetX+lineDifference, -offsetY).Add(_center), _z, _color)
// 	_sp.DrawLine(NewVector2f(offsetX, offsetY+lineDifference).Add(_center), NewVector2f(offsetX, -offsetY-lineDifference).Add(_center), _z, _color)
// 	_sp.DrawLine(NewVector2f(-offsetX, offsetY+lineDifference).Add(_center), NewVector2f(-offsetX, -offsetY-lineDifference).Add(_center), _z, _color)
// }

// func (_sp *ShapeBatch) DrawRectRotated(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8, _angle float32) {
// 	offsetX := _dimensions.Scale(0.5).X
// 	offsetY := _dimensions.Scale(0.5).Y

// 	lineDifference := _sp.LineWidth

// 	_sp.DrawLine(NewVector2f(-offsetX-lineDifference, offsetY).Add(_center).Rotate(_angle, _center), NewVector2f(offsetX+lineDifference, offsetY).Add(_center).Rotate(_angle, _center), _z, _color)
// 	_sp.DrawLine(NewVector2f(-offsetX-lineDifference, -offsetY).Add(_center).Rotate(_angle, _center), NewVector2f(offsetX+lineDifference, -offsetY).Add(_center).Rotate(_angle, _center), _z, _color)
// 	_sp.DrawLine(NewVector2f(offsetX, offsetY+lineDifference).Add(_center).Rotate(_angle, _center), NewVector2f(offsetX, -offsetY-lineDifference).Add(_center).Rotate(_angle, _center), _z, _color)
// 	_sp.DrawLine(NewVector2f(-offsetX, offsetY+lineDifference).Add(_center).Rotate(_angle, _center), NewVector2f(-offsetX, -offsetY-lineDifference).Add(_center).Rotate(_angle, _center), _z, _color)
// }

// func (_sp *ShapeBatch) DrawTriangle(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8) {
// 	const numOfVertices = 3

// 	pos := [numOfVertices]Vector2f{}
// 	for i := 0; i < numOfVertices; i++ {
// 		angle := (float32(i) / float32(numOfVertices)) * 2.0 * PI
// 		pos[i] = NewVector2f(_center.X+float32(math.Cos(float64(angle)))*_dimensions.Y, _center.Y+float32(math.Sin(float64(angle)))*_dimensions.X)
// 	}

// 	for i := 0; i < numOfVertices-1; i++ {
// 		_sp.DrawLine(pos[i], pos[i+1], _z, _color)
// 	}
// 	_sp.DrawLine(pos[numOfVertices-1], pos[0], _z, _color)
// }

// func (_sp *ShapeBatch) DrawTriangleRotated(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8, _rotation float32) {
// 	numOfVertices := 3

// 	pos := [3]Vector2f{}
// 	for i := 0; i < numOfVertices; i++ {
// 		angle := (float32(i) / float32(3.0)) * 2.0 * PI
// 		pos[i] = NewVector2f(_center.X+float32(math.Cos(float64(angle)))*_dimensions.Y, _center.Y+float32(math.Sin(float64(angle)))*_dimensions.X)
// 		pos[i] = pos[i].Rotate(_rotation, _center)
// 	}

// 	for i := 0; i < numOfVertices-1; i++ {
// 		_sp.DrawLine(pos[i], pos[i+1], _z, _color)
// 	}
// 	_sp.DrawLine(pos[numOfVertices-1], pos[0], _z, _color)
// }

// func (_sp *ShapeBatch) DrawFillTriangle(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8) {
// 	m_Offset := _dimensions.Scale(0.5)

// 	vertsSize := _sp.Vertices.Count()

// 	_new_z := remapZForView(_z)
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.Subtract(m_Offset), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.AddXY(m_Offset.X, -m_Offset.Y), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.AddXY(0, m_Offset.Y), Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))

// 	_sp.NumberOfInstances++
// }

// func (_sp *ShapeBatch) DrawFillTriangleRotated(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8, _rotation float32) {
// 	m_Offset := _dimensions.Scale(0.5)

// 	vertsSize := _sp.Vertices.Count()

// 	_new_z := remapZForView(_z)
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.Subtract(m_Offset).Rotate(_rotation, _center), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.AddXY(m_Offset.X, -m_Offset.Y).Rotate(_rotation, _center), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.AddXY(0, m_Offset.Y).Rotate(_rotation, _center), Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))

// 	_sp.NumberOfInstances++
// }

// func (_sp *ShapeBatch) DrawCircle(_center Vector2f, _z, _radius float32, _color RGBA8) {
// 	numOfVertices := 16

// 	pos := [16]Vector2f{}
// 	for i := 0; i < numOfVertices; i++ {
// 		angle := (float32(i) / float32(numOfVertices)) * 2.0 * PI
// 		pos[i] = NewVector2f(_center.X+float32(math.Cos(float64(angle)))*_radius, _center.Y+float32(math.Sin(float64(angle)))*_radius)
// 	}

// 	for i := 0; i < numOfVertices-1; i++ {
// 		_sp.DrawLine(pos[i], pos[i+1], _z, _color)
// 	}
// 	_sp.DrawLine(pos[numOfVertices-1], pos[0], _z, _color)
// }

// func (_sp *ShapeBatch) DrawFillRect(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8) {
// 	m_Offset := _dimensions.Scale(0.5)

// 	vertsSize := _sp.Vertices.Count()

// 	_new_z := remapZForView(_z)
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.Subtract(m_Offset), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.SubtractXY(m_Offset.X, -m_Offset.Y), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.AddXY(m_Offset.X, -m_Offset.Y), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.Add(m_Offset), Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 3))

// 	_sp.NumberOfInstances++
// }

// func (_sp *ShapeBatch) DrawFillRect_Rect(_rect Rect, _z float32, _color RGBA8) {

// 	vertsSize := _sp.Vertices.Count()

// 	_new_z := remapZForView(_z)
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _rect.Position, Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _rect.Position.AddXY(0.0, _rect.Size.Y), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _rect.Position.AddXY(_rect.Size.X, 0), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _rect.Size, Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 3))

// 	_sp.NumberOfInstances++

// }

// func (_sp *ShapeBatch) DrawFillRectBottom(_bottom Vector2f, _z float32, _dimensions Vector2f, _color RGBA8) {
// 	m_Offset := NewVector2f(_dimensions.X/2.0, 0.0)

// 	vertsSize := _sp.Vertices.Count()

// 	_new_z := remapZForView(_z)
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.Subtract(m_Offset), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.SubtractXY(m_Offset.X, -m_Offset.Y), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.AddXY(m_Offset.X, -m_Offset.Y), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.Add(m_Offset), Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 3))

// 	_sp.NumberOfInstances++

// }

// func (_sp *ShapeBatch) DrawFillRectRotated(_center Vector2f, _z float32, _dimensions Vector2f, _color RGBA8, _rotation float32) {
// 	m_Offset := _dimensions.Scale(0.5)

// 	vertsSize := _sp.Vertices.Count()

// 	_new_z := remapZForView(_z)
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.Subtract(m_Offset).Rotate(_rotation, _center), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.SubtractXY(m_Offset.X, -m_Offset.Y).Rotate(_rotation, _center), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.AddXY(m_Offset.X, -m_Offset.Y).Rotate(_rotation, _center), Z: _new_z, Color: _color})
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _center.Add(m_Offset).Rotate(_rotation, _center), Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 3))

// 	_sp.NumberOfInstances++

// }

// func (_sp *ShapeBatch) DrawFillRectBottomRotated(_bottom Vector2f, _z float32, _dimensions Vector2f, _color RGBA8, _rotation float32) {
// 	// m_Offset := NewVector2f(_dimensions.X/2.0, 0.0)
// 	m_Offset := _dimensions.Scale(0.5)

// 	vertsSize := _sp.Vertices.Count()
// 	_new_z := remapZForView(_z)
// 	//Bottom Left
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.SubtractXY(m_Offset.X, 0.0).Rotate(_rotation, _bottom), Z: _new_z, Color: _color})
// 	//Top Left
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.SubtractXY(m_Offset.X, -m_Offset.Y*2).Rotate(_rotation, _bottom), Z: _new_z, Color: _color})

// 	//Bottom Right
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.AddXY(m_Offset.X, 0).Rotate(_rotation, _bottom), Z: _new_z, Color: _color})

// 	//Top Right
// 	_sp.Vertices.PushBack(Vertex{Coordinates: _bottom.AddXY(m_Offset.X, m_Offset.Y*2.0).Rotate(_rotation, _bottom), Z: _new_z, Color: _color})

// 	_sp.Indices.PushBack(int32(vertsSize + 0))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 2))
// 	_sp.Indices.PushBack(int32(vertsSize + 1))
// 	_sp.Indices.PushBack(int32(vertsSize + 3))

// 	_sp.NumberOfInstances++

// }

// func (_sp *ShapeBatch) finalize() {
// 	if _sp.Vertices.Count() == 0 {
// 		_sp.NumberOfElements = 0
// 		return
// 	}

// 	jsVerts := vertexBufferToJsVertexBuffer(_sp.Vertices.Data)
// 	jsElem := int32BufferToJsInt32Buffer(_sp.Indices.Data)

// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _sp.m_Vbo)
// 	canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), jsVerts, canvasContext.Get("STATIC_DRAW"))

// 	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), _sp.m_Ibo)
// 	canvasContext.Call("bufferData", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), jsElem, canvasContext.Get("STATIC_DRAW"))

// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), js.Null())
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), js.Null())

// 	_sp.NumberOfElements = _sp.Indices.Count()

// 	jsVerts.Set("delete", nil)
// 	jsElem.Set("delete", nil)

// 	_sp.Vertices.Clear()
// 	_sp.Indices.Clear()
// }

// func (_sp *ShapeBatch) Render(_cam *Camera2D) {
// 	_sp.finalize()

// 	UseShader(&_sp.Shader)

// 	viewMatrix := _cam.m_ViewMatrix.Data()

// 	viewBuffer := new(bytes.Buffer)
// 	{
// 		err := binary.Write(viewBuffer, binary.LittleEndian, viewMatrix[:])
// 		Assert(err == nil, "Error writing viewBuffer")
// 	}

// 	byte_slice2 := viewBuffer.Bytes()

// 	b2 := js.Global().Get("Uint8Array").New(len(viewMatrix) * 4)
// 	js.CopyBytesToJS(b2, byte_slice2)

// 	viewMatrixJS := js.Global().Get("Float32Array").New(b2.Get("buffer"), b2.Get("byteOffset"), b2.Get("byteLength").Int()/4)
// 	for i := 0; i < len(viewMatrix); i++ {
// 		viewMatrixJS.SetIndex(i, js.ValueOf(viewMatrix[i]))
// 	}

// 	viewmatrix_loc := canvasContext.Call("getUniformLocation", _sp.Shader.ShaderProgramID, "view_matrix")
// 	canvasContext.Call("uniformMatrix4fv", viewmatrix_loc, false, viewMatrixJS)

// 	canvasContext.Call("bindVertexArray", _sp.m_Vao)
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), _sp.m_Ibo)
// 	canvasContext.Call("drawElements", canvasContext.Get("TRIANGLES"), _sp.NumberOfElements, canvasContext.Get("UNSIGNED_INT"), 0)
// 	_sp.NumberOfInstances = 0
// 	UnuseShader()
// }

// /*
// ##############################################################
// ################ Sprite Batch - Sprite Batch #################
// ##############################################################
// */

// const SPRITES_SHADER_VERTEX = `#version 300 es

// precision mediump float;

// in vec2 coordinates;
// in float z;
// in vec4 colors;
// in vec2 uv;

// out vec4 vertex_FragColor;
// out vec2 vertex_UV;

// uniform mat4 projection_matrix;
// uniform mat4 view_matrix;

// void main(void) {
// 	vec4 global_position = vec4(0.0);
// 	global_position = view_matrix * vec4(coordinates, 0.0, 1.0);
// 	global_position.z = z;
// 	global_position.w = 1.0;
// 	gl_Position = global_position;

// 	vertex_FragColor = colors;
// 	vertex_UV = uv;
// }`
// const SPRITES_SHADER_FRAGMENT = `#version 300 es

// precision mediump float;

// in vec4 vertex_FragColor;
// in vec2 vertex_UV;

// uniform sampler2D genericSampler;

// out vec4 fragColor;

// void main(void) {
// 	ivec2 texSize = textureSize(genericSampler, 0);
// 	float s_offset = (1.0 / (float(texSize.x))) * 0.0;
// 	float t_offset = (1.0 / (float(texSize.y))) * 0.0;
// 	vec2 tc_final =  vertex_UV + vec2(s_offset, t_offset);

// 	vec4 thisColor = vertex_FragColor * texture(genericSampler, tc_final);
// 	thisColor.r *= vertex_FragColor.a;
// 	thisColor.g *= vertex_FragColor.a;
// 	thisColor.b *= vertex_FragColor.a;
// 	fragColor = thisColor;
// }`

// type SpriteGlyph struct {
// 	m_Bottomleft, m_Topleft, m_Topright, m_Bottomright Vertex
// 	m_Texture                                          *Texture2D
// }

// func NewSpriteGlyph(_pos, _dimensions, _uv1 Vector2f, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8) SpriteGlyph {
// 	var tempGlyph SpriteGlyph

// 	halfDim := _dimensions.Scale(0.5)

// 	tempGlyph.m_Bottomleft = NewVertex(_pos.Subtract(halfDim), _z, NewVector2f(_uv1.X, _uv2.Y), _tint)
// 	tempGlyph.m_Topleft = NewVertex(_pos.Add(NewVector2f(-halfDim.X, halfDim.Y)), _z, _uv1, _tint)
// 	tempGlyph.m_Topright = NewVertex(_pos.Add(halfDim), _z, NewVector2f(_uv2.X, _uv1.Y), _tint)
// 	tempGlyph.m_Bottomright = NewVertex(_pos.Add(NewVector2f(halfDim.X, -halfDim.Y)), _z, _uv2, _tint)

// 	tempGlyph.m_Texture = _texture

// 	return tempGlyph
// }

// func NewSpriteGlyphRotated(_pos, _dimensions, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8, _rotation float32) SpriteGlyph {
// 	var tempGlyph SpriteGlyph

// 	halfDim := _dimensions.Scale(0.5)

// 	tempGlyph.m_Bottomleft = NewVertex(_pos.Subtract(halfDim).Rotate(_rotation, _pos), _z, NewVector2f(_uv1.X, _uv2.Y), _tint)
// 	tempGlyph.m_Topleft = NewVertex(_pos.Add(NewVector2f(-halfDim.X, halfDim.Y)).Rotate(_rotation, _pos), _z, _uv1, _tint)
// 	tempGlyph.m_Topright = NewVertex(_pos.Add(halfDim).Rotate(_rotation, _pos), _z, NewVector2f(_uv2.X, _uv1.Y), _tint)
// 	tempGlyph.m_Bottomright = NewVertex(_pos.Add(NewVector2f(halfDim.X, -halfDim.Y)).Rotate(_rotation, _pos), _z, _uv2, _tint)

// 	tempGlyph.m_Texture = _texture

// 	return tempGlyph
// }

// type RenderBatch struct {
// 	m_Offset, m_NumberOfVertices, m_NumOfInstances int
// 	m_Texture                                      *Texture2D
// }

// func NewRenderBatch(_offset, _numberOfVertices, _numOfInstances int, _texture *Texture2D) RenderBatch {
// 	return RenderBatch{_offset, _numberOfVertices, _numOfInstances, _texture}
// }

// type SpriteBatch struct {
// 	m_Vbo js.Value
// 	m_Vao js.Value

// 	m_Shader ShaderProgram

// 	m_RenderBatches []RenderBatch
// 	m_SpriteGlyphs  []SpriteGlyph

// 	RenderCam *Camera2D
// }

// func (self *SpriteBatch) Reset() {
// 	self.m_RenderBatches = self.m_RenderBatches[:0]
// 	self.m_SpriteGlyphs = self.m_SpriteGlyphs[:0]
// 	//`self.Init("")
// }

// func (self *SpriteBatch) Init(_shader_path string) {

// 	self.m_RenderBatches = make([]RenderBatch, 0)
// 	self.m_SpriteGlyphs = make([]SpriteGlyph, 0)

// 	self.m_Vao = canvasContext.Call("createVertexArray")
// 	canvasContext.Call("bindVertexArray", self.m_Vao)

// 	self.m_Vbo = canvasContext.Call("createBuffer")
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), self.m_Vbo)

// 	canvasContext.Call("enableVertexAttribArray", 0)
// 	canvasContext.Call("enableVertexAttribArray", 1)
// 	canvasContext.Call("enableVertexAttribArray", 2)
// 	canvasContext.Call("enableVertexAttribArray", 3)

// 	canvasContext.Call("vertexAttribPointer", 0, 2, canvasContext.Get("FLOAT"), false, VertexSize, 0)
// 	canvasContext.Call("vertexAttribPointer", 1, 1, canvasContext.Get("FLOAT"), false, VertexSize, 8)
// 	canvasContext.Call("vertexAttribPointer", 2, 4, canvasContext.Get("UNSIGNED_BYTE"), true, VertexSize, 12)
// 	canvasContext.Call("vertexAttribPointer", 3, 2, canvasContext.Get("FLOAT"), false, VertexSize, 16)

// 	canvasContext.Call("bindVertexArray", js.Null())
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), js.Null())

// 	canvasContext.Call("disableVertexAttribArray", 0)
// 	canvasContext.Call("disableVertexAttribArray", 1)
// 	canvasContext.Call("disableVertexAttribArray", 2)
// 	canvasContext.Call("disableVertexAttribArray", 3)

// 	if _shader_path == "" {
// 		self.m_Shader.ParseShader(SPRITES_SHADER_VERTEX, SPRITES_SHADER_FRAGMENT)
// 	} else {
// 		self.m_Shader.ParseShaderFromFile(_shader_path)
// 	}
// 	self.m_Shader.CreateShaderProgram()
// 	self.m_Shader.AddAttribute("coordinates")
// 	self.m_Shader.AddAttribute("colors")
// 	self.m_Shader.AddAttribute("uv")

// }

// func remapZForView(_z float32) float32 {
// 	return MapFloat32(_z, -1, 1000, -0.999, 0.999)
// }

// func (self *SpriteBatch) DrawSprite(_center, _dimensions, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyph(_center, _dimensions, _uv1, _uv2, remapZForView(_z), _texture, _tint))
// }

// func (self *SpriteBatch) DrawSpriteOrigin(_center, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyph(_center, NewVector2f(float32(_texture.Width), float32(_texture.Height)), _uv1, _uv2, remapZForView(_z), _texture, _tint))

// }
// func (self *SpriteBatch) DrawSpriteOriginScaled(_center, _uv1, _uv2 Vector2f, _scale float32, _z float32, _texture *Texture2D, _tint RGBA8) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyph(_center, NewVector2f(float32(_texture.Width), float32(_texture.Height)).Scale(_scale), _uv1, _uv2, remapZForView(_z), _texture, _tint))
// }
// func (self *SpriteBatch) DrawSpriteBottomLeft(_pos, _dimensions, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyph(_pos.Add(_dimensions.Scale(0.5)), _dimensions, _uv1, _uv2, remapZForView(_z), _texture, _tint))
// }
// func (self *SpriteBatch) DrawSpriteBottomLeftRotated(_pos, _dimensions, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8, _rotation float32) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyphRotated(_pos.Add(_dimensions.Scale(0.5)), _dimensions, _uv1, _uv2, remapZForView(_z), _texture, _tint, _rotation))
// }
// func (self *SpriteBatch) DrawSpriteBottomRight(_pos, _dimensions, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyph(_pos.AddXY(-_dimensions.Scale(0.5).X, _dimensions.Scale(0.5).Y), _dimensions, _uv1, _uv2, remapZForView(_z), _texture, _tint))
// }
// func (self *SpriteBatch) DrawSpriteBottomLeftOrigin(_pos, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyph(_pos.Subtract(NewVector2f(float32(_texture.Width), float32(_texture.Height)).Scale(0.5)), NewVector2f(float32(_texture.Width), float32(_texture.Height)), _uv1, _uv2, remapZForView(_z), _texture, _tint))

// }

// func (self *SpriteBatch) DrawSpriteOriginRotated(_center, _uv1, _uv2 Vector2f, _z float32, _texture *Texture2D, _tint RGBA8, _rotation float32) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyphRotated(_center, NewVector2f(float32(_texture.Width), float32(_texture.Height)), _uv1, _uv2, remapZForView(_z), _texture, _tint, _rotation))
// }

// func (self *SpriteBatch) DrawSpriteOriginScaledRotated(_center, _uv1, _uv2 Vector2f, _scale float32, _z float32, _texture *Texture2D, _tint RGBA8, _rotation float32) {
// 	_spriteDims := NewVector2f(float32(_texture.Width), float32(_texture.Height)).Scale(_scale)
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyphRotated(_center, _spriteDims, _uv1, _uv2, remapZForView(_z), _texture, _tint, _rotation))
// }
// func (self *SpriteBatch) DrawSpriteOriginScaledDimRotated(_center, _dims, _uv1, _uv2 Vector2f, _scale float32, _z float32, _texture *Texture2D, _tint RGBA8, _rotation float32) {
// 	self.m_SpriteGlyphs = append(self.m_SpriteGlyphs, NewSpriteGlyphRotated(_center, _dims.Scale(_scale), _uv1, _uv2, remapZForView(_z), _texture, _tint, _rotation))
// }

// func (self *SpriteBatch) finalize() {
// 	self.createRenderBatches()

// }

// func (self *SpriteBatch) Render() {
// 	self.finalize()

// 	UseShader(&self.m_Shader)

// 	viewMatrix := self.RenderCam.m_ViewMatrix.Data()

// 	viewBuffer := new(bytes.Buffer)
// 	{
// 		err := binary.Write(viewBuffer, binary.LittleEndian, viewMatrix[:])
// 		Assert(err == nil, "Error writing viewBuffer")
// 	}

// 	byte_slice2 := viewBuffer.Bytes()

// 	b2 := js.Global().Get("Uint8Array").New(len(viewMatrix) * 4)
// 	js.CopyBytesToJS(b2, byte_slice2)

// 	viewMatrixJS := js.Global().Get("Float32Array").New(b2.Get("buffer"), b2.Get("byteOffset"), b2.Get("byteLength").Int()/4)
// 	for i := 0; i < len(viewMatrix); i++ {
// 		viewMatrixJS.SetIndex(i, js.ValueOf(viewMatrix[i]))
// 	}

// 	viewmatrix_loc := canvasContext.Call("getUniformLocation", self.m_Shader.ShaderProgramID, "view_matrix")
// 	canvasContext.Call("uniformMatrix4fv", viewmatrix_loc, false, viewMatrixJS)

// 	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
// 	canvasContext.Call("uniform1i", self.m_Shader.GetUniformLocation("genericSampler"), 0)

// 	canvasContext.Call("bindVertexArray", self.m_Vao)
// 	for i := 0; i < len(self.m_RenderBatches); i++ {
// 		canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), self.m_RenderBatches[i].m_Texture.m_TextureId)
// 		canvasContext.Call("drawArrays", canvasContext.Get("TRIANGLES"), self.m_RenderBatches[i].m_Offset, self.m_RenderBatches[i].m_NumberOfVertices)
// 	}
// 	canvasContext.Call("bindVertexArray", js.Null())

// 	self.m_RenderBatches = self.m_RenderBatches[:0]
// 	self.m_SpriteGlyphs = self.m_SpriteGlyphs[:0]
// }

// func (self *SpriteBatch) createRenderBatches() {
// 	vertices := make([]Vertex, len(self.m_SpriteGlyphs)*6)

// 	if len(vertices) == 0 {
// 		return
// 	}

// 	m_Offset := 0
// 	vertexNum := 0

// 	self.m_RenderBatches = append(self.m_RenderBatches, NewRenderBatch(m_Offset, 6, 1, self.m_SpriteGlyphs[0].m_Texture))
// 	vertices[vertexNum] = self.m_SpriteGlyphs[0].m_Bottomleft
// 	vertexNum++
// 	vertices[vertexNum] = self.m_SpriteGlyphs[0].m_Topright
// 	vertexNum++
// 	vertices[vertexNum] = self.m_SpriteGlyphs[0].m_Bottomright
// 	vertexNum++
// 	vertices[vertexNum] = self.m_SpriteGlyphs[0].m_Bottomleft
// 	vertexNum++
// 	vertices[vertexNum] = self.m_SpriteGlyphs[0].m_Topright
// 	vertexNum++
// 	vertices[vertexNum] = self.m_SpriteGlyphs[0].m_Topleft
// 	vertexNum++
// 	m_Offset += 6

// 	for i := 1; i < len(self.m_SpriteGlyphs); i++ {
// 		if !self.m_SpriteGlyphs[i-1].m_Texture.m_TextureId.Equal(self.m_SpriteGlyphs[i].m_Texture.m_TextureId) {
// 			self.m_RenderBatches = append(self.m_RenderBatches, NewRenderBatch(m_Offset, 6, 1, self.m_SpriteGlyphs[i].m_Texture))
// 		} else {
// 			self.m_RenderBatches[len(self.m_RenderBatches)-1].m_NumberOfVertices += 6
// 		}

// 		vertices[vertexNum] = self.m_SpriteGlyphs[i].m_Bottomleft
// 		vertexNum++
// 		vertices[vertexNum] = self.m_SpriteGlyphs[i].m_Topright
// 		vertexNum++
// 		vertices[vertexNum] = self.m_SpriteGlyphs[i].m_Bottomright
// 		vertexNum++
// 		vertices[vertexNum] = self.m_SpriteGlyphs[i].m_Bottomleft
// 		vertexNum++
// 		vertices[vertexNum] = self.m_SpriteGlyphs[i].m_Topright
// 		vertexNum++
// 		vertices[vertexNum] = self.m_SpriteGlyphs[i].m_Topleft
// 		vertexNum++
// 		m_Offset += 6
// 		self.m_RenderBatches[len(self.m_RenderBatches)-1].m_NumOfInstances += 1
// 	}
// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), self.m_Vbo)

// 	jsVerts := vertexBufferToJsVertexBuffer(vertices)

// 	canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), jsVerts, canvasContext.Get("STATIC_DRAW"))

// 	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), js.Null())

// }
///////// OLD RENDERER //////////////////////////////////////////////

var canvasContext js.Value

type RGBA8 struct {
	m_R, m_G, m_B, m_A uint8
}

func (rgba8 *RGBA8) SetColorR(m_R uint8) {
	rgba8.m_R = m_R
}

func (rgba8 *RGBA8) SetColorRFloat32(m_R float32) {
	rgba8.m_R = Float1ToUint8_255(m_R)
}

func (rgba8 *RGBA8) SetColorG(m_G uint8) {
	rgba8.m_G = m_G
}

func (rgba8 *RGBA8) SetColorGFloat32(m_G float32) {
	rgba8.m_G = Float1ToUint8_255(m_G)
}

func (rgba8 *RGBA8) SetColorB(m_B uint8) {
	rgba8.m_B = m_B
}

func (rgba8 *RGBA8) SetColorBFloat32(m_B float32) {
	rgba8.m_B = Float1ToUint8_255(m_B)
}

func (rgba8 *RGBA8) SetColorA(m_A uint8) {
	rgba8.m_A = m_A
}

func (rgba8 *RGBA8) SetColorAFloat32(m_A float32) {
	rgba8.m_A = Float1ToUint8_255(m_A)
}

func (rgba8 *RGBA8) GetColorRFloat32() float32 {
	return Uint8ToFloat1(rgba8.m_R)
}
func (rgba8 *RGBA8) GetColorGFloat32() float32 {
	return Uint8ToFloat1(rgba8.m_G)
}
func (rgba8 *RGBA8) GetColorBFloat32() float32 {
	return Uint8ToFloat1(rgba8.m_B)
}
func (rgba8 *RGBA8) GetColorAFloat32() float32 {
	return Uint8ToFloat1(rgba8.m_A)
}

func NewRGBA8(m_R, m_G, m_B, m_A uint8) RGBA8 {
	return RGBA8{m_R, m_G, m_B, m_A}
}

func NewRGBA8Float(m_R, m_G, m_B, m_A float32) RGBA8 {
	return RGBA8{Float1ToUint8_255(m_R), Float1ToUint8_255(m_G), Float1ToUint8_255(m_B), Float1ToUint8_255(m_A)}
}

func GetRandomRGBA8() RGBA8 {
	return NewRGBA8(Float1ToUint8_255(rand.Float32()), Float1ToUint8_255(rand.Float32()), Float1ToUint8_255(rand.Float32()), 255)
}

var (
	CLEAR   = RGBA8{0, 0, 0, 0}
	WHITE   = RGBA8{255, 255, 255, 255}
	BLACK   = RGBA8{0, 0, 0, 255}
	GREY    = RGBA8{125, 125, 125, 255}
	RED     = RGBA8{255, 0, 0, 255}
	GREEN   = RGBA8{0, 255, 0, 255}
	BLUE    = RGBA8{0, 0, 255, 255}
	YELLOW  = RGBA8{255, 255, 0, 255}
	CYAN    = RGBA8{0, 255, 255, 255}
	MAGENTA = RGBA8{255, 0, 255, 255}
	ORANGE  = RGBA8{255, 180, 0, 255}
)

var lineWidth float32 = 0.1

func SetLineWidth(_newWidth float32) {
	lineWidth = _newWidth
}

// These are wrapper functions for the Renderer, I don't want Renderer2D to be cluttered with little functions that change some functionality, like I did with the SpriteBatch struct above.
func DrawLine(_p1, _p2 Vector2f, _c RGBA8, _z int) {
	DrawLineWRenderer(&Renderer, _p1, _p2, _c, _z)
}

func DrawLineWRenderer(_renderer *Renderer2D, _p1, _p2 Vector2f, _c RGBA8, _z int) {
	rotatingOffset := _p2.Subtract(_p1)
	rotatingOffset = rotatingOffset.Perpendicular()
	rotatingOffset = rotatingOffset.Normalize()
	rotatingOffset = rotatingOffset.Scale(lineWidth)
	_renderer.InsertQuadWithCoordinates(customtypes.Pair[Vector2f, Vector2f]{First: _p1.Subtract(rotatingOffset), Second: _p1.Add(rotatingOffset)},
		customtypes.Pair[Vector2f, Vector2f]{First: _p2.Subtract(rotatingOffset), Second: _p2.Add(rotatingOffset)},
		_z, _c)
}

func DrawRect(_center, _dimensions Vector2f, _c RGBA8, _z int, _r float32) {
	DrawRectWRenderer(&Renderer, _center, _dimensions, _c, _z, _r)
}

func DrawRectWRenderer(_renderer *Renderer2D, _center, _dimensions Vector2f, _c RGBA8, _z int, _r float32) {
	lineDifference := lineWidth

	lbPoint := _center.AddXY(-(_dimensions.X / 2.0), -(_dimensions.Y / 2.0))
	ltPoint := _center.AddXY(-(_dimensions.X / 2.0), (_dimensions.Y / 2.0))
	rtPoint := _center.AddXY((_dimensions.X / 2.0), (_dimensions.Y / 2.0))
	rbPoint := _center.AddXY((_dimensions.X / 2.0), -(_dimensions.Y / 2.0))

	DrawLineWRenderer(_renderer, lbPoint.AddXY(0.0, lineDifference).Rotate(_r, _center), ltPoint.AddXY(0.0, -lineDifference).Rotate(_r, _center), _c, _z)
	DrawLineWRenderer(_renderer, ltPoint.AddXY(-lineDifference, 0.0).Rotate(_r, _center), rtPoint.AddXY(lineDifference, 0.0).Rotate(_r, _center), _c, _z)
	DrawLineWRenderer(_renderer, rtPoint.AddXY(0.0, -lineDifference).Rotate(_r, _center), rbPoint.AddXY(0.0, lineDifference).Rotate(_r, _center), _c, _z)
	DrawLineWRenderer(_renderer, rbPoint.AddXY(lineDifference, 0.0).Rotate(_r, _center), lbPoint.AddXY(-lineDifference, 0.0).Rotate(_r, _center), _c, _z)
}

func DrawCircle(_center Vector2f, _radius float32, _c RGBA8, _z int) {
	DrawCircleWRenderer(&Renderer, _center, _radius, _c, _z)
}

func DrawCircleWRenderer(_renderer *Renderer2D, _center Vector2f, _radius float32, _c RGBA8, _z int) {
	numOfVertices := 16

	pos := [16]Vector2f{}
	for i := 0; i < numOfVertices; i++ {
		angle := (float32(i) / float32(numOfVertices)) * 2.0 * PI
		pos[i] = NewVector2f(_center.X+float32(math.Cos(float64(angle)))*_radius, _center.Y+float32(math.Sin(float64(angle)))*_radius)
	}

	for i := 0; i < numOfVertices-1; i++ {
		DrawLineWRenderer(_renderer, pos[i], pos[i+1], _c, _z)
	}
	DrawLineWRenderer(_renderer, pos[numOfVertices-1], pos[0], _c, _z)
}

/////////////////////////////////////////////////////

type renderObjectFuncType = func(v1, v2, v3, v4 Vector2f, c RGBA8, r float32, z int, t *Texture2D)

var (
	QUAD_RENDEROBJECTTYPEFUNC renderObjectFuncType = func(v1, v2, v3, v4 Vector2f, c RGBA8, r float32, z int, t *Texture2D) {
		// sh.DrawFillRectRotated(v1, z, v2, c, r)
		Renderer.InsertQuadRotated(v1.Subtract(v2.Scale(0.5)), v2, z, c, r)
	}
	// TRI_RENDEROBJECTTYPEFUNC renderObjectFuncType = func(sh *ShapeBatch, sp *SpriteBatch, v1, v2, v3, v4 Vector2f, c RGBA8, r float32, z int, t *Texture2D) {
	// 	sh.DrawFillTriangleRotated(v1, z, v2, c, r)
	// }
	// CIRCLE_RENDEROBJECTTYPEFUNC renderObjectFuncType = func(v1, v2, v3, v4 Vector2f, c RGBA8, r, z float32, t *Texture2D) {
	// 	sh.DrawCircle(v1, z, v2.X, c)
	// }
	LINE_RENDEROBJECTTYPEFUNC renderObjectFuncType = func(v1, v2, v3, v4 Vector2f, c RGBA8, r float32, z int, t *Texture2D) {
		// DrawLine(v1, v1.Add(v2), c, z)
	}
	SPRITE_RENDEROBJECTTYPEFUNC renderObjectFuncType = func(v1, v2, v3, v4 Vector2f, c RGBA8, r float32, z int, t *Texture2D) {
		// sp.DrawSpriteOriginScaledRotated(v1, v3, v4, v2.X, z, t, c, r)
		// Sprites.DrawSpriteBottomLeft(v.Position, v.Dimensions, v.UV1, v.UV2, &tmc.tileset.m_Texture, WHITE)
		// sp.DrawSpriteBottomLeft(v1, v2, v3, v4, t, c)
		// sp.DrawSpriteBottomLeftRotated(v1, v2, v3, v4, z, t, c, r)
		Renderer.InsertQuadTex(v1, v2, customtypes.Pair[Vector2f, Vector2f]{First: v3, Second: v4}, z, c, t)

	}
	SPRITECENTER_RENDEROBJECTTYPEFUNC renderObjectFuncType = func(v1, v2, v3, v4 Vector2f, c RGBA8, r float32, z int, t *Texture2D) {
		// sp.DrawSpriteOriginScaledDimRotated(v1, v2.Multp(t.m_PixelsToMeterDimensions), v3, v4, 1.0, z, t, c, r)
		newDim := v2.Multp(t.m_PixelsToMeterDimensions)
		Renderer.InsertQuadTexRotated(v1.Subtract(newDim.Scale(0.5)), newDim, customtypes.Pair[Vector2f, Vector2f]{First: v3, Second: v4}, z, c, t, r)
	}
)

type RenderObject struct {
	m_EntId      EntId
	m_ObjectType renderObjectFuncType
	m_Texture    *Texture2D
}

func newRenderObject(_entId EntId, _objectType renderObjectFuncType) RenderObject {
	return RenderObject{
		m_EntId:      _entId,
		m_ObjectType: _objectType,
	}
}

func NewQuadComponent(_thisScene *Scene, _entId EntId, vt VisualTransform, _static bool) FillRectRenderComponent {
	renderObj := newRenderObject(_entId, QUAD_RENDEROBJECTTYPEFUNC)
	renderObj.m_EntId = _entId
	if _static {
		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{vt, renderObj}, Rect{Position: vt.Position.Subtract(vt.Dimensions.Scale(0.5)), Size: vt.Dimensions})
	} else {
		// _index = DynamicRenderQuadTreeContainer.Insert(renderObj, Rect{Position: t.Position.Subtract(t.Dimensions.Scale(0.5)), Size: t.Dimensions})
		DynamicRenderQuadTreeContainer.InsertWithIndex(renderObj, Rect{Position: vt.Position.Subtract(vt.Dimensions.Scale(0.5)), Size: vt.Dimensions}, int64(_entId))

	}

	return FillRectRenderComponent{}
}

// func NewTriangleComponent(_thisScene *Scene, _entId EntId, _vt VisualTransform, _static bool) FillTriangleRenderComponent {
// 	renderObj := newRenderObject(_entId, TRI_RENDEROBJECTTYPEFUNC)
// 	if _static {
// 		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{_vt, renderObj}, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(0.5)), Size: _vt.Dimensions})

// 	} else {
// 		DynamicRenderQuadTreeContainer.Insert(renderObj, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(0.5)), Size: _vt.Dimensions})

// 	}

// 	return FillTriangleRenderComponent{}
// }

// func NewCircleComponent(_thisScene *Scene, _entId EntId, _vt VisualTransform, _static bool) CircleRenderComponent {
// 	renderObj := newRenderObject(_entId, CIRCLE_RENDEROBJECTTYPEFUNC)
// 	_vt.Dimensions.X /= 2.0
// 	_vt.Dimensions.Y /= 2.0
// 	if _static {
// 		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{_vt, renderObj}, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(0.5)), Size: _vt.Dimensions})

// 	} else {
// 		DynamicRenderQuadTreeContainer.Insert(renderObj, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(0.5)), Size: _vt.Dimensions})
// 	}

// 	return CircleRenderComponent{}
// }

func NewLineComponent(_thisScene *Scene, _entId EntId, _vt VisualTransform, _static bool) LineRenderComponent {
	renderObj := newRenderObject(_entId, LINE_RENDEROBJECTTYPEFUNC)
	if _static {
		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{_vt, renderObj}, Rect{Position: _vt.Position, Size: _vt.Dimensions})

	} else {
		DynamicRenderQuadTreeContainer.Insert(renderObj, Rect{Position: _vt.Position, Size: _vt.Dimensions})
	}

	return LineRenderComponent{}
}

func NewSpriteComponentBL(_thisScene *Scene, _entId EntId, _vt VisualTransform, _texture *Texture2D, _static bool) SpriteComponent {
	renderObj := newRenderObject(_entId, SPRITE_RENDEROBJECTTYPEFUNC)
	renderObj.m_Texture = _texture
	if _static {
		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{_vt, renderObj}, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(.5)), Size: _vt.Dimensions})

	} else {
		DynamicRenderQuadTreeContainer.Insert(renderObj, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(0.5)), Size: _vt.Dimensions})
	}

	return SpriteComponent{}
}

func NewSpriteComponent(_thisScene *Scene, _entId EntId, _vt VisualTransform, _texture *Texture2D, _static bool) SpriteComponent {
	renderObj := newRenderObject(_entId, SPRITECENTER_RENDEROBJECTTYPEFUNC)
	renderObj.m_Texture = _texture
	if _static {
		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{_vt, renderObj}, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(.5)), Size: _vt.Dimensions})

	} else {
		DynamicRenderQuadTreeContainer.Insert(renderObj, Rect{Position: _vt.Position.Subtract(_vt.Dimensions.Scale(0.5)), Size: _vt.Dimensions})
	}

	return SpriteComponent{}
}
