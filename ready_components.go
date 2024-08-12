package chai

import "github.com/mhamedGd/chai/ecs"

type VisualTransform struct {
	Position   Vector2f
	Dimensions Vector2f
	Rotation   float32
	Scale      float32
	Tint       RGBA8
	UV1        Vector2f
	UV2        Vector2f
}

type SpriteComponent struct {
	Texture Texture2D
	Tint    RGBA8
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
			Shapes.DrawTriangleRotated(t.Position, tri.Dimensions.Scale(t.Scale), tri.Tint, t.Rotation)
		}
	})

	queryRect := ecs.Query2[VisualTransform, RectRenderComponent](GetCurrentScene().Ecs_World)
	queryRect.MapId(func(id ecs.Id, t *VisualTransform, rect *RectRenderComponent) {
		if Cam.IsBoxInView(t.Position, rect.Dimensions.Scale(t.Scale)) {
			Shapes.DrawRectRotated(t.Position, rect.Dimensions.Scale(t.Scale), rect.Tint, t.Rotation)
		}
	})

	queryFillRectBottom := ecs.Query2[VisualTransform, FillRectBottomRenderComponent](GetCurrentScene().Ecs_World)
	queryFillRectBottom.MapId(func(id ecs.Id, t *VisualTransform, rect *FillRectBottomRenderComponent) {
		rectDims := rect.Dimensions.Scale(t.Scale)
		if Cam.IsBoxInView(t.Position.Subtract(rectDims.Scale(0.5)), rectDims) {
			Shapes.DrawFillRectBottomRotated(t.Position, rectDims, rect.Tint, t.Rotation)
		}
	})
}

type LineRenderComponent struct {
	Tint RGBA8
}

type TriangleRenderComponent struct {
	Dimensions  Vector2f
	OffsetPivot Vector2f
	Tint        RGBA8
}

type FillTriangleRenderComponent struct {
	Tint RGBA8
}

type RectRenderComponent struct {
	Dimensions Vector2f
	Tint       RGBA8
}

type FillRectRenderComponent struct {
	Tint          RGBA8
	QuadTreeIndex int64
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
		frc.Fontbatch_atlas.DrawString(frc.Text, t.Position.Add(frc.Offset), frc.Scale, frc.Tint)
	})
}

//////////////////////////////////////////////////////////

type TweenValue[T any] struct {
	timeStep float32
	value    T
}

type AnimationComponent[T any] struct {
	Animations Map[string, *TweenAnimation[T]]
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
		Animations: NewMap[string, *TweenAnimation[int]](),
	}
}

func NewAnimationComponentFloat32() AnimationComponent[float32] {
	return AnimationComponent[float32]{
		Animations: NewMap[string, *TweenAnimation[float32]](),
	}
}

func NewAnimationComponentVector2f() AnimationComponent[Vector2f] {
	return AnimationComponent[Vector2f]{
		Animations: NewMap[string, *TweenAnimation[Vector2f]](),
	}
}

func NewAnimationComponentVector2i() AnimationComponent[Vector2i] {
	return AnimationComponent[Vector2i]{
		Animations: NewMap[string, *TweenAnimation[Vector2i]](),
	}
}

type TweenAnimation[T any] struct {
	KeyframeValues  List[TweenValue[T]]
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
		KeyframeValues: NewList[TweenValue[int]](),
		timeStepFactor: 0.0,
	})
}

func (anim *AnimationComponent[float32]) NewTweenAnimationFloat32(animationName string, loop bool) {
	anim.Animations.Set(animationName, &TweenAnimation[float32]{
		KeyframeValues: NewList[TweenValue[float32]](),
		timeStepFactor: 0.0,
		Loop:           loop,
	})
}

func (anim *AnimationComponent[Vector2f]) NewTweenAnimationVector2f(animationName string, loop bool) {
	anim.Animations.Set(animationName, &TweenAnimation[Vector2f]{
		KeyframeValues: NewList[TweenValue[Vector2f]](),
		timeStepFactor: 0.0,
	})
}

func (anim *AnimationComponent[Vector2i]) NewTweenAnimationVector2i(animationName string) {
	anim.Animations.Set(animationName, &TweenAnimation[Vector2i]{
		KeyframeValues: NewList[TweenValue[Vector2i]](),
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
