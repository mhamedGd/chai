package chai

import (
	"image"
	"image/png"
	"io"
	"syscall/js"
)

type TextureFilter = js.Value

var TEXTURE_FILTER_LINEAR TextureFilter
var TEXTURE_FILTER_NEAREST TextureFilter

func initTextures() {
	TEXTURE_FILTER_LINEAR = canvasContext.Get("LINEAR")
	TEXTURE_FILTER_NEAREST = canvasContext.Get("NEAREST")
}

type TileSet struct {
	texture                   Texture2D
	totalRows, totalColumns   int
	spriteWidth, spriteHeight int
	startPosition             Vector2f
}

func NewTileSet(_startPosition Vector2f, _texture Texture2D, _columns, _rows int) TileSet {
	return TileSet{
		texture:       _texture,
		totalRows:     _rows,
		totalColumns:  _columns,
		spriteWidth:   _texture.Width / _columns,
		spriteHeight:  _texture.Height / _rows,
		startPosition: _startPosition,
	}
}

type TextureSettings struct {
	Filter TextureFilter
}

type Texture2D struct {
	Width, Height, bpp int
	textureId          js.Value
}

type Pixel struct {
	RGBA RGBA8
}

func New(r, g, b, a uint8) Pixel {
	var pixel Pixel
	pixel.RGBA = RGBA8{r, g, b, a}
	return pixel
}

func LoadPng(_filePath string, _textureSettings *TextureSettings) Texture2D {

	var tempTexture Texture2D

	pngChannel := make(chan io.ReadCloser)
	go LoadResponseBody(_filePath, pngChannel)
	pngChannelBody := <-pngChannel
	img, err := png.Decode(pngChannelBody)
	if err != nil {
		LogF("%v", err.Error())
	}
	defer pngChannelBody.Close()

	tempTexture.Width = img.Bounds().Dx()
	tempTexture.Height = img.Bounds().Dy()

	pixels := make([]Pixel, tempTexture.Height*tempTexture.Width)

	for y := 0; y < tempTexture.Height; y++ {
		for x := 0; x < tempTexture.Width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[y*tempTexture.Width+x] = New(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
		}
	}

	tempTexture.textureId = canvasContext.Call("createTexture")
	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), tempTexture.textureId)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), _textureSettings.Filter)
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), _textureSettings.Filter)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)

	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

	return tempTexture
}

func LoadTextureFromImg(img image.Image) Texture2D {
	var tempTexture Texture2D

	tempTexture.Width = img.Bounds().Dx()
	tempTexture.Height = img.Bounds().Dy()

	if tempTexture.Height <= 0 || tempTexture.Width <= 0 {
		LogF("Loaded Image has zero dimensions")
	}

	pixels := make([]Pixel, tempTexture.Height*tempTexture.Width)

	for y := 0; y < tempTexture.Height; y++ {
		for x := 0; x < tempTexture.Width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[y*tempTexture.Width+x] = New(uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	tempTexture.textureId = canvasContext.Call("createTexture")
	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), tempTexture.textureId)

	canvasContext.Call("pixelStorei", canvasContext.Get("UNPACK_ALIGNMENT"), 1)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), canvasContext.Get("NEAREST"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), canvasContext.Get("NEAREST"))

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)
	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA8"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

	return tempTexture
}
