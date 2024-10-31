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

// func NewVisualTransform(_position Vector2f, _z float32, _dimensions Vector2f) VisualTransform {
// 	return VisualTransform{
// 		Position:   _position,
// 		Z:          _z,
// 		Dimensions: _dimensions.Multp(pixelsPerMeterDimensions),
// 		Rotation:   0.0,
// 		Scale:      1.0,
// 		Tint:       WHITE,
// 		UV1:        Vector2fZero,
// 		UV2:        Vector2fOne,
// 	}
// }
// func (vt *VisualTransform) WithRotation(_rotation float32) {
// 	vt.Rotation = _rotation
// }
// func (vt *VisualTransform) WithScale(_scale float32) {
// 	vt.Scale = _scale
// }
// func (vt *VisualTransform) WithTint(_tint RGBA8) {
// 	vt.Tint = _tint
// }
// func (vt *VisualTransform) WithUVs(_uv1, _uv2 Vector2f) {
// 	vt.UV1 = _uv1
// 	vt.UV2 = _uv2
// }

type SpriteComponent struct {
	Texture Texture2D
}

func SpriteRenderSystem(_this_scene *Scene, _dt float32) {
	// query2 := ecs.Query2[VisualTransform, SpriteComponent](GetCurrentScene().Ecs_World)
	// query2.MapId(func(id ecs.Id, t *VisualTransform, s *SpriteComponent) {
	// 	newOffset := _render.Offset.Rotate(t.Rotation, Vector2fZero)
	// 	halfDim := NewVector2f(newOffset.X*float32(s.Texture.Width)/2.0, newOffset.Y*float32(s.Texture.Height)/2.0)
	// 	_render.Sprites.DrawSpriteOriginScaledRotated(t.Position.Add(halfDim), Vector2fZero, Vector2fOne, _render.Scale, &s.Texture, s.Tint, t.Rotation)
	// })
}

func ShapesDrawingSystem(_this_scene *Scene, dt float32) {
	queryTri := ecs.Query2[VisualTransform, TriangleRenderComponent](GetCurrentScene().Ecs_World)
	queryTri.MapId(func(id ecs.Id, t *VisualTransform, tri *TriangleRenderComponent) {
		if Cam.IsBoxInView(t.Position, tri.Dimensions.Scale(t.Scale)) {
			Shapes.DrawTriangleRotated(t.Position, t.Z, tri.Dimensions.Scale(t.Scale), tri.Tint, t.Rotation)
		}
	})

	queryRect := ecs.Query2[VisualTransform, RectRenderComponent](GetCurrentScene().Ecs_World)
	queryRect.MapId(func(id ecs.Id, t *VisualTransform, rect *RectRenderComponent) {
		if Cam.IsBoxInView(t.Position, rect.Dimensions.Scale(t.Scale)) {
			Shapes.DrawRectRotated(t.Position, t.Z, rect.Dimensions.Scale(t.Scale), rect.Tint, t.Rotation)
		}
	})

	queryFillRectBottom := ecs.Query2[VisualTransform, FillRectBottomRenderComponent](GetCurrentScene().Ecs_World)
	queryFillRectBottom.MapId(func(id ecs.Id, t *VisualTransform, rect *FillRectBottomRenderComponent) {
		rectDims := rect.Dimensions.Scale(t.Scale)
		if Cam.IsBoxInView(t.Position.Subtract(rectDims.Scale(0.5)), rectDims) {
			Shapes.DrawFillRectBottomRotated(t.Position, t.Z, rectDims, rect.Tint, t.Rotation)
		}
	})
}

type LineRenderComponent struct {
}

type TriangleRenderComponent struct {
	Dimensions  Vector2f
	OffsetPivot Vector2f
	Tint        RGBA8
}

type FillTriangleRenderComponent struct {
}

type RectRenderComponent struct {
	Dimensions Vector2f
	Tint       RGBA8
}

type FillRectRenderComponent struct {
}

type FillRectBottomRenderComponent struct {
	Dimensions Vector2f
	Tint       RGBA8
}

type CircleRenderComponent struct {
	Tint RGBA8
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

func FontRenderSystem(_this_scene *Scene, _dt float32) {
	Iterate2[VisualTransform, FontRenderComponent](func(i ecs.Id, t *VisualTransform, frc *FontRenderComponent) {
		frc.Fontbatch_atlas.DrawString(frc.Text, t.Position.Add(frc.Offset), frc.Scale, t.Z, frc.Tint)
	})
}

//////////////////////////////////////////////////////////

type TweenValue[T any] struct {
	timeStep float32
	value    T
}

type AnimationComponent[T any] struct {
	Animations customtypes.Map[string, *TweenAnimation[T]]
}

func (a *AnimationComponent[T]) Play(animationName string) {
	for _, val := range a.Animations.AllItems() {
		val.timeStepFactor = 0.0
		a.Animations.Get(animationName).HasFinished = true
	}

	a.Animations.Get(animationName).HasFinished = false
	a.Animations.Get(animationName).timeStepFactor = 1.0
}

func (a *AnimationComponent[T]) PlaySimultaneous(animationNames ...string) {
	for i := range animationNames {
		a.Animations.Get(animationNames[i]).timeStepFactor = 1.0
	}
}

func (a *AnimationComponent[T]) Stop(animationNames ...string) {
	for i := range animationNames {
		a.Animations.Get(animationNames[i]).timeStepFactor = 0.0
	}
}

func (a *AnimationComponent[T]) StopAll() {
	for _, val := range a.Animations.AllItems() {
		val.timeStepFactor = 0.0
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
	KeyframeValues  customtypes.List[TweenValue[T]]
	currentValue    T
	currentIndex    int
	Length          float32
	CurrentTimestep float32
	timeStepFactor  float32
	Loop            bool
	HasFinished     bool
}

func (comp *AnimationComponent[T]) GetCurrentValue(animationName string) T {
	return comp.Animations.Get(animationName).currentValue
}

func (comp *AnimationComponent[T]) HasFinished(animationName string) bool {
	return comp.Animations.Get(animationName).HasFinished
}

func (comp TweenAnimation[T]) IsPlaying() bool {
	return comp.timeStepFactor != 0.0
}

func (anim *AnimationComponent[int]) NewTweenAnimationInt(animationName string) {
	anim.Animations.Set(animationName, &TweenAnimation[int]{
		KeyframeValues: customtypes.NewList[TweenValue[int]](),
		timeStepFactor: 0.0,
	})
}

func (anim *AnimationComponent[float32]) NewTweenAnimationFloat32(animationName string, loop bool) {
	anim.Animations.Set(animationName, &TweenAnimation[float32]{
		KeyframeValues: customtypes.NewList[TweenValue[float32]](),
		timeStepFactor: 0.0,
		Loop:           loop,
	})
}

func (anim *AnimationComponent[Vector2f]) NewTweenAnimationVector2f(animationName string, loop bool) {
	anim.Animations.Set(animationName, &TweenAnimation[Vector2f]{
		KeyframeValues: customtypes.NewList[TweenValue[Vector2f]](),
		timeStepFactor: 0.0,
		Loop:           loop,
	})
}

func (anim *AnimationComponent[Vector2i]) NewTweenAnimationVector2i(animationName string) {
	anim.Animations.Set(animationName, &TweenAnimation[Vector2i]{
		KeyframeValues: customtypes.NewList[TweenValue[Vector2i]](),
		timeStepFactor: 0.0,
	})
}

func (comp *AnimationComponent[T]) RegisterKeyframe(animationName string, timeStep float32, value T) {
	anim := comp.Animations.Get(animationName)
	// anim.KeyframeValues = append(anim.KeyframeValues, TweenValue[T]{timeStep: timeStep, value: value})
	anim.KeyframeValues.PushBack(TweenValue[T]{timeStep: timeStep, value: value})

	lowest := anim.KeyframeValues.Data[0].timeStep
	for _, val := range anim.KeyframeValues.Data {
		if lowest > val.timeStep {
			lowest = val.timeStep
		}
	}

	anim.currentValue = anim.KeyframeValues.Data[0].value
	comp.Animations.Set(animationName, anim)
}

func TweenAnimatorSystem(_this_scene *Scene, _dt float32) {
	// Float Tween Animation
	/////////////////////////////////////////////////////
	Iterate1[AnimationComponent[float32]](func(i EntId, an *AnimationComponent[float32]) {
		for _, tween := range an.Animations.AllItems() {
			if !tween.IsPlaying() || tween.HasFinished {
				continue
			}
			tween.CurrentTimestep += _dt * tween.timeStepFactor
			tween.currentValue = LerpFloat32(tween.KeyframeValues.Data[tween.currentIndex].value, tween.KeyframeValues.Data[tween.currentIndex+1].value, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.currentIndex].timeStep)/(tween.KeyframeValues.Data[tween.currentIndex+1].timeStep-tween.KeyframeValues.Data[tween.currentIndex].timeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.currentIndex+1].timeStep {
				tween.currentIndex++
				if tween.currentIndex == tween.KeyframeValues.Count()-1 {
					tween.currentValue = tween.KeyframeValues.Data[tween.currentIndex].value
					tween.currentIndex = 0
					if tween.Loop {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].timeStep
					} else {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].timeStep
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
			tween.CurrentTimestep += _dt * tween.timeStepFactor
			tween.currentValue.X = LerpFloat32(tween.KeyframeValues.Data[tween.currentIndex].value.X, tween.KeyframeValues.Data[tween.currentIndex+1].value.X, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.currentIndex].timeStep)/(tween.KeyframeValues.Data[tween.currentIndex+1].timeStep-tween.KeyframeValues.Data[tween.currentIndex].timeStep))
			tween.currentValue.Y = LerpFloat32(tween.KeyframeValues.Data[tween.currentIndex].value.Y, tween.KeyframeValues.Data[tween.currentIndex+1].value.Y, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.currentIndex].timeStep)/(tween.KeyframeValues.Data[tween.currentIndex+1].timeStep-tween.KeyframeValues.Data[tween.currentIndex].timeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.currentIndex+1].timeStep {
				tween.currentIndex++
				if tween.currentIndex >= tween.KeyframeValues.Count()-1 {
					tween.currentValue.X = tween.KeyframeValues.Data[tween.currentIndex].value.X
					tween.currentValue.Y = tween.KeyframeValues.Data[tween.currentIndex].value.Y

					tween.currentIndex = 0
					if tween.Loop {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].timeStep
					} else {
						tween.CurrentTimestep = tween.KeyframeValues.Data[0].timeStep
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

			// passingTime := tween.KeyframeValues[tween.currentIndex+1].timeStep - tween.KeyframeValues[tween.currentIndex].timeStep
			tween.CurrentTimestep += _dt * tween.timeStepFactor
			tween.currentValue = LerpInt(tween.KeyframeValues.Data[tween.currentIndex].value, tween.KeyframeValues.Data[tween.currentIndex+1].value, (tween.CurrentTimestep-tween.KeyframeValues.Data[tween.currentIndex].timeStep)/(tween.KeyframeValues.Data[tween.currentIndex+1].timeStep-tween.KeyframeValues.Data[tween.currentIndex].timeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.currentIndex+1].timeStep {
				tween.currentIndex++
				if tween.currentIndex == tween.KeyframeValues.Count()-1 {
					tween.currentIndex = 0
					tween.CurrentTimestep = tween.KeyframeValues.Data[0].timeStep
				}
			}
		}
	})

	// Vector2i Tween Animation
	//////////////////////////////////////////////////////
	Iterate1[AnimationComponent[Vector2i]](func(i EntId, ac *AnimationComponent[Vector2i]) {

		for _, tween := range ac.Animations.AllItems() {
			tween.CurrentTimestep += _dt * tween.timeStepFactor

			if tween.CurrentTimestep >= tween.KeyframeValues.Data[tween.currentIndex+1].timeStep {
				tween.currentIndex++
				if tween.currentIndex == tween.KeyframeValues.Count()-1 {
					tween.currentIndex = 0
					tween.CurrentTimestep = tween.KeyframeValues.Data[0].timeStep
				}
				tween.currentValue = tween.KeyframeValues.Data[tween.currentIndex].value
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

func NewSpriteSheet(_tex Texture2D, _tile_width, _tile_height int) SpriteSheet {
	_sprite_sheet := SpriteSheet{}
	_sprite_sheet.TextureCoordinates = customtypes.NewList[customtypes.Pair[Vector2f, Vector2f]]()

	_coloumns := _tex.Width / _tile_width
	_rows := _tex.Height / _tile_height
	for x := 0; x < _coloumns; x++ {
		for y := 0; y < _rows; y++ {
			_uv1 := NewVector2f(float32(x)/float32(_coloumns), float32(y)/float32(_rows))
			_uv2 := NewVector2f(float32(x+1)/float32(_coloumns), float32(y+1)/float32(_rows))
			_sprite_sheet.TextureCoordinates.PushBack(customtypes.Pair[Vector2f, Vector2f]{_uv1, _uv2})
		}
	}

	_sprite_sheet.Coloumns = _coloumns
	_sprite_sheet.Rows = _rows
	_sprite_sheet.TileWidth = _tile_width
	_sprite_sheet.TileHeight = _tile_height
	_sprite_sheet.Texture = _tex

	return _sprite_sheet
}

type SpriteAnimaionComponent struct {
	CurrentAnimation string
	CurrentTimestep  float32
	CurrentAnimStep  uint16
	AnimationSpeed   int
	Animations       customtypes.Map[string, customtypes.List[Vector2i]]
	spriteSheet      *SpriteSheet
}

func NewSpriteAnimationComponent(_sprite_sheet *SpriteSheet) SpriteAnimaionComponent {
	return SpriteAnimaionComponent{
		CurrentAnimation: "",
		CurrentTimestep:  0.0,
		CurrentAnimStep:  0,
		Animations:       customtypes.NewMap[string, customtypes.List[Vector2i]](),
		spriteSheet:      _sprite_sheet,
	}
}

func (_spa *SpriteAnimaionComponent) NewAnimation(_anim_name string) {
	_spa.Animations.Set(_anim_name, customtypes.NewList[Vector2i]())
}

func (_spa *SpriteAnimaionComponent) RegisterFrame(_anim_name string, _value Vector2i) {
	_anim_list := _spa.Animations.Get(_anim_name)
	_anim_list.PushBack(_value)
	_spa.Animations.Set(_anim_name, _anim_list)
}

func (_spa *SpriteAnimaionComponent) RegisterFrames(_anim_name string, _values []Vector2i) {
	_anim_list := _spa.Animations.Get(_anim_name)
	for i := 0; i < len(_values); i++ {
		_anim_list.PushBack(_values[i])
	}
	_spa.Animations.Set(_anim_name, _anim_list)
}

func SpriteAnimationSystem(_this_scene *Scene, _dt float32) {
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
		t.UV1 = b.spriteSheet.TextureCoordinates.Data[b.spriteSheet.Coloumns*_uv_coords.Y+_uv_coords.X].First
		t.UV2 = b.spriteSheet.TextureCoordinates.Data[b.spriteSheet.Coloumns*_uv_coords.Y+_uv_coords.X].Second
	})
}

// Debug Draw
//////////////////////////////////////////////////////

func DebugBodyDrawSystem(_this_scene *Scene, _dt float32) {
	_original := Shapes.LineWidth
	_color := NewRGBA8(0, 255, 0, 150)
	_z := float32(-1)
	Shapes.LineWidth = 0.1 / float32(pixelsPerMeter)
	Iterate1[DynamicBodyComponent](func(i ecs.Id, dbc *DynamicBodyComponent) {
		if dbc.settings.ColliderShape == SHAPE_RECTBODY {
			Shapes.DrawRectRotated(dbc.GetPosition(), _z, dbc.settings.StartDimensions, _color, dbc.GetRotation())
		} else {
			Shapes.DrawCircle(dbc.GetPosition(), _z, dbc.settings.StartDimensions.X/2.0, _color)
		}
	})
	Iterate1[StaticBodyComponent](func(i ecs.Id, sbc *StaticBodyComponent) {
		if sbc.settings.ColliderShape == SHAPE_RECTBODY {
			Shapes.DrawRectRotated(sbc.GetPosition(), _z, sbc.settings.StartDimensions, _color, sbc.GetRotation())
		} else {
			Shapes.DrawCircle(sbc.GetPosition(), _z, sbc.settings.StartDimensions.X/2.0, _color)
		}
	})
	Iterate1[KinematicBodyComponent](func(i ecs.Id, kbc *KinematicBodyComponent) {
		if kbc.settings.ColliderShape == SHAPE_RECTBODY {
			Shapes.DrawRectRotated(kbc.GetPosition(), _z, kbc.settings.StartDimensions, _color, kbc.GetRotation())
		} else {
			Shapes.DrawCircle(kbc.GetPosition(), _z, kbc.settings.StartDimensions.X/2.0, _color)
		}
	})

	Shapes.LineWidth = _original
}
