package chai

import (
	"image"
	"image/png"
	"io"
	"syscall/js"

	. "github.com/mhamedGd/chai/math"
)

type TextureFilter = js.Value

var TEXTURE_FILTER_LINEAR TextureFilter
var TEXTURE_FILTER_NEAREST TextureFilter

func initTextures() {
	TEXTURE_FILTER_LINEAR = canvasContext.Get("LINEAR")
	TEXTURE_FILTER_NEAREST = canvasContext.Get("NEAREST")
}

type TileSet struct {
	m_Texture                     Texture2D
	m_TotalRows, m_TotalColumns   int
	m_SpriteWidth, m_SpriteHeight float32
}

func NewTileSet(_texture Texture2D, _columns, _rows int) TileSet {
	if _texture.m_Spwidth == 0 {
		return TileSet{
			m_Texture:      _texture,
			m_TotalRows:    _rows,
			m_TotalColumns: _columns,
			m_SpriteWidth:  float32(_texture.Width) / float32(_columns),
			m_SpriteHeight: float32(_texture.Height) / float32(_rows),
		}
	}

	return TileSet{
		m_Texture:      _texture,
		m_TotalRows:    _rows,
		m_TotalColumns: _columns,
		m_SpriteWidth:  float32(_texture.m_Spwidth),
		m_SpriteHeight: float32(_texture.m_Spheight),
	}
}

func NewTileSetBySize(_texture Texture2D, _spriteWidth, _spriteHeight int) TileSet {
	if _texture.m_Spwidth == 0 {
		return TileSet{
			m_Texture:      _texture,
			m_TotalRows:    _texture.Height / _spriteHeight,
			m_TotalColumns: _texture.Width / _spriteWidth,
			m_SpriteWidth:  float32(_spriteWidth),
			m_SpriteHeight: float32(_spriteHeight),
		}
	}

	return TileSet{
		m_Texture:      _texture,
		m_TotalRows:    _texture.Height / _spriteHeight,
		m_TotalColumns: _texture.Width / _spriteWidth,
		m_SpriteWidth:  float32(_texture.m_Spwidth),
		m_SpriteHeight: float32(_texture.m_Spheight),
	}
}

type TextureSettings struct {
	Filter TextureFilter
}

type Texture2D struct {
	Width, Height             int
	m_TextureId               js.Value
	m_Spwidth, m_Spheight     int
	m_PixelsToMeterDimensions Vector2f
	m_LocalPixelsPerMeter     int
}

func (t *Texture2D) OverridePixelsPerMeter(_newPPM int) {
	t.m_LocalPixelsPerMeter = _newPPM
	t.resizePixelsToMeters()
}

func (t *Texture2D) resizePixelsToMeters() {
	if t.m_LocalPixelsPerMeter == 0 {
		t.m_PixelsToMeterDimensions = NewVector2f(float32(t.Width)/float32(pixelsPerMeter), float32(t.Height)/float32(pixelsPerMeter))
	} else {
		t.m_PixelsToMeterDimensions = NewVector2f(float32(t.Width)/float32(t.m_LocalPixelsPerMeter), float32(t.Height)/float32(t.m_LocalPixelsPerMeter))
	}
}

type Pixel struct {
	RGBA RGBA8
}

func NewPixel(_r, _g, _b, _a uint8) Pixel {
	return Pixel{RGBA: RGBA8{_r, _g, _b, _a}}
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

	tempTexture.m_TextureId = canvasContext.Call("createTexture")
	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), tempTexture.m_TextureId)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), canvasContext.Get("NEAREST"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), canvasContext.Get("NEAREST"))

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)

	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), js.Null())

	tempTexture.m_PixelsToMeterDimensions = NewVector2f(float32(pixelsPerMeter), float32(pixelsPerMeter))
	return tempTexture
}

func LoadPng(_filePath string, _textureSettings TextureSettings) Texture2D {

	var tempTexture Texture2D

	pngChannel := make(chan io.ReadCloser)
	go LoadResponseBody(_filePath, pngChannel)
	pngChannelBody := <-pngChannel
	_img, err := png.Decode(pngChannelBody)
	if err != nil {
		LogF("%v", err.Error())
	}
	defer pngChannelBody.Close()

	tempTexture.Width = _img.Bounds().Dx()
	tempTexture.Height = _img.Bounds().Dy()

	pixels := make([]Pixel, tempTexture.Height*tempTexture.Width)

	for y := 0; y < tempTexture.Height; y++ {
		for x := 0; x < tempTexture.Width; x++ {
			r, g, b, a := _img.At(x, y).RGBA()
			pixels[y*tempTexture.Width+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
		}
	}

	tempTexture.m_TextureId = canvasContext.Call("createTexture")
	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), tempTexture.m_TextureId)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), _textureSettings.Filter)
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), _textureSettings.Filter)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)

	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), js.Null())

	// tempTexture.m_PixelsToMeterDimensions = math.NewVector2f(float32(tempTexture.Width)/float32(pixelsPerMeter), float32(tempTexture.Height)/float32(pixelsPerMeter))
	tempTexture.resizePixelsToMeters()
	return tempTexture
}

func LoadPngByTileset(_filePath string, _textureSettings TextureSettings, _spriteWidth, _spriteHeight int) Texture2D {

	var tempTexture Texture2D
	tempTexture.m_Spwidth = _spriteWidth
	tempTexture.m_Spheight = _spriteHeight

	pngChannel := make(chan io.ReadCloser)
	go LoadResponseBody(_filePath, pngChannel)
	pngChannelBody := <-pngChannel
	_img, err := png.Decode(pngChannelBody)
	if err != nil {
		LogF("%v", err.Error())
	}
	defer pngChannelBody.Close()

	tempTexture.Width = _img.Bounds().Dx()
	tempTexture.Height = _img.Bounds().Dy()

	numofCols := tempTexture.Width / _spriteWidth
	numofRows := tempTexture.Height / _spriteHeight

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
			if passed_pixelsx >= _spriteWidth && passed_pixelsy < _spriteHeight {
				if passed_pixelsx == _spriteWidth+1 {
					//Left Side of the tile
					////////////////////////////////////
					r, g, b, a := _img.At(x-passed_spsx*2, y-passed_spsy*2).RGBA()
					pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))

					passed_pixelsx = 0
					passed_spsx++
				} else {
					//Right Side of the tile
					////////////////////////////////////
					passed_pixelsx++
					r, g, b, a := _img.At(x-passed_spsx*2-1, y-passed_spsy*2).RGBA()
					pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
				}
			} else if passed_pixelsy >= _spriteHeight || passed_pixelsy >= _spriteHeight+1 {
				if passed_pixelsy >= _spriteHeight+1 {
					// Top Side of the tile
					////////////////////////////////////
					// pixels[y*(tempTexture.Width)+x] = border_pixel

					// Check if corner
					if passed_pixelsx == _spriteWidth {
						passed_pixelsx++
						pixels[y*(tempTexture.Width)+x] = debug_pixel
						passed_spsx++

					} else if passed_pixelsx == _spriteWidth+1 {
						passed_pixelsx = 0
						pixels[y*(tempTexture.Width)+x] = debug_pixel
					} else {
						r, g, b, a := _img.At(x-passed_spsx*2, y-passed_spsy*2).RGBA()
						pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
						passed_pixelsx++
					}

				} else {
					// Bottom Side of the tile
					////////////////////////////////////

					// Check if Corner

					// pixels[y*(tempTexture.Width)+x] = debug_pixel
					if passed_pixelsx == _spriteWidth {
						passed_pixelsx++
						pixels[y*(tempTexture.Width)+x] = debug_pixel
						passed_spsx++

					} else if passed_pixelsx == _spriteWidth+1 {
						passed_pixelsx = 0
						pixels[y*(tempTexture.Width)+x] = debug_pixel
					} else {
						r, g, b, a := _img.At(x-passed_spsx*2, y-passed_spsy*2-1).RGBA()
						pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
						passed_pixelsx++
					}
				}

			} else {

				r, g, b, a := _img.At(x-passed_spsx*2, y-passed_spsy*2).RGBA()
				pixels[y*(tempTexture.Width)+x] = NewPixel(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a))
				passed_pixelsx++
			}
		}

		if passed_pixelsy >= _spriteHeight+1 {
			passed_pixelsy = 0
			passed_spsy++
		} else {
			passed_pixelsy++
		}
		passed_spsx = 0
		passed_pixelsx = 0
	}

	tempTexture.m_TextureId = canvasContext.Call("createTexture")
	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), tempTexture.m_TextureId)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), _textureSettings.Filter)
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), _textureSettings.Filter)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)

	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), js.Null())

	// tempTexture.m_PixelsToMeterDimensions = math.NewVector2f(float32(tempTexture.Width)/float32(pixelsPerMeter), float32(tempTexture.Height)/float32(pixelsPerMeter))
	// tempTexture.resizePixelsToMeters()
	if tempTexture.m_LocalPixelsPerMeter == 0 {
		tempTexture.m_PixelsToMeterDimensions = NewVector2f(float32(tempTexture.m_Spwidth)/float32(pixelsPerMeter), float32(tempTexture.m_Spheight)/float32(pixelsPerMeter))
	} else {
		tempTexture.m_PixelsToMeterDimensions = NewVector2f(float32(tempTexture.m_Spwidth)/float32(tempTexture.m_LocalPixelsPerMeter), float32(tempTexture.m_Spheight)/float32(tempTexture.m_LocalPixelsPerMeter))
	}
	return tempTexture
}

func LoadTextureFromImg(_img image.Image) Texture2D {
	var tempTexture Texture2D

	tempTexture.Width = _img.Bounds().Dx()
	tempTexture.Height = _img.Bounds().Dy()

	if tempTexture.Height <= 0 || tempTexture.Width <= 0 {
		LogF("Loaded Image has zero dimensions")
	}

	pixels := make([]Pixel, tempTexture.Height*tempTexture.Width)

	for y := 0; y < tempTexture.Height; y++ {
		for x := 0; x < tempTexture.Width; x++ {
			r, g, b, a := _img.At(x, y).RGBA()
			pixels[y*tempTexture.Width+x] = NewPixel(uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	tempTexture.m_TextureId = canvasContext.Call("createTexture")
	canvasContext.Call("activeTexture", canvasContext.Get("TEXTURE0"))
	canvasContext.Call("bindTexture", canvasContext.Get("TEXTURE_2D"), tempTexture.m_TextureId)

	canvasContext.Call("pixelStorei", canvasContext.Get("UNPACK_ALIGNMENT"), 1)

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), canvasContext.Get("NEAREST"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), canvasContext.Get("NEAREST"))

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)
	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA8"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

	tempTexture.resizePixelsToMeters()
	return tempTexture
}
