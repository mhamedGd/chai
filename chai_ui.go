package chai

import (
	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

type UIInteractionsComponent struct {
	InteractionBox Vector2f
	OnCursorEnter  customtypes.ChaiEvent1[EntId]
	OnCursorExit   customtypes.ChaiEvent1[EntId]
	OnClick        customtypes.ChaiEvent1[EntId]
	OnRelease      customtypes.ChaiEvent1[EntId]
	m_JustEntered  bool
	Disabled       bool
}

func UIInteractionSystem(_thisScene *Scene, _dt float32) {
	Iterate2[VisualTransform, UIInteractionsComponent](func(i EntId, t *VisualTransform, uc *UIInteractionsComponent) {
		if uc.Disabled {
			return
		}
		if PointVsRect(GetMouseScreenPosition(), t.Position.Subtract(uc.InteractionBox.Scale(0.5)), t.Position.Add(uc.InteractionBox.Scale(0.5))) {
			if !uc.m_JustEntered {
				uc.OnCursorEnter.Invoke(i)
				uc.m_JustEntered = true
			}
			if IsMouseJustPressed() || IsJustTouched(1) {
				uc.OnClick.Invoke(i)
			} else if IsMouseJustReleased() || IsJustTouchReleased(1) {
				uc.OnRelease.Invoke(i)
			}
		} else {
			if uc.m_JustEntered {
				uc.OnCursorExit.Invoke(i)
				uc.m_JustEntered = false
			}
		}
	})
}
