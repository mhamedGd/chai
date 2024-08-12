package chai

type UIInteractionsComponent struct {
	InteractionBox Vector2f
	OnCursorEnter  ChaiEvent1[EntId]
	OnCursorExit   ChaiEvent1[EntId]
	OnClick        ChaiEvent1[EntId]
	OnRelease      ChaiEvent1[EntId]
	justEntered    bool
	Disabled       bool
}

func UIInteractionSystem(_this_scene *Scene, _dt float32) {
	Iterate2[VisualTransform, UIInteractionsComponent](func(i EntId, t *VisualTransform, uc *UIInteractionsComponent) {
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
