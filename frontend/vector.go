package frontend

import "fmt"

type Vector struct {
	X, Y float64
}

var (
	VectorZero = &Vector{0, 0}
)

func (v *Vector) Add(v2 *Vector) *Vector {
	return &Vector{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v *Vector) Sub(v2 *Vector) *Vector {
	return &Vector{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v *Vector) Multiply(scale float64) *Vector {
	return &Vector{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func (v *Vector) String() string {
	return fmt.Sprintf("{X: %f, Y: %f}", v.X, v.Y)
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
	PivotCenter      = &Pivot{0.5, 0.5}
	PivotBottomLeft  = &Pivot{0, 1}
	PivotBottomRight = &Pivot{1, 1}
)
