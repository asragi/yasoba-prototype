package drawing

import (
	"image"

	"github.com/asragi/yasoba-prototype/frontend"
)

type Image interface {
	Bounds() image.Rectangle
	SubImage(image.Rectangle) Image
	DrawImage(Image, *DrawOptions)
}

type DrawOptions struct {
	opacity  float64
	position *frontend.Vector
}

func NewDrawOptions() *DrawOptions {
	return &DrawOptions{
		opacity:  1.0,
		position: &frontend.Vector{X: 0, Y: 0},
	}
}

func (o *DrawOptions) Translate(x, y float64) {
	o.position.X += x
	o.position.Y += y
}

func (o *DrawOptions) SetScale(sx, sy float64) {
	// noop for now
	_ = sx
	_ = sy
}

type DrawArgFunc func(Image)
type DrawFunc func(DrawArgFunc, Depth)
type DrawEndFunc func(Image)
