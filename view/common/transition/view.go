package transition

import (
	"image/color"

	"github.com/asragi/yasoba-prototype/drawing"
)

type View struct {
	overlay OverlayImage
	depth   drawing.Depth
}

type NewTransitionViewFunc func() *View

type OverlayImage interface {
	Fill(color color.Color)
	Draw(target drawing.Image, op *drawing.DrawOptions)
}

type CreateImageFunc func(width, height int) OverlayImage

func CreateNewView(
	width, height int, depth drawing.Depth, createImage CreateImageFunc,
) NewTransitionViewFunc {
	if width <= 0 {
		panic("transition: width must be positive")
	}
	if height <= 0 {
		panic("transition: height must be positive")
	}
	if createImage == nil {
		panic("transition: createImage must not be nil")
	}
	return func() *View {
		img := createImage(width, height)
		if img == nil {
			panic("transition: createImage returned nil")
		}
		img.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		return &View{
			overlay: img,
			depth:   depth,
		}
	}
}

func (v *View) Draw(drawFunc drawing.DrawFunc, rate float64) {
	if v == nil || v.overlay == nil {
		return
	}
	if drawFunc == nil {
		return
	}
	alpha := clampRate(rate)
	if alpha <= 0 {
		return
	}

	drawFunc(func(screen drawing.Image) {
		op := drawing.NewDrawOptions()
		op.SetOpacity(alpha)
		v.overlay.Draw(screen, op)
	}, v.depth)
}

func clampRate(rate float64) float64 {
	if rate < 0 {
		return 0
	}
	if rate > 1 {
		return 1
	}
	return rate
}
