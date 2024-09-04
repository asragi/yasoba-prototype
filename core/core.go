package core

type Vector struct {
	X, Y float64
}

type Pivot struct {
	X, Y float64
}

func (p *Pivot) ApplyToSize(size *Vector) *Vector {
	return &Vector{
		X: size.X * p.X,
		Y: size.Y * p.Y,
	}
}

var (
	PivotTopLeft     = &Pivot{0, 0}
	PivotTopCenter   = &Pivot{0.5, 0}
	PivotTopRight    = &Pivot{1, 0}
	PivotBottomLeft  = &Pivot{0, 1}
	PivotBottomRight = &Pivot{1, 1}
)
