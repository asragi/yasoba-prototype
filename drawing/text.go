package drawing

import "image/color"

type textFace interface{}

type DrawTextFunc func(Image, string, textFace, *TextDrawOptions)

type TextDrawOptions struct {
	Position     *Vector
	Scale        *Vector
	OutlineColor color.Color
}

func NewTextDrawOptions() *TextDrawOptions {
	return &TextDrawOptions{
		Position: &Vector{X: 0, Y: 0},
		Scale:    &Vector{X: 1, Y: 1},
	}
}

func (o *TextDrawOptions) Translate(x, y float64) {
	o.Position.X += x
	o.Position.Y += y
}

func (o *TextDrawOptions) SetScale(sx, sy float64) {
	o.Scale.X = sx
	o.Scale.Y = sy
}

func (o *TextDrawOptions) SetColorScale(color color.Color) {
	o.OutlineColor = color
}
