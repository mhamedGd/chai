package chai

type TweenValue[T any] struct {
	timeStep float32
	value    T
}

type AnimationComponent[T any] struct {
	Component
	Animations map[string]*TweenAnimation[T]
}

func (a *AnimationComponent[T]) Play(animationName string) {
	for _, val := range a.Animations {
		val.timeStepFactor = 0.0
	}

	a.Animations[animationName].timeStepFactor = 1.0
}

func (a *AnimationComponent[T]) PlaySimultaneous(animationNames ...string) {
	for i := range animationNames {
		a.Animations[animationNames[i]].timeStepFactor = 1.0
	}
}

func (a *AnimationComponent[T]) Stop(animationNames ...string) {
	for i := range animationNames {
		a.Animations[animationNames[i]].timeStepFactor = 0.0
	}
}

func (a *AnimationComponent[T]) StopAll() {
	for _, val := range a.Animations {
		val.timeStepFactor = 0.0
	}
}

func NewAnimationComponentInt() AnimationComponent[int] {
	return AnimationComponent[int]{
		Animations: make(map[string]*TweenAnimation[int]),
	}
}

func NewAnimationComponentFloat32() AnimationComponent[float32] {
	return AnimationComponent[float32]{
		Animations: make(map[string]*TweenAnimation[float32]),
	}
}

func NewAnimationComponentVector2i() AnimationComponent[Vector2i] {
	return AnimationComponent[Vector2i]{
		Animations: make(map[string]*TweenAnimation[Vector2i]),
	}
}

type TweenAnimation[T any] struct {
	KeyframeValues  []TweenValue[T]
	currentValue    T
	currentIndex    int
	Length          float32
	CurrentTimestep float32
	timeStepFactor  float32
	Loop            bool
	HasFinished     bool
}

func (comp *AnimationComponent[T]) GetCurrentValue(animationName string) T {
	return comp.Animations[animationName].currentValue
}

func (comp *AnimationComponent[T]) HasFinished(animationName string) bool {
	return comp.Animations[animationName].HasFinished
}

func (comp TweenAnimation[T]) IsPlaying() bool {
	return comp.timeStepFactor != 0.0
}

func (t *AnimationComponent[T]) ComponentSet(val interface{}) {
	*t = val.(AnimationComponent[T])
}

func (anim *AnimationComponent[int]) NewTweenAnimationInt(animationName string) {
	anim.Animations[animationName] = &TweenAnimation[int]{
		KeyframeValues: make([]TweenValue[int], 0),
		timeStepFactor: 0.0,
	}
}

func (anim *AnimationComponent[float32]) NewTweenAnimationFloat32(animationName string, loop bool) {
	anim.Animations[animationName] = &TweenAnimation[float32]{
		KeyframeValues: make([]TweenValue[float32], 0),
		timeStepFactor: 0.0,
		Loop:           loop,
	}
}

func (anim *AnimationComponent[Vector2i]) NewTweenAnimationVector2i(animationName string) {
	anim.Animations[animationName] = &TweenAnimation[Vector2i]{
		KeyframeValues: make([]TweenValue[Vector2i], 0),
		timeStepFactor: 0.0,
	}
}

// func NewTweenAnimationUint8() TweenAnimation[uint8] {
// 	return TweenAnimation[uint8]{
// 		KeyframeValues: make([]TweenValue[uint8], 0),
// 		timeStepFactor: 1.0,
// 	}
// }

func (comp *AnimationComponent[T]) RegisterKeyframe(animationName string, timeStep float32, value T) {
	anim := comp.Animations[animationName]
	anim.KeyframeValues = append(anim.KeyframeValues, TweenValue[T]{timeStep: timeStep, value: value})

	lowest := anim.KeyframeValues[0].timeStep
	for _, val := range anim.KeyframeValues {
		if lowest > val.timeStep {
			lowest = val.timeStep
		}
	}

	comp.Animations[animationName] = anim
	comp.Animations[animationName].currentValue = comp.Animations[animationName].KeyframeValues[0].value
}

type TweenAnimatorSystemFloat32 struct {
	EcsSystemImpl
}

func (ks *TweenAnimatorSystemFloat32) Update(dt float32) {
	EachEntity(AnimationComponent[float32]{}, func(entity *EcsEntity, a interface{}) {
		anims := a.(AnimationComponent[float32])
		for _, tween := range anims.Animations {
			if !tween.IsPlaying() || tween.HasFinished {
				continue
			}
			tween.CurrentTimestep += dt * tween.timeStepFactor
			tween.currentValue = LerpFloat32(tween.KeyframeValues[tween.currentIndex].value, tween.KeyframeValues[tween.currentIndex+1].value, (tween.CurrentTimestep-tween.KeyframeValues[tween.currentIndex].timeStep)/(tween.KeyframeValues[tween.currentIndex+1].timeStep-tween.KeyframeValues[tween.currentIndex].timeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues[tween.currentIndex+1].timeStep {
				tween.currentIndex++
				if tween.currentIndex == len(tween.KeyframeValues)-1 {
					tween.currentIndex = 0
					if tween.Loop {
						tween.CurrentTimestep = tween.KeyframeValues[0].timeStep
					} else {
						tween.HasFinished = true
					}
				}
			}
		}
		WriteComponent(ks.GetEcsEngine(), entity, anims)
		// passingTime := tween.KeyframeValues[tween.currentIndex+1].timeStep - tween.KeyframeValues[tween.currentIndex].timeStep
	})
}

type TweenAnimatorSystemInt struct {
	EcsSystemImpl
}

func (ks *TweenAnimatorSystemInt) Update(dt float32) {
	EachEntity(AnimationComponent[int]{}, func(entity *EcsEntity, a interface{}) {
		anims := a.(AnimationComponent[int])
		for _, tween := range anims.Animations {

			// passingTime := tween.KeyframeValues[tween.currentIndex+1].timeStep - tween.KeyframeValues[tween.currentIndex].timeStep
			tween.CurrentTimestep += dt * tween.timeStepFactor
			tween.currentValue = LerpInt(tween.KeyframeValues[tween.currentIndex].value, tween.KeyframeValues[tween.currentIndex+1].value, (tween.CurrentTimestep-tween.KeyframeValues[tween.currentIndex].timeStep)/(tween.KeyframeValues[tween.currentIndex+1].timeStep-tween.KeyframeValues[tween.currentIndex].timeStep))
			if tween.CurrentTimestep >= tween.KeyframeValues[tween.currentIndex+1].timeStep {
				tween.currentIndex++
				if tween.currentIndex == len(tween.KeyframeValues)-1 {
					tween.currentIndex = 0
					tween.CurrentTimestep = tween.KeyframeValues[0].timeStep
				}
			}
		}
		WriteComponent(ks.GetEcsEngine(), entity, anims)
	})
}

type TweenAnimatorSystemVector2i struct {
	EcsSystemImpl
}

func (va *TweenAnimatorSystemVector2i) Update(dt float32) {
	EachEntity(AnimationComponent[Vector2i]{}, func(entity *EcsEntity, a interface{}) {
		anims := a.(AnimationComponent[Vector2i])
		for _, tween := range anims.Animations {
			tween.CurrentTimestep += dt * tween.timeStepFactor

			if tween.CurrentTimestep >= tween.KeyframeValues[tween.currentIndex+1].timeStep {
				tween.currentIndex++
				if tween.currentIndex == len(tween.KeyframeValues)-1 {
					tween.currentIndex = 0
					tween.CurrentTimestep = tween.KeyframeValues[0].timeStep
				}
				tween.currentValue = tween.KeyframeValues[tween.currentIndex].value
			}
		}
		WriteComponent(va.GetEcsEngine(), entity, anims)
	})
}

// type TweenAnimatorSystemUint8 struct {
// 	EcsSystemImpl
// }

// func (ks *TweenAnimatorSystemUint8) Update(dt float32) {
// 	EachEntity(TweenAnimation[uint8]{}, func(entity *EcsEntity, a interface{}) {
// 		tween := a.(TweenAnimation[uint8])

//			// passingTime := tween.KeyframeValues[tween.currentIndex+1].timeStep - tween.KeyframeValues[tween.currentIndex].timeStep
//			tween.CurrentTimestep += dt * tween.timeStepFactor
//			tween.currentValue = LerpUint8(tween.KeyframeValues[tween.currentIndex].value, tween.KeyframeValues[tween.currentIndex+1].value, (tween.CurrentTimestep-tween.KeyframeValues[tween.currentIndex].timeStep)/(tween.KeyframeValues[tween.currentIndex+1].timeStep-tween.KeyframeValues[tween.currentIndex].timeStep))
//			if tween.CurrentTimestep >= tween.KeyframeValues[tween.currentIndex+1].timeStep {
//				tween.currentIndex++
//				if tween.currentIndex == len(tween.KeyframeValues)-1 {
//					tween.currentIndex = 0
//					tween.CurrentTimestep = tween.KeyframeValues[0].timeStep
//				}
//			}
//			WriteComponent(ks.GetEcsEngine(), entity, tween)
//		})
//	}

type SpriteAnimation struct {
	Component
	CurrentAnimation string
	StartingSprite   Vector2i
}

func (t *SpriteAnimation) ComponentSet(val interface{}) {
	*t = val.(SpriteAnimation)
}

type SpriteAnimationSystem struct {
	EcsSystemImpl
	TileSet     TileSet
	Sprites     *SpriteBatch
	SpriteScale float32
	Offset      Vector2f
}

func (sa *SpriteAnimationSystem) Update(dt float32) {
	EachEntity(SpriteAnimation{}, func(entity *EcsEntity, a interface{}) {
		spAnim := a.(SpriteAnimation)
		anim := AnimationComponent[Vector2i]{}
		ReadComponent(sa.GetEcsEngine(), entity, &anim)

		_animValue := anim.GetCurrentValue(spAnim.CurrentAnimation)
		_uv1 := NewVector2f(0.0, 0.0)
		_uv1.X = float32(_animValue.X) / float32(sa.TileSet.totalColumns)
		_uv1.Y = float32(_animValue.Y) / float32(sa.TileSet.totalRows)

		_uv2 := NewVector2f(0.0, 0.0)
		_uv2.X = _uv1.X + float32(sa.TileSet.spriteWidth)/float32(sa.TileSet.texture.Width)
		_uv2.Y = _uv1.Y + float32(sa.TileSet.spriteHeight)/float32(sa.TileSet.texture.Height)

		sa.Sprites.DrawSpriteOriginScaledRotated(entity.Pos.Add(sa.Offset).Rotate(entity.Rot, entity.Pos), _uv1, _uv2, sa.SpriteScale, &sa.TileSet.texture, WHITE, entity.Rot)
	})

}
