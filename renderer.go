package chai

import (
	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"

	"bytes"
	"encoding/binary"
	"reflect"
	"sort"
	"syscall/js"
	"unsafe"
)

const vERTEX_SHADER = `#version 300 es

precision mediump float;

in vec2 vertexPosition;
in vec2 vertexUv;
in vec4 vertexColor;
in float texIndex;

out vec2 frag_VertexUV;
out vec4 frag_VertexColor;
out float frag_TexIndex;

uniform mat4 u_View_Matrix;

void main()
{
	
	vec4 global_position = vec4(0.0);
	global_position = u_View_Matrix * vec4(vertexPosition, 0.0, 1.0);
	global_position.w = 1.0;		
	gl_Position = global_position;

    frag_VertexUV = vertexUv;
    // frag_VertexUV.y = 1.0 - frag_VertexUV.y;
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
    vec4 _finalColor = TextureColor(_index)*frag_VertexColor;
    
    fragColor = _finalColor * frag_VertexColor.a;
}`

type Renderer2D struct {
	m_RenderBatches customtypes.List[renderBatch2D]
	m_MaximumQuads  int
	m_App           *App
	m_Camera        *Camera2D
}

func NewRenderer2D(_maxQuads int, _app *App, _camera *Camera2D) Renderer2D {
	return Renderer2D{
		m_RenderBatches: customtypes.NewList[renderBatch2D](),
		m_MaximumQuads:  _maxQuads,
		m_App:           _app,
		m_Camera:        _camera,
	}
}

// //////////////////////
func (_r *Renderer2D) InsertQuad(_origin, _dims Vector2f, _z int, _tint RGBA8) {
	for i := 0; i < _r.m_RenderBatches.Count(); i++ {
		if _z == _r.m_RenderBatches.Data[i].m_Z && _r.m_RenderBatches.Data[i].HasRoom() {
			_r.m_RenderBatches.Data[i].insertQuad(_origin, _dims, _tint)
			return
		}
	}

	// If it didn't find a RenderBatch with this Z
	_r.m_RenderBatches.PushBack(newRenderer(_r.m_App, _r.m_MaximumQuads, _z, _r.m_Camera))
	_r.m_RenderBatches.Back().insertQuad(_origin, _dims, _tint)
}

func (_r *Renderer2D) InsertQuadRotated(_origin, _dims Vector2f, _z int, _tint RGBA8, _rotation float32) {
	for i := 0; i < _r.m_RenderBatches.Count(); i++ {
		if _z == _r.m_RenderBatches.Data[i].m_Z && _r.m_RenderBatches.Data[i].HasRoom() {
			_r.m_RenderBatches.Data[i].insertQuadRotated(_origin, _dims, _tint, _rotation)
			return
		}
	}

	// If it didn't find a RenderBatch with this Z
	_r.m_RenderBatches.PushBack(newRenderer(_r.m_App, _r.m_MaximumQuads, _z, _r.m_Camera))
	_r.m_RenderBatches.Back().insertQuadRotated(_origin, _dims, _tint, _rotation)
}

func (_r *Renderer2D) InsertQuadWithCoordinates(_leftBottomTop, _rightBottomTop customtypes.Pair[Vector2f, Vector2f], _z int, _tint RGBA8) {
	for i := 0; i < _r.m_RenderBatches.Count(); i++ {
		if _z == _r.m_RenderBatches.Data[i].m_Z && _r.m_RenderBatches.Data[i].HasRoom() {
			_r.m_RenderBatches.Data[i].insertQuadWithCoordinates(_leftBottomTop, _rightBottomTop, _tint)
			return
		}
	}

	// If it didn't find a RenderBatch with this Z
	_r.m_RenderBatches.PushBack(newRenderer(_r.m_App, _r.m_MaximumQuads, _z, _r.m_Camera))
	_r.m_RenderBatches.Back().insertQuadWithCoordinates(_leftBottomTop, _rightBottomTop, _tint)
}

func (_r *Renderer2D) InsertQuadTexWithCoordinates(_leftBottomTop, _rightBottomTop, _uv customtypes.Pair[Vector2f, Vector2f], _z int, _tint RGBA8, _texture *Texture2D) {
	for i := 0; i < _r.m_RenderBatches.Count(); i++ {
		if _z == _r.m_RenderBatches.Data[i].m_Z && _r.m_RenderBatches.Data[i].HasRoom() {
			_r.m_RenderBatches.Data[i].insertQuadTexWithCoordinates(_leftBottomTop, _rightBottomTop, _uv, _tint, _texture)
			return
		}
	}

	// If it didn't find a RenderBatch with this Z
	_r.m_RenderBatches.PushBack(newRenderer(_r.m_App, _r.m_MaximumQuads, _z, _r.m_Camera))
	_r.m_RenderBatches.Back().insertQuadTexWithCoordinates(_leftBottomTop, _rightBottomTop, _uv, _tint, _texture)
}

func (_r *Renderer2D) InsertQuadTexWithCoordinatesRotated(_leftBottomTop, _rightBottomTop, _uv customtypes.Pair[Vector2f, Vector2f], _z int, _tint RGBA8, _texture *Texture2D, _rotation float32) {
	for i := 0; i < _r.m_RenderBatches.Count(); i++ {
		if _z == _r.m_RenderBatches.Data[i].m_Z && _r.m_RenderBatches.Data[i].HasRoom() {
			_r.m_RenderBatches.Data[i].insertQuadTexWithCoordinatesRotated(_leftBottomTop, _rightBottomTop, _uv, _tint, _texture, _rotation)
			return
		}
	}

	// If it didn't find a RenderBatch with this Z
	_r.m_RenderBatches.PushBack(newRenderer(_r.m_App, _r.m_MaximumQuads, _z, _r.m_Camera))
	_r.m_RenderBatches.Back().insertQuadTexWithCoordinatesRotated(_leftBottomTop, _rightBottomTop, _uv, _tint, _texture, _rotation)
}

func (_r *Renderer2D) InsertQuadTex(_origin, _dims Vector2f, _uv customtypes.Pair[Vector2f, Vector2f], _z int, _tint RGBA8, _texture *Texture2D) {
	for i := 0; i < _r.m_RenderBatches.Count(); i++ {
		if _z == _r.m_RenderBatches.Data[i].m_Z && _r.m_RenderBatches.Data[i].HasRoom() {
			_r.m_RenderBatches.Data[i].insertQuadTex(_origin, _dims, _uv, _tint, _texture)
			return
		}
	}

	// If it didn't find a RenderBatch with this Z
	_r.m_RenderBatches.PushBack(newRenderer(_r.m_App, _r.m_MaximumQuads, _z, _r.m_Camera))
	_r.m_RenderBatches.Back().insertQuadTex(_origin, _dims, _uv, _tint, _texture)
}

func (_r *Renderer2D) InsertQuadTexRotated(_origin, _dims Vector2f, _uv customtypes.Pair[Vector2f, Vector2f], _z int, _tint RGBA8, _texture *Texture2D, _rotation float32) {
	for i := 0; i < _r.m_RenderBatches.Count(); i++ {
		if _z == _r.m_RenderBatches.Data[i].m_Z && _r.m_RenderBatches.Data[i].HasRoom() {
			_r.m_RenderBatches.Data[i].insertQuadTexRotated(_origin, _dims, _uv, _tint, _texture, _rotation)
			return
		}
	}

	// If it didn't find a RenderBatch with this Z
	_r.m_RenderBatches.PushBack(newRenderer(_r.m_App, _r.m_MaximumQuads, _z, _r.m_Camera))
	_r.m_RenderBatches.Back().insertQuadTexRotated(_origin, _dims, _uv, _tint, _texture, _rotation)
}

////////////////////////

func (_r *Renderer2D) SortBatches() {
	sort.Slice(_r.m_RenderBatches.Data, func(i, j int) bool {
		return _r.m_RenderBatches.Data[i].m_Z < _r.m_RenderBatches.Data[j].m_Z
	})
}

func (_r *Renderer2D) Begin() {
	_r.SortBatches()
	for _rb := range _r.m_RenderBatches.Data {
		_r.m_RenderBatches.Data[_rb].begin()
	}
}

func (_r *Renderer2D) End() {
	// for _rb := range _r.m_RenderBatches.Data {
	// 	_r.m_RenderBatches.Data[_rb].end()
	// }
}

func (_r *Renderer2D) Render() {
	for _rb := range _r.m_RenderBatches.Data {
		_r.m_RenderBatches.Data[_rb].end()
		_r.m_RenderBatches.Data[_rb].render()
	}

	// _r.m_RenderBatches.Clear()
}

/// Render Batch ////////////////////

const (
	MAX_TEXTURES = 6
	NVERTEX_SIZE = 24
)

type NVertex struct {
	m_Coordinates, m_TexCoordinates Vector2f
	m_Color                         RGBA8
	m_textureIndex                  float32
}

type renderBatch2D struct {
	m_Z                 int
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
	m_Camera *Camera2D
	m_App    *App
}

func (_renderer *renderBatch2D) HasRoom() bool {
	return _renderer.m_NumOfVertices/4 < _renderer.m_MaximumQuads
}

func (_renderer *renderBatch2D) insertQuadWithCoordinates(_leftBottomTop, _rightBottomTop customtypes.Pair[Vector2f, Vector2f], _tint RGBA8) {
	// if _renderer.m_NumOfQuads >= _renderer.m_MaximumQuads {
	// 	_renderer.end()
	// 	_renderer.render()
	// 	_renderer.begin()
	// }
	_vertex_index := _renderer.m_NumOfVertices
	_renderer.m_Vertices[_vertex_index+0].m_Coordinates = _leftBottomTop.First
	_renderer.m_Vertices[_vertex_index+0].m_TexCoordinates = Vector2fZero
	_renderer.m_Vertices[_vertex_index+0].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+0].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_Vertices[_vertex_index+1].m_Coordinates = _leftBottomTop.Second
	_renderer.m_Vertices[_vertex_index+1].m_TexCoordinates = Vector2fUp
	_renderer.m_Vertices[_vertex_index+1].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+1].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_Vertices[_vertex_index+2].m_Coordinates = _rightBottomTop.Second
	_renderer.m_Vertices[_vertex_index+2].m_TexCoordinates = Vector2fOne
	_renderer.m_Vertices[_vertex_index+2].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+2].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_Vertices[_vertex_index+3].m_Coordinates = _rightBottomTop.First
	_renderer.m_Vertices[_vertex_index+3].m_TexCoordinates = Vector2fRight
	_renderer.m_Vertices[_vertex_index+3].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+3].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_NumOfElements += 6
	_renderer.m_NumOfVertices += 4
	_renderer.m_NumOfQuads += 1
}

func (_renderer *renderBatch2D) insertQuadWithCoordinatesRotated(_leftBottomTop, _rightBottomTop customtypes.Pair[Vector2f, Vector2f], _tint RGBA8, _rotation float32) {
	// if _renderer.m_NumOfQuads >= _renderer.m_MaximumQuads {
	// 	_renderer.end()
	// 	_renderer.render()
	// 	_renderer.begin()
	// }
	midPoint := Vector2fMidpoint(_leftBottomTop.First, _rightBottomTop.Second)

	_vertex_index := _renderer.m_NumOfVertices
	_renderer.m_Vertices[_vertex_index+0].m_Coordinates = _leftBottomTop.First.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+0].m_TexCoordinates = Vector2fZero
	_renderer.m_Vertices[_vertex_index+0].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+0].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_Vertices[_vertex_index+1].m_Coordinates = _leftBottomTop.Second.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+1].m_TexCoordinates = Vector2fUp
	_renderer.m_Vertices[_vertex_index+1].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+1].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_Vertices[_vertex_index+2].m_Coordinates = _rightBottomTop.Second.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+2].m_TexCoordinates = Vector2fOne
	_renderer.m_Vertices[_vertex_index+2].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+2].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_Vertices[_vertex_index+3].m_Coordinates = _rightBottomTop.First.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+3].m_TexCoordinates = Vector2fRight
	_renderer.m_Vertices[_vertex_index+3].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+3].m_textureIndex = float32(_renderer.m_WhiteTextureSlotIndex)

	_renderer.m_NumOfElements += 6
	_renderer.m_NumOfVertices += 4
	_renderer.m_NumOfQuads += 1
}

func (_renderer *renderBatch2D) insertQuad(_origin, _dims Vector2f, _tint RGBA8) {
	_renderer.insertQuadWithCoordinates(
		customtypes.Pair[Vector2f, Vector2f]{First: _origin, Second: _origin.AddXY(0.0, _dims.Y)},
		customtypes.Pair[Vector2f, Vector2f]{First: _origin.AddXY(_dims.X, 0.0), Second: _origin.Add(_dims)}, _tint)

}

func (_renderer *renderBatch2D) insertQuadRotated(_origin, _dims Vector2f, _tint RGBA8, _rotation float32) {
	_renderer.insertQuadWithCoordinatesRotated(
		customtypes.Pair[Vector2f, Vector2f]{First: _origin, Second: _origin.AddXY(0.0, _dims.Y)},
		customtypes.Pair[Vector2f, Vector2f]{First: _origin.AddXY(_dims.X, 0.0), Second: _origin.Add(_dims)},
		_tint, _rotation)
}

func (_renderer *renderBatch2D) insertQuadTexWithCoordinates(_leftBottomTop, _rightBottomTop, _uv customtypes.Pair[Vector2f, Vector2f], _tint RGBA8, _texture *Texture2D) {
	// if _renderer.m_NumOfQuads >= _renderer.m_MaximumQuads || _renderer.m_TextureSlotIndex > MAX_TEXTURES-1 {
	// 	_renderer.end()
	// 	_renderer.render()
	// 	_renderer.begin()
	// }
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

	_vertex_index := _renderer.m_NumOfVertices

	_renderer.m_Vertices[_vertex_index+0].m_Coordinates = _leftBottomTop.First
	_renderer.m_Vertices[_vertex_index+0].m_TexCoordinates = NewVector2f(_uv.First.X, _uv.Second.Y)
	_renderer.m_Vertices[_vertex_index+0].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+0].m_textureIndex = _texture_index

	_renderer.m_Vertices[_vertex_index+1].m_Coordinates = _leftBottomTop.Second
	_renderer.m_Vertices[_vertex_index+1].m_TexCoordinates = _uv.First
	_renderer.m_Vertices[_vertex_index+1].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+1].m_textureIndex = _texture_index

	_renderer.m_Vertices[_vertex_index+2].m_Coordinates = _rightBottomTop.Second
	_renderer.m_Vertices[_vertex_index+2].m_TexCoordinates = NewVector2f(_uv.Second.X, _uv.First.Y)
	_renderer.m_Vertices[_vertex_index+2].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+2].m_textureIndex = _texture_index

	_renderer.m_Vertices[_vertex_index+3].m_Coordinates = _rightBottomTop.First
	_renderer.m_Vertices[_vertex_index+3].m_TexCoordinates = _uv.Second
	_renderer.m_Vertices[_vertex_index+3].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+3].m_textureIndex = _texture_index

	_renderer.m_NumOfElements += 6
	_renderer.m_NumOfVertices += 4
	_renderer.m_NumOfQuads += 1
}

func (_renderer *renderBatch2D) insertQuadTexWithCoordinatesRotated(_leftBottomTop, _rightBottomTop, _uv customtypes.Pair[Vector2f, Vector2f], _tint RGBA8, _texture *Texture2D, _rotation float32) {
	// if _renderer.m_NumOfQuads >= _renderer.m_MaximumQuads || _renderer.m_TextureSlotIndex > MAX_TEXTURES-1 {
	// 	_renderer.end()
	// 	_renderer.render()
	// 	_renderer.begin()
	// }
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

	midPoint := Vector2fMidpoint(_leftBottomTop.First, _rightBottomTop.Second)
	_vertex_index := _renderer.m_NumOfVertices

	_renderer.m_Vertices[_vertex_index+0].m_Coordinates = _leftBottomTop.First.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+0].m_TexCoordinates = NewVector2f(_uv.First.X, _uv.Second.Y)
	_renderer.m_Vertices[_vertex_index+0].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+0].m_textureIndex = _texture_index

	_renderer.m_Vertices[_vertex_index+1].m_Coordinates = _leftBottomTop.Second.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+1].m_TexCoordinates = _uv.First
	_renderer.m_Vertices[_vertex_index+1].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+1].m_textureIndex = _texture_index

	_renderer.m_Vertices[_vertex_index+2].m_Coordinates = _rightBottomTop.Second.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+2].m_TexCoordinates = NewVector2f(_uv.Second.X, _uv.First.Y)
	_renderer.m_Vertices[_vertex_index+2].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+2].m_textureIndex = _texture_index

	_renderer.m_Vertices[_vertex_index+3].m_Coordinates = _rightBottomTop.First.Rotate(_rotation, midPoint)
	_renderer.m_Vertices[_vertex_index+3].m_TexCoordinates = _uv.Second
	_renderer.m_Vertices[_vertex_index+3].m_Color = _tint
	_renderer.m_Vertices[_vertex_index+3].m_textureIndex = _texture_index

	_renderer.m_NumOfElements += 6
	_renderer.m_NumOfVertices += 4
	_renderer.m_NumOfQuads += 1
}

func (_renderer *renderBatch2D) insertQuadTex(_origin, _dims Vector2f, _uv customtypes.Pair[Vector2f, Vector2f], _tint RGBA8, _texture *Texture2D) {
	_renderer.insertQuadTexWithCoordinates(
		customtypes.Pair[Vector2f, Vector2f]{First: _origin, Second: _origin.AddXY(0.0, _dims.Y)},
		customtypes.Pair[Vector2f, Vector2f]{First: _origin.AddXY(_dims.X, 0.0), Second: _origin.Add(_dims)},
		_uv, _tint, _texture)
}
func (_renderer *renderBatch2D) insertQuadTexRotated(_origin, _dims Vector2f, _uv customtypes.Pair[Vector2f, Vector2f], _tint RGBA8, _texture *Texture2D, _rotation float32) {
	_renderer.insertQuadTexWithCoordinatesRotated(
		customtypes.Pair[Vector2f, Vector2f]{First: _origin, Second: _origin.AddXY(0.0, _dims.Y)},
		customtypes.Pair[Vector2f, Vector2f]{First: _origin.AddXY(_dims.X, 0.0), Second: _origin.Add(_dims)},
		_uv, _tint, _texture, _rotation)
}

func (_renderer *renderBatch2D) begin() {
	_renderer.m_NumOfVertices = 0
	_renderer.m_NumOfElements = 0
	_renderer.m_NumOfQuads = 0

	_renderer.m_TextureSlotIndex = 1
	for i := 1; i < MAX_TEXTURES; i++ {
		_renderer.m_TextureSlots[i] = js.ValueOf(0)
	}
}

func (_renderer *renderBatch2D) end() {
	canvasContext.Call("bindVertexArray", _renderer.m_Vao)

	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _renderer.m_Vbo)
	// _jsNVerts := nVertexBufferToJsNVertexBuffer(_renderer.m_Vertices)
	setNVertexBufferAndJsBuffer(&_renderer.m_JsVertsBuffer, _renderer.m_BytesVertsBuffer, _renderer.m_Vertices)

	canvasContext.Call("bufferSubData", canvasContext.Get("ARRAY_BUFFER"), 0, _renderer.m_JsVertsBuffer)

	// canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), _jsNVerts, canvasContext.Get("DYNAMIC_DRAW"))
	// _jsNVerts.Set("delete", nil)
}

func (_renderer *renderBatch2D) render() {
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
		canvasContext.Call("activeTexture", 33984+i)
		canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), _renderer.m_TextureSlots[i])
	}

	// canvasContext.Call("bindVertexArray", _renderer.m_Vao)
	canvasContext.Call("drawElements", canvasContext.Get("TRIANGLES"), _renderer.m_NumOfElements, canvasContext.Get("UNSIGNED_INT"), 0)
	canvasContext.Call("bindVertexArray", js.Null())
}

func newRenderer(_app *App, _maxQuads, _z int, _cam *Camera2D) renderBatch2D {
	_renderer := renderBatch2D{}

	_renderer.m_Z = _z

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
	canvasContext.Call("vertexAttribPointer", 2, 4, canvasContext.Get("UNSIGNED_BYTE"), true, NVERTEX_SIZE, 16)
	canvasContext.Call("enableVertexAttribArray", 3)
	canvasContext.Call("vertexAttribPointer", 3, 1, canvasContext.Get("FLOAT"), false, NVERTEX_SIZE, 20)

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

	// _renderer.m_Camera.Init(*_app)
	// _renderer.m_Camera.m_CenterOffset = NewVector2f(float32(_app.Width)/2.0, float32(_app.Height)/2.0)
	// _renderer.m_Camera.Update(*_app)
	_renderer.m_Camera = _cam

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
