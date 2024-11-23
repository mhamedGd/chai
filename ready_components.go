package chai

import (
	"github.com/mhamedGd/chai/customtypes"
	"github.com/mhamedGd/chai/ecs"
	. "github.com/mhamedGd/chai/math"
)

type VisualTransform struct {
	Position   Vector2f
	Z          float32
	Dimensions Vector2f
	Rotation   float32
	Scale      float32
	Tint       RGBA8
	UV1        Vector2f
	UV2        Vector2f
}

type SpriteComponent struct {
	Texture Texture2D
}

type LineRenderComponent struct {
}

type RectRenderComponent struct {
	Dimensions Vector2f
	Tint       RGBA8
}

type FillRectRenderComponent struct {
}

/*
###################################################################
###################################################################
*/

// func (frs *FontRenderSystem) SetFont(_font *FontBatchAtlas) {
// 	frs.fontbatch_atlas = _font
// }

type FontRenderComponent struct {
	Text            string
	Scale           float32
	Offset          Vector2f
	Tint            RGBA8
	Fontbatch_atlas *FontBatchAtlas
	FontSettings    *FontBatchSettings
}

func FontRenderSystem(_thisScene *Scene, _dt float32) {
	Iterate2[VisualTransform, FontRenderComponent](func(i ecs.Id, t *VisualTransform, frc *FontRenderComponent) {
		frc.Fontbatch_atlas.DrawString(frc.Text, t.Position.Add(frc.Offset), frc.Scale, t.Z, frc.Tint)
	})
}

//////////////////////////////////////////////////////////

type TweenValue[T any] struct {
	m_TimeStep float32
	value      T
}

type AnimationComponent[T any] struct {
	Animations customtypes.Map[string, *TweenAnimation[T]]
}

func (a *AnimationComponent[T]) Play(_animationName string) {
	for _, val := range a.Animations.AllItems() {
		val.m_TimeStepFactor = 0.0
		a.Animations.Get(_animationName).HasFinished = true
	}

	a.Animations.Get(_animationName).HasFinished = false
	a.Animations.Get(_animationName).m_TimeStepFactor = 1.0
}

func (a *AnimationComponent[T]) PlaySimultaneous(_animationNames ...string) {
	for i := range _animationNames {
		a.Animations.Get(_animationNames[i]).m_TimeStepFactor = 1.0
	}
}

func (a *AnimationComponent[T]) Stop(_animationNames ...string) {
	for i := range _animationNames {
		a.Animations.Get(_animationNames[i]).m_TimeStepFactor = 0.0
	}
}

func (a *AnimationComponent[T]) StopAll() {
	for _, val := range a.Animations.AllItems() {
		val.m_TimeStepFactor = 0.0
	}
}

func NewAnimationComponentInt() AnimationComponent[int] {
	return AnimationComponent[int]{
		Animations: customtypes.NewMap[string, *TweenAnimation[int]](),
	}
}

func NewAnimationComponentFloat32() AnimationComponent[float32] {
	return AnimationComponent[float32]{
		Animations: customtypes.NewMap[string, *TweenAnimation[float32]](),
	}
}

func NewAnimationComponentVector2f() AnimationComponent[Vector2f] {
	return AnimationComponent[Vector2f]{
		Animations: customtypes.NewMap[string, *TweenAnimation[Vector2f]](),
	}
}

func NewAnimationComponentVector2i() AnimationComponent[Vector2i] {
	return AnimationComponent[Vector2i]{
		Animations: customtypes.NewMap[string, *TweenAnimation[Vector2i]](),
	}
}

type TweenAnimation[T any] struct {
	KeyframeValues   customtypes.List[TweenValue[T]]
	m_CurrentValue   T
	m_CurrentIndex   int
	Length           float32
	CurrentTimestep  float32
	m_TimeStepFactor float32
	Loop             bool
	HasFinished      bool
}

func (comp *AnimationComponent[T]) GetCurrentValue(_animationName string) T {
	return comp.Animations.Get(_animationName).m_CurrentValue
}

func (comp *AnimationComponent[T]) HasFinished(_animationName string) bool {
	return comp.Animations.Get(_animationName).HasFinished
}

func (comp TweenAnimation[T]) IsPlaying() bool {
	return comp.m_TimeStepFactor != 0.0
}

func (anim *AnimationComponent[int]) NewTweenAnimationInt(_animationName string) {
	anim.Animations.Set(_animationName, &TweenAnimation[int]{
		KeyframeValues:   customtypes.NewList[TweenValue[int]](),
		m_TimeStepFactor: 0.0,
	})
}

func (anim *AnimationComponent[float32]) NewTweenAnimationFloat32(_animationName string, loop bool) {
	anim.Animations.Set(_animationName, &TweenAnimation[float32]{
		KeyframeValues:   customtypes.NewList[TweenValue[float32]](),
		m_TimeStepFactor: 0.0,
		Loop:             loop,
	})
}

func (anim *AnimationComponent[Vector2f]) NewTweenAnimationVector2f(_animationName string, loop bool) {
	anim.Animations.Set(_animationName, &TweenAnimation[Vector2f]{
		KeyframeValues:   customtypes.NewList[TweenValue[Vector2f]](),
		m_TimeStepFactor: 0.0,
		Loop:             loop,
	})
}

func (anim *AnimationComponent[Vector2i]) NewTweenAnimationVector2i(_animationName string) {
	anim.Animations.Set(_animationName, &TweenAnimation[Vector2i]{
		KeyframeValues:   customtypes.NewList[TweenValue[Vector2i]](),
		m_TimeStepFactor: 0.0,
	})
}

func (comp *AnimationComponent[T]) RegisterKeyframe(_animationName string, m_TimeStep float32, value T) {
	anim := comp.Animations.Get(_animationName)
	// anim.KeyframeValues = append(anim.KeyframeValues, TweenValue[T]{m_TimeStep: m_TimeStep, value: value})
	anim.KeyframeValues.PushBack(TweenValue[T]{m_TimeStep: m_TimeStep, value: value})

	lowest := anim.KeyframeValues.Data[0].m_TimeStep
	for _, val := range anim.KeyframeValues.Data {
		if lowest > val.m_TimeStep {
			lowest = val.m_TimeStep
		}
	}

	anim.m_CurrentValue = anim.KeyframeValues.Data[0].value
	comp.Animations.Set(_animationName, anim)
}

func TweenAnimatorSystem(_thisScene *Scene, _dt float32) {
	// Float Tween Animation
	/////////////////////////////////////////////////////
	Iterate1[AnimationComponent[float32]](func(i EntId, an *AnimationComponent[float32]) {
		for _, tween := range an.Animations.AllItems() {
			if !tween.IsPlaying() || tween.HasFinished {
				continue
			}
			tween.CurrentTimestep += _dt * tween.m_TimeStepFactor
			tween.m_CurrentValue = LerpFloat32(tween.KeyframeValues.Data[tween.m_CurrentIndex].value, tween.KeyframeValues.Data[tween.m_CurrentIndex+1].value, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep)/(tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep {
				tween.m_CurrentIndex++
				if tween.m_CurrentIndex == tween.KeyframeValues.Count()-1 {
					tween.m_CurrentValue = tween.KeyframeValues.Data[tween.m_CurrentIndex].value
					tween.m_CurrentIndex = 0
					if tween.Loop {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].m_TimeStep
					} else {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].m_TimeStep
						tween.HasFinished = true
					}
				}
			}
		}
	})

	// Vector2f Tween Animation
	/////////////////////////////////////////////////////
	Iterate1[AnimationComponent[Vector2f]](func(i EntId, an *AnimationComponent[Vector2f]) {
		for _, tween := range an.Animations.AllItems() {
			if !tween.IsPlaying() || tween.HasFinished {
				continue
			}
			tween.CurrentTimestep += _dt * tween.m_TimeStepFactor
			tween.m_CurrentValue.X = LerpFloat32(tween.KeyframeValues.Data[tween.m_CurrentIndex].value.X, tween.KeyframeValues.Data[tween.m_CurrentIndex+1].value.X, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep)/(tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep))
			tween.m_CurrentValue.Y = LerpFloat32(tween.KeyframeValues.Data[tween.m_CurrentIndex].value.Y, tween.KeyframeValues.Data[tween.m_CurrentIndex+1].value.Y, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep)/(tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep {
				tween.m_CurrentIndex++
				if tween.m_CurrentIndex >= tween.KeyframeValues.Count()-1 {
					tween.m_CurrentValue.X = tween.KeyframeValues.Data[tween.m_CurrentIndex].value.X
					tween.m_CurrentValue.Y = tween.KeyframeValues.Data[tween.m_CurrentIndex].value.Y

					tween.m_CurrentIndex = 0
					if tween.Loop {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].m_TimeStep
					} else {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].m_TimeStep
						tween.HasFinished = true
					}
				}
			}
		}
	})

	// Int Tween Animation
	/////////////////////////////////////////////////////
	Iterate1[AnimationComponent[int]](func(i EntId, ac *AnimationComponent[int]) {
		for _, tween := range ac.Animations.AllItems() {

			// passingTime := tween.KeyframeValues[tween.m_CurrentIndex+1].m_TimeStep - tween.KeyframeValues[tween.m_CurrentIndex].m_TimeStep
			tween.CurrentTimestep += _dt * tween.m_TimeStepFactor
			tween.m_CurrentValue = LerpInt(tween.KeyframeValues.Data[tween.m_CurrentIndex].value, tween.KeyframeValues.Data[tween.m_CurrentIndex+1].value, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep)/(tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep-tween.KeyframeValues.Data[tween.m_CurrentIndex].m_TimeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep {
				tween.m_CurrentIndex++
				if tween.m_CurrentIndex == tween.KeyframeValues.Count()-1 {
					tween.m_CurrentIndex = 0
					tween.CurrentTimestep = tween.KeyframeValues.Data[0].m_TimeStep
				}
			}
		}
	})

	// Vector2i Tween Animation
	//////////////////////////////////////////////////////
	Iterate1[AnimationComponent[Vector2i]](func(i EntId, ac *AnimationComponent[Vector2i]) {

		for _, tween := range ac.Animations.AllItems() {
			tween.CurrentTimestep += _dt * tween.m_TimeStepFactor

			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.m_CurrentIndex+1].m_TimeStep {
				tween.m_CurrentIndex++
				if tween.m_CurrentIndex == tween.KeyframeValues.Count()-1 {
					tween.m_CurrentIndex = 0
					tween.CurrentTimestep = tween.KeyframeValues.Data[0].m_TimeStep
				}
				tween.m_CurrentValue = tween.KeyframeValues.Data[tween.m_CurrentIndex].value
			}
		}
	})
}

// Sprite Animation
//////////////////////////////////////////////////////

type SpriteSheet struct {
	Texture               Texture2D
	Coloumns, Rows        int
	TileWidth, TileHeight int
	TextureCoordinates    customtypes.List[customtypes.Pair[Vector2f, Vector2f]]
}

func NewSpriteSheet(_tex Texture2D, _tileWidth, _tileHeight int) SpriteSheet {
	_sprite_sheet := SpriteSheet{}
	_sprite_sheet.TextureCoordinates = customtypes.NewList[customtypes.Pair[Vector2f, Vector2f]]()

	_coloumns := _tex.Width / _tileWidth
	_rows := _tex.Height / _tileHeight
	for x := 0; x < _coloumns; x++ {
		for y := 0; y < _rows; y++ {
			_uv1 := NewVector2f(float32(x)/float32(_coloumns), float32(y)/float32(_rows))
			_uv2 := NewVector2f(float32(x+1)/float32(_coloumns), float32(y+1)/float32(_rows))
			_sprite_sheet.TextureCoordinates.PushBack(customtypes.Pair[Vector2f, Vector2f]{_uv1, _uv2})
		}
	}

	_sprite_sheet.Coloumns = _coloumns
	_sprite_sheet.Rows = _rows
	_sprite_sheet.TileWidth = _tileWidth
	_sprite_sheet.TileHeight = _tileHeight
	_sprite_sheet.Texture = _tex

	return _sprite_sheet
}

type SpriteAnimaionComponent struct {
	CurrentAnimation string
	CurrentTimestep  float32
	CurrentAnimStep  uint16
	AnimationSpeed   int
	Animations       customtypes.Map[string, customtypes.List[Vector2i]]
	m_SpriteSheet    *SpriteSheet
}

func NewSpriteAnimationComponent(_spriteSheet *SpriteSheet) SpriteAnimaionComponent {
	return SpriteAnimaionComponent{
		CurrentAnimation: "",
		CurrentTimestep:  0.0,
		CurrentAnimStep:  0,
		Animations:       customtypes.NewMap[string, customtypes.List[Vector2i]](),
		m_SpriteSheet:    _spriteSheet,
	}
}

func (_spa *SpriteAnimaionComponent) NewAnimation(_animationName string) {
	_spa.Animations.Set(_animationName, customtypes.NewList[Vector2i]())
}

func (_spa *SpriteAnimaionComponent) RegisterFrame(_animationName string, _value Vector2i) {
	_anim_list := _spa.Animations.Get(_animationName)
	_anim_list.PushBack(_value)
	_spa.Animations.Set(_animationName, _anim_list)
}

func (_spa *SpriteAnimaionComponent) RegisterFrames(_animationName string, _values []Vector2i) {
	_anim_list := _spa.Animations.Get(_animationName)
	for i := 0; i < len(_values); i++ {
		_anim_list.PushBack(_values[i])
	}
	_spa.Animations.Set(_animationName, _anim_list)
}

func SpriteAnimationSystem(_thisScene *Scene, _dt float32) {
	Iterate2(func(i ecs.Id, t *VisualTransform, b *SpriteAnimaionComponent) {
		if !b.Animations.Has(b.CurrentAnimation) {
			return
		}

		// _animation_speed := 1.0 / float32(b.AnimationSpeed)
		b.CurrentTimestep += _dt * float32(b.AnimationSpeed)
		_animation_coords := b.Animations.Get(b.CurrentAnimation)
		if b.CurrentTimestep >= 1.0 {
			b.CurrentTimestep = 0.0
			b.CurrentAnimStep += 1
			if b.CurrentAnimStep >= uint16(_animation_coords.Count()) {
				b.CurrentAnimStep = 0
			}
		}

		_uv_coords := _animation_coords.Data[b.CurrentAnimStep]
		t.UV1 = b.m_SpriteSheet.TextureCoordinates.Data[b.m_SpriteSheet.Coloumns*_uv_coords.Y+_uv_coords.X].First
		t.UV2 = b.m_SpriteSheet.TextureCoordinates.Data[b.m_SpriteSheet.Coloumns*_uv_coords.Y+_uv_coords.X].Second
	})
}

// Debug Draw
//////////////////////////////////////////////////////

func DebugBodyDrawSystem(_thisScene *Scene, _dt float32) {
	_original := lineWidth
	_color := NewRGBA8(0, 255, 0, 255)
	_z := float32(0)
	lineWidth = 0.01
	Iterate1[DynamicBodyComponent](func(i ecs.Id, dbc *DynamicBodyComponent) {
		if dbc.m_Settings.ColliderShape == SHAPE_RECTBODY {
			DrawRect(dbc.GetPosition(), dbc.m_Settings.StartDimensions, _color, _z, dbc.GetRotation())
		} else {
			DrawCircle(dbc.GetPosition(), dbc.m_Settings.StartDimensions.X/2.0, _color, _z)
		}
	})
	Iterate1[StaticBodyComponent](func(i ecs.Id, sbc *StaticBodyComponent) {
		if sbc.m_Settings.ColliderShape == SHAPE_RECTBODY {
			DrawRect(sbc.GetPosition(), sbc.m_Settings.StartDimensions, _color, _z, sbc.GetRotation())

		} else {
			DrawCircle(sbc.GetPosition(), sbc.m_Settings.StartDimensions.X/2.0, _color, _z)
		}
	})
	Iterate1[KinematicBodyComponent](func(i ecs.Id, kbc *KinematicBodyComponent) {
		if kbc.m_Settings.ColliderShape == SHAPE_RECTBODY {
			DrawRect(kbc.GetPosition(), kbc.m_Settings.StartDimensions, _color, _z, kbc.GetRotation())
		} else {
			DrawCircle(kbc.GetPosition(), kbc.m_Settings.StartDimensions.X/2.0, _color, _z)
		}
	})

	for c := range collisionTiles.Data {
		DrawRect(collisionTiles.Data[c].First, collisionTiles.Data[c].Second, _color, _z, 0.0)
	}

	lineWidth = _original
}
