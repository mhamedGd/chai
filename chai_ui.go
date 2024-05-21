package chai

type UIInteractionsComponent struct {
	InteractionBox Vector2f
	OnCursorEnter  ChaiEvent[EntId]
	OnCursorExit   ChaiEvent[EntId]
	OnClick        ChaiEvent[EntId]
	OnRelease      ChaiEvent[EntId]
	justEntered    bool
	Disabled       bool
}

type UIInteractionSystem struct {
	EcsSystem
}

func (uis *UIInteractionSystem) Update(dt float32) {
	Iterate2[Transform, UIInteractionsComponent](func(i EntId, t *Transform, uc *UIInteractionsComponent) {
		if uc.Disabled {
			return
		}
		if PointVsRect(GetMouseScreenPosition(), t.Position.Subtract(uc.InteractionBox.Scale(0.5)), t.Position.Add(uc.InteractionBox.Scale(0.5))) {
			if !uc.justEntered {
				uc.OnCursorEnter.Invoke(i)
				uc.justEntered = true
			}
			if IsMouseJustPressed() || IsJustTouched(1) {
				uc.OnClick.Invoke(i)
			} else if IsMouseJustReleased() || IsJustTouchReleased(1) {
				uc.OnRelease.Invoke(i)
			}
		} else {
			if uc.justEntered {
				uc.OnCursorExit.Invoke(i)
				uc.justEntered = false
			}
		}
	})
}
