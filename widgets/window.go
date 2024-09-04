package widgets

import (
	"github.com/asragi/yasoba-prototype/core"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type WindowOption struct {
	Image      *ebiten.Image
	CornerSize int
	Size       *core.Vector
}

type DrawWindowFunc func(*ebiten.Image)

type Window struct {
	Draw DrawWindowFunc
}

func NewWindow(option *WindowOption) *Window {
	img := option.Image
	type rect struct {
		x0 int
		y0 int
		x1 int
		y1 int
	}
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

	drawWindow := func(screen *ebiten.Image) {
		for i, v := range corners {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(cornerPosition[i].X, cornerPosition[i].Y)
			subImage := img.SubImage(image.Rect(v.x0, v.y0, v.x1, v.y1)).(*ebiten.Image)
			screen.DrawImage(subImage, op)
		}
	}

	return &Window{
		Draw: drawWindow,
	}
}
