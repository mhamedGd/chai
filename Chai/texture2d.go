package chai

import (
	"image"
	"image/png"
	"net/http"
	"syscall/js"
)

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

func LoadPng(_filePath string) Texture2D {

	var tempTexture Texture2D

	resp, err := http.Get(app_url + "/" + _filePath)
	if err != nil {
		LogF("%v", err.Error())
	}
	img, err := png.Decode(resp.Body)
	if err != nil {
		LogF("%v", err.Error())
	}

	resp.Body.Close()

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

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MIN_FILTER"), canvasContext.Get("NEAREST"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_MAG_FILTER"), canvasContext.Get("NEAREST"))

	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_S"), canvasContext.Get("CLAMP_TO_EDGE"))
	canvasContext.Call("texParameteri", canvasContext.Get("TEXTURE_2D"), canvasContext.Get("TEXTURE_WRAP_T"), canvasContext.Get("CLAMP_TO_EDGE"))

	jsPixels := pixelBufferToJsPixelBubffer(pixels)

	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA8"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)

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
