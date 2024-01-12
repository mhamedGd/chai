package chai

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"net/http"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"golang.org/x/image/math/fixed"
)

var ARABIC_UNICODE = []rune{
	// Letter (ﺃ)
	'\uFE83', '\u0623', '\uFE84', '\uFE84',
	// Letter (ﺍ)
	'\uFE8D', '\u0627', '\uFE8E', '\uFE8E',
	// Letter (ﺁ)
	'\uFE81', '\u0622', '\uFE82', '\uFE82',
	// Letter (ﺀ)
	'\uFE80', '\u0621', '\u0621', '\u0621',
	// Letter (ﺅ)
	'\uFE85', '\u0624', '\uFE86', '\uFE86',
	// Letter (ﺇ)
	'\uFE87', '\u0625', '\uFE88', '\uFE88',
	// Letter (ﺉ)
	'\uFE89', '\uFE8B', '\uFE8C', '\uFE8A',
	// Letter (ﺏ)
	'\uFE8F', '\uFE91', '\uFE92', '\uFE90',
	// Letter (ﺕ)
	'\uFE95', '\uFE97', '\uFE98', '\uFE96',
	// Letter (ﺓ)
	'\uFE93', '\u0629', '\u0629', '\uFE94',
	// Letter (ﺙ)
	'\uFE99', '\uFE9B', '\uFE9C', '\uFE9A',
	// Letter (ﺝ)
	'\uFE9D', '\uFE9F', '\uFEA0', '\uFE9E',
	// Letter (ﺡ)
	'\uFEA1', '\uFEA3', '\uFEA4', '\uFEA2',
	// Letter (ﺥ)
	'\uFEA5', '\uFEA7', '\uFEA8', '\uFEA6',
	// Letter (ﺩ)
	'\uFEA9', '\u062F', '\uFEAA', '\uFEAA',
	// Letter (ﺫ)
	'\uFEAB', '\u0630', '\uFEAC', '\uFEAC',
	// Letter (ﺭ)
	'\uFEAD', '\u0631', '\uFEAE', '\uFEAE',
	// Letter (ﺯ)
	'\uFEAF', '\u0632', '\uFEB0', '\uFEB0',
	// Letter (ﺱ)
	'\uFEB1', '\uFEB3', '\uFEB4', '\uFEB2',
	// Letter (ﺵ)
	'\uFEB5', '\uFEB7', '\uFEB8', '\uFEB6',
	// Letter (ﺹ)
	'\uFEB9', '\uFEBB', '\uFEBC', '\uFEBA',
	// Letter (ﺽ)
	'\uFEBD', '\uFEBF', '\uFEC0', '\uFEBE',
	// Letter (ﻁ)
	'\uFEC1', '\uFEC3', '\uFEC4', '\uFEC2',
	// Letter (ﻅ)
	'\uFEC5', '\uFEC7', '\uFEC8', '\uFEC6',
	// Letter (ﻉ)
	'\uFEC9', '\uFECB', '\uFECC', '\uFECA',
	// Letter (ﻍ)
	'\uFECD', '\uFECF', '\uFED0', '\uFECE',
	// Letter (ﻑ)
	'\uFED1', '\uFED3', '\uFED4', '\uFED2',
	// Letter (ﻕ)
	'\uFED5', '\uFED7', '\uFED8', '\uFED6',
	// Letter (ﻙ)
	'\uFED9', '\uFEDB', '\uFEDC', '\uFEDA',
	// Letter (ﻝ)
	'\uFEDD', '\uFEDF', '\uFEE0', '\uFEDE',
	// Letter (ﻡ)
	'\uFEE1', '\uFEE3', '\uFEE4', '\uFEE2',
	// Letter (ﻥ)
	'\uFEE5', '\uFEE7', '\uFEE8', '\uFEE6',
	// Letter (ﻩ)
	'\uFEE9', '\uFEEB', '\uFEEC', '\uFEEA',
	// Letter (ﻭ)
	'\uFEED', '\u0648', '\uFEEE', '\uFEEE',
	// Letter (ﻱ)
	'\uFEF1', '\uFEF3', '\uFEF4', '\uFEF2',
	// Letter (ﻯ)
	'\uFEEF', '\u0649', '\uFEF0', '\uFEF0',
	// Letter (ـ)
	'\u0640', '\u0640', '\u0640', '\u0640',
	// Letter (ﻻ)
	'\uFEFB', '\uFEFB', '\uFEFC', '\uFEFC',
	// Letter (ﻷ)
	'\uFEF7', '\uFEF7', '\uFEF8', '\uFEF8',
	// Letter (ﻹ)
	'\uFEF9', '\uFEF9', '\uFEFA', '\uFEFA',
	// Letter (ﻵ)
	'\uFEF5', '\uFEF5', '\uFEF6', '\uFEF6',
	'٠',
	'١',
	'٢',
	'٣',
	'٤',
	'٥',
	'٦',
	'٧',
	'٨',
	'٩',
	'٪',
}

/*
var ARABIC_UNICODE = []rune{
	'ا',
	'ﺍ',
	'ﺎ',
	'ب',
	'ﺏ',
	'ﺐ',
	'ﺒ',
	'ﺑ',
	'ت',
	'ﺕ',
	'ﺖ',
	'ﺘ',
	'ﺗ',
	'ث',
	'ﺙ',
	'ﺚ',
	'ﺜ',
	'ﺛ',
	'ج',
	'ﺝ',
	'ﺞ',
	'ﺠ',
	'ﺟ',
	'ح',
	'ﺡ',
	'ﺢ',
	'ﺤ',
	'ﺣ',
	'خ',
	'ﺥ',
	'ﺦ',
	'ﺨ',
	'ﺧ',
	'د',
	'ﺩ',
	'ﺪ',
	'ذ',
	'ﺫ',
	'ﺬ',
	'ر',
	'ﺭ',
	'ﺮ',
	'ز',
	'ﺯ',
	'ﺰ',
	'س',
	'ﺱ',
	'ﺲ',
	'ﺴ',
	'ﺳ',
	'ش',
	'ﺵ',
	'ﺶ',
	'ﺸ',
	'ﺷ',
	'ص',
	'ﺹ',
	'ﺺ',
	'ﺼ',
	'ﺻ',
	'ض',
	'ﺽ',
	'ﺾ',
	'ﻀ',
	'ﺿ',
	'ط',
	'ﻁ',
	'ﻂ',
	'ﻄ',
	'ﻃ',
	'ظ',
	'ﻅ',
	'ﻆ',
	'ﻈ',
	'ﻇ',
	'ع',
	'ﻉ',
	'ﻊ',
	'ﻌ',
	'ﻋ',
	'غ',
	'ﻍ',
	'ﻎ',
	'ﻐ',
	'ﻏ',
	'ف',
	'ﻑ',
	'ﻒ',
	'ﻔ',
	'ﻓ',
	'ق',
	'ﻕ',
	'ﻖ',
	'ﻘ',
	'ﻗ',
	'ك',
	'ﻙ',
	'ﻚ',
	'ﻜ',
	'ﻛ',
	'ل',
	'ﻝ',
	'ﻞ',
	'ﻠ',
	'ﻟ',
	'م',
	'ﻡ',
	'ﻢ',
	'ﻤ',
	'ﻣ',
	'ن',
	'ﻥ',
	'ﻦ',
	'ﻨ',
	'ﻧ',
	'ه',
	'ﻩ',
	'ﻪ',
	'ﻬ',
	'ﻫ',
	'و',
	'ﻭ',
	'ﻮ',
	'ي',
	'ﻱ',
	'ﻲ',
	'ﻴ',
	'ﻳ',
	'آ',
	'ﺁ',
	'ﺂ',
	'ة',
	'ﺓ',
	'ﺔ',
	'ى',
	'ﻯ',
	'ﻰ',
	'ء',
	'أ',
	'ﺊ',
	'ﺋ',
	'ﺌ',
	'ؤ',
	'إ',
	'ئ',
	'٠',
	'١',
	'٢',
	'٣',
	'٤',
	'٥',
	'٦',
	'٧',
	'٨',
	'٩',
	'٪',
	'٫',
	'ﻻ',
	'ﻼ',
	'ﻺ',
	'ﻹ',
	'ﻸ',
	'ﻷ',
	'ﻶ',
	'ﻵ',
}
*/

type FontBatchAtlas struct {
	charAtlasSet map[rune]CharAtlasGlyph
	textureAtlas Texture2D
	sPatch       SpriteBatch
	fontSettings FontBatchSettings
}

type CharAtlasGlyph struct {
	uv1           Vector2f
	uv2           Vector2f
	size, bearing Vector2f
	advance       float32
}

type FontBatchSettings struct {
	FontSize, DPI, CharDistance, LineHeight float32
	Arabic                                  bool
}

func (self *FontBatchAtlas) Init() {
	self.charAtlasSet = make(map[rune]CharAtlasGlyph)
	self.sPatch.Init("font.shader")

}

const GLYPH_ATLAS_GAP = 5

func LoadFontToAtlas(_fontPath string, _fontSettings *FontBatchSettings) FontBatchAtlas {
	var tempFont FontBatchAtlas
	tempFont.Init()
	tempFont.fontSettings = *_fontSettings

	resp, err := http.Get(app_url + "/" + _fontPath)
	if err != nil {
		LogF(err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		LogF(err.Error())
	}

	f, err := opentype.Parse(data)

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(_fontSettings.FontSize),
		DPI:     float64(_fontSettings.DPI),
		Hinting: font.HintingFull,
	})
	if err != nil {
		LogF(err.Error())
	}

	max_width, max_height := int(0), int(0)

	dot := fixed.Point26_6{fixed.Int26_6(face.Metrics().Ascent), 0}
	for i := 33; i < 127; i++ {
		char := rune(i)

		_, img, _, _, ok := face.Glyph(dot, char)

		if !ok {
			//LogF("Failed to load rune: %v", char)
		}

		max_width += img.Bounds().Dx() + GLYPH_ATLAS_GAP
		if max_height < img.Bounds().Dy() {
			max_height = img.Bounds().Dy()
		}
	}

	for _, c := range ARABIC_UNICODE {
		_, img, _, _, ok := face.Glyph(dot, c)

		if !ok {
			//LogF("Failed to load rune: %v", c)
		}

		max_width += img.Bounds().Dx() + GLYPH_ATLAS_GAP
		if max_height < img.Bounds().Dy() {
			max_height = img.Bounds().Dy()
		}
	}
	atlas_img := image.NewRGBA(image.Rect(0, 0, max_width, max_height))
	draw.Draw(atlas_img, atlas_img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)

	x_offset := 0
	y_offset := 0

	for i := 32; i < 127; i++ {
		char := rune(i)
		_, img, _, ad, ok := face.Glyph(dot, char)
		if !ok {
			//LogF("Failed To Draw Letter: %v", char)
		}
		y_offset = img.Bounds().Dy()

		if char == ' ' {
			tempFont.charAtlasSet[char] = CharAtlasGlyph{uv1: Vector2fZero, uv2: Vector2fZero, size: Vector2fZero, bearing: Vector2fOne, advance: float32(ad) / float32(1<<6)}
			continue
		}

		bounds, _, ok := face.GlyphBounds(char)
		if !ok {
			//LogF("Failed to load bounds of rune: %v", char)
		}

		tempFont.charAtlasSet[char] = CharAtlasGlyph{
			NewVector2f(float32(x_offset)/float32(max_width), 0),
			NewVector2f(float32(x_offset+img.Bounds().Dx())/float32(max_width), float32(y_offset)/float32(max_height)),
			NewVector2f(float32(img.Bounds().Dx()), float32(img.Bounds().Dy())),
			NewVector2f(float32(bounds.Max.X)/64.0, float32(-bounds.Max.Y)/64.0),
			float32(ad) / 64.0,
		}
		draw.Draw(atlas_img, image.Rect(x_offset, 0, x_offset+img.Bounds().Dx(), y_offset), img, image.ZP, draw.Src)
		x_offset += img.Bounds().Dx() + GLYPH_ATLAS_GAP
	}
	if _fontSettings.Arabic {

		for _, c := range ARABIC_UNICODE {
			_, img, _, ad, ok := face.Glyph(dot, c)
			if !ok {
				//LogF("Failed To Draw Letter: %v", c)
				if c == '\uFEAD' {
					LogF("DETECTED ARR")
				}
				continue
			}
			y_offset = img.Bounds().Dy()

			bounds, _, ok := face.GlyphBounds(c)
			if !ok {
				//LogF("Failed to load bounds of rune: %v", c)
				continue
			}

			tempFont.charAtlasSet[c] = CharAtlasGlyph{
				NewVector2f(float32(x_offset)/float32(max_width), 0),
				NewVector2f(float32(x_offset+img.Bounds().Dx())/float32(max_width), float32(y_offset)/float32(max_height)),
				NewVector2f(float32(img.Bounds().Dx()), float32(img.Bounds().Dy())),
				NewVector2f(float32(bounds.Max.X)/64.0-float32(bounds.Min.X)/64.0, float32(-bounds.Max.Y)/64.0),
				float32(ad) / 64.0,
			}
			draw.Draw(atlas_img, image.Rect(x_offset, 0, x_offset+img.Bounds().Dx(), y_offset), img, image.ZP, draw.Src)
			x_offset += img.Bounds().Dx() + GLYPH_ATLAS_GAP
		}
	}

	tempFont.textureAtlas = LoadTextureFromImg(atlas_img)
	return tempFont
}

func (self *FontBatchAtlas) DrawString(_text string, _position Vector2f, _scale float32, _tint RGBA8) {
	if self.fontSettings.Arabic {
		self.drawStringArabic(_text, _position, _scale, _tint)
	} else {
		self.drawStringEnglish(_text, _position, _scale, _tint)
	}
}

func (self *FontBatchAtlas) drawStringEnglish(_text string, _position Vector2f, _scale float32, _tint RGBA8) {

	originalPos := _position

	for _, v := range _text {
		charglyph := self.charAtlasSet[v]
		if v == ' ' {
			originalPos.X += charglyph.advance * _scale
			continue
		} else if v == '\n' {
			originalPos.X = _position.X
			originalPos.Y += self.fontSettings.LineHeight
		}
		loc_pos := originalPos
		loc_pos.Y += (charglyph.bearing.Y) * _scale

		self.sPatch.DrawSpriteBottomLeft(loc_pos, charglyph.size.Scale(_scale), charglyph.uv1, charglyph.uv2, &self.textureAtlas, _tint)
		originalPos.X += charglyph.advance * _scale
	}

}

func (self *FontBatchAtlas) drawStringArabic(_text string, _position Vector2f, _scale float32, _tint RGBA8) {
	_new_text := ShapeArabic(_text)
	_new_text = reverse(_new_text)
	originalPos := _position

	for _, v := range _new_text {
		charglyph := self.charAtlasSet[v]
		if v == ' ' {
			originalPos.X -= charglyph.advance * 1.25 * _scale
			continue
		} else if v == '\n' {
			originalPos.X = _position.X
			originalPos.Y += self.fontSettings.LineHeight
			continue
		}
		loc_pos := originalPos
		loc_pos.X -= charglyph.bearing.X * _scale
		loc_pos.Y += (charglyph.bearing.Y) * _scale
		self.sPatch.DrawSpriteBottomLeft(loc_pos, charglyph.size.Scale(_scale), charglyph.uv1, charglyph.uv2, &self.textureAtlas, _tint)

		originalPos.X -= charglyph.advance * _scale
	}
}

func (self *FontBatchAtlas) Render() {
	self.sPatch.Render(&Cam)
}

func LoadFontTexture(_fontPath string, _fontSettings *FontBatchSettings) Texture2D {
	var tempFont FontBatchAtlas
	tempFont.Init()
	tempFont.fontSettings = *_fontSettings

	resp, err := http.Get(app_url + "/" + _fontPath)
	if err != nil {
		LogF(err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		LogF(err.Error())
	}

	f, err := opentype.Parse(data)

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(_fontSettings.FontSize),
		DPI:     float64(_fontSettings.DPI),
		Hinting: font.HintingFull,
	})
	if err != nil {
		LogF(err.Error())
	}

	max_width, max_height := int(0), int(0)

	dot := fixed.Point26_6{fixed.Int26_6(face.Metrics().Ascent), 0}
	for i := 33; i < 127; i++ {
		char := rune(i)

		_, img, _, _, ok := face.Glyph(dot, char)

		if !ok {
			//LogF("Failed to load rune: %v", char)
		}

		max_width += img.Bounds().Dx() + GLYPH_ATLAS_GAP
		if max_height < img.Bounds().Dy() {
			max_height = img.Bounds().Dy()
		}
	}

	for _, c := range ARABIC_UNICODE {
		_, img, _, _, ok := face.Glyph(dot, c)

		if !ok {
			//LogF("Failed to load rune: %v", c)
		}

		max_width += img.Bounds().Dx() + GLYPH_ATLAS_GAP
		if max_height < img.Bounds().Dy() {
			max_height = img.Bounds().Dy()
		}
	}

	atlas_img := image.NewRGBA(image.Rect(0, 0, max_width, max_height))
	draw.Draw(atlas_img, atlas_img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)

	x_offset := 0
	y_offset := 0

	for i := 32; i < 127; i++ {
		char := rune(i)
		_, img, _, ad, ok := face.Glyph(dot, char)
		if !ok {
			//LogF("Failed To Draw Letter: %v", char)
		}
		y_offset = img.Bounds().Dy()

		if char == ' ' {
			tempFont.charAtlasSet[char] = CharAtlasGlyph{uv1: Vector2fZero, uv2: Vector2fZero, size: Vector2fZero, bearing: Vector2fOne, advance: float32(ad) / float32(1<<6)}
			continue
		}

		bounds, _, ok := face.GlyphBounds(char)
		if !ok {
			//LogF("Failed to load bounds of rune: %v", char)
		}

		tempFont.charAtlasSet[char] = CharAtlasGlyph{
			NewVector2f(float32(x_offset)/float32(max_width), 0),
			NewVector2f(float32(x_offset+img.Bounds().Dx())/float32(max_width), float32(y_offset)/float32(max_height)),
			NewVector2f(float32(img.Bounds().Dx()), float32(img.Bounds().Dy())),
			NewVector2f(float32(bounds.Max.X)/64.0, float32(-bounds.Max.Y)/64.0),
			float32(ad) / 64.0,
		}
		draw.Draw(atlas_img, image.Rect(x_offset, 0, x_offset+img.Bounds().Dx(), y_offset), img, image.ZP, draw.Src)
		x_offset += img.Bounds().Dx() + GLYPH_ATLAS_GAP
	}
	if _fontSettings.Arabic {

		for _, c := range ARABIC_UNICODE {
			_, img, _, ad, ok := face.Glyph(dot, c)
			if !ok {
				//LogF("Failed To Draw Letter: %v", c)
				continue
			}
			y_offset = img.Bounds().Dy()

			bounds, _, ok := face.GlyphBounds(c)
			if !ok {
				//LogF("Failed to load bounds of rune: %v", c)
				continue
			}

			tempFont.charAtlasSet[c] = CharAtlasGlyph{
				NewVector2f(float32(x_offset)/float32(max_width), 0),
				NewVector2f(float32(x_offset+img.Bounds().Dx())/float32(max_width), float32(y_offset)/float32(max_height)),
				NewVector2f(float32(img.Bounds().Dx()), float32(img.Bounds().Dy())),
				NewVector2f(float32(bounds.Max.X)/64.0-float32(bounds.Min.X)/64.0, float32(-bounds.Max.Y)/64.0),
				float32(ad) / 64.0,
			}
			draw.Draw(atlas_img, image.Rect(x_offset, 0, x_offset+img.Bounds().Dx(), y_offset), img, image.ZP, draw.Src)
			x_offset += img.Bounds().Dx() + GLYPH_ATLAS_GAP
		}
	}

	tempFont.textureAtlas = LoadTextureFromImg(atlas_img)
	return tempFont.textureAtlas
}
