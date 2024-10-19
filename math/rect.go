package math

// Position Bottom-Left
type Rect struct {
	Position, Size Vector2f
}

func (r *Rect) ContainsPoint(p Vector2f) bool {
	return p.X >= r.Position.X && p.X <= (r.Size.X+r.Position.X) && p.Y >= r.Position.Y && p.Y <= (r.Size.Y+r.Position.Y)
}

func (r *Rect) ContainsRect(cr Rect) bool {
	return cr.Position.X >= r.Position.X && cr.Position.Y >= r.Position.Y &&
		(cr.Size.X+cr.Position.X) < (r.Size.X+r.Position.X) && (cr.Size.Y+cr.Position.Y) < (r.Size.Y+r.Position.Y)
}

func (r *Rect) OverlapsRect(cr Rect) bool {
	return r.Position.X < cr.Position.X+cr.Size.X && r.Position.X+r.Size.X >= cr.Position.X && r.Position.Y < cr.Position.Y+cr.Size.Y && r.Position.Y+r.Size.Y >= cr.Position.Y
}

func PointVsRect(_point Vector2f, lower_left, upper_right Vector2f) bool {
	return (_point.X >= lower_left.X && _point.Y >= lower_left.Y && _point.X <= upper_right.X && _point.Y <= upper_right.Y)
}
