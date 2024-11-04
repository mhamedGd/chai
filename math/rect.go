package math

// Position Bottom-Left
type Rect struct {
	Position, Size Vector2f
}

func (r *Rect) ContainsPoint(_p Vector2f) bool {
	return _p.X >= r.Position.X && _p.X <= (r.Size.X+r.Position.X) && _p.Y >= r.Position.Y && _p.Y <= (r.Size.Y+r.Position.Y)
}

func (r *Rect) ContainsRect(_cr Rect) bool {
	return _cr.Position.X >= r.Position.X && _cr.Position.Y >= r.Position.Y &&
		(_cr.Size.X+_cr.Position.X) < (r.Size.X+r.Position.X) && (_cr.Size.Y+_cr.Position.Y) < (r.Size.Y+r.Position.Y)
}

func (r *Rect) OverlapsRect(_cr Rect) bool {
	return r.Position.X < _cr.Position.X+_cr.Size.X && r.Position.X+r.Size.X >= _cr.Position.X && r.Position.Y < _cr.Position.Y+_cr.Size.Y && r.Position.Y+r.Size.Y >= _cr.Position.Y
}

func PointVsRect(_point Vector2f, _lowerLeft, _upperRight Vector2f) bool {
	return (_point.X >= _lowerLeft.X && _point.Y >= _lowerLeft.Y && _point.X <= _upperRight.X && _point.Y <= _upperRight.Y)
}
