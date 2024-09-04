package widgets

import (
	"errors"
	"github.com/asragi/yasoba-prototype/core"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type WindowOption struct {
	Image            *ebiten.Image
	CornerSize       int
	RelativePosition *core.Vector
	Size             *core.Vector
	Depth            core.Depth
	Pivot            core.Pivot
}

func (o *WindowOption) Validation() error {
	if o.Image == nil {
		return errors.New("image is required")
	}
	if o.CornerSize <= 0 {
		return errors.New("corner size must be greater than 0")
	}
	if o.Size == nil {
		return errors.New("size is required")
	}
	if o.Depth == core.Zero {
		return errors.New("depth is required")
	}
	return nil
}

type Window struct {
	MoveTo func(relativePosition *core.Vector)
	Update func(parentPosition *core.Vector)
	Draw   func(core.DrawFunc)
}

func NewWindow(option *WindowOption) *Window {
	if err := option.Validation(); err != nil {
		panic(err)
	}
	img := option.Image
	relativePosition := option.RelativePosition
	type rect struct {
		x0 int
		y0 int
		x1 int
		y1 int
	}
	pivotDiff := option.Pivot.ApplyToSize(option.Size)
	corners := []*rect{
		{0, 0, option.CornerSize, option.CornerSize},
		{img.Bounds().Dx() - option.CornerSize, 0, img.Bounds().Dx(), option.CornerSize},
		{0, img.Bounds().Dy() - option.CornerSize, option.CornerSize, img.Bounds().Dy()},
		{
			img.Bounds().Dx() - option.CornerSize,
			img.Bounds().Dy() - option.CornerSize,
			img.Bounds().Dx(),
			img.Bounds().Dy(),
		},
	}
	cornerPosition := []*core.Vector{
		{0, 0},
		{option.Size.X - float64(option.CornerSize), 0},
		{0, option.Size.Y - float64(option.CornerSize)},
		{option.Size.X - float64(option.CornerSize), option.Size.Y - float64(option.CornerSize)},
	}

	var parentPosition *core.Vector
	moveTo := func(passedPosition *core.Vector) {
		relativePosition = passedPosition
	}

	update := func(passedPosition *core.Vector) {
		parentPosition = passedPosition
	}

	drawWindow := func(drawing core.DrawFunc) {
		for i, v := range corners {
			op := &ebiten.DrawImageOptions{}
			x := cornerPosition[i].X + relativePosition.X - pivotDiff.X + parentPosition.X
			y := cornerPosition[i].Y + relativePosition.Y - pivotDiff.Y + parentPosition.Y
			op.GeoM.Translate(x, y)
			subImage := img.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1)).(*ebiten.Image)
			drawing(
				func(screen *ebiten.Image) {
					screen.DrawImage(subImage, op)
				}, option.Depth,
			)
		}
	}

	return &Window{
		MoveTo: moveTo,
		Update: update,
		Draw:   drawWindow,
	}
}
