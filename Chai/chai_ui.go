package chai

type UIInteractionsComponent struct {
	Component
	InteractionBox Vector2f
	OnCursorEnter  ChaiEvent[*EcsEntity]
	OnCursorExit   ChaiEvent[*EcsEntity]
	OnClick        ChaiEvent[*EcsEntity]
	OnRelease      ChaiEvent[*EcsEntity]
	justEntered    bool
}

func (t *UIInteractionsComponent) ComponentSet(val interface{}) { *t = val.(UIInteractionsComponent) }

type UIInteractionSystem struct {
	EcsSystemImpl
}

func (uis *UIInteractionSystem) Update(dt float32) {
	EachEntity(UIInteractionsComponent{}, func(entity *EcsEntity, a interface{}) {
		uiInter := a.(UIInteractionsComponent)
		if PointVsRect(GetMouseScreenPosition(), entity.Pos.Subtract(uiInter.InteractionBox.Scale(0.5)), entity.Pos.Add(uiInter.InteractionBox.Scale(0.5))) {
			if !uiInter.justEntered {
				uiInter.OnCursorEnter.Invoke(entity)
				uiInter.justEntered = true
			}
			if IsMousejustPressed() || IsJustTouched(1) {
				uiInter.OnClick.Invoke(entity)
			} else if IsMouseJustReleased() || IsJustTouchReleased(1) {
				uiInter.OnRelease.Invoke(entity)
				LogF("Just Released")
			}
		} else {
			if uiInter.justEntered {
				uiInter.OnCursorExit.Invoke(entity)
				uiInter.justEntered = false
			}
		}
		WriteComponent(uis.GetEcsEngine(), entity, uiInter)
	})
}
