package chai

import (
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
in float z;
in vec4 vertexColor;
in float texIndex;

out vec2 frag_VertexUV;
out vec4 frag_VertexColor;
out float frag_TexIndex;

uniform mat4 u_View_Matrix;

void main()
{
	gl_Position.xy = (u_View_Matrix * vec4(vertexPosition.x, vertexPosition.y, z, 1.0)).xy;
    gl_Position.z = 0.0;
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
	coordinates, tex_coordinates Vector2f
	z                            float32
	color                        RGBA8
	texture_index                float32
}

type Renderer2D struct {
	vao, vbo, ebo js.Value
	maximum_quads int
	maximum_z     float32

	num_of_elements  int
	num_of_vertices  int
	num_of_quads     int
	num_of_drawcalls int

	vertices         []NVertex
	jsVertsBuffer    js.Value
	bytesVertsBuffer []byte

	_tex_sampler_jsbuffer js.Value
	_tex_sampler_js       js.Value
	_view_matrix_jsbuffer js.Value
	_view_matrix_js       js.Value

	white_texture          js.Value
	white_textureSlotIndex uint32
	texture_slots          [6]js.Value
	texture_slot_index     uint32

	shader ShaderProgram
	camera Camera2D
	app    *App
}

func (_renderer *Renderer2D) InsertQuad(_center Vector2f, _z float32, _tint RGBA8) {
	if _renderer.num_of_quads >= _renderer.maximum_quads {
		_renderer.End()
		_renderer.Render()
		_renderer.Begin()
	}

	_half_dims := NewVector2f(35.0, 35.0)
	_new_z := RemapFloat32(_z, 0, _renderer.maximum_z, -0.9999, 9.9999)
	_vertex_index := _renderer.num_of_vertices

	_renderer.vertices[_vertex_index+0] = NVertex{_center.Add(_half_dims.Scale(-1.0)), Vector2f{0.0, 0.0}, _new_z, _tint, float32(_renderer.white_textureSlotIndex)}
	_renderer.vertices[_vertex_index+1] = NVertex{_center.Add(NewVector2f(_half_dims.X, -_half_dims.Y)), Vector2f{0.0, 1.0}, _new_z, _tint, float32(_renderer.white_textureSlotIndex)}
	_renderer.vertices[_vertex_index+2] = NVertex{_center.Add(_half_dims.Scale(1.0)), Vector2f{1.0, 1.0}, _new_z, _tint, float32(_renderer.white_textureSlotIndex)}
	_renderer.vertices[_vertex_index+3] = NVertex{_center.Add(NewVector2f(-_half_dims.X, _half_dims.Y)), Vector2f{1.0, 0.0}, _new_z, _tint, float32(_renderer.white_textureSlotIndex)}

	_renderer.num_of_elements += 6
	_renderer.num_of_vertices += 4
	_renderer.num_of_quads += 1
}

func (_renderer *Renderer2D) InsertQuadTex(_center Vector2f, _z float32, _tint RGBA8, _texture *Texture2D) {
	if _renderer.num_of_quads >= _renderer.maximum_quads || _renderer.texture_slot_index > MAX_TEXTURES-1 {
		_renderer.End()
		_renderer.Render()
		_renderer.Begin()
	}
	_texture_index := float32(0.0)
	for i := uint32(1); i < _renderer.texture_slot_index; i++ {
		if _renderer.texture_slots[i].Equal(_texture.textureId) {
			_texture_index = float32(i)
		}
	}
	if _texture_index == 0.0 {
		_texture_index = float32(_renderer.texture_slot_index)
		_renderer.texture_slots[_renderer.texture_slot_index] = _texture.textureId
		_renderer.texture_slot_index += 1
	}

	_half_dims := NewVector2f(35.0, 35.0)
	_new_z := RemapFloat32(_z, 0, _renderer.maximum_z, -0.9999, 9.9999)
	_vertex_index := _renderer.num_of_vertices

	_renderer.vertices[_vertex_index+0] = NVertex{_center.Add(_half_dims.Scale(-1.0)), Vector2f{0.0, 0.0}, _new_z, _tint, float32(_texture_index)}
	_renderer.vertices[_vertex_index+1] = NVertex{_center.Add(NewVector2f(_half_dims.X, -_half_dims.Y)), Vector2f{0.0, 1.0}, _new_z, _tint, float32(_texture_index)}
	_renderer.vertices[_vertex_index+2] = NVertex{_center.Add(_half_dims.Scale(1.0)), Vector2f{1.0, 1.0}, _new_z, _tint, float32(_texture_index)}
	_renderer.vertices[_vertex_index+3] = NVertex{_center.Add(NewVector2f(-_half_dims.X, _half_dims.Y)), Vector2f{1.0, 0.0}, _new_z, _tint, float32(_texture_index)}

	_renderer.num_of_elements += 6
	_renderer.num_of_vertices += 4
	_renderer.num_of_quads += 1
}

func (_renderer *Renderer2D) Begin() {
	_renderer.num_of_vertices = 0
	_renderer.num_of_elements = 0
	_renderer.num_of_quads = 0

	_renderer.texture_slot_index = 1
	for i := 1; i < MAX_TEXTURES; i++ {
		_renderer.texture_slots[i] = js.ValueOf(0)
	}
}

func (_renderer *Renderer2D) End() {
	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _renderer.vbo)
	// _jsNVerts := nVertexBufferToJsNVertexBuffer(_renderer.vertices)
	setNVertexBufferAndJsBuffer(&_renderer.jsVertsBuffer, _renderer.bytesVertsBuffer, _renderer.vertices)
	canvasContext.Call("bufferSubData", canvasContext.Get("ARRAY_BUFFER"), 0, _renderer.jsVertsBuffer)
	// canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), _jsNVerts, canvasContext.Get("DYNAMIC_DRAW"))
	// _jsNVerts.Set("delete", nil)
}

func (_renderer *Renderer2D) Render() {
	_renderer.camera.Update(*_renderer.app)

	UseShader(&_renderer.shader)

	// Sending the samplers to JS
	/////////////////////////////
	_tex_samplers := [MAX_TEXTURES]int32{}
	_tex_samplers_uniloc := canvasContext.Call("getUniformLocation", _renderer.shader.ShaderProgramID, "u_Textures")
	for i := 0; i < len(_tex_samplers); i++ {
		_tex_samplers[i] = int32(i)
		_renderer._tex_sampler_js.SetIndex(i, js.ValueOf(_tex_samplers[i]))
	}
	canvasContext.Call("uniform1iv", _tex_samplers_uniloc, _renderer._tex_sampler_js)
	//////////////////////////////

	// Sending the ViewMatrix to JS
	///////////////////////////////
	viewMatrix := _renderer.camera.viewMatrix.Data()

	for i := 0; i < len(viewMatrix); i++ {
		_renderer._view_matrix_js.SetIndex(i, js.ValueOf(viewMatrix[i]))
	}
	viewmatrix_loc := canvasContext.Call("getUniformLocation", _renderer.shader.ShaderProgramID, "u_View_Matrix")
	canvasContext.Call("uniformMatrix4fv", viewmatrix_loc, false, _renderer._view_matrix_js)
	///////////////////////////////

	for i := uint32(0); i < _renderer.texture_slot_index; i++ {
		canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE1"))
		canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), _renderer.texture_slots[i])
	}

	canvasContext.Call("bindVertexArray", _renderer.vao)
	canvasContext.Call("drawElements", canvasContext.Get("TRIANGLES"), _renderer.num_of_elements, canvasContext.Get("UNSIGNED_INT"), 0)
	canvasContext.Call("bindVertexArray", js.Null())
}

func NewRenderer(_app *App, _max_quads int) Renderer2D {
	_renderer := Renderer2D{}

	_MAX_VERTEX_COUNT := 4 * _max_quads
	_MAX_ELEMENTS_COUNT := 6 * _max_quads

	_renderer.vertices = make([]NVertex, _MAX_VERTEX_COUNT)
	_renderer.jsVertsBuffer, _renderer.bytesVertsBuffer = makeNVertexBufferAndJsBuffer(_renderer.vertices)

	_renderer.maximum_quads = _max_quads
	_renderer.maximum_z = 1000
	_renderer.app = _app

	_renderer.vao = canvasContext.Call("createVertexArray")
	canvasContext.Call("bindVertexArray", _renderer.vao)

	_renderer.vbo = canvasContext.Call("createBuffer")
	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _renderer.vbo)
	// jsNVerts := nVertexBufferToJsNVertexBuffer(_renderer.vertices)
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

	_renderer.ebo = canvasContext.Call("createBuffer")
	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), _renderer.ebo)
	jsElements := int32BufferToJsInt32Buffer(_elements)
	canvasContext.Call("bufferData", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), jsElements, canvasContext.Get("STATIC_DRAW"))

	_renderer.white_texture = loadWhiteTexture().textureId
	_renderer.texture_slots[0] = _renderer.white_texture
	_renderer.texture_slot_index = 1
	for i := 1; i < MAX_TEXTURES; i++ {
		_renderer.texture_slots[i] = js.ValueOf(0)
	}

	canvasContext.Call("bindVertexArray", js.Null())
	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), js.Null())
	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), js.Null())

	canvasContext.Call("disableVertexAttribArray", 0)
	canvasContext.Call("disableVertexAttribArray", 1)
	canvasContext.Call("disableVertexAttribArray", 2)
	canvasContext.Call("disableVertexAttribArray", 3)
	canvasContext.Call("disableVertexAttribArray", 4)

	_renderer.camera.Init(*_app)
	_renderer.camera.centerOffset = NewVector2f(float32(_app.Width)/2.0, float32(_app.Height)/2.0)
	_renderer.camera.Update(*_app)

	_renderer.shader.ParseShader(vERTEX_SHADER, fRAGMENT_SHADER)
	_renderer.shader.CreateShaderProgram()
	_renderer.shader.AddAttribute("vertexPosition")
	_renderer.shader.AddAttribute("vertexUv")
	_renderer.shader.AddAttribute("z")
	_renderer.shader.AddAttribute("vertexColor")
	_renderer.shader.AddAttribute("texIndex")

	////////////////////////////////////////////////
	////////////////////////////////////////////////
	_tex_samplers := [MAX_TEXTURES]int32{}
	_tex_samplers_viewbuffer := new(bytes.Buffer)
	{
		err := binary.Write(_tex_samplers_viewbuffer, binary.LittleEndian, _tex_samplers[:])
		Assert(err == nil, "Error Writing Texture Sampler Unifrom")
	}
	_tex_byte_slice := _tex_samplers_viewbuffer.Bytes()
	_renderer._tex_sampler_jsbuffer = js.Global().Get("Uint8Array").New(len(_tex_samplers) * 4)
	js.CopyBytesToJS(_renderer._tex_sampler_jsbuffer, _tex_byte_slice)
	_renderer._tex_sampler_js = js.Global().Get("Int32Array").New(_renderer._tex_sampler_jsbuffer.Get("buffer"), _renderer._tex_sampler_jsbuffer.Get("byteOffset"), _renderer._tex_sampler_jsbuffer.Get("byteLength").Int()/4)
	///////////////////////////////////////////////
	viewMatrix := _renderer.camera.viewMatrix.Data()
	viewBuffer := new(bytes.Buffer)
	{
		err := binary.Write(viewBuffer, binary.LittleEndian, viewMatrix[:])
		Assert(err == nil, "Error writing viewBuffer")
	}
	byte_slice2 := viewBuffer.Bytes()
	_renderer._view_matrix_jsbuffer = js.Global().Get("Uint8Array").New(len(viewMatrix) * 4)
	js.CopyBytesToJS(_renderer._view_matrix_jsbuffer, byte_slice2)
	_renderer._view_matrix_js = js.Global().Get("Float32Array").New(_renderer._view_matrix_jsbuffer.Get("buffer"), _renderer._view_matrix_jsbuffer.Get("byteOffset"), _renderer._view_matrix_jsbuffer.Get("byteLength").Int()/4)

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
