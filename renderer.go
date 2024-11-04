package chai

import (
	. "github.com/mhamedGd/chai/math"

	"bytes"
	"encoding/binary"
	"reflect"
	"syscall/js"
	"unsafe"
)

const vERTEX_SHADER = `#version 300 es

precision mediump float;

in vec2 vertexPosition;
in vec2 vertexUv;
in float m_Z;
in vec4 vertexColor;
in float texIndex;

out vec2 frag_VertexUV;
out vec4 frag_VertexColor;
out float frag_TexIndex;

uniform mat4 u_View_Matrix;

void main()
{
	gl_Position.xy = (u_View_Matrix * vec4(vertexPosition.x, vertexPosition.y, m_Z, 1.0)).xy;
    gl_Position.m_Z = 0.0;
	gl_Position.w = 1.0;

    frag_VertexUV = vertexUv;
    frag_VertexUV.y = 1.0 - frag_VertexUV.y;
	frag_VertexColor = vertexColor;
    frag_TexIndex = texIndex;
}`
const fRAGMENT_SHADER = `#version 300 es

precision mediump float;

out vec4 fragColor;

in vec2 frag_VertexUV;
in vec4 frag_VertexColor;
in float frag_TexIndex;

uniform sampler2D u_Textures[6];
uniform sampler2D tex01;
uniform sampler2D tex02;

vec4 TextureColor(int _index) {
    switch(_index){
        case 0:
            return texture(u_Textures[0], frag_VertexUV);
        case 1:
            return texture(u_Textures[1], frag_VertexUV);
        case 2:
            return texture(u_Textures[2], frag_VertexUV);
        case 3:
            return texture(u_Textures[3], frag_VertexUV);
        case 4:
            return texture(u_Textures[4], frag_VertexUV);
        case 5:
            return texture(u_Textures[5], frag_VertexUV);
    }
    return vec4(1.0);
}

void main()
{
    int _index = int(frag_TexIndex);
    vec4 _finalColor = vec4(1.0);
    
    fragColor = frag_VertexColor * TextureColor(_index);
}`

const (
	MAX_TEXTURES = 6
	NVERTEX_SIZE = 28
)

type NVertex struct {
	m_Coordinates, m_TexCoordinates Vector2f
	m_Z                             float32
	m_Color                         RGBA8
	m_textureIndex                  float32
}

type Renderer2D struct {
	m_Vao, m_Vbo, m_Ebo js.Value
	m_MaximumQuads      int
	m_MaximumZ          float32

	m_NumOfElements  int
	m_NumOfVertices  int
	m_NumOfQuads     int
	m_NumOfDrawcalls int

	m_Vertices         []NVertex
	m_JsVertsBuffer    js.Value
	m_BytesVertsBuffer []byte

	m_TexSamplerJsbuffer js.Value
	m_TexSamplerJs       js.Value
	m_ViewMatrixJsbuffer js.Value
	m_ViewMatrixJs       js.Value

	m_WhiteTexture          js.Value
	m_WhiteTextureSlotIndex uint32
	m_TextureSlots          [6]js.Value
	m_TextureSlotIndex      uint32

	m_Shader ShaderProgram
	m_Camera Camera2D
	m_App    *App
}

func (_renderer *Renderer2D) InsertQuad(_center Vector2f, _z float32, _tint RGBA8) {
	if _renderer.m_NumOfQuads >= _renderer.m_MaximumQuads {
		_renderer.End()
		_renderer.Render()
		_renderer.Begin()
	}

	_half_dims := NewVector2f(35.0, 35.0)
	_new_z := RemapFloat32(_z, 0, _renderer.m_MaximumZ, -0.9999, 9.9999)
	_vertex_index := _renderer.m_NumOfVertices

	_renderer.m_Vertices[_vertex_index+0] = NVertex{_center.Add(_half_dims.Scale(-1.0)), Vector2f{0.0, 0.0}, _new_z, _tint, float32(_renderer.m_WhiteTextureSlotIndex)}
	_renderer.m_Vertices[_vertex_index+1] = NVertex{_center.Add(NewVector2f(_half_dims.X, -_half_dims.Y)), Vector2f{0.0, 1.0}, _new_z, _tint, float32(_renderer.m_WhiteTextureSlotIndex)}
	_renderer.m_Vertices[_vertex_index+2] = NVertex{_center.Add(_half_dims.Scale(1.0)), Vector2f{1.0, 1.0}, _new_z, _tint, float32(_renderer.m_WhiteTextureSlotIndex)}
	_renderer.m_Vertices[_vertex_index+3] = NVertex{_center.Add(NewVector2f(-_half_dims.X, _half_dims.Y)), Vector2f{1.0, 0.0}, _new_z, _tint, float32(_renderer.m_WhiteTextureSlotIndex)}

	_renderer.m_NumOfElements += 6
	_renderer.m_NumOfVertices += 4
	_renderer.m_NumOfQuads += 1
}

func (_renderer *Renderer2D) InsertQuadTex(_center Vector2f, _z float32, _tint RGBA8, _texture *Texture2D) {
	if _renderer.m_NumOfQuads >= _renderer.m_MaximumQuads || _renderer.m_TextureSlotIndex > MAX_TEXTURES-1 {
		_renderer.End()
		_renderer.Render()
		_renderer.Begin()
	}
	_texture_index := float32(0.0)
	for i := uint32(1); i < _renderer.m_TextureSlotIndex; i++ {
		if _renderer.m_TextureSlots[i].Equal(_texture.m_TextureId) {
			_texture_index = float32(i)
		}
	}
	if _texture_index == 0.0 {
		_texture_index = float32(_renderer.m_TextureSlotIndex)
		_renderer.m_TextureSlots[_renderer.m_TextureSlotIndex] = _texture.m_TextureId
		_renderer.m_TextureSlotIndex += 1
	}

	_half_dims := NewVector2f(35.0, 35.0)
	_new_z := RemapFloat32(_z, 0, _renderer.m_MaximumZ, -0.9999, 9.9999)
	_vertex_index := _renderer.m_NumOfVertices

	_renderer.m_Vertices[_vertex_index+0] = NVertex{_center.Add(_half_dims.Scale(-1.0)), Vector2f{0.0, 0.0}, _new_z, _tint, float32(_texture_index)}
	_renderer.m_Vertices[_vertex_index+1] = NVertex{_center.Add(NewVector2f(_half_dims.X, -_half_dims.Y)), Vector2f{0.0, 1.0}, _new_z, _tint, float32(_texture_index)}
	_renderer.m_Vertices[_vertex_index+2] = NVertex{_center.Add(_half_dims.Scale(1.0)), Vector2f{1.0, 1.0}, _new_z, _tint, float32(_texture_index)}
	_renderer.m_Vertices[_vertex_index+3] = NVertex{_center.Add(NewVector2f(-_half_dims.X, _half_dims.Y)), Vector2f{1.0, 0.0}, _new_z, _tint, float32(_texture_index)}

	_renderer.m_NumOfElements += 6
	_renderer.m_NumOfVertices += 4
	_renderer.m_NumOfQuads += 1
}

func (_renderer *Renderer2D) Begin() {
	_renderer.m_NumOfVertices = 0
	_renderer.m_NumOfElements = 0
	_renderer.m_NumOfQuads = 0

	_renderer.m_TextureSlotIndex = 1
	for i := 1; i < MAX_TEXTURES; i++ {
		_renderer.m_TextureSlots[i] = js.ValueOf(0)
	}
}

func (_renderer *Renderer2D) End() {
	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _renderer.m_Vbo)
	// _jsNVerts := nVertexBufferToJsNVertexBuffer(_renderer.m_Vertices)
	setNVertexBufferAndJsBuffer(&_renderer.m_JsVertsBuffer, _renderer.m_BytesVertsBuffer, _renderer.m_Vertices)
	canvasContext.Call("bufferSubData", canvasContext.Get("ARRAY_BUFFER"), 0, _renderer.m_JsVertsBuffer)
	// canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), _jsNVerts, canvasContext.Get("DYNAMIC_DRAW"))
	// _jsNVerts.Set("delete", nil)
}

func (_renderer *Renderer2D) Render() {
	_renderer.m_Camera.Update(*_renderer.m_App)

	UseShader(&_renderer.m_Shader)

	// Sending the samplers to JS
	/////////////////////////////
	_tex_samplers := [MAX_TEXTURES]int32{}
	_tex_samplers_uniloc := canvasContext.Call("getUniformLocation", _renderer.m_Shader.ShaderProgramID, "u_Textures")
	for i := 0; i < len(_tex_samplers); i++ {
		_tex_samplers[i] = int32(i)
		_renderer.m_TexSamplerJs.SetIndex(i, js.ValueOf(_tex_samplers[i]))
	}
	canvasContext.Call("uniform1iv", _tex_samplers_uniloc, _renderer.m_TexSamplerJs)
	//////////////////////////////

	// Sending the ViewMatrix to JS
	///////////////////////////////
	viewMatrix := _renderer.m_Camera.m_ViewMatrix.Data()

	for i := 0; i < len(viewMatrix); i++ {
		_renderer.m_ViewMatrixJs.SetIndex(i, js.ValueOf(viewMatrix[i]))
	}
	viewmatrix_loc := canvasContext.Call("getUniformLocation", _renderer.m_Shader.ShaderProgramID, "u_View_Matrix")
	canvasContext.Call("uniformMatrix4fv", viewmatrix_loc, false, _renderer.m_ViewMatrixJs)
	///////////////////////////////

	for i := uint32(0); i < _renderer.m_TextureSlotIndex; i++ {
		canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE1"))
		canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), _renderer.m_TextureSlots[i])
	}

	canvasContext.Call("bindVertexArray", _renderer.m_Vao)
	canvasContext.Call("drawElements", canvasContext.Get("TRIANGLES"), _renderer.m_NumOfElements, canvasContext.Get("UNSIGNED_INT"), 0)
	canvasContext.Call("bindVertexArray", js.Null())
}

func NewRenderer(_app *App, _maxQuads int) Renderer2D {
	_renderer := Renderer2D{}

	_MAX_VERTEX_COUNT := 4 * _maxQuads
	_MAX_ELEMENTS_COUNT := 6 * _maxQuads

	_renderer.m_Vertices = make([]NVertex, _MAX_VERTEX_COUNT)
	_renderer.m_JsVertsBuffer, _renderer.m_BytesVertsBuffer = makeNVertexBufferAndJsBuffer(_renderer.m_Vertices)

	_renderer.m_MaximumQuads = _maxQuads
	_renderer.m_MaximumZ = 1000
	_renderer.m_App = _app

	_renderer.m_Vao = canvasContext.Call("createVertexArray")
	canvasContext.Call("bindVertexArray", _renderer.m_Vao)

	_renderer.m_Vbo = canvasContext.Call("createBuffer")
	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _renderer.m_Vbo)
	// jsNVerts := nVertexBufferToJsNVertexBuffer(_renderer.m_Vertices)
	canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), NVERTEX_SIZE*_MAX_VERTEX_COUNT, canvasContext.Get("DYNAMIC_DRAW"))

	canvasContext.Call("enableVertexAttribArray", 0)
	canvasContext.Call("vertexAttribPointer", 0, 2, canvasContext.Get("FLOAT"), false, NVERTEX_SIZE, 0)
	canvasContext.Call("enableVertexAttribArray", 1)
	canvasContext.Call("vertexAttribPointer", 1, 2, canvasContext.Get("FLOAT"), false, NVERTEX_SIZE, 8)
	canvasContext.Call("enableVertexAttribArray", 2)
	canvasContext.Call("vertexAttribPointer", 2, 1, canvasContext.Get("FLOAT"), false, NVERTEX_SIZE, 16)
	canvasContext.Call("enableVertexAttribArray", 3)
	canvasContext.Call("vertexAttribPointer", 3, 4, canvasContext.Get("UNSIGNED_BYTE"), true, NVERTEX_SIZE, 20)
	canvasContext.Call("enableVertexAttribArray", 4)
	canvasContext.Call("vertexAttribPointer", 4, 1, canvasContext.Get("FLOAT"), false, NVERTEX_SIZE, 24)

	_elements := make([]int32, _MAX_ELEMENTS_COUNT)
	_vert_offset := int32(0)
	for e := 0; e < _MAX_ELEMENTS_COUNT; e += 6 {
		_elements[e+0] = _vert_offset + 0
		_elements[e+1] = _vert_offset + 1
		_elements[e+2] = _vert_offset + 2

		_elements[e+3] = _vert_offset + 2
		_elements[e+4] = _vert_offset + 3
		_elements[e+5] = _vert_offset + 0
		_vert_offset += 4
	}

	_renderer.m_Ebo = canvasContext.Call("createBuffer")
	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), _renderer.m_Ebo)
	jsElements := int32BufferToJsInt32Buffer(_elements)
	canvasContext.Call("bufferData", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), jsElements, canvasContext.Get("STATIC_DRAW"))

	_renderer.m_WhiteTexture = loadWhiteTexture().m_TextureId
	_renderer.m_TextureSlots[0] = _renderer.m_WhiteTexture
	_renderer.m_TextureSlotIndex = 1
	for i := 1; i < MAX_TEXTURES; i++ {
		_renderer.m_TextureSlots[i] = js.ValueOf(0)
	}

	canvasContext.Call("bindVertexArray", js.Null())
	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), js.Null())
	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), js.Null())

	canvasContext.Call("disableVertexAttribArray", 0)
	canvasContext.Call("disableVertexAttribArray", 1)
	canvasContext.Call("disableVertexAttribArray", 2)
	canvasContext.Call("disableVertexAttribArray", 3)
	canvasContext.Call("disableVertexAttribArray", 4)

	_renderer.m_Camera.Init(*_app)
	_renderer.m_Camera.m_CenterOffset = NewVector2f(float32(_app.Width)/2.0, float32(_app.Height)/2.0)
	_renderer.m_Camera.Update(*_app)

	_renderer.m_Shader.ParseShader(vERTEX_SHADER, fRAGMENT_SHADER)
	_renderer.m_Shader.CreateShaderProgram()
	_renderer.m_Shader.AddAttribute("vertexPosition")
	_renderer.m_Shader.AddAttribute("vertexUv")
	_renderer.m_Shader.AddAttribute("m_Z")
	_renderer.m_Shader.AddAttribute("vertexColor")
	_renderer.m_Shader.AddAttribute("texIndex")

	////////////////////////////////////////////////
	////////////////////////////////////////////////
	_tex_samplers := [MAX_TEXTURES]int32{}
	_tex_samplers_viewbuffer := new(bytes.Buffer)
	{
		err := binary.Write(_tex_samplers_viewbuffer, binary.LittleEndian, _tex_samplers[:])
		Assert(err == nil, "Error Writing Texture Sampler Unifrom")
	}
	_tex_byte_slice := _tex_samplers_viewbuffer.Bytes()
	_renderer.m_TexSamplerJsbuffer = js.Global().Get("Uint8Array").New(len(_tex_samplers) * 4)
	js.CopyBytesToJS(_renderer.m_TexSamplerJsbuffer, _tex_byte_slice)
	_renderer.m_TexSamplerJs = js.Global().Get("Int32Array").New(_renderer.m_TexSamplerJsbuffer.Get("buffer"), _renderer.m_TexSamplerJsbuffer.Get("byteOffset"), _renderer.m_TexSamplerJsbuffer.Get("byteLength").Int()/4)
	///////////////////////////////////////////////
	viewMatrix := _renderer.m_Camera.m_ViewMatrix.Data()
	viewBuffer := new(bytes.Buffer)
	{
		err := binary.Write(viewBuffer, binary.LittleEndian, viewMatrix[:])
		Assert(err == nil, "Error writing viewBuffer")
	}
	byte_slice2 := viewBuffer.Bytes()
	_renderer.m_ViewMatrixJsbuffer = js.Global().Get("Uint8Array").New(len(viewMatrix) * 4)
	js.CopyBytesToJS(_renderer.m_ViewMatrixJsbuffer, byte_slice2)
	_renderer.m_ViewMatrixJs = js.Global().Get("Float32Array").New(_renderer.m_ViewMatrixJsbuffer.Get("buffer"), _renderer.m_ViewMatrixJsbuffer.Get("byteOffset"), _renderer.m_ViewMatrixJsbuffer.Get("byteLength").Int()/4)

	return _renderer
}

func nVertexBufferToJsNVertexBuffer(_buffer []NVertex) js.Value {
	jsVerts := js.Global().Get("Uint8Array").New(len(_buffer) * NVERTEX_SIZE)
	var verticesBytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&verticesBytes))
	header.Cap = cap(_buffer) * NVERTEX_SIZE
	header.Len = len(_buffer) * NVERTEX_SIZE
	header.Data = uintptr(unsafe.Pointer(&_buffer[0]))

	js.CopyBytesToJS(jsVerts, verticesBytes)
	return jsVerts
}

func makeNVertexBufferAndJsBuffer(_buffer []NVertex) (js.Value, []byte) {
	jsVerts := js.Global().Get("Uint8Array").New(len(_buffer) * NVERTEX_SIZE)
	var verticesBytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&verticesBytes))
	header.Cap = cap(_buffer) * NVERTEX_SIZE
	header.Len = len(_buffer) * NVERTEX_SIZE
	header.Data = uintptr(unsafe.Pointer(&_buffer[0]))

	return jsVerts, verticesBytes
}

func setNVertexBufferAndJsBuffer(_jsBuffer *js.Value, _nvertexBuffer []byte, _buffer []NVertex) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&_nvertexBuffer))
	header.Cap = cap(_buffer) * NVERTEX_SIZE
	header.Len = len(_buffer) * NVERTEX_SIZE
	header.Data = uintptr(unsafe.Pointer(&_buffer[0]))

	js.CopyBytesToJS(*_jsBuffer, _nvertexBuffer)
}

func uint32BufferToJsUInt32Buffer(_buffer []uint32) js.Value {
	jsElements := js.Global().Get("Uint8Array").New(len(_buffer) * 4)
	var elementsBytes []byte
	headerElem := (*reflect.SliceHeader)(unsafe.Pointer(&elementsBytes))
	headerElem.Cap = cap(_buffer) * 4
	headerElem.Len = len(_buffer) * 4
	headerElem.Data = uintptr(unsafe.Pointer(&_buffer[0]))
	js.CopyBytesToJS(jsElements, elementsBytes)

	return jsElements
}
