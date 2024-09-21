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
	spriteWidth, spriteHeight float32
}

func NewTileSet(_texture Texture2D, _columns, _rows int) TileSet {
	if _texture.spwidth == 0 {
		return TileSet{
			texture:      _texture,
			totalRows:    _rows,
			totalColumns: _columns,
			spriteWidth:  float32(_texture.Width) / float32(_columns),
			spriteHeight: float32(_texture.Height) / float32(_rows),
		}
	}

	return TileSet{
		texture:      _texture,
		totalRows:    _rows,
		totalColumns: _columns,
		spriteWidth:  float32(_texture.spwidth),
		spriteHeight: float32(_texture.spheight),
	}
}

func NewTileSetBySize(_texture Texture2D, _spriteWidth, _spriteHeight int) TileSet {
	if _texture.spwidth == 0 {
		return TileSet{
			texture:      _texture,
			totalRows:    _texture.Height / _spriteHeight,
			totalColumns: _texture.Width / _spriteWidth,
			spriteWidth:  float32(_spriteWidth),
			spriteHeight: float32(_spriteHeight),
		}
	}

	return TileSet{
		texture:      _texture,
		totalRows:    _texture.Height / _spriteHeight,
		totalColumns: _texture.Width / _spriteWidth,
		spriteWidth:  float32(_texture.spwidth),
		spriteHeight: float32(_texture.spheight),
	}
}

type TextureSettings struct {
	Filter TextureFilter
}

type Texture2D struct {
	Width, Height, bpp int
	textureId          js.Value
	spwidth, spheight  int
}

type Pixel struct {
	RGBA RGBA8
}

func NewPixel(r, g, b, a uint8) Pixel {
	var pixel Pixel
	pixel.RGBA = RGBA8{r, g, b, a}
	return pixel
}

func loadWhiteTexture() Texture2D {

	var tempTexture Texture2D

	tempTexture.Width = 1
	tempTexture.Height = 1

	pixels := make([]Pixel, tempTexture.Height*tempTexture.Width)

	for y := 0; y < tempTexture.Height; y++ {
		for x := 0; x < tempTexture.Width; x++ {
			pixels[y*tempTexture.Width+x] = NewPixel(255, 255, 255, 255)
		}
	}

	tempTexture.textureId = canvasContext.Call("createTexture")
	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), tempTexture.textureId)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), canvasContext.Get("NEAREST"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), canvasContext.Get("NEAREST"))

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)

	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), js.Null())

	return tempTexture
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
			pixels[y*tempTexture.Width+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
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

	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), js.Null())

	return tempTexture
}

func LoadPngByTileset(_filePath string, _textureSettings *TextureSettings, _spritewidth, _spriteheight int) Texture2D {

	var tempTexture Texture2D
	tempTexture.spwidth = _spritewidth
	tempTexture.spheight = _spriteheight

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

	numofCols := tempTexture.Width / _spritewidth
	numofRows := tempTexture.Height / _spriteheight

	tempTexture.Width += numofCols * 2
	tempTexture.Height += numofRows * 2

	pixels := make([]Pixel, tempTexture.Height*tempTexture.Width)
	passed_pixelsx := int(0)
	passed_pixelsy := int(0)
	passed_spsx := int(0)
	passed_spsy := int(0)

	// border_pixel := NewPixel(100, 0, 0, 255)
	debug_pixel := NewPixel(0, 0, 0, 255)

	for y := 0; y < tempTexture.Height; y++ {
		for x := 0; x < tempTexture.Width; x++ {
			if passed_pixelsx >= _spritewidth && passed_pixelsy < _spriteheight {
				if passed_pixelsx == _spritewidth+1 {
					//Left Side of the tile
					////////////////////////////////////
					r, g, b, a := img.At(x-passed_spsx*2, y-passed_spsy*2).RGBA()
					pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))

					passed_pixelsx = 0
					passed_spsx++
				} else {
					//Right Side of the tile
					////////////////////////////////////
					passed_pixelsx++
					r, g, b, a := img.At(x-passed_spsx*2-1, y-passed_spsy*2).RGBA()
					pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
				}
			} else if passed_pixelsy >= _spriteheight || passed_pixelsy >= _spriteheight+1 {
				if passed_pixelsy >= _spriteheight+1 {
					// Top Side of the tile
					////////////////////////////////////
					// pixels[y*(tempTexture.Width)+x] = border_pixel

					// Check if corner
					if passed_pixelsx == _spritewidth {
						passed_pixelsx++
						pixels[y*(tempTexture.Width)+x] = debug_pixel
						passed_spsx++

					} else if passed_pixelsx == _spritewidth+1 {
						passed_pixelsx = 0
						pixels[y*(tempTexture.Width)+x] = debug_pixel
					} else {
						r, g, b, a := img.At(x-passed_spsx*2, y-passed_spsy*2).RGBA()
						pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
						passed_pixelsx++
					}

				} else {
					// Bottom Side of the tile
					////////////////////////////////////

					// Check if Corner

					// pixels[y*(tempTexture.Width)+x] = debug_pixel
					if passed_pixelsx == _spritewidth {
						passed_pixelsx++
						pixels[y*(tempTexture.Width)+x] = debug_pixel
						passed_spsx++

					} else if passed_pixelsx == _spritewidth+1 {
						passed_pixelsx = 0
						pixels[y*(tempTexture.Width)+x] = debug_pixel
					} else {
						r, g, b, a := img.At(x-passed_spsx*2, y-passed_spsy*2-1).RGBA()
						pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
						passed_pixelsx++
					}
				}

			} else {

				r, g, b, a := img.At(x-passed_spsx*2, y-passed_spsy*2).RGBA()
				pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
				passed_pixelsx++
			}
		}

		if passed_pixelsy >= _spriteheight+1 {
			passed_pixelsy = 0
			passed_spsy++
		} else {
			passed_pixelsy++
		}
		passed_spsx = 0
		passed_pixelsx = 0
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

	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), js.Null())

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
			pixels[y*tempTexture.Width+x] = NewPixel(uint8(r), uint8(g), uint8(b), uint8(a))
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
